package dynamo_test

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/bysidecar/voalarm"
)

// Response represents the data structure needed to create a response
type Response struct {
	Code    int
	Message string `json:"message"`
}

// response sets the params to generate a JSON response
func response(w http.ResponseWriter, ra Response) {
	w.WriteHeader(ra.Code)
}

// responseError generates log, alarm and response when an error occurs
func responseError(w http.ResponseWriter, message string, err error) {
	// e := &errorLogger{message, http.StatusInternalServerError, err, logError(err)}
	// e.sendAlarm()

	ra := Response{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
	response(w, ra)
}

// responseUnprocessable calls response function to inform user of something does not work 100% OK
func responseUnprocessable(w http.ResponseWriter, message string, err error) {
	ra := Response{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
	response(w, ra)
}

// responseOk calls response function with proper data to generate an OK response
func responseOk(w http.ResponseWriter) {
	ra := Response{
		Code: http.StatusOK,
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
