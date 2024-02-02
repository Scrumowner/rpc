package responder

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Responder interface {
	OutputHtml(w http.ResponseWriter, responseData interface{})
	OutputJSON(w http.ResponseWriter, responseData interface{})
	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}
type Respond struct {
	log *zap.SugaredLogger
}

func NewResponder(logger *zap.SugaredLogger) Responder {
	return &Respond{log: logger}
}

func (r *Respond) OutputHtml(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, responseData)

}
func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch responseData.(type) {
	case string:
		fmt.Fprint(w, responseData)
		return
	}
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.Error(w, fmt.Sprintf("%s", err), http.StatusForbidden)
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne Unauthorized", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.Error(w, fmt.Sprintf("%s", err), http.StatusUnauthorized)
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
}
