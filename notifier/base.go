package notifier

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/holgerhuo/gobackup/config"
	"github.com/spf13/viper"
)

type Base struct {
	viper     *viper.Viper
	Name      string
	onSuccess bool
	onFailure bool
}

type Notifier interface {
	notify(title, message string) error
}

var (
	notifyTypeSuccess = 1
	notifyTypeFailure = 2
)

func newNotifier(name string, config config.SubConfig) (Notifier, *Base, error) {
	base := &Base{
		viper: config.Viper,
		Name:  name,
	}
	base.viper.SetDefault("on_success", true)
	base.viper.SetDefault("on_failure", true)

	base.onSuccess = base.viper.GetBool("on_success")
	base.onFailure = base.viper.GetBool("on_failure")

	switch config.Type {
	case "webhook":
		return NewWebhook(base), base, nil
	}
	return nil, nil, fmt.Errorf("Notifier: %s is not supported", name)
}

func notify(model config.ModelConfig, title, message string, notifyType int) {
	notifyTypeStr := "failure"
	if notifyType == notifyTypeSuccess {
		notifyTypeStr = "success"
	}
	
	slog.Info("Running notifiers", 
		"component", "notifier",
		"model", model.Name,
		"count", len(model.Notifiers),
		"type", notifyTypeStr)
	
	for name, config := range model.Notifiers {
		notifier, base, err := newNotifier(name, config)
		if err != nil {
			slog.Error("Failed to initialize notifier", 
				"component", "notifier",
				"model", model.Name,
				"notifier", name,
				"error", err)
			continue
		}

		if notifyType == notifyTypeSuccess {
			if base.onSuccess {
				if err := notifier.notify(title, message); err != nil {
					slog.Error("Notification failed", 
						"component", "notifier",
						"model", model.Name,
						"notifier", name,
						"type", "success",
						"error", err)
				}
			}
		} else if notifyType == notifyTypeFailure {
			if base.onFailure {
				if err := notifier.notify(title, message); err != nil {
					slog.Error("Notification failed", 
						"component", "notifier",
						"model", model.Name,
						"notifier", name,
						"type", "failure",
						"error", err)
				}
			}
		}
	}
}

func Success(model config.ModelConfig) {
	title := fmt.Sprintf("[GoBackup] OK: Backup %s has successfully", model.Name)
	message := fmt.Sprintf("Backup of %s completed successfully at %s", model.Name, time.Now().Local())
	notify(model, title, message, notifyTypeSuccess)
}

func Failure(model config.ModelConfig, reason string) {
	title := fmt.Sprintf("[GoBackup] Err: Backup %s has failed", model.Name)
	message := fmt.Sprintf("Backup of %s failed at %s:\n\n%s", model.Name, time.Now().Local(), reason)

	notify(model, title, message, notifyTypeFailure)
}
