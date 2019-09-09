package presentation_test

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"encoding/json"
	"github.com/cloudfoundry/log-stream-cli/internal/presentation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Envelope Presenter", func() {
	DescribeTable("conversion to presentation", func(env *loggregator_v2.Envelope, expectedResult string) {
		convertedEnvelope, err := presentation.Envelope(env)
		Expect(err).ToNot(HaveOccurred())

		envelopeString, err := json.Marshal(convertedEnvelope)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(envelopeString)).To(Equal(expectedResult))
	},
		Entry("serializes log message payload to base64 decoded string", logEnv, `{"log":{"payload":"hello world"}}`),
		Entry("does not omit empty counter delta and total", emptyCounter, `{"counter":{"name":"counter","delta":0,"total":0}}`),
		Entry("does not omit empty gauge unit and value", emptyGauge, `{"gauge":{"metrics":{"gauge1":{"unit":"dollars","value":100},"gauge2":{"unit":"","value":0}}}}`),
		Entry("does not omit empty timer start and stop", emptyTimer, `{"timer":{"name":"timer","start":0,"stop":0}}`),
		Entry("does not omit empty event title and body", emptyEvent, `{"event":{"title":"","body":""}}`),
	)

	It("errors with unknown envelope type", func() {
		envUnknown := &loggregator_v2.Envelope{}

		_, err := presentation.Envelope(envUnknown)
		Expect(err).To(HaveOccurred())
	})
})

var logEnv = &loggregator_v2.Envelope{
	Message: &loggregator_v2.Envelope_Log{
		Log: &loggregator_v2.Log{
			Payload: []byte("hello world"),
		},
	},
}

var emptyCounter = &loggregator_v2.Envelope{
	Message: &loggregator_v2.Envelope_Counter{
		Counter: &loggregator_v2.Counter{
			Name:  "counter",
			Delta: 0,
			Total: 0,
		},
	},
}

var emptyGauge = &loggregator_v2.Envelope{
	Message: &loggregator_v2.Envelope_Gauge{
		Gauge: &loggregator_v2.Gauge{
			Metrics: map[string]*loggregator_v2.GaugeValue{
				"gauge1": {
					Unit:  "dollars",
					Value: 100,
				},
				"gauge2": {
					Value: 0,
				},
			},
		},
	},
}

var emptyTimer = &loggregator_v2.Envelope{
	Message: &loggregator_v2.Envelope_Timer{
		Timer: &loggregator_v2.Timer{
			Name:  "timer",
			Start: 0,
			Stop:  0,
		},
	},
}

var emptyEvent = &loggregator_v2.Envelope{
	Message: &loggregator_v2.Envelope_Event{
		Event: &loggregator_v2.Event{},
	},
}
