package command_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace/tracefakes"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	"github.com/cloudfoundry/log-stream-cli/internal/testhelpers"
	"github.com/gogo/protobuf/jsonpb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamLogs", func() {
	var (
		ui terminal.UI

		printer      *testhelpers.FakePrinter
		tracePrinter *tracefakes.FakePrinter

		stdin  *testhelpers.SyncedBuffer
		stdout *testhelpers.SyncedBuffer
	)

	BeforeEach(func() {
		stdin = &testhelpers.SyncedBuffer{}
		stdout = &testhelpers.SyncedBuffer{}

		printer = new(testhelpers.FakePrinter)
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			return fmt.Fprintf(stdout, format, a...)
		}
		tracePrinter = new(tracefakes.FakePrinter)

		ui = terminal.NewUI(stdin, stdout, printer, tracePrinter)
	})

	It("connects to the specified gateway host with the correct query params", func() {
		ch := make(chan []byte, 100)
		doer := &fakeDoer{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}
		go command.StreamLogs("https://log-stream.test-minster.cf-app.com", doer, ui)

		Eventually(func() string { return doer.ReqHost }).Should(Equal("log-stream.test-minster.cf-app.com"))
		Eventually(func() url.Values { return doer.ReqQuery }).Should(HaveKeyWithValue("log", []string{""}))
		Eventually(func() url.Values { return doer.ReqQuery }).Should(HaveKeyWithValue("counter", []string{""}))
		Eventually(func() url.Values { return doer.ReqQuery }).Should(HaveKeyWithValue("gauge", []string{""}))
		Eventually(func() url.Values { return doer.ReqQuery }).Should(HaveKeyWithValue("timer", []string{""}))
		Eventually(func() url.Values { return doer.ReqQuery }).Should(HaveKeyWithValue("event", []string{""}))
	})

	It("writes received envelopes to the terminal", func() {
		ch := make(chan []byte, 1000)
		doer := &fakeDoer{
			response: &http.Response{
				Body:       ioutil.NopCloser(channelReader(ch)),
				StatusCode: 200,
			},
		}
		go command.StreamLogs("https://log-stream.test-minster.cf-app.com", doer, ui)

		go func() {
			m := jsonpb.Marshaler{}
			for i := 0; i < 10; i++ {
				s, err := m.MarshalToString(&loggregator_v2.EnvelopeBatch{
					Batch: []*loggregator_v2.Envelope{
						{Timestamp: int64(i)},
					},
				})
				if err != nil {
					panic(err)
				}
				ch <- []byte(fmt.Sprintf("data: %s\n\n", s))
			}
		}()

		Eventually(stdout).Should(ContainSubstring("timestamp"))
	})
})

type fakeDoer struct {
	response *http.Response
	ReqHost  string
	ReqQuery url.Values
	err      error
}

func (d *fakeDoer) Do(req *http.Request) (*http.Response, error) {
	d.ReqHost = req.URL.Host
	d.ReqQuery = req.URL.Query()

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
