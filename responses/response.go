package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"net/http"

)

type ErrorResponse struct {
	Status	int			`json:"status"`
	Msg 	string		`json:"error"`
}

type SuccessResponse struct {
	Status	int			`json:"status"`
	Data   string		`json:"data"`
}

func SendDataResponse(w http.ResponseWriter, statusCode int, dataResult interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if dataResult != nil {
		dr := SuccessResponse{Status: statusCode, Data: dataResult.(string)}
		responseData, err := json.Marshal(dr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(responseData)
		return
	}
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	es := ErrorResponse{Status: statusCode, Msg: msg}
	responseData, err := json.Marshal(es)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseData)
}

func JsonErrorResponse(w http.ResponseWriter, err error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var msg string
	var status int

	switch {
	case errors.As(err, &syntaxError):
		msg = fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		status = http.StatusBadRequest
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg = fmt.Sprintf("request body contains badly-formed JSON")
		status = http.StatusBadRequest
	case errors.As(err, &unmarshalTypeError):
		msg = fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		status = http.StatusBadRequest
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg = fmt.Sprintf("request body contains unknown field %s", fieldName)
		status = http.StatusBadRequest
	case errors.Is(err, io.EOF):
		msg = "request body must not be empty"
		status = http.StatusBadRequest
	default:
		log.Println(err.Error())
		msg = http.StatusText(http.StatusInternalServerError)
		status = http.StatusInternalServerError
	}

	SendDataResponse(w, status, msg)
}

