param(
    [string]$InstallDir = "$env:USERPROFILE\tools\prompt-enhancer"
)

$exePath = Join-Path $InstallDir "prompt-enhancer.exe"
if (-not (Test-Path $exePath)) {
    Write-Host "prompt-enhancer not found in $InstallDir" -ForegroundColor Yellow
    Write-Host "Trying current directory..." -ForegroundColor Gray
    $exePath = Join-Path $PSScriptRoot "prompt-enhancer.exe"
    if (-not (Test-Path $exePath)) {
        Write-Host "prompt-enhancer.exe not found. Run install.ps1 first." -ForegroundColor Red
        exit 1
    }
}

Write-Host "Starting prompt-enhancer in background..." -ForegroundColor Green
$process = Start-Process -WindowStyle Hidden -FilePath $exePath -PassThru
Write-Host "PID: $($process.Id)" -ForegroundColor Cyan
Write-Host "Running in background (hidden window)." -ForegroundColor Green
Write-Host "To stop: Stop-Process -Id $($process.Id)" -ForegroundColor Yellow
