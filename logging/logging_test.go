package logging

import (
	"bytes"
	"log"
	"testing"

	"github.com/makkes/assert"
)

func TestLevelFromStringShouldBeCaseInsensitive(t *testing.T) {
	a := assert.NewAssert(t)
	l, _ := LevelFromString("debug")
	a.Equal(*l, DEBUG, "Wrong level returned")

	l, _ = LevelFromString("DEBUG")
	a.Equal(*l, DEBUG, "Wrong level returned")
}

func TestLevelFromStringShouldReturnNilWithUnknownString(t *testing.T) {
	a := assert.NewAssert(t)
	l, _ := LevelFromString("oink")
	a.Nil(l, "Level is not nil")
}

func TestLevelFromStringShouldReturnNilWithEmptyString(t *testing.T) {
	a := assert.NewAssert(t)
	l, _ := LevelFromString("")
	a.Nil(l, "Level is not nil")
}

func TestDefaultDebugShouldLogWithDebugLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	Debug("did that")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[DEBUG\\] \\[logging_test.go:\\d*\\] did that\n$", buf.String(), "Unexpected log output")
}

func TestLoggerDebugShouldLogWithDebugLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	logger := NewLogger(DEBUG, "")
	log.SetOutput(buf)

	logger.Debug("did that")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[DEBUG\\] \\[logging_test.go:\\d*\\] did that\n$", buf.String(), "Unexpected log output")
}

func TestDefaultInfoShouldLogWithInfoLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	Info("did it")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[INFO\\] \\[logging_test.go:\\d*\\] did it\n$", buf.String(), "Unexpected log output")
}

func TestLoggerWarnShouldLogWithWarnLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	logger := NewLogger(WARN, "")
	log.SetOutput(buf)

	logger.Warn("warning")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[WARN\\] \\[logging_test.go:\\d*\\] warning\n$", buf.String(), "Unexpected log output")
}

func TestDefaultWarnShouldLogWithWarnLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	Warn("warning")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[WARN\\] \\[logging_test.go:\\d*\\] warning\n$", buf.String(), "Unexpected log output")
}

func TestLoggerErrorShouldLogWithErrorLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	logger := NewLogger(ERROR, "")
	log.SetOutput(buf)

	logger.Error("erroaarrr!!")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[ERROR\\] \\[logging_test\\.go:\\d*\\] erroaarrr!!\n$", buf.String(), "Unexpected log output")
}

func TestDefaultErrorShouldLogWithErrorLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	Error("erroaarrr!!")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[ERROR\\] \\[logging_test\\.go:\\d*\\] erroaarrr!!\n$", buf.String(), "Unexpected log output")
}

func TestDefaultSetLevelShouldSetLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	SetLevel(ERROR)
	Warn("Warp reactor core primary coolant failure")

	a.Match("^$", buf.String(), "Unexpected log output")
}

func TestNewLogger(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	logger := NewLogger(DEBUG, "PREFIX")
	logger.Info("What a nice day, innit?")

	a.Match("^\\d\\d\\d\\d\\/\\d\\d\\/\\d\\d \\d\\d:\\d\\d:\\d\\d \\[INFO\\] \\[PREFIX\\] \\[logging_test\\.go:\\d*\\] What a nice day, innit\\?\n$", buf.String(), "Unexpected log output")
}

func TestNewLoggerWithLevelShouldSuppressLowerLevelLogs(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	logger := NewLogger(ERROR, "PREFIX")
	logger.Info("What a nice day, innit?")

	a.Match("^$", buf.String(), "Unexpected log output")
}

func TestNewLoggerWithDefaultLevelShouldSetCorrectLevel(t *testing.T) {
	a := assert.NewAssert(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	SetLevel(ERROR)
	logger := NewDefaultLevelLogger("PREFIX")
	logger.Info("What a nice day, innit?")

	a.Match("^$", buf.String(), "Unexpected log output")
}
