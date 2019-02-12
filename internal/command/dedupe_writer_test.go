package command_test

import (
	"github.com/cloudfoundry/log-stream-cli/internal/command"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DedupeWriter", func() {
	var (
		fw     *fakeWriter
		writer *command.DedupeWriter
	)

	BeforeEach(func() {
		fw = &fakeWriter{}
		writer = command.NewDedupeWriter(fw)
	})

	It("writes the message", func() {
		writer.Write([]byte("this is the message"))

		Eventually(func() string { return fw.calls[0] }).Should(Equal("this is the message"))
	})

	It("only writes a message if it's different than the previous", func() {
		writer.Write([]byte("this is the message"))
		writer.Write([]byte("this is the message"))

		Expect(len(fw.calls)).To(Equal(1))

		writer.Write([]byte("this is another message"))
		Expect(len(fw.calls)).To(Equal(2))
	})
})

type fakeWriter struct {
	calls []string
}

func (fw *fakeWriter) Write(p []byte) (n int, err error) {
	fw.calls = append(fw.calls, string(p))
	return 0, nil
}
