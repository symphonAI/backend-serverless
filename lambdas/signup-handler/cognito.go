package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func SaveUserToCognito(id string, email string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

	// Create a new CognitoIdentityProvider client.
	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	// Specify the user pool ID and client ID.
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")

	// Set up the user attributes.
	userAttributes := []types.AttributeType{
		{
			Name:  aws.String("email"),
			Value: aws.String(email),
		},
	}

	// Create the user input.
	userInput := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:          aws.String(userPoolID),
		UserAttributes:      userAttributes,
		DesiredDeliveryMediums: types.DeliveryMediumTypeEmail.Values(),
		ForceAliasCreation: false,
	}

	// Create the user in the user pool.
	_, err = cognitoClient.AdminCreateUser(context.TODO(), userInput)
	if err != nil {
		return err
	}

	fmt.Println("User created successfully.")
	return nil
}