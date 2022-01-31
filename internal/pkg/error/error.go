package errno

import (
	"encoding/json"
	"net/http"
)

const (
	InternalError = 10002
)

type Error struct {
	Code    int64
	Message string
}

func GenErrorResp(w http.ResponseWriter, code int64, message string) {
	json.NewEncoder(w).Encode(Error{
		Code:    code,
		Message: message,
	})
}
