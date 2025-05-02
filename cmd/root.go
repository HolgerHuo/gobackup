package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	version    = "1.1.0-fork"
	configFile string
	verbose    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gobackup",
	Short:   "Easy full stack backup operations on UNIX-like systems",
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define persistent flags for the root command
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Special a config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	
	// Initialize logger with cobra's PreRun hook to ensure flags are parsed
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	logLevel := new(slog.LevelVar)
	
	// Set log level based on verbose flag
	if verbose {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	slog.SetDefault(slog.New(handler))
}
