package util

import (
	"errors"
	"fmt"
	"strings"
)

// PrettyPrintError print error as pretty
func PrettyPrintError(err error) string {
	messages := extractMessagesFromError(err)
	return joinErrorMessages(messages)
}

func joinErrorMessages(messages []string) (message string) {
	for i := len(messages) - 1; i >= 0; i-- {
		prefix := "  "
		if i == len(messages)-1 {
			prefix = "Error: "
		}
		message = message + fmt.Sprintln(prefix+strings.TrimSuffix(messages[i], ": "))
	}
	return
}

func extractMessagesFromError(err error) (messages []string) {
	errs := unwrapErrors(err)
	beforeErrMsg := ""
	for i := len(errs) - 1; i >= 0; i-- {
		e := errs[i]
		eMsg := ""
		if beforeErrMsg == "" {
			eMsg = e.Error()
		} else {
			eMsgs := strings.Split(e.Error(), beforeErrMsg)
			eMsg = eMsgs[0]
		}
		messages = append(messages, eMsg)
		beforeErrMsg = e.Error()
	}
	return
}

func unwrapErrors(err error) (errs []error) {
	for {
		errs = append(errs, err)
		if e := errors.Unwrap(err); e == nil {
			break
		} else {
			err = e
		}
	}
	return
}
