package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func SaveUserToCognito(username, password, email string) error {
	// Create a new session using your AWS credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"), 
	})
	if err != nil {
		return err
	}

	// Create a new CognitoIdentityProvider client.
	cognitoClient := cognitoidentityprovider.New(sess)

	// Specify the user pool ID and client ID.
	userPoolID := "your_user_pool_id"
	clientID := "your_client_id"

	// Set up the user attributes.
	userAttributes := []*cognitoidentityprovider.AttributeType{
		{
			Name:  aws.String("email"),
			Value: aws.String(email),
		},
	}

	// Create the user input.
	userInput := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:          aws.String(userPoolID),
		Username:            aws.String(username),
		TemporaryPassword:   aws.String(password),
		UserAttributes:      userAttributes,
		DesiredDeliveryMediums: []*string{
			aws.String("EMAIL"),
		},
		ForceAliasCreation: aws.Bool(true),
	}

	// Create the user in the user pool.
	_, err = cognitoClient.AdminCreateUser(userInput)
	if err != nil {
		return err
	}

	fmt.Println("User created successfully.")
	return nil
}