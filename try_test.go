package golib

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestTryZeroTimesShouldCallFunctionIndefinitelyAndEventuallyReturn(t *testing.T) {
	cnt := 0
	rand.Seed(time.Now().UnixNano())
	successAfterNTries := rand.Intn(100000)
	err := Try(0, 0, func() error {
		cnt++
		if cnt == successAfterNTries {
			return nil
		}
		return errors.New("Some error")
	})
	assert.NoError(t, err, "Try returned with an error")
	assert.Equal(t, cnt, successAfterNTries, "The callback wasn't called enough times")
}

func TestTryTwoTimesShouldCallPassingFunctionOnlyOnce(t *testing.T) {
	called := 0
	err := Try(2, 0, func() error {
		called++
		return nil
	})
	assert.Equal(t, 1, called, fmt.Sprintf("f has been called %d times", called))
	assert.NoError(t, err, "Try returned with an error")
}

func TestTryTwoTimesShouldCallFailingFunctionTwoTimes(t *testing.T) {
	called := 0
	err := Try(2, 0, func() error {
		called++
		if called == 1 {
			return errors.New("Warning. Warp reactor core primary coolant failure")
		}
		return nil
	})
	assert.Equal(t, 2, called, "f has been called %d times", called)
	assert.NoError(t, err, "Try returned with an error")
}
