package server

import (
	"encoding/json"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
	"io/ioutil"
	"net"
	"net/http"
)

type Server struct {
	listener  net.Listener
	reporters []ginkgo.Reporter
}

func New() (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
	}, nil
}

func (server *Server) Start() {
	httpServer := &http.Server{}
	mux := http.NewServeMux()
	httpServer.Handler = mux

	mux.HandleFunc("/SpecSuiteWillBegin", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		body, _ := ioutil.ReadAll(request.Body)
		server.SpecSuiteWillBegin(body)
		writer.WriteHeader(200)
	})

	mux.HandleFunc("/ExampleWillRun", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		body, _ := ioutil.ReadAll(request.Body)
		server.ExampleWillRun(body)
		writer.WriteHeader(200)
	})

	mux.HandleFunc("/ExampleDidComplete", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		body, _ := ioutil.ReadAll(request.Body)
		server.ExampleDidComplete(body)
		writer.WriteHeader(200)
	})

	mux.HandleFunc("/SpecSuiteDidEnd", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		body, _ := ioutil.ReadAll(request.Body)
		server.SpecSuiteDidEnd(body)
		writer.WriteHeader(200)
	})

	go httpServer.Serve(server.listener)
}

func (server *Server) Stop() {
	server.listener.Close()
}

func (server *Server) Address() string {
	return server.listener.Addr().String()
}

func (server *Server) RegisterReporters(reporters ...ginkgo.Reporter) {
	server.reporters = reporters
}

func (server *Server) SpecSuiteWillBegin(body []byte) {
	var data struct {
		Config  config.GinkgoConfigType `json:"config"`
		Summary *types.SuiteSummary     `json:"suite-summary"`
	}

	json.Unmarshal(body, &data)

	for _, reporter := range server.reporters {
		reporter.SpecSuiteWillBegin(data.Config, data.Summary)
	}
}

func (server *Server) ExampleWillRun(body []byte) {
	var exampleSummary *types.ExampleSummary
	json.Unmarshal(body, &exampleSummary)

	for _, reporter := range server.reporters {
		reporter.ExampleWillRun(exampleSummary)
	}
}

func (server *Server) ExampleDidComplete(body []byte) {
	var exampleSummary *types.ExampleSummary
	json.Unmarshal(body, &exampleSummary)

	for _, reporter := range server.reporters {
		reporter.ExampleDidComplete(exampleSummary)
	}
}

func (server *Server) SpecSuiteDidEnd(body []byte) {
	var suiteSummary *types.SuiteSummary
	json.Unmarshal(body, &suiteSummary)

	for _, reporter := range server.reporters {
		reporter.SpecSuiteDidEnd(suiteSummary)
	}
}
