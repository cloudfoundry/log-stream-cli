package command_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	"github.com/gogo/protobuf/jsonpb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamLogs", func() {
	var (
		buf *bytes.Buffer
	)

	BeforeEach(func() {
		buf = bytes.NewBuffer([]byte{})
	})

	It("connects to the specified gateway host with the correct query params", func() {
		ch := make(chan []byte, 100)
		doer := &fakeDoer{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}
		go command.StreamLogs("https://log-stream.test-minster.cf-app.com", doer, buf)

		Eventually(func() string { return doer.host }).Should(Equal("log-stream.test-minster.cf-app.com"))
		Eventually(func() url.Values { return doer.query }).Should(HaveKeyWithValue("log", []string{""}))
		Eventually(func() url.Values { return doer.query }).Should(HaveKeyWithValue("counter", []string{""}))
		Eventually(func() url.Values { return doer.query }).Should(HaveKeyWithValue("gauge", []string{""}))
		Eventually(func() url.Values { return doer.query }).Should(HaveKeyWithValue("timer", []string{""}))
		Eventually(func() url.Values { return doer.query }).Should(HaveKeyWithValue("event", []string{""}))
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

		go command.StreamLogs("https://log-stream.test-minster.cf-app.com", doer, buf)

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

		Eventually(func() string { return buf.String() }).Should(Equal(
			"{\"log\":{\"payload\":\"aGVsbG8sIHdvcmxk\"}}\n{\"log\":{\"payload\":\"Z29vZGJ5ZSwgd29ybGQ=\"}}\n"))
	})

	Context("when there is an error", func() {
		It("writes the error", func() {
			doer := &fakeDoer{
				response: &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader(`{"message": "there was an error"}`)),
					StatusCode: 404,
				},
			}

			go command.StreamLogs("https://log-stream.test-minster.cf-app.com", doer, buf)

			Eventually(func() string { return buf.String() }).Should(ContainSubstring(`{"message": "there was an error"}`))
		})
	})
})

type fakeDoer struct {
	response *http.Response
	host     string
	query    url.Values
	err      error
}

func (d *fakeDoer) Do(req *http.Request) (*http.Response, error) {
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
