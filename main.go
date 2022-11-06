package main

import (
	"github.com/Devenlabo123/connectRN/api"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type handler struct {
	apiHandler     http.Handler
}


func main() {
	formatter := &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyMsg:  "message",
		},
	}
	formatter.TimestampFormat = time.RFC3339Nano
	log.SetFormatter(formatter)

	log.SetLevel(log.DebugLevel)
	log.Info("service starting")

	apiHandler, err := api.CreateHandler()
	if err != nil {
		log.WithError(err).Fatal("error encountered creating api handler")
	}

	h := handler{
		apiHandler: apiHandler,
	}
	srv := &http.Server{Addr: ":8080", Handler: h}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go startServer(&wg, srv)

	log.Info("service initialized")
	wg.Wait()

}

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.Index(request.URL.Path, "/api") == 0 {
		h.apiHandler.ServeHTTP(writer, request)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func startServer(wg *sync.WaitGroup, srv *http.Server) {
	defer wg.Done()
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.WithError(err).Fatal("error encountered in http server")
	}
}