package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
)

type ReflectResponse struct {
	Code          int
	Message       string
	Data          map[string]interface{}
	Headers       map[string][]string
	Host          string
	RemoteHost    string
	ContentLength int64
}

var port int

func main() {
	flag.IntVar(&port, "p", 8080, "Port to run reflect server on")
	flag.Parse()

	http.HandleFunc("/reflect", func(w http.ResponseWriter, r *http.Request) {
		kvs := make(map[string]interface{})
		var vals url.Values
		switch r.Method {
		case "POST":
			vals = r.PostForm
		case "GET":
			vals = r.URL.Query()
		}
		for k, v := range vals {
			kvs[k] = v[0]
		}

		response := ReflectResponse{0, "ok", kvs, r.Header, r.Host, r.RemoteAddr, r.ContentLength}

		fmt.Printf("[%s] Response: %+v\n", r.RemoteAddr, response)

		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			errStr := fmt.Sprintf("{\"code\":1,\"message\":\"%v\"}", err)
			w.Write([]byte(errStr))
		}
	})

	http.ListenAndServe(":"+fmt.Sprintf("%d", port), nil)
}
