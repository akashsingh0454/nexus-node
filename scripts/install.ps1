<#
.SYNOPSIS
One-line installer for Nexus Node Agent on Windows.

.DESCRIPTION
Downloads the latest release from GitHub, extracts it to C:\Program Files\NexusNode, and registers it.

.EXAMPLE
powershell -ExecutionPolicy Bypass -File install.ps1
#>

param(
    [string]$ConfigUrl = ""
)

$Repo = "YOUR_USERNAME/nexus-node" # PLACEHOLDER: Update this after forking!
$InstallDir = "C:\Program Files\NexusNode"
$Version = "latest"

Write-Host "🕷️  Installing Nexus Node Agent..." -ForegroundColor Cyan

# 1. Determine Architecture
$Arch = "amd64" # Expanding support later, default to amd64 for Windows
$AssetPattern = "nexus-windows-$Arch.zip"

# 2. Find Release URL
try {
    if ($Version -eq "latest") {
        $ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
    } else {
        $ApiUrl = "https://api.github.com/repos/$Repo/releases/tags/$Version"
    }
    
    $Release = Invoke-RestMethod -Uri $ApiUrl
    $Asset = $Release.assets | Where-Object { $_.name -like $AssetPattern } | Select-Object -First 1
    
    if (-not $Asset) {
        Write-Error "Could not find asset matching '$AssetPattern' in release $($Release.tag_name)"
        exit 1
    }
    
    $DownloadUrl = $Asset.browser_download_url
    Write-Host "   - Found version: $($Release.tag_name)"
    Write-Host "   - Downloading from: $DownloadUrl"
} catch {
    Write-Error "Failed to fetch release info: $_"
    exit 1
}

# 3. Download & Extract
$ZipPath = "$env:TEMP\nexus_installer.zip"
Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

Expand-Archive -Path $ZipPath -DestinationPath $InstallDir -Force
Remove-Item $ZipPath

# 4. Configure
Write-Host "   - Installed to: $InstallDir"

# Custom Config Override
if ($ConfigUrl -ne "") {
    Write-Host "   - Fetching custom configuration from: $ConfigUrl"
    try {
        Invoke-WebRequest -Uri $ConfigUrl -OutFile "$InstallDir\config.yaml"
    } catch {
        Write-Warning "Failed to download custom config. Using default."
    }
}

$ExePath = "$InstallDir\nexus.exe"

if (-not (Test-Path $ExePath)) {
    Write-Error "Installation failed: nexus.exe not found in extraction."
    exit 1
}

# 5. Service Registration
$ServiceName = "NexusNode"
if (Get-Service $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "   - Service already exists. stopping..."
    Stop-Service $ServiceName
    sc.exe delete $ServiceName | Out-Null
    Start-Sleep -Seconds 2
}

Write-Host "   - Registering Windows Service..."
# Use sc.exe for robust service creation
$BinPath = "`"$ExePath`" --config `"$InstallDir\config.yaml`""
sc.exe create $ServiceName binPath= $BinPath start= auto displayname= "Nexus Node Agent" | Out-Null
sc.exe description $ServiceName "Unified Security & Observability Agent" | Out-Null
sc.exe failure $ServiceName reset= 86400 actions= restart/60000/restart/60000/restart/60000 | Out-Null

Write-Host "   - Starting Service..."
Start-Service $ServiceName

Write-Host "✅ Installation Complete!" -ForegroundColor Green
Write-Host "   Agent is running as a service: $ServiceName"
Write-Host "   Status Check: & '$ExePath' --status" -ForegroundColor Yellow
