package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func saveUserAndRefreshTokenToDb(userId string, email string, refreshToken string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

	ddb := dynamodb.NewFromConfig(cfg)

	user := User{
		PartitionKey: "USER",
		SortKey: userId,
		RefreshToken: refreshToken,
		Email: email,
	}
	
	ddbitem, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal User: %w", err)
	}

	// Create the DynamoDB input parameters
	params := &dynamodb.PutItemInput{
		Item:      ddbitem,
		TableName: aws.String("YourDynamoDBTableName"), // Replace with your DynamoDB table name
	}

	// Save the item to DynamoDB
	_, err = ddb.PutItem(context.TODO(), params)
	if err != nil {
		fmt.Println("Failed to save item to DynamoDB:", err)
		return err
	}
	
	fmt.Println("Item saved to DynamoDB successfully")
	return nil
}