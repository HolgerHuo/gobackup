package cmd

import (
	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/model"
	"github.com/spf13/cobra"
)

var (
	modelName string
)

var performCmd = &cobra.Command{
	Use:   "perform",
	Short: "perform backup",
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

	performCmd.Flags().StringVarP(&modelName, "model", "m", "", "the model to perform")
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
	modelConfig := config.GetModelByName(modelName)
	if modelConfig != nil {
		m := model.Model{
			Config: *modelConfig,
		}
		m.Perform()
	}
}
