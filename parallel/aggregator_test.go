package parallel_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/ginkgo/parallel"
	st "github.com/onsi/ginkgo/stenographer"
	"github.com/onsi/ginkgo/types"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregator", func() {
	var (
		aggregator     *Aggregator
		reporterConfig config.DefaultReporterConfigType
		stenographer   *st.FakeStenographer

		ginkgoConfig1 config.GinkgoConfigType
		ginkgoConfig2 config.GinkgoConfigType

		suiteSummary1 *types.SuiteSummary
		suiteSummary2 *types.SuiteSummary

		exampleSummary1 *types.ExampleSummary
		exampleSummary2 *types.ExampleSummary

		suiteDescription string
	)

	BeforeEach(func() {
		config = config.DefaultReporterConfigType{
			NoColor:           false,
			SlowSpecThreshold: 0.1,
			NoisyPendings:     true,
			Succinct:          false,
			Verbose:           true,
		}
		stenographer = st.NewFakeStenographer()
		aggregator = NewAggregator(2, config, stenographer)

		//
		// now set up some fixture data
		//

		ginkgoConfig1 = config.GinkgoConfigType{
			RandomSeed:        1138,
			RandomizeAllSpecs: true,
			ParallelNode:      1,
			ParallelTotal:     2,
		}

		ginkgoConfig2 = config.GinkgoConfigType{
			RandomSeed:        1138,
			RandomizeAllSpecs: true,
			ParallelNode:      2,
			ParallelTotal:     2,
		}

		suiteDescription = "My Parallel Suite"

		suiteSummary1 = &types.SuiteSummary{
			SuiteDescription: suiteDescription,

			NumberOfExamplesBeforeParallelization: 30,
			NumberOfTotalExamples:                 17,
			NumberOfExamplesThatWillBeRun:         15,
			NumberOfPendingExamples:               1,
			NumberOfSkippedExamples:               1,
		}

		suiteSummary2 = &types.SuiteSummary{
			SuiteDescription: suiteDescription,

			NumberOfExamplesBeforeParallelization: 30,
			NumberOfTotalExamples:                 13,
			NumberOfExamplesThatWillBeRun:         8,
			NumberOfPendingExamples:               2,
			NumberOfSkippedExamples:               3,
		}

		exampleSummary1 = &types.ExampleSummary{}
		exampleSummary2 = &types.ExampleSummary{}
	})

	call := func(method string, args ...interface{}) st.FakeStenographerCall {
		return st.NewFakeStenographerCall(method, args...)
	}

	Describe("Announcing the beginning of the suite", func() {
		Context("When one of the parallel-suites starts", func() {
			BeforeEach(func() {
				aggregator.SpecSuiteWillBegin(ginkgoConfig2, suiteSummary2)
			})

			It("should be silent", func() {
				Ω(stenographer.Calls).Should(BeEmpty())
			})
		})

		Context("once all of the parallel-suites have started", func() {
			BeforeEach(func() {
				aggregator.SpecSuiteWillBegin(ginkgoConfig2, suiteSummary2)
				aggregator.SpecSuiteWillBegin(ginkgoConfig1, suiteSummary1)
			})

			It("should announce the beginning of the suite", func() {
				Ω(stenographer.Calls).Should(HaveLen(2))
				Ω(stenographer.Calls[0]).Should(Equal(call("AnnounceSuite", suiteDescription, ginkgoConfig1.RandomSeed, true)))
				//TODO: ANNOUNCE SOMETHING ABOUT PARALLELIZATION HERE
				Ω(stenographer.Calls[1]).Should(Equal(call("AnnounceNumberOfSpecs", 23, 30)))
			})
		})
	})

	Describe("Announcing examples", func() {
		Context("when the parallel-suites have not all started", func() {
			BeforeEach(func() {
				exampleSummary1.State = types.ExampleStatePassed
				aggregator.ExampleDidComplete(exampleSummary1)
			})

			It("should not announce any examples", func() {
				Ω(stenographer.Calls).Should(BeEmpty())
			})

			Context("when the parallel-suites subsequently start", func() {
				BeforeEach(func() {
					aggregator.SpecSuiteWillBegin(ginkgoConfig2, suiteSummary2)
					aggregator.SpecSuiteWillBegin(ginkgoConfig1, suiteSummary1)
				})

				It("should announce the examples", func() {
					//MAKE AN ASSERTION THAT MAKES SENSE HERE!
					Ω(stenographer.Calls).Should(Equal())
				})
			})
		})

		Context("When the parallel-suites have all started", func() {
			BeforeEach(func() {
				aggregator.SpecSuiteWillBegin(ginkgoConfig2, suiteSummary2)
				aggregator.SpecSuiteWillBegin(ginkgoConfig1, suiteSummary1)
				stenographer.Reset()
			})

			Context("When an example completes", func() {
				BeforeEach(func() {
					exampleSummary.State = types.ExampleStatePassed
					aggregator.ExampleDidComplete(exampleSummary1)
				})

				It("should announce that the example will run (when in verbose mode)", func() {
					Ω(stenographer.Calls[0]).Should(Equal(call("AnnounceExampleWillRun", exampleSummary1)))
				})

				It("should announce the captured stdout of the example", func() {
					Ω(stenographer.Calls[1]).Should(Equal(call("AnnounceCapturedOutput", exampleSummary1)))
				})

				It("should announce completion", func() {
					Ω(stenographer.Calls[2]).Should(Equal(call("AnnounceSuccesfulExample", exampleSummary1)))
				})
			})
		})
	})

	Describe("Announcing the end of the suite", func() {
		Context("When one of the parallel-suites ends", func() {
			BeforeEach(func() {
				aggregator.SpecSuiteDidEnd(ginkgoConfig2, suiteSummary2)
			})

			It("should be silent", func() {
				Ω(stenographer.Calls).Should(BeEmpty())
			})
		})

		Context("once all of the parallel-suites end", func() {
			BeforeEach(func() {
				aggregator.SpecSuiteDidEnd(ginkgoConfig2, suiteSummary2)
				aggregator.SpecSuiteDidEnd(ginkgoConfig1, suiteSummary1)
			})

			It("should announce the end of the suite", func() {
				//TODO: break this out in stenographer
			})
		})

		Context("when all the parallel-suites pass", func() {
			It("should notify the channel that it succeded", func() {

			})
		})

		Context("when one of the parallel-suites fails", func() {
			It("should notify the channel that it failed", func() {

			})
		})
	})
})
