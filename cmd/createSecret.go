package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
)

var createSecretCmd = &cobra.Command{
	Use:   "createSecret",
	Short: "Used to convert a .env file into an AWS Secret",
	Long:  `Converts a .env file into a JSON string and creates or updates a given aws secret with that value`,
	Run:   CreateSecret,
}

func BuildJSONStringFromEnv(data []string) string {
	jsonEnvFile := "{"
	for i, entry := range data {
		splitEntry := strings.Split(entry, "=")
		jsonEnvFile += fmt.Sprintf("\"%s\":\"%s\"", splitEntry[0], strings.Join(splitEntry[1:], ""))
		if i != len(data)-1 {
			jsonEnvFile += ","
		} else {
			jsonEnvFile += "}"
		}
	}

	return jsonEnvFile
}

func CreateSecret(cmd *cobra.Command, args []string) {
	region, _ := cmd.Flags().GetString("region")

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	client := secretsmanager.NewFromConfig(cfg)

	secretName, _ := cmd.Flags().GetString("secretName")

	envFilePath, _ := cmd.Flags().GetString("envFile")

	workingDir, _ := os.Getwd()
	file, err := os.Open(workingDir + "/" + envFilePath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var perline string
	envFileData := make([]string, 0)

	for {

		_, err := fmt.Fscanf(file, "%v\n", &perline) // give a patter to scan

		if err != nil {

			if err == io.EOF {
				break // stop reading the file
			}
			fmt.Println(err)
			os.Exit(1)
		}

		envFileData = append(envFileData, perline)
	}

	stringEnvFile := BuildJSONStringFromEnv(envFileData)

	newSecret, _ := cmd.Flags().GetBool("newSecret")

	if newSecret {
		log.Printf("Creating secret %s", secretName)
		_, err := client.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
			Name:         &secretName,
			SecretString: &stringEnvFile,
		})
		if err != nil {
			log.Fatalf("Error creating secret %v", err)
		}
	} else {
		log.Printf("Updating secret %s", secretName)
		_, err = client.PutSecretValue(context.TODO(), &secretsmanager.PutSecretValueInput{
			SecretId:     &secretName,
			SecretString: &stringEnvFile,
		})
		if err != nil {
			log.Fatalf("Error updating secret %v", err)
		}
	}

}

func init() {
	cobra.OnInitialize()
	rootCmd.AddCommand(createSecretCmd)
	createSecretCmd.Flags().StringP("region", "r", "us-east-1", "The AWS region the Secret is located")
	createSecretCmd.Flags().StringP("secretName", "s", "test", "The Secret name")
	createSecretCmd.Flags().BoolP("newSecret", "n", false, "If the Secret is new or not")
	createSecretCmd.Flags().StringP("envFile", "e", ".env", "The name of the env file")

}
