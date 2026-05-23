package main

import (
	"fmt"
	"os/exec"
)

type Notifier interface {
	Notify(title, message string) error
}

type ToastNotifier struct{}

func (t *ToastNotifier) Notify(title, message string) error {
	return showWindowsToast(title, message)
}

func showWindowsToast(title, message string) error {
	script := fmt.Sprintf(`
$xml = [Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$textNodes = $xml.GetElementsByTagName("text")
$textNodes.Item(0).AppendChild($xml.CreateTextNode("%s")) > $null
$textNodes.Item(1).AppendChild($xml.CreateTextNode("%s")) > $null
$toast = [Windows.UI.Notifications.ToastNotification]::New($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Prompt Enhancer").Show($toast)
`, escapePowershell(title), escapePowershell(message))

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("toast notification failed: %w", err)
	}
	return nil
}

func escapePowershell(s string) string {
	escaped := ""
	for _, c := range s {
		switch c {
		case '"':
			escaped += "`\""
		case '`':
			escaped += "``"
		case '\n':
			escaped += " "
		case '\r':
			escaped += ""
		default:
			escaped += string(c)
		}
	}
	return escaped
}
