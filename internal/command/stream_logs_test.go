package command_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	"github.com/gogo/protobuf/jsonpb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamLogs", func() {
	var (
		writer *syncedWriter
		fc     *fakeClient
		ch     chan []byte
	)

	BeforeEach(func() {
		writer = &syncedWriter{
			buf: bytes.NewBuffer([]byte{}),
		}

		ch = make(chan []byte, 1000)
		fc = &fakeClient{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}
	})

	It("connects to the specified gateway host with the correct query params", func() {
		go command.StreamLogs("https://log-stream.test-minster.cf-app.com", fc, writer)

		Eventually(fc.Host).Should(Equal("log-stream.test-minster.cf-app.com"))
		Eventually(fc.Query).Should(HaveKeyWithValue("log", []string{""}))
		Eventually(fc.Query).Should(HaveKeyWithValue("counter", []string{""}))
		Eventually(fc.Query).Should(HaveKeyWithValue("gauge", []string{""}))
		Eventually(fc.Query).Should(HaveKeyWithValue("timer", []string{""}))
		Eventually(fc.Query).Should(HaveKeyWithValue("event", []string{""}))
	})

	It("filters by source id", func() {
		go command.StreamLogs(
			"https://log-stream.test-minster.cf-app.com",
			fc,
			writer,
			command.WithSourceIDs([]string{"some-source-id", "another-source-id"}),
		)

		Eventually(fc.Query).Should(HaveKeyWithValue("source_id", []string{"some-source-id", "another-source-id"}))
	})

	It("filters by metric type", func() {
		go command.StreamLogs(
			"https://log-stream.test-minster.cf-app.com",
			fc,
			writer,
			command.WithMetricTypes([]string{"gauge", "timer"}),
		)

		Eventually(fc.Query).Should(HaveKeyWithValue("gauge", []string{""}))
		Eventually(fc.Query).Should(HaveKeyWithValue("timer", []string{""}))
		Consistently(fc.Query).ShouldNot(HaveKeyWithValue("event", []string{""}))
		Consistently(fc.Query).ShouldNot(HaveKeyWithValue("log", []string{""}))
		Consistently(fc.Query).ShouldNot(HaveKeyWithValue("counter", []string{""}))
	})

	It("writes received envelopes to the terminal", func() {
		envelopeOne := &loggregator_v2.Envelope{
			Message: &loggregator_v2.Envelope_Log{
				Log: &loggregator_v2.Log{
					Payload: []byte("hello, world"),
				},
			},
		}

		envelopeTwo := &loggregator_v2.Envelope{
			Message: &loggregator_v2.Envelope_Log{
				Log: &loggregator_v2.Log{
					Payload: []byte("goodbye, world"),
				},
			},
		}

		go command.StreamLogs("https://log-stream.test-minster.cf-app.com", fc, writer)

		go func() {
			m := jsonpb.Marshaler{}
			for i := 0; i < 1; i++ {
				s, err := m.MarshalToString(&loggregator_v2.EnvelopeBatch{
					Batch: []*loggregator_v2.Envelope{
						envelopeOne,
						envelopeTwo,
					},
				})
				if err != nil {
					panic(err)
				}
				ch <- []byte(fmt.Sprintf("data: %s\n\n", s))
			}
		}()

		Eventually(writer.String).Should(Equal(
			"{\"log\":{\"payload\":\"hello, world\"}}\n{\"log\":{\"payload\":\"goodbye, world\"}}\n"))
	})

	It("accepts a shard ID", func() {
		go command.StreamLogs(
			"https://log-stream.test-minster.cf-app.com",
			fc,
			writer,
			command.WithShardID("tralala"),
		)

		Eventually(fc.Query).Should(HaveKeyWithValue("shard_id", []string{"tralala"}))
	})

	Context("when there is an error", func() {
		It("writes the error", func() {
			fc.response.Body = ioutil.NopCloser(strings.NewReader(`{"message": "there was an error"}`))
			fc.response.StatusCode = http.StatusNotFound

			go command.StreamLogs("https://log-stream.test-minster.cf-app.com", fc, writer)

			Eventually(writer.String).Should(ContainSubstring(`{"message": "there was an error"}`))
		})
	})
})

type fakeClient struct {
	response *http.Response
	host     string
	query    url.Values
	err      error

	mu sync.Mutex
}

func (d *fakeClient) Query() url.Values {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.query
}

func (d *fakeClient) Host() string {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.host
}

func (d *fakeClient) Do(req *http.Request) (*http.Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.host = req.URL.Host
	d.query = req.URL.Query()

	if d.err != nil {
		return nil, d.err
	}

	return d.response, nil
}

type channelReader <-chan []byte

func (r channelReader) Read(buf []byte) (int, error) {
	data, ok := <-r
	if !ok {
		return 0, io.EOF
	}
	n := copy(buf, data)
	return n, nil
}

type syncedWriter struct {
	buf *bytes.Buffer

	mu sync.Mutex
}

func (w *syncedWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.Write(p)
}

func (w *syncedWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.String()
}
