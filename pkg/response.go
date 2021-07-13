package dynamoapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/josedelrio85/voalarm"
)

// Response represents the data structure needed to create a response
type Response struct {
	Code    int                    `json:"-"`
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// response sets the params to generate a JSON response
func response(w http.ResponseWriter, ra Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ra.Code)

	json.NewEncoder(w).Encode(ra)
}

// responseError generates a 500 status response
func responseError(w http.ResponseWriter, message string) {
	ra := Response{
		Code:    http.StatusInternalServerError,
		Success: false,
		Message: message,
	}
	response(w, ra)
}

// responseUnprocessable a 422 status response
func responseUnprocessable(w http.ResponseWriter, message string) {
	ra := Response{
		Code:    http.StatusUnprocessableEntity,
		Success: false,
		Message: message,
	}
	response(w, ra)
}

// responseOk calls response function with proper data to generate an OK response
func responseOk(w http.ResponseWriter, data map[string]interface{}) {
	ra := Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}
	if data != nil {
		ra.Data = data
	}
	response(w, ra)
}

// errorLogger is a struct to handle error properties
type errorLogger struct {
	msg    string
	status int
	err    error
	log    string
}

// sendAlarm to VictorOps plattform and format the error for more info
func (e *errorLogger) sendAlarm() {
	e.msg = fmt.Sprintf("Dynamodb -> %s", e.msg)
	log.Println(e.log)

	mstype := voalarm.Acknowledgement
	switch e.status {
	case http.StatusInternalServerError:
		mstype = voalarm.Warning
	case http.StatusUnprocessableEntity:
		mstype = voalarm.Info
	}

	alarm := voalarm.NewClient("")
	_, err := alarm.SendAlarm(e.msg, mstype, e.err)
	if err != nil {
		log.Fatalf(e.msg)
	}
}

// logError obtains a trace of the line and file where the error happens
func logError(err error) string {
	pc, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
}
