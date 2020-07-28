package dynamo_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

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

// HandleTables is a function used to retrive the list of tables
func (h *Handler) HandleTables() ([]string, error) {
	tbllist, err := h.AppContext.Db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		msg := "Error retrieving tables list"
		test(err, msg)
		return nil, err
	}
	tables := []string{}

	for _, tbl := range tbllist.TableNames {
		tables = append(tables, *tbl)
	}
	return tables, err
}

// PrintTables is a function used to print list tables
func (h *Handler) PrintTables() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tables, err := h.HandleTables()
		if err != nil {
			msg := "Error retrieving tables list"
			test(err, msg)
			responseError(w, msg, err)
			return
		}
		log.Println("Tables: ")
		for _, table := range tables {
			log.Println(table)
		}
		responseOk(w)
	})
}

// DescribeTable blabla
func (h *Handler) DescribeTable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tables, err := h.HandleTables()
		if err != nil {
			msg := "Error retrieving tables list"
			test(err, msg)
			responseError(w, msg, err)
			return
		}

		for _, tbl := range tables {
			req := &dynamodb.DescribeTableInput{
				TableName: aws.String(tbl),
			}
			result, err := h.AppContext.Db.DescribeTable(req)
			if err != nil {
				msg := fmt.Sprintf("Error describe table %s", tbl)
				test(err, msg)
				responseError(w, msg, err)
				return
			}
			table := result.Table
			fmt.Printf("done %s", table)
		}

	})
}

// Input example
type Input struct {
	CustomerID string `json:"cust_id"`
	Email      string `json:"email"`
}

// Item lalalal
type Item struct {
	CustomerID             string `json:"customerId"`
	LastName               string `json:"lastname"`
	DateOfBirth            string `json:"date_of_birth"`
	Email                  string `json:"email"`
	IsEligibleForPromotion bool   `json:"is_eligible"`
	Test                   string `json:"test"`
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
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
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
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg, err)
			return
		}

		// item := Item{}
		item := make(map[string]interface{})

		if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
			msg := "Error mapping to item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg, err)
			return
		}

		log.Printf("retrieved item: %v", item)

		responseOk(w)
	})
}

// PutItem blabla
func (h *Handler) PutItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// item := Item{}
		item := make(map[string]interface{})

		rawdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Error parsing body to bytes"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg, err)
			return
		}

		if err := json.Unmarshal(rawdata, &item); err != nil {
			msg := "Error unmarshaling data"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg, err)
			return
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			msg := "Error marshalling item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg, err)
			return
		}

		// Create item in table Movies
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("demo-customer-info"),
		}

		_, err = h.AppContext.Db.PutItem(input)

		if err != nil {
			msg := "Error putting item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg, err)
			return
		}

		log.Printf("putted item: %v", item)

		responseOk(w)
	})
}

//HelperRandstring lalalal
func HelperRandstring(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
