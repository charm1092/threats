package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/info").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateThreat)
	router.Path("/info/{md5}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetByMD5)
	router.Path("/info/{md5}").Methods("PUT").HandlerFunc(s.httpHandlers.HandlerUpdateThreat)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}