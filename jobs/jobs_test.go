package jobs

import (
	"testing"
	"time"

	"github.com/arschles/assert"
)

type testPeriodic struct {
	t    *testing.T
	err  error
	i    int
	freq time.Duration
}

func (t *testPeriodic) Do() error {
	t.t.Logf("testPeriodic Do at %s", time.Now())
	t.i++
	return t.err
}

func (t testPeriodic) Frequency() time.Duration {
	return t.freq
}

func TestDoPeriodic(t *testing.T) {
	interval := time.Duration(100) * time.Millisecond
	p := &testPeriodic{t: t, err: nil, freq: interval}
	closeCh1 := DoPeriodic([]Periodic{p})
	time.Sleep(interval * 2)
	assert.True(t, p.i >= 1, "the periodic wasn't called at least once")
	close(closeCh1)
}
