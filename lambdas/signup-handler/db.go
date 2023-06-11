package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func saveUserAndRefreshTokenToDb(userId string, email string, refreshToken string) error {
	fmt.Println("Saving user and refresh token to DB...")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

	ddb := dynamodb.NewFromConfig(cfg)

	user := User{
		PartitionKey: "USER",
		SortKey: email,
		RefreshToken: refreshToken,
		Username: userId,
	}
	
	ddbitem, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal User: %w", err)
	}

	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
	// Create the DynamoDB input parameters
	params := &dynamodb.PutItemInput{
		Item:      ddbitem,
		TableName: aws.String(ddbTableName), // Replace with your DynamoDB table name
	}

	// Save the item to DynamoDB
	_, err = ddb.PutItem(context.TODO(), params)
	if err != nil {
		fmt.Println("Failed to save user data to DynamoDB:", err)
		return err
	}
	
	fmt.Println("User data saved to DynamoDB successfully")
	return nil
}