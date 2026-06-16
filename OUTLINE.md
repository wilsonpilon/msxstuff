# MSX Stuffs - Guia e Instruções de Desenvolvimento (OUTLINE)

Este documento contém as especificações técnicas, decisões de arquitetura e instruções passo a passo para que qualquer Desenvolvedor Humano ou Inteligência Artificial (IA) possa recriar, compilar ou continuar o desenvolvimento do projeto **MSX Stuffs** a partir deste ponto, independentemente da IDE ou do ambiente de desenvolvimento utilizado.

---

## 1. Visão Geral do Projeto

O **MSX Stuffs** é um catálogo visual moderno para jogos e programas de MSX, inspirado no pacote clássico de 10 CDs lançado pela *Nemesis Informática* nos anos 2000. O objetivo é ler os índices textuais dos CDs (`MSXSTUFF1.txt` a `MSXSTUFF10.txt`), estruturá-los em um banco de dados local SQLite e oferecer uma interface gráfica rápida, responsiva e com suporte a temas visuais, que permita filtrar títulos e executá-los em emuladores modernos (como openMSX, fMSX, blueMSX e ruMSX).

---

## 2. Pilha Tecnológica (Tech Stack)

* **Linguagem**: Go (Golang) 1.20+.
* **Interface Gráfica**: [Fyne v2](https://fyne.io/) (Framework multiplataforma nativo em Go).
* **Banco de Dados**: SQLite, acessado via driver Pure Go `modernc.org/sqlite` (dispensa a necessidade de GCC/CGO para compilar o SQLite).
* **CLI (Interface de Linha de Comando)**: [Cobra](https://github.com/spf13/cobra) para tratamento e expansão futura de comandos de console.
* **Automação de Build**: PowerShell 7 (`build.ps1`).
* **Recursos do Windows**: Ferramenta `rsrc` para embutir ícones (`.ico`) no executável final no Windows.

---

## 3. Estrutura do Projeto e Arquivos

O projeto está modularizado de forma limpa, separando responsabilidades em múltiplos arquivos na raiz da pasta `goMSXStuffs`:

```
e:\goMSXStuffs
│   build.ps1                 # Script PowerShell de build
│   cli.go                    # Configuração de linha de comando (Cobra)
│   db.go                     # Lógica de banco de dados e parser de catálogos
│   go.mod                    # Declaração do módulo Go e dependências
│   go.sum                    # Sommas de verificação das dependências Go
│   gui.go                    # Janela principal Fyne e loop gráfico
│   i18n.go                   # Recursos de internacionalização e tradução (i18n)
│   LICENSE                   # Licença GNU GPL 3.0
│   main.go                   # Ponto de entrada (Entrypoint)
│   OUTLINE.md                # Este manual
│   README.md                 # Documentação e histórico do projeto
│   rsrc_windows_386.syso     # Recurso de ícone embutido para Windows 32 bits
│   rsrc_windows_amd64.syso   # Recurso de ícone embutido para Windows 64 bits
│   settings.go               # Tela e abas de Configurações
│   themes.go                 # Paletas de cores e definições de Temas Fyne
│
├───data
│       msxstuff.db           # Banco de dados gerado
│       MSXSTUFF1.txt...10    # Arquivos originais de texto com listas de arquivos
│
├───images
│       Icon.ico              # Ícone original do Windows
│       Icon.png              # Ícone convertido para uso na interface gráfica Fyne
│       MSXStuffs.png         # Imagem de fundo principal
│
└───pictures                  # Pasta que conterá as screenshots dos jogos
```

---

## 4. Detalhamento e Regras dos Arquivos de Código

### 4.1. [main.go](file:///e:/goMSXStuffs/main.go)
Inicia o aplicativo delegando a execução para a CLI (Cobra):
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
```

### 4.2. [cli.go](file:///e:/goMSXStuffs/cli.go)
Configura o comando raiz (`rootCmd`). Se nenhum subcomando for fornecido, executa o loop gráfico padrão via `runGUI()`:
```go
package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "msxstuffs",
	Short: "MSX Stuffs - Catálogo Visual",
	Run: func(cmd *cobra.Command, args []string) {
		runGUI()
	},
}
```

### 4.3. [themes.go](file:///e:/goMSXStuffs/themes.go)
Define a estrutura `CustomFyneTheme` implementando `fyne.Theme`. O sistema é baseado em paletas de cores estáticas (ex: `ThemeOneDark`, `ThemeDracula`, `ThemeMSXClassic`, `ThemeMSXExpert`, etc.).
* **Regra Importante**: Para que os ícones nativos e componentes do Fyne tenham contraste adequado, a função `Color(...)` deve retornar a variante `theme.VariantDark` ou `theme.VariantLight` condizente com o campo `IsDark` definido na paleta correspondente.

### 4.4. [db.go](file:///e:/goMSXStuffs/db.go)
Gerencia todas as operações do SQLite (`data/msxstuff.db`):
* **`InitDB()`**:
  1. Cria/atualiza a tabela `configuracoes` (armazenando pares chave/valor de caminhos e executáveis).
  2. Apaga e recria a tabela consolidada `msxstuffs` e as individuais `disco1` a `disco10`.
  3. Varre e processa cada arquivo `MSXSTUFF[1-10].txt` no diretório `data/`.
  4. Trata variações de maiúsculas/minúsculas no nome dos arquivos e ignora a linha de cabeçalho `disco=descricao` caso exista.
  5. Extrai extensão e calcula campos derivados:
     - `ext := filepath.Ext(discoVal)`
     - `raiz := strings.TrimSuffix(discoVal, ext)`
     - `tipo := strings.TrimPrefix(ext, ".")` (ROM, DSK, etc.)
  6. **Lógica de Opções padrão de Emulador**:
     - Se o tipo for `ROM` (ou `.ROM`), define a coluna `options` com: `"-machine Gradiente_Expert_GPC-1"`.
     - Se o tipo for `DSK` (ou `.DSK`), define a coluna `options` com: `"-machine Gradiente_Expert_GPC-1 -ext DDX_3.0"`.
     - Qualquer outra extensão deixa a coluna `options` em branco (`""`).
  7. Insere os registros em lote dentro de uma transação rápida (`tx.Commit()`).

### 4.5. [settings.go](file:///e:/goMSXStuffs/settings.go)
Contém a janela de configurações (`650x480`), composta por quatro abas:
1. **Aparência**: Seletor de Tema Visual e Seletor do "Volume Inicial" (calculado a partir dos volumes cadastrados no banco ou fallback de 1 a 10).
2. **Caminhos**: Gerencia diretórios para Arquivos ROM, DSK, Imagens (Screenshots) e do Banco de Dados. Utiliza caixas de texto com botões de navegação lateral (`dialog.ShowFolderOpen` / `dialog.ShowFileOpen`).
3. **Executáveis**: Configura os caminhos absolutos dos executáveis dos emuladores: *openMSX* (`openmsx.exe`), *fMSX* (`fmsx.exe`), *blueMSX* (`bluemsx.exe`) e *ruMSX* (`rumsx.exe`).
4. **Banco de Dados**: Botão "Inicializar DB" que executa `InitDB()` limpando e populando o banco de dados.

### 4.6. [gui.go](file:///e:/goMSXStuffs/gui.go)
Contrói a interface gráfica principal (`1024x768`):
* Aplica o ícone da aplicação (`images/Icon.png`).
* Aplica o tema salvo nas preferências.
* Cria o menu superior:
  - **Arquivo**: Status (abre a janela de logs do sistema e saídas dos emuladores) e Sair.
  - **Configuração**: Configurações (exibe a tela `settings.go`, controlando a instância de forma única para evitar múltiplas abas abertas).
  - **Ajuda**: Sobre.
* Exibe a imagem de fundo e desenha o layout do catálogo: barra lateral esquerda contendo o cabeçalho, os 10 botões segmentados de volume, o campo de busca de jogos e a lista dinâmica (`widget.NewList`); e o painel direito contendo detalhes do jogo selecionado (título, screenshot e botões "Jogar" / "Reconfigurar").
* Desenha e inicia a tela de splash sem bordas com fade-out de 1 segundo após 2 segundos de exibição.

### 4.7. [i18n.go](file:///e:/goMSXStuffs/i18n.go)
Gerencia os recursos de localização para suporte multilíngue:
* **Estrutura de Tradução**: Define um mapa de strings associando chaves de texto para 5 idiomas (`en`, `pt`, `it`, `es`, `nl`).
* **`T(key string, args ...any)`**: Função global para buscar e formatar strings traduzidas com base na linguagem ativa (`CurrentLanguage`), contendo fallback automático para a versão em inglês caso o termo não seja encontrado.

---

## 5. Como Compilar e Empacotar o Projeto

### 5.1. Pré-requisitos
1. **Go (Golang)** instalado e configurado no PATH do sistema.
2. **GCC (MinGW)**: Embora o banco SQLite dispense CGO, o framework gráfico Fyne exige CGO ativado para compilar componentes nativos que fazem chamadas à GPU (OpenGL/GLFW) no Windows e Linux.
3. **PowerShell 7**: Para executar os scripts de build.

### 5.2. Embutindo o Ícone no Executável (Windows)
Para garantir que o ícone oficial (`images/Icon.ico`) apareça no Windows Explorer, na barra de tarefas e nos atalhos, é utilizada a ferramenta `rsrc`:
```powershell
# Instalar a ferramenta de recursos
go install github.com/akavel/rsrc@latest

# Gerar o arquivo .syso para arquitetura 64 bits (executar no diretório do projeto)
rsrc -manifest msxstuffs.manifest -ico images/Icon.ico -o rsrc_windows_amd64.syso

# Gerar o arquivo .syso para arquitetura 32 bits (opcional)
rsrc -manifest msxstuffs.manifest -ico images/Icon.ico -o rsrc_windows_386.syso
```
*Nota*: O arquivo `.syso` nomeado especificamente com sufixo do OS/Arq impede que o Go tente fazer o link de recursos do Windows quando compilado no Linux, evitando quebras de build multiplataforma.

### 5.3. Executando o Build
Utilize o script PowerShell para compilar ou rodar o aplicativo de forma automática:
```powershell
# Compilar versão de lançamento (Release - limpa os símbolos de debug)
.\build.ps1 -OS windows -Mode Release

# Compilar e executar imediatamente
.\build.ps1 -OS windows -Mode Debug -Run
```

---

## 6. Diretrizes para Continuidade (Próximos Passos)

Com a arquitetura básica de navegação, visualização de screenshots, configuração de emuladores e suporte multilíngue concluída, os próximos passos sugeridos para evolução do projeto são:

### 6.1. Suporte a Arquivos Compactados (ZIP/7Z)
* **Objetivo**: Permitir executar jogos que estejam compactados sem precisar descompactar previamente.
* **Sugestão**:
  1. Modificar o carregamento do jogo para detectar se o arquivo termina em `.zip`, `.gz` ou `.7z`.
  2. Extrair temporariamente o arquivo na pasta temporária do sistema antes de repassar o caminho do executável ao emulador, limpando o arquivo temporário após o encerramento do processo.

### 6.2. Download Automático de Assets (Screenshots)
* **Objetivo**: Baixar screenshots ausentes diretamente da internet.
* **Sugestão**:
  1. Adicionar uma URL base de repositório de mídia do MSX (ou API pública) nas configurações.
  2. No fluxo de carregamento da imagem em `gui.go`, caso o arquivo local não exista, fazer uma chamada assíncrona (HTTP Get) para baixar a imagem correspondente a `<raiz>.png` e salvá-la em cache local.

### 6.3. Mapeamento e Suporte a Gamepads
* **Objetivo**: Permitir a navegação na interface e inicialização de jogos usando controles físicos.
* **Sugestão**:
  1. Integrar uma biblioteca de detecção de gamepad em Go.
  2. Mapear o D-pad / Analógico para mudar a seleção da lista de jogos, botões superiores (L/R) para trocar de volume, e o botão de ação (A/Confirmar) para executar o jogo selecionado.

### 6.4. Estatísticas de Uso Avançadas
* **Objetivo**: Ampliar as informações e contagens do banco de dados na UI.
* **Sugestão**:
  1. Aproveitar a tabela `game_emulacao` (que já possui a coluna `execucoes`) para criar uma tela ou painel de "Estatísticas".
  2. Exibir o ranking dos 10 jogos mais jogados no catálogo, tempo acumulado de jogo, ou data da última execução.
