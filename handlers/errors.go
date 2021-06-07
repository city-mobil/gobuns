package handlers

import (
	"encoding/json"
	"net/http"
)

type RequestError struct {
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (r *RequestError) Error() string {
	return r.Message
}

func (r *RequestError) getStatus() int {
	return r.Status
}

func (r *RequestError) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RequestErrorWithBody struct {
	Status  int    `json:"-"`
	Message string `json:"-"`
	Body    interface{}
}

func (r *RequestErrorWithBody) Error() string {
	return r.Message
}

func (r *RequestErrorWithBody) marshal() ([]byte, error) {
	return json.Marshal(r.Body)
}

func (r *RequestErrorWithBody) getStatus() int {
	return r.Status
}

type requestErrorWrapper interface {
	getStatus() int
	marshal() ([]byte, error)
}

func writeRequestError(w http.ResponseWriter, werr requestErrorWrapper) {
	status := werr.getStatus()
	var dt []byte
	dt, err := werr.marshal()
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(dt)
}
