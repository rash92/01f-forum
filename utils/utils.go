package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func PrintErrOnCommandLine(err error) {
	if err != nil {
		fmt.Println("unable to open logfile", err)
	}
}

func WriteMessageToLogFile(message interface{}) {
	now := time.Now()
	formatTime := now.Format(time.UnixDate)
	_, _, functionName := trace()
	stringMessage := string(AssertString(message)) + " in " + functionName
	MessageWithFormatTime := formatTime + ": " + stringMessage + "\n"
	WriteToLogFile(MessageWithFormatTime)
}

func HandleError(message string, err error) {
	if err != nil {
		now := time.Now()
		formatTime := now.Format(time.UnixDate)
		_, _, functionName := trace()
		errorMessage := fmt.Sprintf("%s: %v in %s", message, err, functionName)
		errorMessageWithFormatTime := formatTime + ": " + errorMessage + "\n"
		WriteToLogFile(errorMessageWithFormatTime)
	}
}

func WriteToLogFile(message string) {
	file, err := os.OpenFile("./logfile.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	PrintErrOnCommandLine(err)

	messageWithNewline := message

	n, err := file.Write([]byte(messageWithNewline))
	PrintErrOnCommandLine(err)
	if n != len(messageWithNewline) {
		fmt.Println("message length not the same")
	}
}

func trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "?"
	}

	return file, line, fn.Name()
}

func AssertString(val interface{}) string {
	v := val.(string)
	return v
}
