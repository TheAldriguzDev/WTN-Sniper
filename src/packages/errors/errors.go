package error_manager

import (
	console_logger "WTN-Sniper/src/packages/logger"
	"fmt"
	"time"
)

type ErrorInformation struct {
	Error     error
	TimeStamp string
}

type Errors interface {
	BadStatus(status int) ErrorInformation
	Marshal(err error) ErrorInformation
	UnMarshal(err error) ErrorInformation
	Session(err error) ErrorInformation
	NewRequest(err error) ErrorInformation
	DoRequest(err error) ErrorInformation
	RotateProxy(err error) ErrorInformation
	ReadBody(err error) ErrorInformation
}

var l = console_logger.New("Error", true)

func BadStatus(status int, message string) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("Bad status code: %v. [%s]", status, message),
		TimeStamp: getTime(),
	}

	return tempErr
}

func Marshal(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while marhsalling the json: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func UnMarshal(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while unmarhsalling the json: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func Session(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while creating the wtn session: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func NewRequest(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while creating a new request: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func DoRequest(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while performing a request: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func RotateProxy(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while rotating the proxy: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func ReadBody(err error) ErrorInformation {
	tempErr := ErrorInformation{
		Error:     fmt.Errorf("An error occured while reading the response body: %v", err),
		TimeStamp: getTime(),
	}

	return tempErr
}

func getTime() string {
	now := time.Now()
	return now.Format("15:04:05")
}
