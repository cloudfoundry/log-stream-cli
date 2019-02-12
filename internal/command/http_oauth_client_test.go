package command_test

import (
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DoerProvider", func() {
	It("correctly sets insecureSkipVerify", func() {
		auth := command.NewAuthClient("foo", true)

		client := auth.Client
		Expect(client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify).To(BeTrue())
	})

	It("correctly sets the auth token", func() {
		auth := command.NewAuthClient("foo", true)

		req := &http.Request{
			Header: make(http.Header),
		}
		auth.Do(req)

		Expect(req.Header.Get("Authorization")).To(Equal("foo"))
	})
})
