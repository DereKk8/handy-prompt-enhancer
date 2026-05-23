package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/atotto/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() {
	mainthread.Init(run)
}

func run() {
	cfg, err := LoadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, ErrNoAPIKey) {
			cfg = &Config{Model: "gemini-2.0-flash"}
			showConfigSetup(cfg)
			return
		}
		log.Fatalf("Config error: %v", err)
	}

	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}, hotkey.KeyF12)
	if err := hk.Register(); err != nil {
		log.Fatalf("Failed to register hotkey: %v", err)
	}
	log.Println("prompt-enhancer running — press Ctrl+Alt+F12 to enhance clipboard")
	log.Println("Press Ctrl+C to quit")

	notifier := &ToastNotifier{}
	provider := &GeminiProvider{
		APIKey: cfg.APIKey,
		Model:  cfg.Model,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for {
		select {
		case <-hk.Keydown():
			handleEnhance(provider, notifier)
		case <-sigCh:
			log.Println("Shutting down...")
			return
		}
	}
}

func handleEnhance(provider LLMProvider, notifier Notifier) {
	raw, err := clipboard.ReadAll()
	if err != nil {
		log.Printf("Clipboard read failed: %v", err)
		notifier.Notify("Prompt Enhancer", fmt.Sprintf("Clipboard read failed: %v", err))
		return
	}
	if raw == "" {
		notifier.Notify("Prompt Enhancer", "Clipboard is empty — copy text first")
		return
	}

	log.Println("Enhancing prompt...")
	enhanced, err := provider.Enhance(raw)
	if err != nil {
		log.Printf("Enhancement failed: %v", err)
		notifier.Notify("Prompt Enhancer", fmt.Sprintf("Enhancement failed: %v", err))
		return
	}

	if err := clipboard.WriteAll(enhanced); err != nil {
		log.Printf("Clipboard write failed: %v", err)
		notifier.Notify("Prompt Enhancer", fmt.Sprintf("Clipboard write failed: %v", err))
		return
	}

	log.Println("Prompt enhanced and copied to clipboard")
	notifier.Notify("Prompt Enhancer", "Prompt enhanced and copied to clipboard")
}

func showConfigSetup(cfg *Config) {
	fmt.Println("=== Prompt Enhancer Setup ===")
	fmt.Println()
	fmt.Println("Config file not found or incomplete.")
	fmt.Println("Create %APPDATA%\\prompt-enhancer\\config.toml with:")
	fmt.Println()
	fmt.Println(`  api_key = "your-gemini-api-key"`)
	fmt.Println(`  model = "gemini-2.0-flash"`)
	fmt.Println()
	fmt.Println("Get a Gemini API key at: https://aistudio.google.com/app/apikey")
	fmt.Println()
	fmt.Print("Enter your Gemini API key: ")

	var key string
	fmt.Scanln(&key)
	if key == "" {
		log.Fatal("API key required")
	}
	cfg.APIKey = key
	cfg.Model = "gemini-2.0-flash"

	if err := WriteConfig(cfg); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}
	fmt.Println("Config saved. Restart prompt-enhancer.")
}
