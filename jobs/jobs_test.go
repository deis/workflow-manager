package jobs

import (
	"testing"
	"time"

	"github.com/arschles/assert"
)

type testPeriodic struct {
	t   *testing.T
	err error
	i   int
}

func (t *testPeriodic) Do() error {
	t.t.Logf("testPeriodic Do at %s", time.Now())
	t.i++
	return t.err
}

func TestRunJobs(t *testing.T) {
	p := &testPeriodic{t: t, err: nil}
	runJobs([]Periodic{p})
	assert.Equal(t, p.i, 1, "number of invocations")
}

func TestDoPeriodic(t *testing.T) {
	interval := time.Duration(100) * time.Millisecond
	p := &testPeriodic{t: t, err: nil}
	closeCh1 := DoPeriodic([]Periodic{p}, interval)
	time.Sleep(interval * 2)
	assert.True(t, p.i >= 1, "the periodic wasn't called at least once")
	close(closeCh1)
}
