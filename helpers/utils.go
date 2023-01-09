package utils

import "log"

func HandleError(message string, err error) {
	if err != nil {
		log.Println(message, err.Error())
	}
}
