package render

import (
	"fmt"
	"net/http"
)

type Err struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, err error, status int, message string) {

	//err = unpack(err)
	err = fmt.Errorf("%s: %w", message, err)

	JSON(w, Err{Error: err.Error()}, status)
}

//func unpack(err error) error {
//	for {
//		e := errors.Unwrap(err)
//		if e == nil {
//			break
//		}
//
//		err = e
//	}
//
//	return err
//}
