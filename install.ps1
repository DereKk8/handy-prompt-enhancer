param(
    [string]$InstallDir = "$env:USERPROFILE\tools\prompt-enhancer"
)

$exePath = Join-Path $PSScriptRoot "prompt-enhancer.exe"
if (-not (Test-Path $exePath)) {
    Write-Host "prompt-enhancer.exe not found next to this script." -ForegroundColor Red
    exit 1
}

New-Item -Path $InstallDir -ItemType Directory -Force | Out-Null
Copy-Item -Path $exePath -Destination $InstallDir -Force
Write-Host "Copied prompt-enhancer.exe to $InstallDir" -ForegroundColor Green

$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    $newPath = "$currentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added $InstallDir to your PATH (user-wide)." -ForegroundColor Green
    Write-Host "Restart your terminal for the change to take full effect." -ForegroundColor Yellow
} else {
    Write-Host "$InstallDir is already in your PATH." -ForegroundColor Cyan
}

Write-Host ""
Write-Host "Done. From any terminal, run:" -ForegroundColor White
Write-Host "  prompt-enhancer.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "To run in the background (hidden window):" -ForegroundColor White
Write-Host "  Start-Process -WindowStyle Hidden -FilePath prompt-enhancer.exe" -ForegroundColor Cyan
