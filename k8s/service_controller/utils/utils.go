package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	fmt.Fprintf(w, string(data))
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

// KubeCreateResource takes input of resource definitions in the form
// of a reader. The intermingled output of stdout and stderr is
// returned as a string. It exists until we vendor a k8s client or
// figure out how to authenticate directly to apiserver.
func KubeCreateResource(r io.Reader) (string, error) {
	c := exec.Command("kubectl", "create", "-oname", "-f", "-")
	c.Stdin = r
	b, e := c.CombinedOutput()
	s := string(b)
	return s, e
}
