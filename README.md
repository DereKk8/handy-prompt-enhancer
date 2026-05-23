# prompt-enhancer

Takes raw speech-to-text (from [Handy](https://github.com/cjpais/Handy) or any source), sends it to Gemini, rewrites it into a structured coding-agent prompt, and copies the result to your clipboard — all with a global hotkey.

## How it works

1. Press your Handy shortcut, speak, Handy pastes raw text
2. Select the raw text and press `Ctrl+C`
3. Press `Ctrl+Alt+F12`
4. App reads clipboard → Gemini rewrites → structured prompt copied back
5. Paste into Cursor / Claude Code / Codex / Copilot

## Build

### From WSL2 / Linux (cross-compile for Windows)

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o prompt-enhancer.exe
```

### Native on Windows

```bash
go build -o prompt-enhancer.exe
```

## Install (Windows)

Copies the exe to `%USERPROFILE%\tools\prompt-enhancer\` and adds it to your user PATH so you can run `prompt-enhancer` from any terminal.

```powershell
# Run from the directory containing prompt-enhancer.exe
.\install.ps1

# Restart your terminal afterwards
```

## Usage

### Foreground (for testing)

```powershell
prompt-enhancer.exe
```

Press `Ctrl+Alt+F12` to enhance clipboard contents. Press `Ctrl+C` to quit.

### Background (for daily use)

```powershell
.\run-bg.ps1
```

This launches `prompt-enhancer.exe` in a hidden window. The hotkey still works. To stop it:

```powershell
Stop-Process -Name prompt-enhancer
```

### On startup (optional)

Place `run-bg.ps1` or a shortcut in `shell:startup` to auto-launch on login.

### Update API key

```powershell
prompt-enhancer.exe --setup
```

## API Key

1. Get a Gemini API key: https://aistudio.google.com/app/apikey
2. Run once — it will prompt for your key and save to `%APPDATA%\prompt-enhancer\config.toml`

```toml
# %APPDATA%\prompt-enhancer\config.toml
api_key = "your-gemini-api-key"
model = "gemini-flash-latest"
```

## Files

| File | Purpose |
|---|---|
| `main.go` | Hotkey loop, clipboard, orchestration |
| `llm.go` | LLM provider interface + Gemini |
| `config.go` | Config loading from `%APPDATA%` |
| `notify.go` | Windows toast notification |
| `install.ps1` | Adds exe to PATH and copies to tools dir |
| `run-bg.ps1` | Launches exe in background (hidden window) |

## Troubleshooting

### Toast notifications not showing
The app uses PowerShell to show toasts. If they don't appear, the app still logs to its terminal output.

### Clipboard issues
Make sure the target text is already copied before pressing `Ctrl+Alt+F12`. The app reads whatever is currently on the clipboard.

### Hotkey conflicts
If `Ctrl+Alt+F12` is used by another app, you can change the hotkey by editing `main.go` and recompiling.
