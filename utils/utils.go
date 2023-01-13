package utils

import (
	"fmt"
)

func HandleError(message string, err error) {
	if err != nil {
		fmt.Println(message, err.Error())
	}
}

func AssertString(val interface{}) string {
	v := val.(string)
	return v
}
