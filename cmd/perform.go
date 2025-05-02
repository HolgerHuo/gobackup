package cmd

import (
	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/model"
	"github.com/spf13/cobra"
)

var (
	modelName string
)

// performCmd represents the perform command
var performCmd = &cobra.Command{
	Use:   "perform",
	Short: "Perform backup operations",
	Run: func(cmd *cobra.Command, args []string) {
		config.Init(configFile)

		if len(modelName) == 0 {
			performAll()
		} else {
			performOne(modelName)
		}
	},
}

func init() {
	rootCmd.AddCommand(performCmd)

	// Define flags for the perform command
	performCmd.Flags().StringVarP(&modelName, "model", "m", "", "Model name that you want execute")
}

func performAll() {
	for _, modelConfig := range config.Models {
		m := model.Model{
			Config: modelConfig,
		}
		m.Perform()
	}
}

func performOne(modelName string) {
	for _, modelConfig := range config.Models {
		if modelConfig.Name == modelName {
			m := model.Model{
				Config: modelConfig,
			}
			m.Perform()
			return
		}
	}
}
