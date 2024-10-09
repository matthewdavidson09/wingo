package uptime

import (
	"embed"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/sys/windows"
	"gopkg.in/toast.v1"
)

// Configuration holds the application settings
type Configuration struct {
	UptimeThresholdWarning  int    `mapstructure:"uptime_threshold_warning"`  // Days before warning
	UptimeThresholdCritical int    `mapstructure:"uptime_threshold_critical"` // Days before forced restart
	CountdownMinutes        int    `mapstructure:"countdown_minutes"`         // Countdown before forced restart
	RestartCommand          string `mapstructure:"restart_command"`           // Command to execute for restart
}

//go:embed config.yaml
var configFile embed.FS

// LoadConfig loads configuration from the embedded config.yaml file
func LoadConfig() (*Configuration, error) {
	f, err := configFile.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error opening embedded config file: %w", err)
	}
	defer f.Close()

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(f) // Read directly from the opened file
	if err != nil {
		return nil, fmt.Errorf("error reading embedded config file: %w", err)
	}

	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}

// GetSystemUptime gets the system uptime in milliseconds using GetTickCount64
func GetSystemUptime() (time.Duration, error) {
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetTickCount64")
	ret, _, err := proc.Call()
	if ret == 0 {
		return 0, fmt.Errorf("GetTickCount64 failed: %w", err)
	}

	uptime := time.Duration(ret) * time.Millisecond
	return uptime, nil
}

// ForceRestart forces a system restart with a countdown notification
func ForceRestart(config *Configuration) {
	for remaining := config.CountdownMinutes; remaining > 0; remaining-- {
		message := fmt.Sprintf("Your system has been running for over %d days. "+
			"It will restart in %d minutes.", config.UptimeThresholdCritical, remaining)
		log.Println(message)

		notification := toast.Notification{
			AppID:   "Restart Notification",
			Title:   "Restart Notification - Urgent",
			Message: message,
		}

		if err := notification.Push(); err != nil {
			log.Printf("Failed to push notification: %v", err) // Log error and continue
		}

		time.Sleep(1 * time.Minute) // Wait 1 minute between updates
	}

	// Restart immediately after countdown
	parts := strings.Split(config.RestartCommand, " ")
	cmdName := parts[0]
	cmdArgs := parts[1:]

	// Run the PowerShell command
	cmd := exec.Command(cmdName, cmdArgs...)
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr
	if err != nil {
		log.Fatalf("Failed to restart the system: %v, Output: %s", err, output)
	} else {
		log.Printf("Restart command executed successfully. Output: %s", output)
	}
}

// PushNotification sends a toast notification with an option to restart
func PushNotification(config *Configuration, message string) {
	notification := toast.Notification{
		AppID:   "System Reminder",
		Title:   "System Uptime Warning",
		Message: message,
		Actions: []toast.Action{
			{Type: "protocol", Label: "Later", Arguments: "later"},
			{Type: "protocol", Label: "Restart Now", Arguments: "restart"},
		},
	}

	if err := notification.Push(); err != nil {
		log.Printf("Failed to push notification: %v", err)
	}
}

// CheckUptime checks the system uptime and decides whether to notify or restart
func CheckUptime(config *Configuration) {
	uptime, err := GetSystemUptime()
	if err != nil {
		log.Fatalf("Error getting system uptime: %v", err)
	}

	// Convert uptime to days
	uptimeInDays := int(uptime.Hours() / 24)
	fmt.Printf("System uptime: %d days\n", uptimeInDays)

	if uptimeInDays >= config.UptimeThresholdWarning && uptimeInDays < config.UptimeThresholdCritical {
		// Show toast notification offering to snooze or restart
		log.Println("Uptime is between warning and critical thresholds, offering to restart.")
		PushNotification(config, fmt.Sprintf("Your system has been up for over %d days. "+
			"Would you like to restart now or later?", config.UptimeThresholdWarning))
	} else if uptimeInDays >= config.UptimeThresholdCritical {
		// Notify and force a restart with countdown
		log.Println("Uptime is at or above critical threshold, starting countdown and forcing restart.")
		ForceRestart(config)
	} else {
		// Do nothing if uptime is less than the warning threshold
		log.Printf("Uptime is less than %d days, doing nothing.\n", config.UptimeThresholdWarning)
	}
}

func main() {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Check uptime and take action accordingly
	CheckUptime(config)
}
