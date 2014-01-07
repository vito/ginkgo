package parallel

import (
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/stenographer"
	"github.com/onsi/ginkgo/types"
)

type Aggregator struct {
	nodeCount    int
	config       config.DefaultReporterConfigType
	stenographer stenographer.Stenographer
}

func NewAggregator(nodeCount int, config config.DefaultReporterConfigType, stenographer stenographer.Stenographer) *Aggregator {
	return &Aggregator{
		nodeCount:    int,
		config:       config,
		stenographer: stenographer,
	}
}

func (aggregator *Aggregator) SpecSuiteWillBegin(config config.GinkgoConfigType, summary *types.SuiteSummary) {
	aggregator.stenographer.AnnounceSuite(summary.SuiteDescription, config.RandomSeed, config.RandomizeAllSpecs)
	aggregator.stenographer.AnnounceNumberOfSpecs(summary.NumberOfExamplesThatWillBeRun, summary.NumberOfTotalExamples)
}

func (aggregator *Aggregator) ExampleWillRun(exampleSummary *types.ExampleSummary) {
}

func (aggregator *Aggregator) ExampleDidComplete(exampleSummary *types.ExampleSummary) {
	//put this in a channel to serialize!
	if aggregator.config.Verbose && exampleSummary.State != types.ExampleStatePending && exampleSummary.State != types.ExampleStateSkipped {
		aggregator.stenographer.AnnounceExampleWillRun(exampleSummary)
	}

	// print out captured output

	switch exampleSummary.State {
	case types.ExampleStatePassed:
		if exampleSummary.IsMeasurement {
			aggregator.stenographer.AnnounceSuccesfulMeasurement(exampleSummary, aggregator.config.Succinct)
		} else if exampleSummary.RunTime.Seconds() >= aggregator.config.SlowSpecThreshold {
			aggregator.stenographer.AnnounceSuccesfulSlowExample(exampleSummary, aggregator.config.Succinct)
		} else {
			aggregator.stenographer.AnnounceSuccesfulExample(exampleSummary)
		}
	case types.ExampleStatePending:
		aggregator.stenographer.AnnouncePendingExample(exampleSummary, aggregator.config.NoisyPendings, aggregator.config.Succinct)
	case types.ExampleStateSkipped:
		aggregator.stenographer.AnnounceSkippedExample(exampleSummary)
	case types.ExampleStateTimedOut:
		aggregator.stenographer.AnnounceExampleTimedOut(exampleSummary, aggregator.config.Succinct)
	case types.ExampleStatePanicked:
		aggregator.stenographer.AnnounceExamplePanicked(exampleSummary, aggregator.config.Succinct)
	case types.ExampleStateFailed:
		aggregator.stenographer.AnnounceExampleFailed(exampleSummary, aggregator.config.Succinct)
	}
}

func (aggregator *Aggregator) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	aggregator.stenographer.AnnounceSpecRunCompletion(summary)
}
