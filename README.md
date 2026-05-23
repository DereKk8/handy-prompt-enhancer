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

## Setup

1. Get a Gemini API key: https://aistudio.google.com/app/apikey
2. Run `prompt-enhancer.exe` once — it will prompt for your API key and create the config file
3. Or manually create `%APPDATA%\prompt-enhancer\config.toml`:

```toml
api_key = "your-gemini-api-key"
model = "gemini-flash-latest"
```

## Usage

```bash
prompt-enhancer.exe
```

The app runs in the foreground. Press `Ctrl+Alt+F12` to enhance clipboard contents. Press `Ctrl+C` to quit.

To update your API key or model, run: `prompt-enhancer.exe --setup`

## Files

| File | Purpose |
|---|---|
| `main.go` | Hotkey loop, clipboard, orchestration |
| `llm.go` | LLM provider interface + Gemini |
| `config.go` | Config loading from `%APPDATA%` |
| `notify.go` | Windows toast notification |

## Dependencies

- [golang.design/x/hotkey](https://golang.design/x/hotkey) — global hotkey
- [github.com/atotto/clipboard](https://github.com/atotto/clipboard) — clipboard
- [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml) — TOML parsing
- [golang.org/x/sys](https://golang.org/x/sys) — Windows API calls

## Troubleshooting

### Toast notifications not showing
The app uses PowerShell to show toasts. If they don't appear, the app still logs to its terminal output.

### Clipboard issues
Make sure the target text is already copied before pressing `Ctrl+Alt+P`. The app reads whatever is currently on the clipboard.

### Hotkey conflicts
If `Ctrl+Alt+F12` is used by another app, you can change the hotkey by editing `main.go` and recompiling.
