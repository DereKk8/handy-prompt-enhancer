package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func runWizard(cfg *Config) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Prompt Enhancer Setup ===")
	fmt.Println()

	if cfg.APIKey == "" {
		fmt.Print("Enter your Gemini API key: ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		if key == "" {
			fmt.Println("No key entered. Exiting.")
			return
		}
		cfg.APIKey = key
	}

	installDir := promptInstallDir(reader)
	copyExe(installDir)
	promptPathInstall(reader, installDir)
	promptStartupInstall(reader, installDir)
	promptLaunchNow(reader, installDir)

	if err := WriteConfig(cfg); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}
	fmt.Println("Setup complete!")
}

func promptInstallDir(reader *bufio.Reader) string {
	defaultDir := filepath.Join(os.Getenv("USERPROFILE"), "tools", "prompt-enhancer")
	fmt.Printf("Install directory [%s]: ", defaultDir)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultDir
	}
	return input
}

func copyExe(dir string) {
	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("  Cannot determine executable path: %v\n", err)
		return
	}
	os.MkdirAll(dir, 0755)
	dst := filepath.Join(dir, "prompt-enhancer.exe")
	if err := copyFile(exe, dst); err != nil {
		fmt.Printf("  Failed to copy exe: %v\n", err)
		return
	}
	fmt.Printf("  Installed to %s\n", dir)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}

func promptPathInstall(reader *bufio.Reader, installDir string) {
	fmt.Print("Add to system PATH so you can run it from any terminal? (Y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "n" || input == "no" {
		fmt.Println("  Skipped PATH installation.")
		return
	}
	if err := addToPath(installDir); err != nil {
		fmt.Printf("  Failed to update PATH: %v\n", err)
		return
	}
	fmt.Println("  Added to PATH (user-wide). Restart terminal to apply.")
}

func promptStartupInstall(reader *bufio.Reader, installDir string) {
	fmt.Print("Launch on Windows startup in the background? (Y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "n" || input == "no" {
		fmt.Println("  Skipped startup registration.")
		return
	}
	if err := addToStartup(installDir); err != nil {
		fmt.Printf("  Failed to register startup: %v\n", err)
		return
	}
	fmt.Println("  Added to Windows startup.")
}

func promptLaunchNow(reader *bufio.Reader, installDir string) {
	fmt.Print("Launch now in the background? (Y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "n" || input == "no" {
		fmt.Println("  Skipped launch.")
		return
	}
	launchBackground(installDir)
}

func addToPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("open registry: %w", err)
	}
	defer k.Close()

	current, _, err := k.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("read PATH: %w", err)
	}

	entries := strings.Split(current, ";")
	for _, e := range entries {
		if strings.EqualFold(strings.TrimSpace(e), dir) {
			return nil
		}
	}

	newPath := current
	if newPath != "" && !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir

	if err := k.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("write PATH: %w", err)
	}
	return nil
}

func addToStartup(dir string) error {
	startupDir := filepath.Join(os.Getenv("APPDATA"), `Microsoft\Windows\Start Menu\Programs\Startup`)
	exePath := filepath.Join(dir, "prompt-enhancer.exe")
	vbsPath := filepath.Join(startupDir, "prompt-enhancer.vbs")

	vbsContent := fmt.Sprintf(
		`Set WshShell = CreateObject("WScript.Shell")
WshShell.Run "%s", 0, False
`, exePath)

	if err := os.WriteFile(vbsPath, []byte(vbsContent), 0644); err != nil {
		return fmt.Errorf("write startup script: %w", err)
	}
	return nil
}

func launchBackground(dir string) {
	exePath := filepath.Join(dir, "prompt-enhancer.exe")
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Start-Process -WindowStyle Hidden -FilePath "%s"`, exePath))
	cmd.Start()
	fmt.Println("  Launched in background (hidden window).")
}

func removeStartup() error {
	startupDir := filepath.Join(os.Getenv("APPDATA"), `Microsoft\Windows\Start Menu\Programs\Startup`)
	vbsPath := filepath.Join(startupDir, "prompt-enhancer.vbs")
	if err := os.Remove(vbsPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func removeFromPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("open registry: %w", err)
	}
	defer k.Close()

	current, _, err := k.GetStringValue("Path")
	if err != nil {
		return nil
	}

	entries := strings.Split(current, ";")
	var kept []string
	for _, e := range entries {
		if !strings.EqualFold(strings.TrimSpace(e), dir) {
			kept = append(kept, e)
		}
	}
	newPath := strings.Join(kept, ";")
	return k.SetStringValue("Path", newPath)
}
