package lorg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog_Fatal_ExitWithCode1(t *testing.T) {
	testcaseEnv := os.Getenv("TESTCASE")

	switch testcaseEnv {
	case "log_fatal":
		log := NewLog()
		log.SetFormat(new(mockFormat))

		log.Fatal("TESTCASE: ", testcaseEnv)

		assert.Fail(t, "Fatal does not exit program")

		return

	case "log_fatalf":
		log := NewLog()
		log.SetFormat(new(mockFormat))

		log.Fatalf("TESTCASE: %s", testcaseEnv)

		assert.Fail(t, "Fatalf does not exit program")

		return
	}

	testcases := []string{
		"log_fatal",
		"log_fatalf",
	}

	for _, testcase := range testcases {
		output, status, err := runLogFatalTestcase(
			testcase,
		)
		if err != nil {
			assert.FailNow(t, "%s testcase failed: %s", testcase, err)
		}

		assert.Equal(t, 1, status)
		assert.Equal(t, "[testcase] TESTCASE: "+testcase+"\n", output)
	}
}

func runLogFatalTestcase(
	testcase string,
) (output string, status int, err error) {
	cmd := exec.Command(os.Args[0], "-test.run=TestLog_Fatal_ExitWithCode1")
	cmd.Env = append(os.Environ(), "TESTCASE="+testcase)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", 0, err
	}

	err = cmd.Start()
	if err != nil {
		return "", 0, err
	}

	stderrData, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", 0, fmt.Errorf("can't read stderr: %s", err)
	}

	err = cmd.Wait()
	if err == nil {
		return "", 0, fmt.Errorf("process exited with status 0")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return "", 0, fmt.Errorf("can't run process: %s", err)
	}

	return string(stderrData),
		exitErr.Sys().(syscall.WaitStatus).ExitStatus(),
		nil
}
