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

func TestLogFatal(t *testing.T) {
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
		err := assertLogFatalExit(
			t,
			"TESTCASE="+testcase,
			"[testcase] TESTCASE: "+testcase+"\n",
		)
		if err != nil {
			t.Fatalf("%s testcase failed: %s", testcase, err)
		}
	}
}

func assertLogFatalExit(
	t *testing.T, env string, expectedOutput string,
) error {
	cmd := exec.Command(os.Args[0], "-test.run=TestLogFatal")
	cmd.Env = append(os.Environ(), env)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	stderrData, err := ioutil.ReadAll(stderr)
	if err != nil {
		return fmt.Errorf("can't read stderr: %s", err)
	}

	err = cmd.Wait()
	if err == nil {
		return fmt.Errorf("process exited with status 0")
	}

	exiterr, ok := err.(*exec.ExitError)
	if !ok {
		return fmt.Errorf("can't run process: %s", err)
	}

	status := exiterr.Sys().(syscall.WaitStatus).ExitStatus()
	assert.Equal(t, 1, status)

	assert.Equal(t, expectedOutput, string(stderrData))

	return nil
}
