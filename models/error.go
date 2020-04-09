package models

import "fmt"

type Error struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("error_code: %s, error_message: %s", e.Code, e.Message)
}
