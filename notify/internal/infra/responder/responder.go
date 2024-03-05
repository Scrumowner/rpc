package responder

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})
	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

func NewResponder(logger *zap.Logger) Responder {
	return &Respond{
		logger: logger,
		json:   jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

type Respond struct {
	logger *zap.Logger
	json   jsoniter.API
}

func (r *Respond) OutputJSON(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	switch resp.(type) {
	case string:
		_, err := fmt.Fprint(w, resp)
		if err != nil {
			r.logger.Error("Can't response json invalid format of message")
		}
	}
	err := r.json.NewEncoder(w).Encode(resp)
	if err != nil {
		r.logger.Error("Can't response json invalid format of message")
	}

}
func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.Error(w, err.Error(), http.StatusUnauthorized)

}
func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

}
func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.Error(w, err.Error(), http.StatusForbidden)
}
func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
