package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
)

// createEnvCmd represents the createEnv command
var createEnvCmd = &cobra.Command{
	Use:   "createEnv",
	Short: "Creates a .env file from a Secret",
	Long: `Creates a .env file from a Secret taking in the region and secret name.
          In the event that you already have a .env created in the location you choose
          the new one to be it WILL override that so be careful.`,
	Run: CreateEnvFile,
}

func BuildEnvFileString(secretMap map[string]interface{}) string {
	result := ""
	for k := range secretMap {
		secretValue := k + "=" + secretMap[k].(string) + "\n"
		result = result + secretValue
	}

	return result
}

var osGetWd = os.Getwd
var osCreate = os.Create

func InitEnvFile(path string, envFileName string) (*os.File, error) {

	wd, err := osGetWd()

	if err != nil {
		panic(err)
	}

	if !bytes.HasSuffix([]byte(path), []byte("/")) {
		path = path + "/"
	}

	envFile, err := osCreate(wd + "/" + path + envFileName)

	return envFile, err
}

func CreateEnvFile(cmd *cobra.Command, args []string) {
	region, _ := cmd.Flags().GetString("region")

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	client := secretsmanager.NewFromConfig(cfg)

	secretName, _ := cmd.Flags().GetString("secretName")

	secVal, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})

	if err != nil {
		log.Fatalf("Error fetching secret value: %v", err)
	}

	secretString := aws.ToString(secVal.SecretString)

	sec := map[string]interface{}{}

	if err := json.Unmarshal([]byte(secretString), &sec); err != nil {
		log.Fatalf("Error converting secret value to JSON: %v", err)
	}

	envString := BuildEnvFileString(sec)

	path, _ := cmd.Flags().GetString("path")
	envFileName, _ := cmd.Flags().GetString("envFile")

	envFile, err := InitEnvFile(path, envFileName)

	if err != nil {
		log.Fatalf("Error creating the env file with path %v file name %v: Error %v", path, envFileName, err)
	}

	_, err = envFile.WriteString(envString)

	if err != nil {
		log.Fatalf("Error writing to env file: %v", err)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.AddCommand(createEnvCmd)
	createEnvCmd.Flags().StringP("region", "r", "us-east-1", "The AWS region the Secret is located")
	createEnvCmd.Flags().StringP("secretName", "s", "test", "The Secret name")
	createEnvCmd.Flags().StringP("envFile", "e", ".env", "The name of the env file")
	createEnvCmd.Flags().StringP("path", "p", "./", "The relative path to place the env file in")
}
