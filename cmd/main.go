package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	client "github.com/bysidecar/dynamo_test/pkg"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	log.Println("DynamoDB with Golang test started!")

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
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost)
	subrouter := router.PathPrefix("/dynamo").Subrouter()
	subrouter.Handle("/test", handler.HandleFunction())
	subrouter.Handle("/tables", handler.HandleTables())
	subrouter.Handle("/item", handler.GetItem())

	log.Println("starting web server...")
	log.Fatal(http.ListenAndServe(":9001", cors.Default().Handler(subrouter)))
}
