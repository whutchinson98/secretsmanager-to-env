package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "secretsmanager-to-env",
	Short: "Utilities to work with AWS SecretsManager and Env Files",
	Long: `This package is used to quickly pull down a JSON secret from AWS and create a .env file for it for your
         projects that do not use the secrets via aws-sdk`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
