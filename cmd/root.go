package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	version    = "dev"
	configFile string
	debug      bool
	jsonLog    bool
)

var rootCmd = &cobra.Command{
	Use:     "gobackup",
	Short:   "Easy backup solution on Linux",
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to config file")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable verbose log")
	rootCmd.PersistentFlags().BoolVar(&jsonLog, "json", false, "output logs in json format")
	
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	logLevel := new(slog.LevelVar)

	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	var handler slog.Handler
	if jsonLog {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
			AddSource: debug,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
			AddSource: debug,
		})
	}
	slog.SetDefault(slog.New(handler))
}
