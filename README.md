# secretsmanager-to-env

Generates a `.env` file using an AWS SecretsManager secret

## Installation

`go install github.com/whutchinson98/secretsmanager-to-env@latest`

## Local Setup

`go build` - to build the application within the repo

`go install secretsmanager-to-env` to install the application inside of your go packages to make it available throughout your system

## Example Usage

`secretsmanager-to-env createEnv -s api-secret-production -r us-east-2 -e .env.production -p ../`
