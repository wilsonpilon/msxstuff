# Create-DskFiles.ps1

param(
    [string]$InputFile = "games.txt",
    [string]$OutputDir = "."
)

if (-not (Test-Path $InputFile)) {
    Write-Error "Arquivo nao encontrado: $InputFile"
    exit 1
}

# Cria pastas A-Z corretamente
65..90 | ForEach-Object {
    $letter = [char]$_
    $dir = Join-Path $OutputDir $letter
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir | Out-Null
    }
}

# Pasta para nomes que nao comecam com letra
$otherDir = Join-Path $OutputDir "_outros"
if (-not (Test-Path $otherDir)) {
    New-Item -ItemType Directory -Path $otherDir | Out-Null
}

$created = 0
$skipped = 0
$lineNum = 0

Get-Content $InputFile | ForEach-Object {
    $lineNum++
    $line = $_.Trim()

    if ([string]::IsNullOrWhiteSpace($line)) { return }

    $fields = $line -split ','
    if ($fields.Count -lt 3) {
        Write-Warning "Linha $lineNum ignorada (menos de 3 campos): $line"
        $skipped++
        return
    }

    # Terceiro campo (indice 2) e o nome do jogo
    $gameName = $fields[2].Trim()

    if ([string]::IsNullOrWhiteSpace($gameName)) {
        Write-Warning "Linha $lineNum ignorada (nome vazio)"
        $skipped++
        return
    }

    $safeName = $gameName -replace '[\\/:*?"<>|]', '_'

    $firstChar = $safeName.Substring(0, 1).ToUpper()
    if ($firstChar -match '^[A-Z]$') {
        $targetDir = Join-Path $OutputDir $firstChar
    } else {
        $targetDir = $otherDir
    }

    $filePath = Join-Path $targetDir "$safeName.dsk"

    if (-not (Test-Path $filePath)) {
        New-Item -ItemType File -Path $filePath | Out-Null
        Write-Host "Criado: $filePath"
        $created++
    } else {
        $skipped++
    }
}

Write-Host ""
Write-Host "Concluido. Arquivos criados: $created | Ignorados/duplicados: $skipped"
