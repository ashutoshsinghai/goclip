//go:build windows

package tray

import "os/exec"

func openPicker() {
	exe := goclipExe()
	exec.Command("cmd", "/C", "start", "cmd", "/K", exe+" pick").Start() //nolint:errcheck
}

func notifyCopied(preview string) {
	// PowerShell toast notification (Windows 10+)
	script := `[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType=WindowsRuntime] | Out-Null
$t = [Windows.UI.Notifications.ToastTemplateType]::ToastText02
$xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent($t)
$xml.GetElementsByTagName('text')[0].AppendChild($xml.CreateTextNode('goclip')) | Out-Null
$xml.GetElementsByTagName('text')[1].AppendChild($xml.CreateTextNode('Copied to clipboard')) | Out-Null
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('goclip').Show($toast)`
	exec.Command("powershell", "-Command", script).Start() //nolint:errcheck
}
