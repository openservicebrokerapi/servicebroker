package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Log(str string) {
	fmt.Printf("%s\n", str)
}

func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	var err error
	var data []byte

	if str, ok := object.(string); ok {
		data = []byte(str)
	} else {
		data, err = json.Marshal(object)
		if err != nil {
			code = http.StatusInternalServerError
			data = []byte(fmt.Sprintf("%s\n", err))
		}
	}

	w.WriteHeader(code)
	fmt.Fprintf(w, string(data))
}

func getCatalog(w http.ResponseWriter, r *http.Request) {
	Log("Got: " + r.URL.String())
	res := `{ "services": [] }`

	res += "\n"
	WriteResponse(w, 200, res)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", getCatalog).Methods("GET")

	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	fmt.Println("Broker started on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err.Error())
}
