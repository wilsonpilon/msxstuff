# Manual do Usuário e Desenvolvedor - MSX Stuffs

Este manual contém as instruções básicas sobre como baixar, compilar, configurar e executar o **MSX Stuffs - Catálogo Visual**, bem como a descrição das funcionalidades que já estão implementadas no sistema.

---

## 1. Requisitos do Sistema

Antes de iniciar, certifique-se de que seu ambiente atende aos seguintes requisitos:

* **Go (Golang)**: Versão 1.20 ou superior instalada.
* **Git**: Para clonar o repositório.
* **Compilador C (GCC)**: O framework gráfico **Fyne v2** necessita do CGO habilitado para ligar as bibliotecas gráficas nativas (OpenGL/GLFW) no Windows e Linux. No Windows, recomenda-se instalar o [MSYS2](https://www.msys2.org/) (com o pacote `mingw-w64-x86_64-toolchain`) ou o [w64devkit](https://github.com/skeeto/w64devkit/releases).
* **PowerShell 7** (Opcional): Para execução do script simplificado de compilação ([`build.ps1`](file:///e:/goMSXStuffs/build.ps1)).

---

## 2. Como Baixar o Projeto

Abra o seu terminal (Prompt de Comando ou PowerShell) e execute o comando abaixo para clonar o repositório:

```bash
git clone https://github.com/wilsonpilon/goMSXStuffs.git
cd goMSXStuffs
```

> [!NOTE]
> Certifique-se de que a estrutura das pastas `data/` (com os catálogos `MSXSTUFF1.txt` a `MSXSTUFF10.txt`) e `images/` (contendo os ícones e a imagem de fundo) estão presentes após o clone.

---

## 3. Como Compilar o Programa

O projeto possui um script PowerShell automatizado que organiza as dependências e gera o binário otimizado de forma simples.

### Método 1: Usando o Script de Build (PowerShell)
Abra o PowerShell no diretório do projeto e execute:

```powershell
# Compilação padrão para o Windows em modo de distribuição (Release)
.\build.ps1 -OS windows -Mode Release
```

Isso gerará o arquivo executável `msxstuffs.exe` na raiz do projeto.

### Método 2: Compilação Manual (Go CLI)
Se preferir compilar diretamente pela linha de comando do Go:

1. Baixe e organize as dependências:
   ```bash
   go mod tidy
   ```
2. Compile o executável removendo símbolos de depuração para diminuir o tamanho do arquivo:
   ```bash
   go build -ldflags="-s -w" -o msxstuffs.exe
   ```

### Embutindo o Ícone no Executável (Windows)
Se você precisar regenerar o ícone embutido no executável para que ele apareça no gerenciador de arquivos (Windows Explorer):

1. Instale a ferramenta `rsrc`:
   ```bash
   go install github.com/akavel/rsrc@latest
   ```
2. Gere o arquivo de recurso do Windows (faça isso antes de rodar o `go build`):
   ```bash
   rsrc -ico images/Icon.ico -o rsrc_windows_amd64.syso
   ```

---

## 4. Como Executar

### Executar o Binário Compilado
Após compilar, você pode rodar o programa dando um duplo clique no arquivo `msxstuffs.exe` ou executando via console:

```bash
.\msxstuffs.exe
```

### Compilar e Executar Imediatamente via Script
Você pode utilizar o script de build para compilar e inicializar o sistema de uma só vez:

```powershell
.\build.ps1 -OS windows -Mode Debug -Run
```

---

## 5. Funcionalidades Implementadas

O projeto de modernização já conta com as seguintes funcionalidades ativas:

### 5.1. Inicialização Automatizada do Banco de Dados
Na aba de configurações do banco de dados, o sistema lê de forma inteligente os arquivos de listagem originais (`MSXSTUFF1.txt` a `MSXSTUFF10.txt` localizados na pasta `data/`). Ele cria e popula:
* **Tabelas de Volumes**: Tabelas `disco1` a `disco10` contendo os registros brutos.
* **Tabela Consolidada (`msxstuffs`)**: Uma tabela unificada contendo todos os títulos indexados, calculando dinamicamente:
  * **Raiz**: Nome do arquivo sem extensão.
  * **Tipo**: Extensão do arquivo de MSX (ROM, DSK, etc.).
  * **Opções de Inicialização do Emulador**: 
    * Se for `.ROM`: configura automaticamente `"-machine Gradiente_Expert_GPC-1"`.
    * Se for `.DSK`: configura automaticamente `"-machine Gradiente_Expert_GPC-1 -ext DDX_3.0"`.
    * Outras extensões ficam com as opções em branco.

### 5.2. Personalização Visual com Temas
O sistema conta com um motor de temas nativo que customiza completamente a janela gráfica. Os seguintes esquemas de cores pré-configurados estão disponíveis para seleção instantânea na tela de configurações:
* **One Dark** (Padrão moderno escuro)
* **Dracula**
* **Monokai**
* **GitHub Light** (Tema claro)
* **Solarized Light**
* **VS Code Light**
* **Antigravity Dark**
* **MSX Classic** (Inspirado na clássica tela azul do MSX 1)
* **MSX Expert** (Inspirado no visual verde/cinza do Gradiente Expert)
* **MSX Cyber** (Visual futurista magenta/púrpura)

### 5.3. Gerenciamento Completo de Configurações (Aba de Preferências)
A tela de Configurações (`700x520`), acessível pelo menu superior *Configuração -> Programa*, é composta pelas seguintes abas funcionais:
* **Aparência**: Onde é possível escolher o Tema Visual, definir o **Volume Inicial** e escolher o **Idioma** da interface.
* **Caminhos**: Inputs interativos com assistentes de seleção de diretórios e arquivos para pastas do sistema (`Diretório Raiz`, `Imagens`, `DSK`, `ROM` e localização do `Banco de Dados`).
* **Executáveis**: Configuração dos executáveis locais dos emuladores compatíveis (*openMSX*, *fMSX*, *blueMSX* e *ruMSX*), permitindo buscar o arquivo `.exe` diretamente no Windows Explorer.
* **openMSX**: Definição da máquina padrão do openMSX (ex: `Gradiente_Expert_GPC-1`), opções livres de linha de comando e até 4 extensões padrão (ex: `DDX_3.0`).
* **fMSX / blueMSX / ruMSX**: Opções adicionais livres de linha de comando específicas para cada um desses emuladores.
* **Banco de Dados**: Botão interativo para reconstruir e resetar todas as tabelas do banco de dados a partir dos arquivos originais.

### 5.4. Interface Principal de Catálogo (Pesquisa e Navegação)
* **Seletor de Volumes**: Conjunto de 10 botões segmentados no topo do painel esquerdo que permite filtrar instantaneamente os jogos do volume selecionado. O botão do volume ativo recebe um destaque visual de alta importância.
* **Filtro Reativo**: Um campo de busca reativo que filtra a lista de jogos do volume selecionado em tempo real por descrição ou nome de arquivo.
* **Painel de Detalhes**: Exibição central/direita do jogo selecionado contendo a descrição com fonte em negrito e quebra de linha inteligente, visualização automática de screenshot contida na estrutura `pictures/XX/<raiz>.BMP` (ou `.bmp`) e botões interativos de ação rápida.

### 5.5. Configuração Individual de Emulação por Jogo
Ao selecionar um jogo e clicar em "Reconfigurar", uma janela de configurações individual é exibida:
* **Escolha de Emulador**: Dropdown para definir qual emulador será associado a este jogo específico.
* **Configuração Específica**:
  * Para o **openMSX**: permite selecionar uma máquina de MSX personalizada e até 4 extensões específicas (ex: mapeador MegaROM, drives de disco adicionais). O sistema busca interativamente a lista de máquinas e extensões instaladas no diretório `share/machines` e `share/extensions` do openMSX.
  * Para **fMSX/blueMSX/ruMSX**: inputs para opções de linha de comando customizadas.
* **Preview de Comando**: Exibe em tempo real o comando de terminal completo que será disparado ao iniciar o jogo.
* **Contador de Execuções**: Exibe a quantidade de vezes que o jogo foi executado a partir do catálogo e incrementa este número automaticamente a cada nova inicialização.

### 5.6. Suporte Multilíngue (i18n)
O sistema possui localização nativa para os seguintes idiomas:
* **Inglês** (Padrão)
* **Português (Brasil)**
* **Italiano**
* **Español**
* **Holandês**
A alteração é feita na aba Aparência e solicita que o usuário reinicie a aplicação para que os menus principais e janelas do sistema sejam completamente atualizados na nova linguagem.

### 5.7. Status do Sistema e Histórico de Logs
Acessível através do menu superior *Arquivo -> Status*, esta tela captura em tempo real:
* A inicialização do sistema e alteração de parâmetros.
* Detalhes de depuração sobre o comando disparado ao iniciar o jogo.
* A saída padrão (`stdout`) e de erros (`stderr`) dos emuladores em execução.
* Botões rápidos para "Copiar Tudo" para a área de transferência ou "Limpar" o histórico de logs.

### 5.8. Tela de Inicialização Animada (Splash Screen)
Ao inicializar, se suportado pelo driver de desktop do sistema operacional, uma janela splash sem bordas exibe a imagem promocional [`splash.png`](file:///e:/goMSXStuffs/images/splash.png) por 2 segundos e executa uma suave animação de fade-out antes de abrir a janela principal do catálogo.
