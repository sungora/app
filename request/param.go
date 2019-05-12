package request

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Param struct {
	request       *http.Request
	requestParams map[string][]string
}

func (par Param) UUID(key string) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(par.request, key))
}

func (par Param) Int(key string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(par.request, key), 10, 64)
}

func (par Param) Uint(key string) (uint64, error) {
	return strconv.ParseUint(chi.URLParam(par.request, key), 10, 64)
}

func (par Param) Float(key string) (float64, error) {
	return strconv.ParseFloat(chi.URLParam(par.request, key), 64)
}

// GetParams получение пармаметра запроса ("application/x-www-form-urlencoded").
func (par Param) Key(key string) ([]string, error) {
	if par.requestParams == nil {
		_, _ = par.All()
	}
	if _, ok := par.requestParams[key]; ok == false {
		return nil, errors.New("param not found")
	}
	return par.requestParams[key], nil
}

// GetParams получение всех данных запроса ("application/x-www-form-urlencoded").
func (par Param) All() (map[string][]string, error) {
	if par.requestParams == nil {
		par.requestParams, _ = url.ParseQuery(par.request.URL.Query().Encode())
		if err := par.request.ParseForm(); err != nil {
			return nil, err
		}
		for i, v := range par.request.Form {
			par.requestParams[i] = v
		}
	}
	return par.requestParams, nil
}
