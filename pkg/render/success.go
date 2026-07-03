package render

import "net/http"

type Scs struct {
	Success string `json:"success"`
}

func Success(w http.ResponseWriter, status int, msg string) {
	JSON(w, Scs{Success: msg}, status)
}
