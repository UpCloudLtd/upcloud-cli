package config

import (
	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
	"github.com/jedib0t/go-pretty/v6/text"
)

func GetProgressOutputConfig(cfg *Config) *progress.OutputConfig {
	config := progress.GetDefaultOutputConfig()

	config.StatusColorMap = map[messages.MessageStatus]messages.Color{
		messages.MessageStatusSuccess: text.FgHiGreen,
		messages.MessageStatusWarning: text.FgHiYellow,
		messages.MessageStatusError:   text.FgHiRed,
		messages.MessageStatusStarted: text.FgHiBlue,
		messages.MessageStatusPending: text.FgHiCyan,
		messages.MessageStatusSkipped: text.FgHiMagenta,
	}
	config.ColorMessage = true

	// We can not update progress animations when debug messages are outputted to the same output stream. Thus, disable animations when debug output is enabled.
	config.DisableAnimations = cfg.GlobalFlags.Debug

	return config
}
