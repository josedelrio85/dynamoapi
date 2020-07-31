package dynamo_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Dev        bool
}

// AppContext is a stuct to handle DynamoDB context
type AppContext struct {
	Sess *session.Session
	Db   *dynamodb.DynamoDB
}

// HandleFunction is a function used to manage all received requests.
func (h *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseOk(w, nil)
	})
}

// HandleTables is a function used to retrive the list of tables
func (h *Handler) HandleTables() ([]string, error) {
	tbllist, err := h.AppContext.Db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		msg := "Error retrieving tables list"
		e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
		e.sendAlarm()
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
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		logtxt := "Tables \n"
		tablelist := make(map[string]interface{})
		for idx, table := range tables {
			i := fmt.Sprintf("%d", idx)
			tablelist[i] = table
			logtxt += fmt.Sprintf("%s \n", table)
		}
		h.PrintLog(logtxt)
		responseOk(w, tablelist)
	})
}

// DescribeTable blabla
func (h *Handler) DescribeTable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tables, err := h.HandleTables()
		if err != nil {
			msg := "Error retrieving tables list"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		for _, tbl := range tables {
			req := &dynamodb.DescribeTableInput{
				TableName: aws.String(tbl),
			}
			result, err := h.AppContext.Db.DescribeTable(req)
			if err != nil {
				msg := fmt.Sprintf("Error describe table %s", tbl)
				e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
				e.sendAlarm()
				responseError(w, msg)
				return
			}
			table := result.Table
			h.PrintLog(fmt.Sprintf("done %s", table))
		}

	})
}

// Input represents structure to retrieve an element
type Input struct {
	Table string `json:"table"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Item represents structure to put an item.
// You must set table name and its data
type Item struct {
	Table string                 `json:"table"`
	Data  map[string]interface{} `json:"data"`
}

// GetItem retrives an item from the tablename using key value as index
func (h *Handler) GetItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		input := Input{}
		rawdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Error parsing body to bytes"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		if err := json.Unmarshal(rawdata, &input); err != nil {
			msg := "Error unmarshaling data"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		h.PrintLog(fmt.Sprintf("input: %v", input))

		result, err := h.AppContext.Db.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(input.Table), // "demo-customer-info"
			Key: map[string]*dynamodb.AttributeValue{
				input.Key: { // "customerId"
					S: aws.String(input.Value), // input.CustomerID
				},
			},
		})
		if err != nil {
			msg := "Error retrieving item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		item := make(map[string]interface{})
		if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
			msg := "Error mapping to item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}
		h.PrintLog(fmt.Sprintf("retrieved item: %v", item))

		responseOk(w, item)
	})
}

// PutItem saves posted data into param table
func (h *Handler) PutItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item := Item{}
		rawdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Error parsing body to bytes"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		if err := json.Unmarshal(rawdata, &item); err != nil {
			msg := "Error unmarshaling data"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		av, err := dynamodbattribute.MarshalMap(item.Data)
		if err != nil {
			msg := "Error marshalling item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}

		// Create item in table
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(item.Table), // "demo-customer-info"
		}

		_, err = h.AppContext.Db.PutItem(input)

		if err != nil {
			msg := "Error putting item"
			e := &errorLogger{msg, http.StatusInternalServerError, err, logError(err)}
			e.sendAlarm()
			responseError(w, msg)
			return
		}
		h.PrintLog(fmt.Sprintf("putted item: %v", item))

		responseOk(w, nil)
	})
}

// PrintLog prints log info when we are in dev mode
func (h *Handler) PrintLog(txt string) {
	if h.Dev {
		now := time.Now().Format("2006-01-02 15-04-05")
		log.Printf("%s - %s", now, txt)
	}
}
