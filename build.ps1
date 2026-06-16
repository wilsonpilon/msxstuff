# build.ps1
# Script de automação para compilar e executar o projeto MSX Stuffs em Go.

[CmdletBinding(PositionalBinding=$false)]
param(
    [ValidateSet("windows", "linux")]
    [string]$OS = "windows",

    [ValidateSet("Release", "Debug")]
    [string]$Mode = "Release",

    [switch]$Run,

    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$AppArgs
)

# Habilitar strict mode e interromper em caso de erros
$ErrorActionPreference = "Stop"

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host " MSX Stuffs - Script de Compilação" -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# Priorizar compiladores modernos no Windows se estiverem instalados
if ($OS -eq "windows") {
    $MSYS2Paths = @(
        "C:\msys64\ucrt64\bin",
        "C:\msys64\mingw64\bin"
    )
    foreach ($Path in $MSYS2Paths) {
        if (Test-Path $Path) {
            $env:PATH = "$Path;" + $env:PATH
            break
        }
    }
}

# 1. Organizar go.mod e go.sum
Write-Host "[1/3] Organizando dependências do Go (go mod tidy)..." -ForegroundColor Yellow
go mod tidy
if ($LastExitCode -ne 0) {
    Write-Error "A organização das dependências (go mod tidy) falhou."
    exit $LastExitCode
}

# 2. Configurar compilação cruzada
$BinaryName = if ($OS -eq "windows") { "msxstuffs.exe" } else { "msxstuffs" }
Write-Host "[2/3] Configurando ambiente para OS: $OS | Modo: $Mode..." -ForegroundColor Yellow

$env:GOOS = $OS
# Fyne exige CGO para bindings gráficos (OpenGL/GLFW), portanto mantemos o padrão do Go.

# Configurar flags de compilação
$BuildFlags = @()
if ($Mode -eq "Release") {
    Write-Host "      -> Compilando em modo Release (removendo símbolos de depuração)..." -ForegroundColor Gray
    if ($OS -eq "windows") {
        $BuildFlags += "-ldflags=-s -w -H=windowsgui"
    } else {
        $BuildFlags += "-ldflags=-s -w"
    }
} else {
    Write-Host "      -> Compilando em modo Debug..." -ForegroundColor Gray
    if ($OS -eq "windows") {
        $BuildFlags += "-ldflags=-H=windowsgui"
    }
}

# Compilar
Write-Host "      -> Executando 'go build'..." -ForegroundColor Gray
go build $BuildFlags -o $BinaryName
if ($LastExitCode -ne 0) {
    Write-Error "A compilação falhou com código de saída $LastExitCode."
    exit $LastExitCode
}

Write-Host "[3/3] Compilação concluída com sucesso! Binário gerado: $BinaryName" -ForegroundColor Green

# 3. Executar o programa se o parâmetro -Run estiver presente
if ($Run) {
    if ($OS -ne "windows") {
        Write-Warning "Não é possível executar um binário do Linux ($OS) diretamente no host Windows."
        Write-Warning "Por favor, execute o binário manualmente em um ambiente compatível."
    } else {
        Write-Host "==================================================" -ForegroundColor Cyan
        Write-Host " Executando programa: $BinaryName" -ForegroundColor Cyan
        if ($AppArgs) {
            Write-Host " Parâmetros passados: $AppArgs" -ForegroundColor Gray
        }
        Write-Host "==================================================" -ForegroundColor Cyan
        
        # Executa o executável passando os argumentos adicionais
        & ".\$BinaryName" $AppArgs
    }
}
