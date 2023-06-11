package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func saveUserToCognito(id string, email string) error {
	fmt.Println("Saving user to cognito...")

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
		Username: &email, // TODO Should we use ID instead?
		DesiredDeliveryMediums: []types.DeliveryMediumType{types.DeliveryMediumTypeEmail},
		ForceAliasCreation: false,
	}

	// Create the user in the user pool.
	_, err = cognitoClient.AdminCreateUser(context.TODO(), userInput)
	if err != nil {
		// Don't mind if the username exists
		var usernameExistsErr *types.UsernameExistsException
		if errors.As(err, &usernameExistsErr) {
			fmt.Println("Username already exists. Ignoring error.")
			return nil // Return nil to indicate success without error.
		}
		fmt.Println("Error saving user in Cognito user pool:", err)
		return err
	}

	fmt.Println("User created successfully in Cognito.")
	return nil
}