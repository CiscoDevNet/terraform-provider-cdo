package acctest

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type CdoSecretManager struct {
	region string
	client *secretsmanager.Client
}

func NewCdoSecretManager(region string) *CdoSecretManager {

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	return &CdoSecretManager{
		region: region,
		client: svc,
	}
}

func (manager *CdoSecretManager) getSecretValue(inp *secretsmanager.GetSecretValueInput) (string, error) {

	result, err := manager.client.GetSecretValue(context.TODO(), inp)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", err
	}

	return *result.SecretString, nil
}

func (manager *CdoSecretManager) getCurrentSecretValue(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	return manager.getSecretValue(input)
}
