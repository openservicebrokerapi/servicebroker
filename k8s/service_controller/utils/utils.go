package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	var data []byte
	var err error

	if str, ok := object.(string); ok {
		data = []byte(str)
	} else if err, ok = object.(error); ok {
		if jerr, ok := err.(*json.SyntaxError); ok {
			data = []byte(fmt.Sprintf("%s - offset: %d", err, jerr.Offset))
		} else {
			data = []byte(err.Error())
		}
	} else {
		data, err = json.Marshal(object)
		if err != nil {
			code = http.StatusInternalServerError
			data = []byte(fmt.Sprintf("%s", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, string(data)+"\n")
}

func BodyToObject(r *http.Request, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

func ResponseBodyToObject(r *http.Response, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

func ExtractVarFromRequest(r *http.Request, varName string) string {
	return mux.Vars(r)[varName]
}
