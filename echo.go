package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/go-yaml/yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	edn "gopkg.in/edn.v1"
)

var FORMAT_PARAMETER = "__fmt"

type ReflectResponse struct {
	Code    int
	Message string
	Request ReflectData
}

type ReflectData struct {
	Parameters    map[string]interface{}
	Body          []byte
	Headers       map[string][]string
	Host          string
	RemoteHost    string
	ContentLength int64
	URI           string
	Method        string
	Protocol      string
}

var port int

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("HTTP Echo server - v1.0")

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	go func() {
		for sig := range s {
			log.Error().Msgf("%+v caught, exiting...", sig)
			os.Exit(0)
		}
	}()

	flag.IntVar(&port, "p", 8080, "Port to run reflect server on")
	flag.Parse()

	log.Info().Msg("Building HTTP handler")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		bytes := make([]byte, r.ContentLength)

		if r.Method == "POST" {
			_, err := r.Body.Read(bytes)

			if err != nil {
				writeErr(err, w)
				return
			}
		}

		data := ReflectData{
			Parameters:    kvs,
			Headers:       r.Header,
			Host:          r.Host,
			RemoteHost:    r.RemoteAddr,
			ContentLength: r.ContentLength,
			URI:           r.RequestURI,
			Method:        r.Method,
			Body:          bytes,
			Protocol:      r.Proto,
		}

		response := ReflectResponse{
			Code:    0,
			Message: "ok",
			Request: data,
		}

		log.Info().Msgf("[%s] Response: %+v\n", r.RemoteAddr, response)
		var err error
		if format, ok := kvs[FORMAT_PARAMETER]; ok {

			switch format {
			case "edn":
				var data []byte
				data, err = edn.Marshal(&response)
				w.Write(data)
			case "yaml":
				var data []byte
				data, err = yaml.Marshal(&response)
				w.Write(data)
			case "json":
				fallthrough
			default:
				err = json.NewEncoder(w).Encode(&response)
			}
		} else {
			err = json.NewEncoder(w).Encode(&response)
		}

		if err != nil {
			writeErr(err, w)
		}
	})

	log.Info().Msgf("Listening on port %d", port)
	log.Error().Msgf("Exit: %v", http.ListenAndServe(":"+fmt.Sprintf("%d", port), nil))
}

func writeErr(err error, writer http.ResponseWriter) {
	errStr := fmt.Sprintf("{\"code\":1,\"message\":\"%v\"}", err)
	writer.Write([]byte(errStr))
}
