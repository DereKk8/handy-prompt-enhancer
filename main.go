package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/atotto/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() {
	setup := flag.Bool("setup", false, "Re-run configuration wizard")
	flag.Parse()
	if *setup {
		mainthread.Init(runWizardMain)
		return
	}
	mainthread.Init(run)
}

func runWizardMain() {
	cfg := &Config{Model: "gemini-flash-latest"}
	runWizard(cfg)
}

func run() {
	cfg, err := LoadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, ErrNoAPIKey) {
			cfg = &Config{Model: "gemini-flash-latest"}
			runWizard(cfg)
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
		notifier.Notify("Prompt Enhancer", "Clipboard read failed — see terminal for details")
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
		notifier.Notify("Prompt Enhancer", "Enhancement failed — see terminal for details")
		return
	}

	if err := clipboard.WriteAll(enhanced); err != nil {
		log.Printf("Clipboard write failed: %v", err)
		notifier.Notify("Prompt Enhancer", "Clipboard write failed — see terminal for details")
		return
	}

	log.Println("Prompt enhanced and copied to clipboard")
	notifier.Notify("Prompt Enhancer", "Prompt enhanced and copied to clipboard")
}
