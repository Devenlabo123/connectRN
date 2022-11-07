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
	userApiHandler http.Handler
	imageApiHandler http.Handler
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

	userApiHandler, err := api.CreateUserHandler()
	if err != nil {
		log.WithError(err).Fatal("error encountered creating user api handler")
	}

	imageApiHandler, err := api.CreateImageHandler()
	if err != nil {
		log.WithError(err).Fatal("error encountered creating image api handler")
	}


	h := handler{
		userApiHandler: userApiHandler,
		imageApiHandler: imageApiHandler,
	}
	srv := &http.Server{Addr: ":8080", Handler: h}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go startServer(&wg, srv)

	log.Info("service initialized")
	wg.Wait()
}

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.Index(request.URL.Path, "/api/createUser") == 0 {
		h.userApiHandler.ServeHTTP(writer, request)
	} else if strings.Index(request.URL.Path, "/api/images") == 0 {
		h.imageApiHandler.ServeHTTP(writer, request)
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