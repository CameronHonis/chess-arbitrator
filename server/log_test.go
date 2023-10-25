package server_test

import (
	"bytes"
	"github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

func ReadStdout(stdoutWriter *os.File, stdoutReader *os.File) string {
	var stdoutBuffer bytes.Buffer
	_ = stdoutWriter.Close()
	_, _ = io.Copy(&stdoutBuffer, stdoutReader)
	return stdoutBuffer.String()
}

var _ = Describe("Log", func() {
	var oldStdout *os.File
	var stdoutWriter *os.File
	var stdoutReader *os.File
	BeforeEach(func() {
		oldStdout = os.Stdout
		stdoutReader, stdoutWriter, _ = os.Pipe()
		os.Stdout = stdoutWriter
	})
	Describe("::Log", func() {
		It("logs the message", func() {
			server.GetLogManager().Log("TEST", "test message")
			stdout := ReadStdout(stdoutWriter, stdoutReader)
			Expect(stdout).To(ContainSubstring("[TEST] test message\n"))
		})
		When("multiple strings are passed in", func() {
			It("logs the message", func() {
				server.GetLogManager().Log("TEST", "test message", " other test message", " 123")
				stdout := ReadStdout(stdoutWriter, stdoutReader)
				Expect(stdout).To(ContainSubstring("[TEST] test message other test message 123\n"))
			})
		})
	})
	Describe("::LogRed", func() {
		It("logs the message in color", func() {
			server.GetLogManager().LogRed("TEST", "test message")
			stdout := ReadStdout(stdoutWriter, stdoutReader)
			Expect(stdout).To(ContainSubstring("\x1b[31m[TEST] test message\x1b[0m\n"))
		})
	})
	Describe("::LogGreen", func() {
		It("logs the message in color", func() {
			server.GetLogManager().LogGreen("TEST", "test message")
			stdout := ReadStdout(stdoutWriter, stdoutReader)
			Expect(stdout).To(ContainSubstring("\x1b[32m[TEST] test message\x1b[0m\n"))
		})
	})
	AfterEach(func() {
		_ = stdoutWriter.Close()
		_ = stdoutReader.Close()
		os.Stdout = oldStdout
		server.GetLogManager().LogGreen("TEST", "resetting stdout")
	})
})
