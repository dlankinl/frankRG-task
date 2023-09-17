package util

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Error bool        `json:"error"`
	Msg   string      `json:"message"`
	Data  interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JSONResponse
	payload.Error = true
	payload.Msg = err.Error()
	WriteJSON(w, statusCode, payload)
}
