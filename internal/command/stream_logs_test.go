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
	)

	BeforeEach(func() {
		writer = &syncedWriter{
			buf: bytes.NewBuffer([]byte{}),
		}
	})

	It("connects to the specified gateway host with the correct query params", func() {
		ch := make(chan []byte, 100)
		doer := &fakeDoer{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}
		go command.StreamLogs([]string{}, "https://log-stream.test-minster.cf-app.com", doer, writer)

		Eventually(doer.Host).Should(Equal("log-stream.test-minster.cf-app.com"))
		Eventually(doer.Query).Should(HaveKeyWithValue("log", []string{""}))
		Eventually(doer.Query).Should(HaveKeyWithValue("counter", []string{""}))
		Eventually(doer.Query).Should(HaveKeyWithValue("gauge", []string{""}))
		Eventually(doer.Query).Should(HaveKeyWithValue("timer", []string{""}))
		Eventually(doer.Query).Should(HaveKeyWithValue("event", []string{""}))
	})

	It("passes source_id when present", func() {
		ch := make(chan []byte, 100)
		doer := &fakeDoer{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}
		go command.StreamLogs([]string{"some-source-id", "another-source-id"}, "https://log-stream.test-minster.cf-app.com", doer, writer)

		Eventually(doer.Query).Should(HaveKeyWithValue("source_id", []string{"some-source-id", "another-source-id"}))
	})

	It("writes received envelopes to the terminal", func() {
		ch := make(chan []byte, 1000)
		doer := &fakeDoer{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}

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

		go command.StreamLogs([]string{}, "https://log-stream.test-minster.cf-app.com", doer, writer)

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

	Context("when there is an error", func() {
		It("writes the error", func() {
			doer := &fakeDoer{
				response: &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader(`{"message": "there was an error"}`)),
					StatusCode: 404,
				},
			}

			go command.StreamLogs([]string{}, "https://log-stream.test-minster.cf-app.com", doer, writer)

			Eventually(writer.String).Should(ContainSubstring(`{"message": "there was an error"}`))
		})
	})
})

type fakeDoer struct {
	response *http.Response
	host     string
	query    url.Values
	err      error

	mu sync.Mutex
}

func (d *fakeDoer) Query() url.Values {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.query
}

func (d *fakeDoer) Host() string {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.host
}

func (d *fakeDoer) Do(req *http.Request) (*http.Response, error) {
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
