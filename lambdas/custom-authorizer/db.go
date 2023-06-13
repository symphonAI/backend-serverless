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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func getUserFromDB(us string) (user *User, err error){
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
	fmt.Println("User email:", us)

	ddb := dynamodb.NewFromConfig(cfg)
	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")

	params := &dynamodb.GetItemInput{
		TableName: aws.String(ddbTableName), 
		Key: map[string]types.AttributeValue{
            "PartitionKey": &types.AttributeValueMemberS{Value: "USER"},
			"SortKey": &types.AttributeValueMemberS{Value: us},
		},
	}

	// Save the item to DynamoDB
	resp, err := ddb.GetItem(context.TODO(), params)
	if err != nil {
		fmt.Println("Failed to retrieve user data from DynamoDB:", err)
		return nil, err
	}
	item := resp.Item
	if item == nil {
		errorStr := "User not found in DB"
		e := fmt.Errorf(errorStr)
		fmt.Println(errorStr)
		return nil, e
	}
	u, err := mapDynamoDbItemToUser(item)
	if err != nil {
		fmt.Println("Error mapping dynamoDB item to user:", err)
		return nil, err
	}
	return u, nil
}

func mapDynamoDbItemToUser(m map[string]types.AttributeValue) (u *User, err error){
	u = &User{}
	e := attributevalue.UnmarshalMap(m, u)
	if e != nil {
		return nil, e
	}
	return u, nil
}