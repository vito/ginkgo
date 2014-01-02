package parallel

import (
	"bytes"
	"encoding/json"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
	"io"
	"net/http"
)

type Poster interface {
	Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
}

//The ForwardingReporter is automatically used when running parallel tests.
//You shouldn' need to use this reporter in your own code.
type ForwardingReporter struct {
	serverHost        string
	poster            Poster
	outputInterceptor OutputInterceptor
}

func NewForwardingReporter(serverHost string, poster Poster, outputInterceptor OutputInterceptor) *ForwardingReporter {
	return &ForwardingReporter{
		serverHost:        serverHost,
		poster:            poster,
		outputInterceptor: outputInterceptor,
	}
}

func (reporter *ForwardingReporter) post(path string, data interface{}) {
	encoded, _ := json.Marshal(data)
	buffer := bytes.NewBuffer(encoded)
	reporter.poster.Post("http://"+reporter.serverHost+path, "application/json", buffer)
}

func (reporter *ForwardingReporter) SpecSuiteWillBegin(conf config.GinkgoConfigType, summary *types.SuiteSummary) {
	data := struct {
		Config  config.GinkgoConfigType `json:"config"`
		Summary *types.SuiteSummary     `json:"suite-summary"`
	}{
		conf,
		summary,
	}

	reporter.post("/SpecSuiteWillBegin", data)
}

func (reporter *ForwardingReporter) ExampleWillRun(exampleSummary *types.ExampleSummary) {
	reporter.outputInterceptor.StartInterceptingOutput()
	reporter.post("/ExampleWillRun", exampleSummary)
}

func (reporter *ForwardingReporter) ExampleDidComplete(exampleSummary *types.ExampleSummary) {
	output, _ := reporter.outputInterceptor.StopInterceptingAndReturnOutput()
	exampleSummary.CapturedOutput = output
	reporter.post("/ExampleDidComplete", exampleSummary)
}

func (reporter *ForwardingReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	reporter.post("/SpecSuiteDidEnd", summary)
}
