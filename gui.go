package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/image/bmp" // Registra o decodificador de arquivos BMP

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var currentCategory = "MSX Stuffs"

var (
	msxManiaTimerCancel chan struct{}
	msxManiaTimerMutex  sync.Mutex
	debugMode           bool
)

func logDebug(format string, args ...any) {
	if debugMode {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("[%s] [DEBUG] %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	}
}

func stopMSXManiaSlideshow() {
	msxManiaTimerMutex.Lock()
	defer msxManiaTimerMutex.Unlock()
	if msxManiaTimerCancel != nil {
		close(msxManiaTimerCancel)
		msxManiaTimerCancel = nil
	}
}

func findImageFile(dir, categoryDir, name string) string {
	if dir == "" {
		return ""
	}
	// Strip suffix like " %A", " %B", etc. at the end
	if len(name) >= 3 && name[len(name)-2] == '%' {
		if name[len(name)-3] == ' ' {
			name = name[:len(name)-3]
		}
	}
	logDebug("findImageFile: Buscando imagem para '%s' em '%s/%s'", name, dir, categoryDir)
	// 1. Tenta correspondência exata primeiro (rápido)
	path1 := filepath.Join(dir, categoryDir, name+".png")
	if _, err := os.Stat(path1); err == nil {
		logDebug("findImageFile: Encontrado (exato, minúsculo) em '%s'", path1)
		return path1
	}
	path2 := filepath.Join(dir, categoryDir, name+".PNG")
	if _, err := os.Stat(path2); err == nil {
		logDebug("findImageFile: Encontrado (exato, maiúsculo) em '%s'", path2)
		return path2
	}

	// 2. Se não achar, faz uma busca inteligente no diretório por nomes similares
	searchDir := filepath.Join(dir, categoryDir)
	files, err := os.ReadDir(searchDir)
	if err != nil {
		logDebug("findImageFile: Não foi possível ler diretório '%s': %v", searchDir, err)
		return ""
	}

	removeParenthesesAndBrackets := func(s string) string {
		var sb strings.Builder
		inParen := 0
		inBracket := 0
		for _, r := range s {
			if r == '(' {
				inParen++
			} else if r == '[' {
				inBracket++
			} else if r == ')' {
				if inParen > 0 {
					inParen--
				}
			} else if r == ']' {
				if inBracket > 0 {
					inBracket--
				}
			} else if inParen == 0 && inBracket == 0 {
				sb.WriteRune(r)
			}
		}
		return sb.String()
	}

	normalize := func(s string) string {
		cleaned := removeParenthesesAndBrackets(s)
		var sb strings.Builder
		for _, r := range strings.ToLower(cleaned) {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				sb.WriteRune(r)
			}
		}
		return sb.String()
	}

	normTarget := normalize(name)
	if normTarget == "" {
		return ""
	}
	logDebug("findImageFile: Busca exata falhou. Tentando busca similar para termo normalizado '%s'...", normTarget)

	// Primeiro critério: começa com o nome procurado (normalizado)
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		ext := filepath.Ext(entryName)
		if strings.ToLower(ext) != ".png" {
			continue
		}
		nameWithoutExt := strings.TrimSuffix(entryName, ext)
		normEntry := normalize(nameWithoutExt)

		if strings.HasPrefix(normEntry, normTarget) {
			matchedPath := filepath.Join(searchDir, entryName)
			logDebug("findImageFile: Encontrado similar por prefixo: '%s' -> '%s'", entryName, matchedPath)
			return matchedPath
		}
	}

	// Segundo critério: contém o nome procurado (normalizado)
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		ext := filepath.Ext(entryName)
		if strings.ToLower(ext) != ".png" {
			continue
		}
		nameWithoutExt := strings.TrimSuffix(entryName, ext)
		normEntry := normalize(nameWithoutExt)

		if strings.Contains(normEntry, normTarget) {
			matchedPath := filepath.Join(searchDir, entryName)
			logDebug("findImageFile: Encontrado similar por conter termo: '%s' -> '%s'", entryName, matchedPath)
			return matchedPath
		}
	}

	logDebug("findImageFile: Nenhuma imagem similar encontrada para '%s' em '%s'", name, searchDir)
	return ""
}

func runGUI() {
	CurrentLanguage = getConfigWithFallback("language", "en")
	debugMode = getConfigWithFallback("debug", "false") == "true"
	initOutputCapture()
	logDebug("Modo debug ativado! Inicializando interface gráfica...")
	fmt.Println(T("log_started", time.Now().Format("2006-01-02 15:04:05")))

	// Inicializa tabelas de emulação específicas do jogo se necessário
	_ = InitGameEmulacaoTables()

	// Inicializa a aplicação Fyne
	myApp := app.NewWithID("com.nemesis.msxstuffs")

	// Carrega e aplica o ícone do programa (barra de tarefas e janela principal)
	if iconRes, err := fyne.LoadResourceFromPath("images/Icon.png"); err == nil {
		myApp.SetIcon(iconRes)
	}

	// Carrega e aplica o tema salvo
	savedTheme := myApp.Preferences().StringWithFallback("theme", string(ThemeOneDark))
	myApp.Settings().SetTheme(GetTheme(savedTheme))

	myWindow := myApp.NewWindow(T("title_main"))

	// Barra de Status inferior
	statusBar := widget.NewLabel(T("status_theme_active", savedTheme))
	statusBar.Alignment = fyne.TextAlignLeading

	// Ações dos Menus
	sairItem := fyne.NewMenuItem(T("menu_exit"), func() {
		logDebug("Menu 'Sair' acionado. Encerrando aplicação...")
		myApp.Quit()
	})
	sairItem.Icon = theme.LogoutIcon()

	statusItem := fyne.NewMenuItem(T("menu_status"), func() {
		logDebug("Menu 'Status/Logs' acionado. Abrindo janela de status...")
		showStatusWindow(myApp)
	})
	statusItem.Icon = theme.InfoIcon()

	programaItem := fyne.NewMenuItem(T("menu_program"), func() {
		logDebug("Menu 'Configuração -> Programa' acionado. Abrindo janela de configurações...")
		showSettings(myApp, statusBar)
	})
	programaItem.Icon = theme.SettingsIcon()

	sobreItem := fyne.NewMenuItem(T("menu_about"), func() {
		logDebug("Menu 'Ajuda -> Sobre' acionado. Exibindo janela informativa...")
		dialog.ShowInformation(
			T("about_title"),
			T("about_text"),
			myWindow,
		)
	})
	sobreItem.Icon = theme.InfoIcon()

	// Menus principais
	arquivoMenu := fyne.NewMenu(T("menu_file"), statusItem, fyne.NewMenuItemSeparator(), sairItem)
	configMenu := fyne.NewMenu(T("menu_config"), programaItem)
	ajudaMenu := fyne.NewMenu(T("menu_help"), sobreItem)

	// Cria o menu superior
	mainMenu := fyne.NewMainMenu(arquivoMenu, configMenu, ajudaMenu)
	myWindow.SetMainMenu(mainMenu)

	// Inicialização do estado local
	volInicialStr := getConfigWithFallback("volume_inicial", "1")
	var currentVolume int
	fmt.Sscanf(volInicialStr, "%d", &currentVolume)
	if currentVolume < 1 || currentVolume > 10 {
		currentVolume = 1
	}

	var selectedGame *Game
	var filteredGames []Game

	var listWidget *widget.List
	rightContainer := container.NewStack()

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder(T("search_placeholder"))

	// Função para atualizar a lista filtrada de jogos
	reloadGames := func() {
		logDebug("reloadGames: Iniciando busca com categoria='%s', volume=%d, termo de busca='%s'", currentCategory, currentVolume, searchEntry.Text)
		games, err := GetGamesByVolume(currentCategory, currentVolume, searchEntry.Text)
		if err != nil {
			logDebug("reloadGames: Erro na busca do banco: %v", err)
			filteredGames = []Game{}
			if listWidget != nil {
				listWidget.Refresh()
			}
			return
		}
		logDebug("reloadGames: Busca concluída. Registros retornados: %d", len(games))
		filteredGames = games
		if listWidget != nil {
			listWidget.Refresh()
			listWidget.UnselectAll()
		}
		selectedGame = nil
	}

	var updateGameDetails func()
 
	updateGameDetails = func() {
		stopMSXManiaSlideshow()
		if selectedGame == nil {
			logDebug("updateGameDetails: Nenhum jogo selecionado. Limpando painel de detalhes.")
			// Mostra a área direita totalmente limpa quando não há jogo selecionado
			rightContainer.Objects = []fyne.CanvasObject{}
			rightContainer.Refresh()
			return
		}
		logDebug("updateGameDetails: Exibindo detalhes para o jogo '%s' (CD %d, Arquivo: %s)", selectedGame.Descricao, selectedGame.CdNumero, selectedGame.Disco)
 
		var detailImg fyne.CanvasObject
		if currentCategory == "MSX Mania" || currentCategory == "Good MSX 1" {
			rootDir, err := filepath.Abs(".")
			if err != nil {
				rootDir = "."
			}
			msxmaniaPicsDir := getConfigWithFallback("msxmania_pictures", filepath.Join(rootDir, "pictures", "msxmania", "MSX"))
			logDebug("updateGameDetails: Buscando imagens em '%s'", msxmaniaPicsDir)
			titlePath := findImageFile(msxmaniaPicsDir, "Title", selectedGame.Descricao)
			snapPath := findImageFile(msxmaniaPicsDir, "Snap", selectedGame.Descricao)
 
			if titlePath != "" && snapPath != "" {
				logDebug("updateGameDetails: Encontrado titlePath='%s' e snapPath='%s'. Iniciando alternância de 2 segundos.", titlePath, snapPath)
				img := canvas.NewImageFromFile(titlePath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
 
				// Inicia a alternância das imagens
				msxManiaTimerMutex.Lock()
				cancelChan := make(chan struct{})
				msxManiaTimerCancel = cancelChan
				msxManiaTimerMutex.Unlock()
 
				go func(img *canvas.Image, title, snap string, cancel chan struct{}) {
					isTitle := true
					ticker := time.NewTicker(2 * time.Second)
					defer ticker.Stop()
					for {
						select {
						case <-ticker.C:
							isTitle = !isTitle
							if isTitle {
								logDebug("Slideshow: Alternando para tela de título: '%s'", title)
								img.File = title
							} else {
								logDebug("Slideshow: Alternando para tela de jogo: '%s'", snap)
								img.File = snap
							}
							img.Refresh()
						case <-cancel:
							logDebug("Slideshow: Alternador parado/cancelado.")
							return
						}
					}
				}(img, titlePath, snapPath, cancelChan)
			} else if titlePath != "" {
				logDebug("updateGameDetails: Apenas tela de título encontrada: '%s'", titlePath)
				img := canvas.NewImageFromFile(titlePath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else if snapPath != "" {
				logDebug("updateGameDetails: Apenas tela de jogo (snap) encontrada: '%s'", snapPath)
				img := canvas.NewImageFromFile(snapPath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else if box2dPath := findImageFile(msxmaniaPicsDir, "2DBoxes", selectedGame.Descricao); box2dPath != "" {
				logDebug("updateGameDetails: Encontrada 2D Box: '%s'", box2dPath)
				img := canvas.NewImageFromFile(box2dPath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else if box2dzxPath := findImageFile(msxmaniaPicsDir, "2DBoxeszx", selectedGame.Descricao); box2dzxPath != "" {
				logDebug("updateGameDetails: Encontrada 2D Box ZX: '%s'", box2dzxPath)
				img := canvas.NewImageFromFile(box2dzxPath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else if box3dPath := findImageFile(msxmaniaPicsDir, "3DBoxes", selectedGame.Descricao); box3dPath != "" {
				logDebug("updateGameDetails: Encontrada 3D Box: '%s'", box3dPath)
				img := canvas.NewImageFromFile(box3dPath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else if logosPath := findImageFile(msxmaniaPicsDir, "Logos", selectedGame.Descricao); logosPath != "" {
				logDebug("updateGameDetails: Encontrada Logos: '%s'", logosPath)
				img := canvas.NewImageFromFile(logosPath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else if logoPath := findImageFile(msxmaniaPicsDir, "Logo", selectedGame.Descricao); logoPath != "" {
				logDebug("updateGameDetails: Encontrada Logo: '%s'", logoPath)
				img := canvas.NewImageFromFile(logoPath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else {
				logDebug("updateGameDetails: Nenhuma imagem de jogo encontrada para '%s' nos diretórios do MSX Mania.", selectedGame.Descricao)
				spacer := container.NewGridWrap(fyne.NewSize(380, 280))
				lbl := widget.NewLabelWithStyle(T("no_screenshot"), fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
				detailImg = container.NewStack(spacer, container.NewCenter(lbl))
			}
		} else {
			picsPath := getConfigWithFallback("pictures", "pictures")
			dirName := fmt.Sprintf("%02d", selectedGame.CdNumero)
			imagePath := filepath.Join(picsPath, dirName, selectedGame.Raiz+".BMP")
			logDebug("updateGameDetails: Buscando screenshot em '%s'", imagePath)
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				imagePath = filepath.Join(picsPath, dirName, selectedGame.Raiz+".bmp")
				logDebug("updateGameDetails: BMP maiúsculo não encontrado. Tentando minúsculo: '%s'", imagePath)
			}
 
			if _, err := os.Stat(imagePath); err == nil {
				logDebug("updateGameDetails: Screenshot encontrada: '%s'", imagePath)
				img := canvas.NewImageFromFile(imagePath)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(380, 280))
				detailImg = img
			} else {
				logDebug("updateGameDetails: Nenhuma screenshot encontrada para '%s' (%s).", selectedGame.Descricao, imagePath)
				spacer := container.NewGridWrap(fyne.NewSize(380, 280))
				lbl := widget.NewLabelWithStyle(T("no_screenshot"), fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
				detailImg = container.NewStack(spacer, container.NewCenter(lbl))
			}
		}
 
		playBtn := widget.NewButtonWithIcon(T("btn_play"), theme.MediaPlayIcon(), func() {
			logDebug("Ação: Botão 'Jogar' (Play) clicado para o jogo '%s'", selectedGame.Descricao)
			launchGame(selectedGame, myWindow)
		})
		playBtn.Importance = widget.HighImportance
 
		reconfigBtn := widget.NewButtonWithIcon(T("btn_reconfig"), theme.SettingsIcon(), func() {
			logDebug("Ação: Botão 'Reconfigurar' clicado para o jogo '%s'", selectedGame.Descricao)
			showGameConfig(selectedGame, myWindow, myApp)
		})
 
		// Espaçadores para empurrar a imagem (ajuste fino de alinhamento)
		var topOffset float32 = 135
		var leftOffset float32 = 109
		if currentCategory == "MSX Mania" || currentCategory == "Good MSX 1" {
			topOffset = 172  // desce mais 2 pixels (total 37)
			leftOffset = 116 // joga mais 2 pixels para a direita (total 7)
		}
		imgSpacerTop := container.NewGridWrap(fyne.NewSize(1, topOffset))
		imgSpacerLeft := container.NewGridWrap(fyne.NewSize(leftOffset, 1))
 
		imgLayout := container.NewBorder(
			nil,
			nil,
			imgSpacerLeft,
			nil,
			detailImg,
		)
 
		// Espaçador para empurrar os botões um pouco mais para baixo
		playBtnSpacerTop := container.NewGridWrap(fyne.NewSize(1, 20))
		// Define o tamanho de cada botão para mantê-los proporcionais
		playBtnSized := container.NewGridWrap(fyne.NewSize(140, 38), playBtn)
		reconfigBtnSized := container.NewGridWrap(fyne.NewSize(140, 38), reconfigBtn)
		
		// Agrupa os botões lado a lado no centro
		buttonsRow := container.NewHBox(playBtnSized, reconfigBtnSized)
		buttonsContainer := container.NewCenter(buttonsRow)
 
		detailsVBox := container.NewVBox(
			imgSpacerTop,
			imgLayout,
			playBtnSpacerTop,
			buttonsContainer,
		)
 
		scrollContainer := container.NewScroll(detailsVBox)
		rightContainer.Objects = []fyne.CanvasObject{scrollContainer}
		rightContainer.Refresh()
	}

	// Linha de 10 botões (Seletor Segmentado Customizado)
	volumeButtons := make([]*widget.Button, 10)
	var updateVolumeSelection func(int)

	updateVolumeSelection = func(selectedVol int) {
		for idx, btn := range volumeButtons {
			if currentCategory == "MSX Mania" || currentCategory == "Good MSX 1" {
				btn.Disable()
			} else {
				btn.Enable()
				if idx+1 == selectedVol {
					btn.Importance = widget.HighImportance
				} else {
					btn.Importance = widget.MediumImportance
				}
			}
			btn.Refresh()
		}
	}

	for i := 0; i < 10; i++ {
		volNum := i + 1
		btn := widget.NewButton(fmt.Sprintf("%d", volNum), func() {
			logDebug("Ação: Botão de Volume '%d' selecionado", volNum)
			currentVolume = volNum
			updateVolumeSelection(currentVolume)
			searchEntry.SetText("")
			reloadGames()
			updateGameDetails()
		})
		volumeButtons[i] = btn
	}

	// Agrupa os botões em um grid horizontal de 10 colunas
	segmentedControl := container.NewGridWithColumns(10,
		volumeButtons[0], volumeButtons[1], volumeButtons[2], volumeButtons[3], volumeButtons[4],
		volumeButtons[5], volumeButtons[6], volumeButtons[7], volumeButtons[8], volumeButtons[9],
	)

	listBg := canvas.NewRectangle(color.White)

	listWidget = widget.NewList(
		func() int {
			return len(filteredGames)
		},
		func() fyne.CanvasObject {
			diskLabel := canvas.NewText("[000]", color.RGBA{R: 120, G: 120, B: 120, A: 255})
			diskLabel.TextStyle = fyne.TextStyle{Bold: true}
			diskLabel.TextSize = 13

			nameLabel := canvas.NewText("Nome do Jogo", color.RGBA{R: 20, G: 20, B: 20, A: 255})
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			nameLabel.TextSize = 13

			return container.NewHBox(diskLabel, nameLabel)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(filteredGames) {
				return
			}
			game := filteredGames[id]
			box := item.(*fyne.Container)
			diskText := box.Objects[0].(*canvas.Text)
			nameText := box.Objects[1].(*canvas.Text)

			if currentCategory == "MSX Mania" {
				if game.CdNumero == 131 {
					diskText.Text = "[SP1]"
				} else if game.CdNumero == 132 {
					diskText.Text = "[SP2]"
				} else {
					diskText.Text = fmt.Sprintf("[%03d]", game.CdNumero)
				}
				diskText.Show()
			} else {
				diskText.Hide()
			}

			nameText.Text = game.Descricao

			diskText.Refresh()
			nameText.Refresh()
			box.Refresh()
		},
	)

	listWidget.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(filteredGames) {
			logDebug("Ação: Item selecionado na listagem de jogos (ID: %d, Descrição: '%s')", id, filteredGames[id].Descricao)
			selectedGame = &filteredGames[id]
			updateGameDetails()
		}
	}

	searchEntry.OnChanged = func(val string) {
		logDebug("Ação: Campo de texto de busca modificado para: '%s'", val)
		reloadGames()
	}

	// Inicializa a lista e o controle segmentado
	reloadGames()
	updateVolumeSelection(currentVolume)
	updateGameDetails()

	// Espaçadores externos de margem (topo e esquerda)
	sidebarSpacerTop := container.NewGridWrap(fyne.NewSize(1, 140))
	sidebarSpacerLeft := container.NewGridWrap(fyne.NewSize(15, 1))

	// Título do painel lateral (letra branca)
	sidebarHeader := canvas.NewText(T("sidebar_title"), color.White)
	sidebarHeader.Alignment = fyne.TextAlignCenter
	sidebarHeader.TextStyle = fyne.TextStyle{Bold: true}
	sidebarHeader.TextSize = 14

	// Borda fina cinza ao redor do listbox (combobox) para melhor definição visual
	listBorder := canvas.NewRectangle(color.Transparent)
	listBorder.StrokeColor = color.RGBA{R: 180, G: 180, B: 180, A: 255}
	listBorder.StrokeWidth = 1

	// Carrega as categorias do banco
	cats, err := GetCategories()
	if err != nil {
		cats = []Category{
			{1, "MSX Stuffs", 1},
			{2, "MSX Mania", 0},
			{3, "CAS Collection", 0},
			{4, "Good MSX 1", 1},
			{5, "Wave Games", 0},
			{6, "MSX Tools", 0},
			{7, "Nemesis Diskpack", 0},
		}
	}

	var selectOptions []string
	categoryMap := make(map[string]Category)

	for _, c := range cats {
		label := c.Nome
		if c.Ativo == 0 {
			label = c.Nome + " (" + T("disabled_label") + ")"
		}
		selectOptions = append(selectOptions, label)
		categoryMap[label] = c
	}

	var categorySelect *widget.Select
	var lastSelectedCategory = selectOptions[0]

	var onCategoryChange func(string)
	onCategoryChange = func(selected string) {
		logDebug("Ação: Seleção de categoria modificada para '%s'", selected)
		cat, ok := categoryMap[selected]
		if !ok || cat.Ativo == 0 {
			logDebug("onCategoryChange: Categoria '%s' inativa ou inexistente. Revertendo...", selected)
			displayName := selected
			if ok {
				displayName = cat.Nome
			}
			dialog.ShowInformation(T("msg_disabled_title"), T("msg_category_disabled", displayName), myWindow)

			// Reverte a seleção sem disparar callback recursivamente
			categorySelect.OnChanged = nil
			categorySelect.SetSelected(lastSelectedCategory)
			categorySelect.OnChanged = onCategoryChange
			return
		}
		lastSelectedCategory = selected
		currentCategory = cat.Nome
		reloadGames()
		updateVolumeSelection(currentVolume)
		updateGameDetails()
	}

	categorySelect = widget.NewSelect(selectOptions, onCategoryChange)
	categorySelect.SetSelected(lastSelectedCategory)

	categoryBox := container.NewVBox(
		widget.NewLabelWithStyle(T("label_collection"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		categorySelect,
		container.NewGridWrap(fyne.NewSize(1, 10)), // Espaçamento inferior
	)

	// Agrupa a lista rolável com o fundo branco e a borda cinza
	listContainer := container.NewStack(listBg, listWidget, listBorder)

	listContainerLayout := container.NewBorder(
		nil,
		categoryBox,
		nil,
		nil,
		listContainer,
	)

	// Layout interno da barra lateral
	sidebarLayout := container.NewBorder(
		container.NewVBox(
			sidebarHeader,
			segmentedControl,
			searchEntry,
		),
		nil,
		nil,
		nil,
		listContainerLayout,
	)

	// Adiciona pequeno espaçamento (padding) interno mais justo (6px) nas bordas do painel
	paddedSidebar := container.NewBorder(
		container.NewGridWrap(fyne.NewSize(1, 6)),
		container.NewGridWrap(fyne.NewSize(1, 6)),
		container.NewGridWrap(fyne.NewSize(6, 1)),
		container.NewGridWrap(fyne.NewSize(6, 1)),
		sidebarLayout,
	)

	// Moldura externa (borda) para o painel de navegação
	panelBorder := canvas.NewRectangle(color.Transparent)
	panelBorder.StrokeColor = color.RGBA{R: 160, G: 160, B: 160, A: 255}
	panelBorder.StrokeWidth = 1.5

	// Fundo escuro translúcido para dar destaque e contraste sobre a imagem original
	panelBg := canvas.NewRectangle(color.RGBA{R: 24, G: 28, B: 37, A: 160})

	// Junta fundo, borda e conteúdo
	sidebarFrame := container.NewStack(
		panelBg,
		panelBorder,
		paddedSidebar,
	)

	// Monta o layout da barra lateral aplicando as margens externas (topo e esquerda)
	sidebarLayoutWithMargins := container.NewBorder(
		sidebarSpacerTop,
		nil,
		sidebarSpacerLeft,
		nil,
		sidebarFrame,
	)

	splitPane := container.NewHSplit(
		sidebarLayoutWithMargins,
		rightContainer,
	)
	splitPane.Offset = 0.32

	// Imagem de fundo mantendo o aspecto original sem espremer
	bgImage := canvas.NewImageFromFile("images/MSXStuffs.png")
	bgImage.FillMode = canvas.ImageFillContain

	// O layout central do programa sobrepõe o fundo original diretamente com o splitPane
	centerLayout := container.NewStack(
		bgImage,
		splitPane,
	)

	// Layout principal da janela
	mainLayout := container.NewBorder(
		nil,        // Sem elemento no topo
		statusBar,  // Status bar na parte inferior
		nil,
		nil,
		centerLayout,
	)

	myWindow.SetContent(mainLayout)
	myWindow.Resize(fyne.NewSize(1024, 768))

	splashShown := false

	// Tenta obter o driver de desktop para criar a janela sem bordas (splash)
	if drv, ok := myApp.Driver().(desktop.Driver); ok {
		splashWin := drv.CreateSplashWindow()

		splashImg := canvas.NewImageFromFile("images/splash.png")
		splashImg.FillMode = canvas.ImageFillStretch

		// Retângulo de fade preto por cima da imagem
		overlay := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 0})

		splashWin.SetContent(container.NewStack(splashImg, overlay))
		splashWin.Resize(fyne.NewSize(640, 480))
		splashWin.Show()
		splashShown = true

		go func() {
			// Exibe o splash estático por 2 segundos
			time.Sleep(2 * time.Second)

			// Executa a animação no thread principal
			fyne.Do(func() {
				fadeAnim := fyne.NewAnimation(1000*time.Millisecond, func(percent float32) {
					alpha := uint8(percent * 255)
					overlay.FillColor = color.RGBA{R: 0, G: 0, B: 0, A: alpha}
					overlay.Refresh()

					if percent >= 1.0 {
						myWindow.Show()
						splashWin.Close()
					}
				})
				fadeAnim.Start()
			})
		}()
	}

	if !splashShown {
		myWindow.Show()
	}

	myApp.Run()
}

func getGamePath(game *Game) string {
	if currentCategory == "MSX Mania" {
		msxmaniaDir := getConfigWithFallback("msxmania", "")
		return filepath.Join(msxmaniaDir, game.Disco)
	}
	if currentCategory == "Good MSX 1" {
		rootDir, err := filepath.Abs(".")
		if err != nil {
			rootDir = "."
		}
		return filepath.Join(rootDir, "Common", "Good_MSX1_Roms", game.Disco)
	}

	tipoUpper := strings.ToUpper(game.Tipo)
	var baseDir string
	if tipoUpper == "ROM" {
		baseDir = getConfigWithFallback("roms", "ROM")
	} else if tipoUpper == "DSK" {
		baseDir = getConfigWithFallback("dsks", "DSK")
	} else {
		baseDir = getConfigWithFallback("raiz", ".")
	}
	dirVol := fmt.Sprintf("%02d", game.CdNumero)
	return filepath.Join(baseDir, dirVol, game.Disco)
}

func launchGame(game *Game, w fyne.Window) {
	// Verificar se o jogo possui configuração personalizada na tabela game_emulacao
	hasCustom, emuConfigured, err := HasGameEmulacao(game.CdNumero, game.Disco)
	if err != nil {
		fmt.Printf("Erro ao verificar configuração customizada de emulação: %v\n", err)
	}

	var emuName string
	if hasCustom {
		emuName = strings.ToLower(emuConfigured)
	} else {
		emuName = strings.ToLower(game.Emulador)
	}

	if emuName == "" {
		emuName = strings.ToLower(getConfigWithFallback("emulador", "openmsx"))
		if emuName == "" {
			emuName = "openmsx"
		}
	}

	emuPath, err := GetConfig(emuName)
	if err != nil || emuPath == "" {
		dialog.ShowError(fmt.Errorf(T("err_emulator_path", emuName)), w)
		return
	}

	gamePath := getGamePath(game)
	if _, err := os.Stat(gamePath); os.IsNotExist(err) {
		dialog.ShowError(fmt.Errorf(T("err_game_not_found", gamePath)), w)
		return
	}

	var args []string
	if !hasCustom && strings.TrimSpace(game.Options) != "" && emuName == "openmsx" {
		// Sem configuração customizada, usa as opções da tabela msxstuffs diretamente
		args = strings.Fields(game.Options)
	} else {
		// Caso exista configuração customizada, ou opções da tabela estejam vazias (fallback para setup padrão)
		if emuName == "openmsx" {
			var maquina string
			if hasCustom {
				maquina, _ = GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "maquina")
			}
			if maquina == "" {
				maquina = getConfigWithFallback("openmsx_maquina", "Gradiente_Expert_GPC-1")
			}
			if maquina != "" {
				args = append(args, "-machine", maquina)
			}

			for i := 1; i <= 4; i++ {
				var ext string
				if hasCustom {
					ext, _ = GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", fmt.Sprintf("extensao%d", i))
				}
				if ext == "" {
					ext = getConfigWithFallback(fmt.Sprintf("openmsx_extensao%d", i), "")
				}
				if ext != "" {
					args = append(args, "-ext", ext)
				}
			}

			var opcoes string
			if hasCustom {
				opcoes, _ = GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "opcoes")
			}
			if opcoes == "" {
				opcoes = getConfigWithFallback("openmsx_opcoes", "")
			}
			if opcoes != "" {
				args = append(args, strings.Fields(opcoes)...)
			}
		} else if emuName == "bluemsx" {
			var maquina string
			if hasCustom {
				maquina, _ = GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "bluemsx", "maquina")
			}
			if maquina == "" {
				maquina = getConfigWithFallback("bluemsx_maquina", "")
			}
			if maquina != "" {
				args = append(args, "-machine", fmt.Sprintf(`"%s"`, maquina))
			}

			var opcoes string
			if hasCustom {
				opcoes, _ = GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "bluemsx", "opcoes")
			}
			if opcoes == "" {
				opcoes = getConfigWithFallback("bluemsx_opcoes", "")
			}
			if opcoes != "" {
				args = append(args, strings.Fields(opcoes)...)
			}
		} else {
			// fmsx, rumsx
			var opcoes string
			if hasCustom {
				opcoes, _ = GetGameEmuladorDetalhe(game.CdNumero, game.Disco, emuName, "opcoes")
			}
			if opcoes == "" {
				opcoes = getConfigWithFallback(emuName+"_opcoes", "")
			}
			if opcoes != "" {
				args = append(args, strings.Fields(opcoes)...)
			}
		}
	}

	companionPath := findCompanionFile(gamePath)
	fileArgs := buildFileArgs(emuName, game.Tipo, gamePath, companionPath)
	args = append(args, fileArgs...)

	fmt.Println(T("log_starting_game", time.Now().Format("2006-01-02 15:04:05"), game.Descricao))
	fmt.Println(T("log_command", time.Now().Format("2006-01-02 15:04:05"), emuPath, strings.Join(args, " ")))

	cmd := exec.Command(emuPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		if err := cmd.Start(); err != nil {
			fmt.Println(T("log_emu_err", time.Now().Format("2006-01-02 15:04:05"), err))
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf(T("err_emulator_start", err)), w)
			})
			return
		}

		// Incrementa a contagem de execuções do jogo
		_ = IncrementGameExecution(game.CdNumero, game.Disco)

		err := cmd.Wait()
		if err != nil {
			fmt.Println(T("log_emu_exit_err", time.Now().Format("2006-01-02 15:04:05"), err))
		} else {
			fmt.Println(T("log_emu_exit_ok", time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
}

func showMachineListDialog(parent fyne.Window, openmsxExe string, onSelect func(string)) {
	if openmsxExe == "" {
		dialog.ShowError(fmt.Errorf(T("err_openmsx_required")), parent)
		return
	}
	baseDir := filepath.Dir(openmsxExe)
	machinesDir := filepath.Join(baseDir, "share", "machines")

	entries, err := os.ReadDir(machinesDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf(T("err_machines_dir", machinesDir)), parent)
		return
	}

	var machines []string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(strings.ToLower(name), ".xml") {
				machineName := strings.TrimSuffix(name, filepath.Ext(name))
				machines = append(machines, machineName)
			}
		}
	}

	if len(machines) == 0 {
		dialog.ShowInformation(T("msg_success"), T("err_no_machines", machinesDir), parent)
		return
	}

	sort.Strings(machines)

	var d dialog.Dialog
	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("Filtrar máquinas...")

	filteredMachines := make([]string, len(machines))
	copy(filteredMachines, machines)

	var list *widget.List
	list = widget.NewList(
		func() int {
			return len(filteredMachines)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Nome da Máquina")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(filteredMachines) {
				return
			}
			item.(*widget.Label).SetText(filteredMachines[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(filteredMachines) {
			onSelect(filteredMachines[id])
			d.Hide()
		}
	}

	filterEntry.OnChanged = func(val string) {
		valLower := strings.ToLower(val)
		filteredMachines = nil
		for _, m := range machines {
			if strings.Contains(strings.ToLower(m), valLower) {
				filteredMachines = append(filteredMachines, m)
			}
		}
		list.Refresh()
		list.UnselectAll()
	}

	content := container.NewBorder(
		filterEntry,
		nil,
		nil,
		nil,
		container.NewGridWrap(fyne.NewSize(350, 250), list),
	)

	d = dialog.NewCustom(T("dialog_select_machine"), T("btn_cancel"), content, parent)
	d.Resize(fyne.NewSize(400, 350))
	d.Show()
}

func showBluemsxMachineListDialog(parent fyne.Window, bluemsxExe string, onSelect func(string)) {
	if bluemsxExe == "" {
		dialog.ShowError(fmt.Errorf(T("err_bluemsx_required")), parent)
		return
	}
	baseDir := filepath.Dir(bluemsxExe)
	machinesDir := filepath.Join(baseDir, "Machines")
	if _, err := os.Stat(machinesDir); err != nil {
		machinesDir = filepath.Join(baseDir, "machines")
	}

	entries, err := os.ReadDir(machinesDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf(T("err_machines_dir", machinesDir)), parent)
		return
	}

	var machines []string
	for _, entry := range entries {
		if entry.IsDir() {
			machines = append(machines, entry.Name())
		}
	}

	if len(machines) == 0 {
		dialog.ShowInformation(T("msg_success"), T("err_no_machines", machinesDir), parent)
		return
	}

	sort.Strings(machines)

	var d dialog.Dialog
	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("Filtrar máquinas...")

	filteredMachines := make([]string, len(machines))
	copy(filteredMachines, machines)

	var list *widget.List
	list = widget.NewList(
		func() int {
			return len(filteredMachines)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Nome da Máquina")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(filteredMachines) {
				return
			}
			item.(*widget.Label).SetText(filteredMachines[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(filteredMachines) {
			onSelect(filteredMachines[id])
			d.Hide()
		}
	}

	filterEntry.OnChanged = func(val string) {
		valLower := strings.ToLower(val)
		filteredMachines = nil
		for _, m := range machines {
			if strings.Contains(strings.ToLower(m), valLower) {
				filteredMachines = append(filteredMachines, m)
			}
		}
		list.Refresh()
		list.UnselectAll()
	}

	content := container.NewBorder(
		filterEntry,
		nil,
		nil,
		nil,
		container.NewGridWrap(fyne.NewSize(350, 250), list),
	)

	d = dialog.NewCustom(T("dialog_select_machine"), T("btn_cancel"), content, parent)
	d.Resize(fyne.NewSize(400, 350))
	d.Show()
}

func showExtensionListDialog(parent fyne.Window, openmsxExe string, onSelect func(string)) {
	if openmsxExe == "" {
		dialog.ShowError(fmt.Errorf(T("err_openmsx_required")), parent)
		return
	}
	baseDir := filepath.Dir(openmsxExe)
	extensionsDir := filepath.Join(baseDir, "share", "extensions")

	entries, err := os.ReadDir(extensionsDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf(T("err_extensions_dir", extensionsDir)), parent)
		return
	}

	var extensions []string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(strings.ToLower(name), ".xml") {
				extName := strings.TrimSuffix(name, filepath.Ext(name))
				extensions = append(extensions, extName)
			}
		}
	}

	if len(extensions) == 0 {
		dialog.ShowInformation(T("msg_success"), T("err_no_extensions", extensionsDir), parent)
		return
	}

	sort.Strings(extensions)

	var d dialog.Dialog
	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("Filtrar extensões...")

	filteredExtensions := make([]string, len(extensions))
	copy(filteredExtensions, extensions)

	var list *widget.List
	list = widget.NewList(
		func() int {
			return len(filteredExtensions)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Nome da Extensão")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(filteredExtensions) {
				return
			}
			item.(*widget.Label).SetText(filteredExtensions[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(filteredExtensions) {
			onSelect(filteredExtensions[id])
			d.Hide()
		}
	}

	filterEntry.OnChanged = func(val string) {
		valLower := strings.ToLower(val)
		filteredExtensions = nil
		for _, e := range extensions {
			if strings.Contains(strings.ToLower(e), valLower) {
				filteredExtensions = append(filteredExtensions, e)
			}
		}
		list.Refresh()
		list.UnselectAll()
	}

	content := container.NewBorder(
		filterEntry,
		nil,
		nil,
		nil,
		container.NewGridWrap(fyne.NewSize(350, 250), list),
	)

	d = dialog.NewCustom(T("dialog_select_extension"), T("btn_cancel"), content, parent)
	d.Resize(fyne.NewSize(400, 350))
	d.Show()
}

var gameConfigWindow fyne.Window

func showGameConfig(game *Game, parent fyne.Window, myApp fyne.App) {
	if gameConfigWindow != nil {
		gameConfigWindow.RequestFocus()
		return
	}

	gameConfigWindow = myApp.NewWindow(T("config_title", game.Descricao))
	gameConfigWindow.SetOnClosed(func() {
		gameConfigWindow = nil
	})

	// 1. Carrega ou cria os dados de emulação do banco
	execucoes, emuEscolhido, err := GetOrCreateGameEmulacao(game.CdNumero, game.Disco)
	if err != nil {
		dialog.ShowError(err, parent)
		gameConfigWindow.Close()
		return
	}

	// 2. Define widgets
	execucoesLabel := widget.NewLabel(T("label_executions", execucoes))
	
	// Dropdown para escolher o emulador
	emuladorSelect := widget.NewSelect([]string{"openMSX", "fMSX", "blueMSX", "ruMSX"}, nil)
	// Normaliza para comparar
	selectedEmuNormal := strings.ToLower(emuEscolhido)
	switch selectedEmuNormal {
	case "openmsx":
		emuladorSelect.SetSelected("openMSX")
	case "fmsx":
		emuladorSelect.SetSelected("fMSX")
	case "bluemsx":
		emuladorSelect.SetSelected("blueMSX")
	case "rumsx":
		emuladorSelect.SetSelected("ruMSX")
	default:
		emuladorSelect.SetSelected("openMSX")
	}

	gamePath := getGamePath(game)

	// 3. Inputs dos emuladores
	// openMSX inputs
	openmsxMachineEntry := widget.NewEntry()
	valMachine, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "maquina")
	openmsxMachineEntry.SetText(valMachine)
	openmsxMachineEntry.SetPlaceHolder(T("placeholder_machine"))

	openmsxMachineBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		openmsxExe, _ := GetConfig("openmsx")
		showMachineListDialog(gameConfigWindow, openmsxExe, func(machine string) {
			openmsxMachineEntry.SetText(machine)
		})
	})
	openmsxMachineContainer := container.NewBorder(nil, nil, nil, openmsxMachineBtn, openmsxMachineEntry)

	openmsxExt1Entry := widget.NewEntry()
	valExt1, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao1")
	openmsxExt1Entry.SetText(valExt1)
	openmsxExt1Entry.SetPlaceHolder(T("placeholder_extension"))

	openmsxExt1Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		openmsxExe, _ := GetConfig("openmsx")
		showExtensionListDialog(gameConfigWindow, openmsxExe, func(ext string) {
			openmsxExt1Entry.SetText(ext)
		})
	})
	openmsxExt1Container := container.NewBorder(nil, nil, nil, openmsxExt1Btn, openmsxExt1Entry)

	openmsxExt2Entry := widget.NewEntry()
	valExt2, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao2")
	openmsxExt2Entry.SetText(valExt2)

	openmsxExt2Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		openmsxExe, _ := GetConfig("openmsx")
		showExtensionListDialog(gameConfigWindow, openmsxExe, func(ext string) {
			openmsxExt2Entry.SetText(ext)
		})
	})
	openmsxExt2Container := container.NewBorder(nil, nil, nil, openmsxExt2Btn, openmsxExt2Entry)

	openmsxExt3Entry := widget.NewEntry()
	valExt3, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao3")
	openmsxExt3Entry.SetText(valExt3)

	openmsxExt3Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		openmsxExe, _ := GetConfig("openmsx")
		showExtensionListDialog(gameConfigWindow, openmsxExe, func(ext string) {
			openmsxExt3Entry.SetText(ext)
		})
	})
	openmsxExt3Container := container.NewBorder(nil, nil, nil, openmsxExt3Btn, openmsxExt3Entry)

	openmsxExt4Entry := widget.NewEntry()
	valExt4, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao4")
	openmsxExt4Entry.SetText(valExt4)

	openmsxExt4Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		openmsxExe, _ := GetConfig("openmsx")
		showExtensionListDialog(gameConfigWindow, openmsxExe, func(ext string) {
			openmsxExt4Entry.SetText(ext)
		})
	})
	openmsxExt4Container := container.NewBorder(nil, nil, nil, openmsxExt4Btn, openmsxExt4Entry)

	openmsxOptionsEntry := widget.NewEntry()
	valOptionsOpenmsx, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "opcoes")
	openmsxOptionsEntry.SetText(valOptionsOpenmsx)
	openmsxOptionsEntry.SetPlaceHolder(T("placeholder_options"))

	openmsxForm := widget.NewForm(
		widget.NewFormItem(T("form_custom_machine"), openmsxMachineContainer),
		widget.NewFormItem(T("form_extension", 1), openmsxExt1Container),
		widget.NewFormItem(T("form_extension", 2), openmsxExt2Container),
		widget.NewFormItem(T("form_extension", 3), openmsxExt3Container),
		widget.NewFormItem(T("form_extension", 4), openmsxExt4Container),
		widget.NewFormItem(T("form_free_options"), openmsxOptionsEntry),
	)

	// fMSX inputs
	fmsxOptionsEntry := widget.NewEntry()
	valOptionsFmsx, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "fmsx", "opcoes")
	fmsxOptionsEntry.SetText(valOptionsFmsx)
	fmsxOptionsEntry.SetPlaceHolder(T("placeholder_options_fmsx"))
	fmsxForm := widget.NewForm(
		widget.NewFormItem(T("form_free_options"), fmsxOptionsEntry),
	)

	// blueMSX inputs
	bluemsxMachineEntry := widget.NewEntry()
	valMachineBluemsx, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "bluemsx", "maquina")
	bluemsxMachineEntry.SetText(valMachineBluemsx)
	bluemsxMachineEntry.SetPlaceHolder("Ex: MSX2")

	bluemsxMachineBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		bluemsxExe, _ := GetConfig("bluemsx")
		showBluemsxMachineListDialog(gameConfigWindow, bluemsxExe, func(machine string) {
			bluemsxMachineEntry.SetText(machine)
		})
	})
	bluemsxMachineContainer := container.NewBorder(nil, nil, nil, bluemsxMachineBtn, bluemsxMachineEntry)

	bluemsxOptionsEntry := widget.NewEntry()
	valOptionsBluemsx, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "bluemsx", "opcoes")
	bluemsxOptionsEntry.SetText(valOptionsBluemsx)
	bluemsxOptionsEntry.SetPlaceHolder(T("placeholder_options_bluemsx"))

	bluemsxForm := widget.NewForm(
		widget.NewFormItem(T("form_custom_machine"), bluemsxMachineContainer),
		widget.NewFormItem(T("form_free_options"), bluemsxOptionsEntry),
	)

	// ruMSX inputs
	rumsxOptionsEntry := widget.NewEntry()
	valOptionsRumsx, _ := GetGameEmuladorDetalhe(game.CdNumero, game.Disco, "rumsx", "opcoes")
	rumsxOptionsEntry.SetText(valOptionsRumsx)
	rumsxOptionsEntry.SetPlaceHolder(T("placeholder_options_rumsx"))
	rumsxForm := widget.NewForm(
		widget.NewFormItem(T("form_free_options"), rumsxOptionsEntry),
	)

	// Linha de Comando Gerada (Preview)
	cmdPreviewEntry := widget.NewMultiLineEntry()
	cmdPreviewEntry.Disable()

	// Função que reconstrói a linha de comando em tempo real para fins de visualização
	updateCommandPreview := func() {
		emuName := strings.ToLower(emuladorSelect.Selected)
		emuPath, _ := GetConfig(emuName)
		if emuPath == "" {
			emuPath = emuName + ".exe"
		}

		var args []string
		if emuName == "openmsx" {
			maquina := openmsxMachineEntry.Text
			if maquina == "" {
				maquina = getConfigWithFallback("openmsx_maquina", "Gradiente_Expert_GPC-1")
			}
			if maquina != "" {
				args = append(args, "-machine", maquina)
			}

			for i := 1; i <= 4; i++ {
				var ext string
				if i == 1 {
					ext = openmsxExt1Entry.Text
				} else if i == 2 {
					ext = openmsxExt2Entry.Text
				} else if i == 3 {
					ext = openmsxExt3Entry.Text
				} else if i == 4 {
					ext = openmsxExt4Entry.Text
				}
				if ext == "" {
					ext = getConfigWithFallback(fmt.Sprintf("openmsx_extensao%d", i), "")
				}
				if ext != "" {
					args = append(args, "-ext", ext)
				}
			}

			opcoes := openmsxOptionsEntry.Text
			if opcoes == "" {
				opcoes = getConfigWithFallback("openmsx_opcoes", "")
			}
			if opcoes != "" {
				args = append(args, strings.Fields(opcoes)...)
			}
		} else if emuName == "bluemsx" {
			maquina := bluemsxMachineEntry.Text
			if maquina == "" {
				maquina = getConfigWithFallback("bluemsx_maquina", "")
			}
			if maquina != "" {
				args = append(args, "-machine", fmt.Sprintf(`"%s"`, maquina))
			}

			opcoes := bluemsxOptionsEntry.Text
			if opcoes == "" {
				opcoes = getConfigWithFallback("bluemsx_opcoes", "")
			}
			if opcoes != "" {
				args = append(args, strings.Fields(opcoes)...)
			}
		} else {
			// fmsx, rumsx
			var opcoes string
			if emuName == "fmsx" {
				opcoes = fmsxOptionsEntry.Text
			} else if emuName == "rumsx" {
				opcoes = rumsxOptionsEntry.Text
			}
			if opcoes == "" {
				opcoes = getConfigWithFallback(emuName+"_opcoes", "")
			}
			if opcoes != "" {
				args = append(args, strings.Fields(opcoes)...)
			}
		}

		companionPath := findCompanionFile(gamePath)
		fileArgs := buildFileArgs(emuName, game.Tipo, gamePath, companionPath)
		args = append(args, fileArgs...)

		cmdLine := fmt.Sprintf("%s %s", emuPath, strings.Join(args, " "))
		cmdPreviewEntry.SetText(cmdLine)
	}

	// Adiciona hooks de atualização reativa em tempo real para os campos editáveis
	openmsxMachineEntry.OnChanged = func(string) { updateCommandPreview() }
	openmsxExt1Entry.OnChanged = func(string) { updateCommandPreview() }
	openmsxExt2Entry.OnChanged = func(string) { updateCommandPreview() }
	openmsxExt3Entry.OnChanged = func(string) { updateCommandPreview() }
	openmsxExt4Entry.OnChanged = func(string) { updateCommandPreview() }
	openmsxOptionsEntry.OnChanged = func(string) { updateCommandPreview() }
	fmsxOptionsEntry.OnChanged = func(string) { updateCommandPreview() }
	bluemsxMachineEntry.OnChanged = func(string) { updateCommandPreview() }
	bluemsxOptionsEntry.OnChanged = func(string) { updateCommandPreview() }
	rumsxOptionsEntry.OnChanged = func(string) { updateCommandPreview() }

	// Função auxiliar para atualizar a edição dos emuladores conforme escolha do dropdown
	updateEditableEmulator := func(selected string) {
		selectedLower := strings.ToLower(selected)
		
		// openMSX toggles
		if selectedLower == "openmsx" {
			openmsxMachineEntry.Enable()
			openmsxMachineBtn.Enable()
			openmsxExt1Entry.Enable()
			openmsxExt1Btn.Enable()
			openmsxExt2Entry.Enable()
			openmsxExt2Btn.Enable()
			openmsxExt3Entry.Enable()
			openmsxExt3Btn.Enable()
			openmsxExt4Entry.Enable()
			openmsxExt4Btn.Enable()
			openmsxOptionsEntry.Enable()
		} else {
			openmsxMachineEntry.Disable()
			openmsxMachineBtn.Disable()
			openmsxExt1Entry.Disable()
			openmsxExt1Btn.Disable()
			openmsxExt2Entry.Disable()
			openmsxExt2Btn.Disable()
			openmsxExt3Entry.Disable()
			openmsxExt3Btn.Disable()
			openmsxExt4Entry.Disable()
			openmsxExt4Btn.Disable()
			openmsxOptionsEntry.Disable()
		}

		// fMSX toggles
		if selectedLower == "fmsx" {
			fmsxOptionsEntry.Enable()
		} else {
			fmsxOptionsEntry.Disable()
		}

		// blueMSX toggles
		if selectedLower == "bluemsx" {
			bluemsxMachineEntry.Enable()
			bluemsxMachineBtn.Enable()
			bluemsxOptionsEntry.Enable()
		} else {
			bluemsxMachineEntry.Disable()
			bluemsxMachineBtn.Disable()
			bluemsxOptionsEntry.Disable()
		}

		// ruMSX toggles
		if selectedLower == "rumsx" {
			rumsxOptionsEntry.Enable()
		} else {
			rumsxOptionsEntry.Disable()
		}
	}

	// Atribui callback ao dropdown
	emuladorSelect.OnChanged = func(selected string) {
		updateEditableEmulator(selected)
		updateCommandPreview()
	}

	// Inicializa os estados e desenha o preview inicial
	updateEditableEmulator(emuladorSelect.Selected)
	updateCommandPreview()

	// Abas dos Emuladores
	tabs := container.NewAppTabs(
		container.NewTabItem("openMSX", openmsxForm),
		container.NewTabItem("fMSX", fmsxForm),
		container.NewTabItem("blueMSX", bluemsxForm),
		container.NewTabItem("ruMSX", rumsxForm),
	)

	generalForm := widget.NewForm(
		widget.NewFormItem(T("form_executed"), execucoesLabel),
		widget.NewFormItem(T("form_active_emulator"), emuladorSelect),
	)

	// Botões de Ação
	salvarBtn := widget.NewButtonWithIcon(T("btn_save"), theme.ConfirmIcon(), func() {
		chosenEmu := strings.ToLower(emuladorSelect.Selected)
		if err := SaveGameEmulacao(game.CdNumero, game.Disco, chosenEmu); err != nil {
			dialog.ShowError(err, gameConfigWindow)
			return
		}

		// Salva os dados do openMSX
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "maquina", openmsxMachineEntry.Text)
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao1", openmsxExt1Entry.Text)
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao2", openmsxExt2Entry.Text)
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao3", openmsxExt3Entry.Text)
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "extensao4", openmsxExt4Entry.Text)
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "openmsx", "opcoes", openmsxOptionsEntry.Text)

		// Salva fMSX
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "fmsx", "opcoes", fmsxOptionsEntry.Text)

		// Salva blueMSX
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "bluemsx", "maquina", bluemsxMachineEntry.Text)
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "bluemsx", "opcoes", bluemsxOptionsEntry.Text)

		// Salva ruMSX
		_ = SetGameEmuladorDetalhe(game.CdNumero, game.Disco, "rumsx", "opcoes", rumsxOptionsEntry.Text)

		dialog.ShowInformation(T("msg_saved_title"), T("msg_saved_success"), gameConfigWindow)
		gameConfigWindow.Close()
	})
	salvarBtn.Importance = widget.HighImportance

	fecharBtn := widget.NewButtonWithIcon(T("btn_close"), theme.CancelIcon(), func() {
		gameConfigWindow.Close()
	})

	actionsRow := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(240, 10)), // Espaçador
		salvarBtn,
		fecharBtn,
	)

	bottomContainer := container.NewVBox(
		widget.NewLabelWithStyle(T("label_cmd_preview"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		cmdPreviewEntry,
		actionsRow,
	)

	mainLayout := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(T("dialog_config_game_title"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(T("dialog_config_game_label", game.Descricao, game.Disco)),
			generalForm,
		),
		bottomContainer,
		nil,
		nil,
		tabs,
	)

	gameConfigWindow.SetContent(mainLayout)
	gameConfigWindow.Resize(fyne.NewSize(540, 520))
	gameConfigWindow.Show()
}

// Global log variables for stdout/stderr capturing
var (
	logBuffer       []string
	logMutex        sync.Mutex
	statusLogUpdate func()
)

func initOutputCapture() {
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, T("log_pipe_err", err))
		return
	}

	// Redireciona Stdout e Stderr
	os.Stdout = w
	os.Stderr = w
	log.SetOutput(w)

	// Thread de leitura do pipe
	go func() {
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadString('\n')
			if len(line) > 0 {
				line = strings.TrimSuffix(line, "\n")
				line = strings.TrimSuffix(line, "\r")

				logMutex.Lock()
				logBuffer = append(logBuffer, line)
				if len(logBuffer) > 1000 {
					logBuffer = logBuffer[len(logBuffer)-1000:]
				}
				logMutex.Unlock()

				if statusLogUpdate != nil {
					fyne.Do(statusLogUpdate)
				}
			}
			if err != nil {
				break
			}
		}
	}()
}

var statusWindow fyne.Window

func showStatusWindow(myApp fyne.App) {
	if statusWindow != nil {
		statusWindow.RequestFocus()
		return
	}

	statusWindow = myApp.NewWindow(T("logs_window_title"))
	statusWindow.SetOnClosed(func() {
		statusLogUpdate = nil
		statusWindow = nil
	})

	// Cria a lista virtualizada de logs
	logList := widget.NewList(
		func() int {
			logMutex.Lock()
			defer logMutex.Unlock()
			return len(logBuffer)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.TextStyle = fyne.TextStyle{Monospace: true}
			return label
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			logMutex.Lock()
			defer logMutex.Unlock()
			if id >= 0 && id < len(logBuffer) {
				item.(*widget.Label).SetText(logBuffer[id])
			}
		},
	)

	// Callback para atualizar a UI
	statusLogUpdate = func() {
		logList.Refresh()
		logMutex.Lock()
		length := len(logBuffer)
		logMutex.Unlock()
		if length > 0 {
			logList.ScrollTo(length - 1)
		}
	}

	// Scroll inicial para o final
	logMutex.Lock()
	initialLength := len(logBuffer)
	logMutex.Unlock()
	if initialLength > 0 {
		logList.ScrollTo(initialLength - 1)
	}

	// Botões
	copyBtn := widget.NewButtonWithIcon(T("btn_copy_all"), theme.ContentCopyIcon(), func() {
		logMutex.Lock()
		allText := strings.Join(logBuffer, "\n")
		logMutex.Unlock()
		statusWindow.Clipboard().SetContent(allText)
	})

	clearBtn := widget.NewButtonWithIcon(T("btn_clear"), theme.DeleteIcon(), func() {
		logMutex.Lock()
		logBuffer = []string{T("log_cleared", time.Now().Format("2006-01-02 15:04:05"))}
		logMutex.Unlock()
		if statusLogUpdate != nil {
			statusLogUpdate()
		}
	})

	closeBtn := widget.NewButtonWithIcon(T("btn_close"), theme.CancelIcon(), func() {
		statusWindow.Close()
	})

	buttonsRow := container.NewHBox(
		copyBtn,
		clearBtn,
		container.NewGridWrap(fyne.NewSize(150, 10)), // Espaçador
		closeBtn,
	)

	content := container.NewBorder(
		widget.NewLabelWithStyle(T("logs_header"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewPadded(buttonsRow),
		nil,
		nil,
		logList,
	)

	statusWindow.SetContent(content)
	statusWindow.Resize(fyne.NewSize(700, 450))
	statusWindow.Show()
}

func findCompanionFile(path string) string {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, ext)

	if len(name) == 0 {
		return ""
	}

	exists := func(n string) (string, bool) {
		p := filepath.Join(dir, n)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
		// Fallback for case sensitivity
		files, err := os.ReadDir(dir)
		if err != nil {
			return "", false
		}
		for _, f := range files {
			if strings.EqualFold(f.Name(), n) {
				return filepath.Join(dir, f.Name()), true
			}
		}
		return "", false
	}

	// Suffix mapping for part 1 -> part 2
	suffixes := []struct {
		search  string
		replace string
	}{
		{"-1", "-2"},
		{"_1", "_2"},
		{"-A", "-B"},
		{"-a", "-b"},
		{"_A", "_B"},
		{"_a", "_b"},
	}

	for _, s := range suffixes {
		if strings.HasSuffix(strings.ToUpper(name), s.search) {
			newName := name[:len(name)-len(s.search)] + s.replace
			if p, ok := exists(newName + ext); ok {
				return p
			}
		}
	}

	// Replace last character directly if it is '1', 'A', 'a'
	lastChar := name[len(name)-1:]
	if lastChar == "1" {
		newName := name[:len(name)-1] + "2"
		if p, ok := exists(newName + ext); ok {
			return p
		}
	} else if lastChar == "A" {
		newName := name[:len(name)-1] + "B"
		if p, ok := exists(newName + ext); ok {
			return p
		}
	} else if lastChar == "a" {
		newName := name[:len(name)-1] + "b"
		if p, ok := exists(newName + ext); ok {
			return p
		}
	}

	// Try appending suffix -2, _2, -B, _B
	appends := []string{"-2", "_2", "-B", "_B", "-b", "_b"}
	for _, app := range appends {
		if p, ok := exists(name + app + ext); ok {
			return p
		}
	}

	return ""
}

func buildFileArgs(emuName string, tipo string, gamePath string, companionPath string) []string {
	tipoUpper := strings.ToUpper(tipo)
	var args []string

	switch emuName {
	case "openmsx":
		if tipoUpper == "DSK" {
			args = append(args, "-diska", gamePath)
		} else if tipoUpper == "ROM" {
			if companionPath != "" {
				args = append(args, "-carta", gamePath, "-cartb", companionPath)
			} else {
				args = append(args, "-carta", gamePath)
			}
		} else if tipoUpper == "CAS" {
			args = append(args, "-cassette", gamePath)
		} else {
			args = append(args, gamePath)
		}

	case "bluemsx":
		if tipoUpper == "DSK" {
			args = append(args, "/diskA", gamePath)
		} else if tipoUpper == "ROM" {
			if companionPath != "" {
				args = append(args, "/rom1", gamePath, "/rom2", companionPath)
			} else {
				args = append(args, "/rom1", gamePath)
			}
		} else if tipoUpper == "CAS" {
			args = append(args, "/cas", gamePath)
		} else {
			args = append(args, gamePath)
		}

	case "fmsx":
		if tipoUpper == "DSK" {
			args = append(args, "-diska", gamePath)
		} else if tipoUpper == "ROM" {
			if companionPath != "" {
				args = append(args, gamePath, companionPath)
			} else {
				args = append(args, gamePath)
			}
		} else if tipoUpper == "CAS" {
			args = append(args, "-cassette", gamePath)
		} else {
			args = append(args, gamePath)
		}

	case "rumsx":
		if tipoUpper == "DSK" {
			args = append(args, "/diskA", gamePath)
		} else if tipoUpper == "ROM" {
			if companionPath != "" {
				args = append(args, "/rom0", gamePath, "/rom1", companionPath)
			} else {
				args = append(args, "/rom0", gamePath)
			}
		} else if tipoUpper == "CAS" {
			args = append(args, "/cas", gamePath)
		} else {
			args = append(args, gamePath)
		}

	default:
		args = append(args, gamePath)
	}

	return args
}
