package log_stream_plugin_test

import (
	"net/http"

	"github.com/cloudfoundry/log-stream-cli/internal/log_stream_plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DoerProvider", func() {
	It("correctly sets insecureSkipVerify", func() {
		auth := log_stream_plugin.NewDoer("foo", true).(*log_stream_plugin.AuthDoer)

		client := auth.Client.(*http.Client)
		Expect(client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify).To(BeTrue())
	})

	It("correctly sets the auth token", func() {
		auth := log_stream_plugin.NewDoer("foo", true).(*log_stream_plugin.AuthDoer)

		req := &http.Request{
			Header: make(http.Header),
		}
		auth.Do(req)

		Expect(req.Header.Get("Authorization")).To(Equal("foo"))
	})
})
