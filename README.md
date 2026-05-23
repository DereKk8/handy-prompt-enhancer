# prompt-enhancer

Takes raw speech-to-text (from [Handy](https://github.com/cjpais/Handy) or any source), sends it to Gemini, rewrites it into a structured coding-agent prompt, and copies the result to your clipboard — all with a global hotkey.

## How it works

1. Press your Handy shortcut, speak, Handy pastes raw text
2. Select the raw text and press `Ctrl+C`
3. Press `Ctrl+Alt+F12`
4. App reads clipboard → Gemini rewrites → structured prompt copied back
5. Paste into Cursor / Claude Code / Codex / Copilot

## Build

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o prompt-enhancer.exe .
```

## First Run

Just run `prompt-enhancer.exe` — it launches an interactive setup wizard that asks:

- Your Gemini API key
- Install directory (default: `%USERPROFILE%\tools\prompt-enhancer`)
- Whether to add to system PATH (run from any terminal)
- Whether to launch on Windows startup (background)
- Whether to launch now in the background

## Usage

```powershell
# Normal (foreground terminal)
prompt-enhancer.exe

# Or if launched in background, hotkey still works silently
```

Press `Ctrl+Alt+F12` to enhance clipboard contents. Press `Ctrl+C` to quit (foreground mode).

### Re-run setup

```powershell
prompt-enhancer.exe --setup
```

## API Key

Get one at https://aistudio.google.com/app/apikey

Config saved to `%APPDATA%\prompt-enhancer\config.toml`:

```toml
api_key = "your-key"
model = "gemini-flash-latest"
```

## Files

| File | Purpose |
|---|---|
| `main.go` | Hotkey loop, clipboard, orchestration |
| `llm.go` | Gemini provider |
| `config.go` | Config loading/saving |
| `notify.go` | Windows toast notification |
| `wizard.go` | Setup wizard (PATH, startup, launch) |

## Troubleshooting

### Toast notifications not showing
The app uses PowerShell to show toasts. If they don't appear, the app still logs to its terminal output.

### Clipboard issues
Make sure the target text is already copied before pressing `Ctrl+Alt+F12`.

### Hotkey conflicts
If `Ctrl+Alt+F12` is used by another app, change it in `main.go` and recompile.
