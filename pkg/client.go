package dynamo_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Handler manages global stuff
type Handler struct {
	AppContext AppContext
}

// AppContext is a stuct to handle DynamoDB context
type AppContext struct {
	Sess *session.Session
	Db   *dynamodb.DynamoDB
}

// HandleFunction is a function used to manage all received requests.
func (h *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseOk(w)
	})
}

// HandleTables is a function used to manage all received requests.
func (h *Handler) HandleTables() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tables, err := h.AppContext.Db.ListTables(&dynamodb.ListTablesInput{})
		if err != nil {
			msg := "Error retrieving tables list"
			test(err, msg)
			responseError(w, msg, err)
			return
		}
		log.Println("Tables: ")
		for _, table := range tables.TableNames {
			log.Println(*table)
		}
		responseOk(w)
	})
}

// Input example
type Input struct {
	CustomerID string `json:"cust_id"`
	Email      string `json:"email"`
}

// Item lalalal
type Item struct {
	CustomerID             string
	LastName               string
	DateOfBirth            string
	Email                  string
	IsEligibleForPromotion bool
}

// GetItem asdfasdf
func (h *Handler) GetItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		input := Input{}
		rawdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Error parsing body to bytes"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		if err := json.Unmarshal(rawdata, &input); err != nil {
			msg := "Error unmarshaling data"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		log.Printf("input: %v", input)

		result, err := h.AppContext.Db.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String("demo-customer-info"),
			Key: map[string]*dynamodb.AttributeValue{
				"customerId": {
					S: aws.String(input.CustomerID),
				},
			},
		})
		if err != nil {
			msg := "Error retrieving item"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		item := Item{}
		if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
			msg := "Error mapping to item"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		log.Printf("item: %v", item)

		responseOk(w)
	})
}

func test(err error, msg string) {
	e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
	e.sendAlarm()
}

// GetItemV2 blablabla
func (h *Handler) GetItemV2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		input := Input{}
		rawdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Error parsing body to bytes"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		if err := json.Unmarshal(rawdata, &input); err != nil {
			msg := "Error unmarshaling data"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		log.Printf("input: %v", input)

		result, err := h.AppContext.Db.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String("demo-customer-info"),
			Key: map[string]*dynamodb.AttributeValue{
				"customerId": {
					S: aws.String(input.CustomerID),
				},
			},
		})
		if err != nil {
			msg := "Error retrieving item"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		item := Item{}
		if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
			msg := "Error mapping to item"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		log.Printf("item: %v", item)

		responseOk(w)
	})
}
