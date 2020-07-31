package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	client "github.com/bysidecar/dynamoapi/pkg"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	log.Println("DynamoDB with Golang test started!")

	dev := false
	devstring := getSetting("DEV")
	if devstring == "true" {
		dev = true
	}
	log.Printf("Are we working on dev? %t", dev)

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:8042"),
		// Credentials: credentials.NewSharedCredentials(""),
	})
	if err != nil {
		log.Fatalf("Error connecting database. Err: %v", err)
		return
	}
	db := dynamodb.New(sess)
	handler := client.Handler{
		AppContext: client.AppContext{
			Sess: sess,
			Db:   db,
		},
		Dev: dev,
	}

	// TODO v2 => hookify this structure (?)
	tablename := "leads"
	req := &dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	}
	_, err = handler.AppContext.Db.DescribeTable(req)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			t := dynamodb.ResourceNotFoundException{}
			switch aerr.Code() {
			case t.Code():
				input := &dynamodb.CreateTableInput{
					TableName: aws.String(tablename),
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String("passport_id"),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("passport_id"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(10),
						WriteCapacityUnits: aws.Int64(10),
					},
				}
				_, err = handler.AppContext.Db.CreateTable(input)
				if err != nil {
					log.Fatalf("Error creating table %s. Err: %v", tablename, err)
					return
				}
			default:
				log.Fatalf("Error describing table %s. Err: %v", tablename, err)
				return
			}
		}
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost)
	subrouter := router.PathPrefix("/dynamo").Subrouter()
	subrouter.Handle("/test", handler.HandleFunction())
	subrouter.Handle("/tables", handler.PrintTables())
	subrouter.Handle("/describe", handler.DescribeTable())
	subrouter.Handle("/item", handler.GetItem())
	subrouter.Handle("/put", handler.PutItem())

	log.Println("starting web server...")
	log.Fatal(http.ListenAndServe(":9001", cors.Default().Handler(subrouter)))
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}
