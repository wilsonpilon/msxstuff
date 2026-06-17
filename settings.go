package main

import (
	"errors"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Controlar abertura única da janela de configurações
var configWindow fyne.Window

func showSettings(myApp fyne.App, statusBar *widget.Label) {
	if configWindow != nil {
		configWindow.RequestFocus()
		return
	}

	configWindow = myApp.NewWindow(T("settings_window_title"))
	configWindow.SetOnClosed(func() {
		configWindow = nil
	})

	// Seletor de Tema
	themeSelect := widget.NewSelect(ThemeList, func(selected string) {
		myApp.Settings().SetTheme(GetTheme(selected))
		myApp.Preferences().SetString("theme", selected)
		statusBar.SetText(T("msg_theme_changed", selected))
	})
	
	// Inicializa com o tema atual
	activeTheme := myApp.Preferences().StringWithFallback("theme", string(ThemeOneDark))
	themeSelect.SetSelected(activeTheme)

	// Seletor de Volume Inicial
	volumes, err := GetVolumes()
	if err != nil {
		volumes = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	}

	volumeSelect := widget.NewSelect(volumes, func(selected string) {
		SetConfig("volume_inicial", selected)
		statusBar.SetText(T("msg_volume_changed", selected))
	})
	volumeSelect.SetSelected(getConfigWithFallback("volume_inicial", "1"))

	// Seletor de Idioma
	var currentLang string = getConfigWithFallback("language", "en")
	var langName string
	switch currentLang {
	case "en":
		langName = "English"
	case "pt":
		langName = "Português (Brasil)"
	case "it":
		langName = "Italiano"
	case "es":
		langName = "Español"
	case "nl":
		langName = "Nederlands"
	default:
		langName = "English"
	}

	langSelect := widget.NewSelect([]string{"English", "Português (Brasil)", "Italiano", "Español", "Nederlands"}, func(selected string) {
		var langCode string
		switch selected {
		case "English":
			langCode = "en"
		case "Português (Brasil)":
			langCode = "pt"
		case "Italiano":
			langCode = "it"
		case "Español":
			langCode = "es"
		case "Nederlands":
			langCode = "nl"
		default:
			langCode = "en"
		}
		SetConfig("language", langCode)
		CurrentLanguage = langCode
		statusBar.SetText(T("status_language_changed", selected))
		dialog.ShowInformation(T("msg_language_title"), T("msg_language_restart"), configWindow)
	})
	langSelect.SetSelected(langName)

	themeForm := widget.NewForm(
		widget.NewFormItem(T("form_visual_theme"), themeSelect),
		widget.NewFormItem(T("form_initial_volume"), volumeSelect),
		widget.NewFormItem(T("form_language"), langSelect),
	)

	rootDir, err := filepath.Abs(".")
	if err != nil {
		rootDir = "."
	}

	raizEntry := widget.NewEntry()
	raizEntry.SetText(getConfigWithFallback("raiz", rootDir))

	picturesEntry := widget.NewEntry()
	picturesEntry.SetText(getConfigWithFallback("pictures", filepath.Join(rootDir, "pictures")))

	dsksEntry := widget.NewEntry()
	dsksEntry.SetText(getConfigWithFallback("dsks", filepath.Join(rootDir, "DSK")))

	romsEntry := widget.NewEntry()
	romsEntry.SetText(getConfigWithFallback("roms", filepath.Join(rootDir, "ROM")))

	msxmaniaEntry := widget.NewEntry()
	msxmaniaEntry.SetText(getConfigWithFallback("msxmania", filepath.Join(rootDir, "Common", "MSX_MANIA")))

	msxmaniaPicturesEntry := widget.NewEntry()
	msxmaniaPicturesEntry.SetText(getConfigWithFallback("msxmania_pictures", filepath.Join(rootDir, "pictures", "msxmania", "MSX")))

	dbEntry := widget.NewEntry()
	dbEntry.SetText(getConfigWithFallback("database", filepath.Join(rootDir, "data", "msxstuff.db")))

	goodmsx1Entry := widget.NewEntry()
	goodmsx1Entry.SetText(getConfigWithFallback("goodmsx1_dir", filepath.Join(rootDir, "Common", "Good_MSX1_Roms")))

	goodmsx2Entry := widget.NewEntry()
	goodmsx2Entry.SetText(getConfigWithFallback("goodmsx2_dir", filepath.Join(rootDir, "Common", "Good_MSX2_Roms")))

	megaramEntry := widget.NewEntry()
	megaramEntry.SetText(getConfigWithFallback("megaram_dir", filepath.Join(rootDir, "Common", "MEGARAM")))

	// Botões de navegação
	raizBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			raizEntry.SetText(uri.Path())
		}, configWindow)
	})

	picturesBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			picturesEntry.SetText(uri.Path())
		}, configWindow)
	})

	dsksBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			dsksEntry.SetText(uri.Path())
		}, configWindow)
	})

	romsBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			romsEntry.SetText(uri.Path())
		}, configWindow)
	})

	msxmaniaBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			msxmaniaEntry.SetText(uri.Path())
		}, configWindow)
	})

	msxmaniaPicturesBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			msxmaniaPicturesEntry.SetText(uri.Path())
		}, configWindow)
	})

	dbBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			dbEntry.SetText(reader.URI().Path())
		}, configWindow)
	})

	goodmsx1BrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			goodmsx1Entry.SetText(uri.Path())
		}, configWindow)
	})

	goodmsx2BrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			goodmsx2Entry.SetText(uri.Path())
		}, configWindow)
	})

	megaramBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			megaramEntry.SetText(uri.Path())
		}, configWindow)
	})

	// Containers acoplando caixas de entrada com botões de navegação
	raizContainer := container.NewBorder(nil, nil, nil, raizBrowseBtn, raizEntry)
	picturesContainer := container.NewBorder(nil, nil, nil, picturesBrowseBtn, picturesEntry)
	dsksContainer := container.NewBorder(nil, nil, nil, dsksBrowseBtn, dsksEntry)
	romsContainer := container.NewBorder(nil, nil, nil, romsBrowseBtn, romsEntry)
	msxmaniaContainer := container.NewBorder(nil, nil, nil, msxmaniaBrowseBtn, msxmaniaEntry)
	msxmaniaPicturesContainer := container.NewBorder(nil, nil, nil, msxmaniaPicturesBrowseBtn, msxmaniaPicturesEntry)
	dbContainer := container.NewBorder(nil, nil, nil, dbBrowseBtn, dbEntry)
	goodmsx1Container := container.NewBorder(nil, nil, nil, goodmsx1BrowseBtn, goodmsx1Entry)
	goodmsx2Container := container.NewBorder(nil, nil, nil, goodmsx2BrowseBtn, goodmsx2Entry)
	megaramContainer := container.NewBorder(nil, nil, nil, megaramBrowseBtn, megaramEntry)

	pathsForm := widget.NewForm(
		widget.NewFormItem(T("form_root_dir"), raizContainer),
		widget.NewFormItem(T("form_pictures_dir"), picturesContainer),
		widget.NewFormItem(T("form_dsks_dir"), dsksContainer),
		widget.NewFormItem(T("form_roms_dir"), romsContainer),
		widget.NewFormItem(T("form_msxmania_dir"), msxmaniaContainer),
		widget.NewFormItem(T("form_msxmania_pictures_dir"), msxmaniaPicturesContainer),
		widget.NewFormItem(T("form_database_file"), dbContainer),
		widget.NewFormItem(T("form_goodmsx1_dir"), goodmsx1Container),
		widget.NewFormItem(T("form_goodmsx2_dir"), goodmsx2Container),
		widget.NewFormItem(T("form_megaram_dir"), megaramContainer),
	)

	savePathsBtn := widget.NewButtonWithIcon(T("btn_save_paths"), theme.ConfirmIcon(), func() {
		var err error
		if err = SetConfig("raiz", raizEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("pictures", picturesEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("dsks", dsksEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("roms", romsEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("msxmania", msxmaniaEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("msxmania_pictures", msxmaniaPicturesEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("database", dbEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("goodmsx1_dir", goodmsx1Entry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("goodmsx2_dir", goodmsx2Entry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("megaram_dir", megaramEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		dialog.ShowInformation(T("msg_success"), T("msg_paths_saved"), configWindow)
	})

	// Executáveis dos Emuladores
	openmsxEntry := widget.NewEntry()
	openmsxEntry.SetText(getConfigWithFallback("openmsx", ""))

	fmsxEntry := widget.NewEntry()
	fmsxEntry.SetText(getConfigWithFallback("fmsx", ""))

	bluemsxEntry := widget.NewEntry()
	bluemsxEntry.SetText(getConfigWithFallback("bluemsx", ""))

	rumsxEntry := widget.NewEntry()
	rumsxEntry.SetText(getConfigWithFallback("rumsx", ""))

	openmsxBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			openmsxEntry.SetText(reader.URI().Path())
		}, configWindow)
	})

	fmsxBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			fmsxEntry.SetText(reader.URI().Path())
		}, configWindow)
	})

	bluemsxBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			bluemsxEntry.SetText(reader.URI().Path())
		}, configWindow)
	})

	rumsxBrowseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			rumsxEntry.SetText(reader.URI().Path())
		}, configWindow)
	})

	openmsxContainer := container.NewBorder(nil, nil, nil, openmsxBrowseBtn, openmsxEntry)
	fmsxContainer := container.NewBorder(nil, nil, nil, fmsxBrowseBtn, fmsxEntry)
	bluemsxContainer := container.NewBorder(nil, nil, nil, bluemsxBrowseBtn, bluemsxEntry)
	rumsxContainer := container.NewBorder(nil, nil, nil, rumsxBrowseBtn, rumsxEntry)

	execsForm := widget.NewForm(
		widget.NewFormItem(T("form_openmsx_exe"), openmsxContainer),
		widget.NewFormItem(T("form_fmsx_exe"), fmsxContainer),
		widget.NewFormItem(T("form_bluemsx_exe"), bluemsxContainer),
		widget.NewFormItem(T("form_rumsx_exe"), rumsxContainer),
	)

	saveExecsBtn := widget.NewButtonWithIcon(T("btn_save_executables"), theme.ConfirmIcon(), func() {
		var err error
		if err = SetConfig("openmsx", openmsxEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("fmsx", fmsxEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("bluemsx", bluemsxEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("rumsx", rumsxEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		dialog.ShowInformation(T("msg_success"), T("msg_executables_saved"), configWindow)
	})

	// Configurações do openMSX
	openmsxMachineEntry := widget.NewEntry()
	openmsxMachineEntry.SetText(getConfigWithFallback("openmsx_maquina", "Gradiente_Expert_GPC-1"))

	openmsxMachineBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		// Busca de máquinas (futuro)
	})
	openmsxMachineContainer := container.NewBorder(nil, nil, nil, openmsxMachineBtn, openmsxMachineEntry)

	openmsxOptionsEntry := widget.NewEntry()
	openmsxOptionsEntry.SetText(getConfigWithFallback("openmsx_opcoes", ""))

	openmsxExt1Entry := widget.NewEntry()
	openmsxExt1Entry.SetText(getConfigWithFallback("openmsx_extensao1", "DDX_3.0"))
	openmsxExt1Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		// Listagem de extensões (futuro)
	})
	openmsxExt1Container := container.NewBorder(nil, nil, nil, openmsxExt1Btn, openmsxExt1Entry)

	openmsxExt2Entry := widget.NewEntry()
	openmsxExt2Entry.SetText(getConfigWithFallback("openmsx_extensao2", ""))
	openmsxExt2Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		// Listagem de extensões (futuro)
	})
	openmsxExt2Container := container.NewBorder(nil, nil, nil, openmsxExt2Btn, openmsxExt2Entry)

	openmsxExt3Entry := widget.NewEntry()
	openmsxExt3Entry.SetText(getConfigWithFallback("openmsx_extensao3", ""))
	openmsxExt3Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		// Listagem de extensões (futuro)
	})
	openmsxExt3Container := container.NewBorder(nil, nil, nil, openmsxExt3Btn, openmsxExt3Entry)

	openmsxExt4Entry := widget.NewEntry()
	openmsxExt4Entry.SetText(getConfigWithFallback("openmsx_extensao4", ""))
	openmsxExt4Btn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		// Listagem de extensões (futuro)
	})
	openmsxExt4Container := container.NewBorder(nil, nil, nil, openmsxExt4Btn, openmsxExt4Entry)

	openmsxForm := widget.NewForm(
		widget.NewFormItem(T("form_default_machine"), openmsxMachineContainer),
		widget.NewFormItem(T("form_free_options"), openmsxOptionsEntry),
		widget.NewFormItem(T("form_extension", 1), openmsxExt1Container),
		widget.NewFormItem(T("form_extension", 2), openmsxExt2Container),
		widget.NewFormItem(T("form_extension", 3), openmsxExt3Container),
		widget.NewFormItem(T("form_extension", 4), openmsxExt4Container),
	)

	saveOpenmsxBtn := widget.NewButtonWithIcon(T("btn_save_openmsx"), theme.ConfirmIcon(), func() {
		var err error
		if err = SetConfig("openmsx_maquina", openmsxMachineEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("openmsx_opcoes", openmsxOptionsEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("openmsx_extensao1", openmsxExt1Entry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("openmsx_extensao2", openmsxExt2Entry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("openmsx_extensao3", openmsxExt3Entry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err = SetConfig("openmsx_extensao4", openmsxExt4Entry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		dialog.ShowInformation(T("msg_success"), T("msg_openmsx_saved"), configWindow)
	})

	// Configurações do fMSX
	fmsxOptionsEntry := widget.NewEntry()
	fmsxOptionsEntry.SetText(getConfigWithFallback("fmsx_opcoes", ""))

	fmsxForm := widget.NewForm(
		widget.NewFormItem(T("form_free_options"), fmsxOptionsEntry),
	)

	saveFmsxBtn := widget.NewButtonWithIcon(T("btn_save_fmsx"), theme.ConfirmIcon(), func() {
		if err := SetConfig("fmsx_opcoes", fmsxOptionsEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		dialog.ShowInformation(T("msg_success"), T("msg_fmsx_saved"), configWindow)
	})

	// Configurações do blueMSX
	bluemsxMachineEntry := widget.NewEntry()
	bluemsxMachineEntry.SetText(getConfigWithFallback("bluemsx_maquina", ""))

	bluemsxMachineBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		bluemsxExe, _ := GetConfig("bluemsx")
		showBluemsxMachineListDialog(configWindow, bluemsxExe, func(machine string) {
			bluemsxMachineEntry.SetText(machine)
		})
	})
	bluemsxMachineContainer := container.NewBorder(nil, nil, nil, bluemsxMachineBtn, bluemsxMachineEntry)

	bluemsxOptionsEntry := widget.NewEntry()
	bluemsxOptionsEntry.SetText(getConfigWithFallback("bluemsx_opcoes", ""))

	bluemsxForm := widget.NewForm(
		widget.NewFormItem(T("form_default_machine"), bluemsxMachineContainer),
		widget.NewFormItem(T("form_free_options"), bluemsxOptionsEntry),
	)

	saveBluemsxBtn := widget.NewButtonWithIcon(T("btn_save_bluemsx"), theme.ConfirmIcon(), func() {
		if err := SetConfig("bluemsx_maquina", bluemsxMachineEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		if err := SetConfig("bluemsx_opcoes", bluemsxOptionsEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		dialog.ShowInformation(T("msg_success"), T("msg_bluemsx_saved"), configWindow)
	})

	// Configurações do ruMSX
	rumsxOptionsEntry := widget.NewEntry()
	rumsxOptionsEntry.SetText(getConfigWithFallback("rumsx_opcoes", ""))

	rumsxForm := widget.NewForm(
		widget.NewFormItem(T("form_free_options"), rumsxOptionsEntry),
	)

	saveRumsxBtn := widget.NewButtonWithIcon(T("btn_save_rumsx"), theme.ConfirmIcon(), func() {
		if err := SetConfig("rumsx_opcoes", rumsxOptionsEntry.Text); err != nil {
			dialog.ShowError(err, configWindow)
			return
		}
		dialog.ShowInformation(T("msg_success"), T("msg_rumsx_saved"), configWindow)
	})

	debugCheck := widget.NewCheck(T("form_debug_mode"), func(checked bool) {
		val := "false"
		if checked {
			val = "true"
		}
		SetConfig("debug", val)
		debugMode = checked
		if checked {
			statusBar.SetText("Debug Mode Enabled")
		} else {
			statusBar.SetText("Debug Mode Disabled")
		}
	})
	debugCheck.SetChecked(getConfigWithFallback("debug", "false") == "true")

	// Abas de Configurações
	tabs := container.NewAppTabs(
		container.NewTabItem(T("tab_appearance"), container.NewVBox(
			widget.NewLabelWithStyle(T("label_customization"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(T("label_customization_desc")),
			themeForm,
			debugCheck,
		)),
		container.NewTabItem(T("tab_paths"), container.NewVBox(
			widget.NewLabelWithStyle(T("label_paths_desc"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			pathsForm,
			savePathsBtn,
		)),
		container.NewTabItem(T("tab_executables"), container.NewVBox(
			widget.NewLabelWithStyle(T("label_executables_desc"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			execsForm,
			saveExecsBtn,
		)),
		container.NewTabItem("openMSX", container.NewVBox(
			widget.NewLabelWithStyle(T("label_openmsx_desc"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			openmsxForm,
			saveOpenmsxBtn,
		)),
		container.NewTabItem("fMSX", container.NewVBox(
			widget.NewLabelWithStyle(T("label_fmsx_desc"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			fmsxForm,
			saveFmsxBtn,
		)),
		container.NewTabItem("blueMSX", container.NewVBox(
			widget.NewLabelWithStyle(T("label_bluemsx_desc"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			bluemsxForm,
			saveBluemsxBtn,
		)),
		container.NewTabItem("ruMSX", container.NewVBox(
			widget.NewLabelWithStyle(T("label_rumsx_desc"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			rumsxForm,
			saveRumsxBtn,
		)),
		container.NewTabItem(T("tab_database"), container.NewVBox(
			widget.NewLabelWithStyle(T("form_db_init_header"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(T("label_db_import_desc")),
			func() *widget.Button {
				btn := widget.NewButtonWithIcon(T("btn_initialize_msxstuffs"), theme.StorageIcon(), func() {
					err := InitDB()
					if err != nil {
						dialog.ShowError(err, configWindow)
					} else {
						dialog.ShowInformation(T("msg_success"), T("msg_db_initialized"), configWindow)
					}
				})
				return btn
			}(),
			func() *widget.Button {
				btn := widget.NewButtonWithIcon(T("btn_initialize_msxmania"), theme.StorageIcon(), func() {
					msxmaniaDir := getConfigWithFallback("msxmania", "")
					if msxmaniaDir == "" {
						dialog.ShowError(errors.New("Diretório do MSX Mania não configurado"), configWindow)
						return
					}
					err := ImportMSXMania(msxmaniaDir)
					if err != nil {
						dialog.ShowError(err, configWindow)
					} else {
						dialog.ShowInformation(T("msg_success"), T("msg_db_initialized"), configWindow)
					}
				})
				return btn
			}(),
			func() *widget.Button {
				btn := widget.NewButtonWithIcon(T("btn_initialize_goodmsx1"), theme.StorageIcon(), func() {
					rootDir := getConfigWithFallback("raiz", "")
					if rootDir == "" {
						var err error
						rootDir, err = filepath.Abs(".")
						if err != nil {
							rootDir = "."
						}
					}
					err := ImportGoodMSX1(rootDir)
					if err != nil {
						dialog.ShowError(err, configWindow)
					} else {
						dialog.ShowInformation(T("msg_success"), T("msg_db_initialized"), configWindow)
					}
				})
				return btn
			}(),
			func() *widget.Button {
				btn := widget.NewButtonWithIcon(T("btn_initialize_goodmsx2"), theme.StorageIcon(), func() {
					rootDir := getConfigWithFallback("raiz", "")
					if rootDir == "" {
						var err error
						rootDir, err = filepath.Abs(".")
						if err != nil {
							rootDir = "."
						}
					}
					err := ImportGoodMSX2(rootDir)
					if err != nil {
						dialog.ShowError(err, configWindow)
					} else {
						dialog.ShowInformation(T("msg_success"), T("msg_db_initialized"), configWindow)
					}
				})
				return btn
			}(),
			func() *widget.Button {
				btn := widget.NewButtonWithIcon(T("btn_initialize_megaram"), theme.StorageIcon(), func() {
					rootDir := getConfigWithFallback("raiz", "")
					if rootDir == "" {
						var err error
						rootDir, err = filepath.Abs(".")
						if err != nil {
							rootDir = "."
						}
					}
					err := ImportMegaram(rootDir)
					if err != nil {
						dialog.ShowError(err, configWindow)
					} else {
						dialog.ShowInformation(T("msg_success"), T("msg_db_initialized"), configWindow)
					}
				})
				return btn
			}(),
		)),
	)

	// Botão de fechar
	fecharBtn := widget.NewButtonWithIcon(T("btn_close"), theme.CancelIcon(), func() {
		configWindow.Close()
	})
	
	bottomContainer := container.NewHBox(
		layoutSpacer(), // Spacer
		fecharBtn,
	)

	configLayout := container.NewBorder(
		nil,
		bottomContainer,
		nil,
		nil,
		tabs,
	)

	configWindow.SetContent(configLayout)
	configWindow.Resize(fyne.NewSize(700, 520))
	configWindow.Show()
}

// Auxiliar para preencher espaço vazio nos layouts horizontais empurrando itens para a direita
func layoutSpacer() fyne.CanvasObject {
	return container.NewGridWrap(fyne.NewSize(580, 10))
}

// Retorna uma configuração do banco com um valor padrão caso falhe
func getConfigWithFallback(key string, fallback string) string {
	val, err := GetConfig(key)
	if err != nil {
		return fallback
	}
	return val
}
