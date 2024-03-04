package cronticker

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// A Ticker holds a channel that delivers “ticks” of a clock
// acording to a crontab schedule.
type Ticker struct {
	C    chan time.Time
	stop chan struct{}

	schedule cron.Schedule
	parser   cron.Parser
}

// NewTicker returns a new Ticker containing a channel that will send
// the current time on the channel according to the crontab schedule.
// The crontab schedule is specified by the spec argument and will be
// parsed using the standard specification.
// Stop the ticker to release associated resources.
func NewTicker(spec string) (*Ticker, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	c := make(chan time.Time, 1)
	stop := make(chan struct{})

	return newTicker(spec, parser, c, stop)
}

// NewTickerWithParser returns a new Ticker containing a channel that will send
// the current time on the channel accordint to the crontab schedule.
// The crontab schedule is specified by the spec argument and will be
// parsed using the provided parser.
// Stop the ticker to release associated resources.
func NewTickerWithParser(spec string, parser cron.Parser) (*Ticker, error) {
	c := make(chan time.Time, 1)
	stop := make(chan struct{})

	return newTicker(spec, parser, c, stop)
}

// newTicker creates a new ticker using the provided spec, paser and channel.
func newTicker(spec string, parser cron.Parser, c chan time.Time, stop chan struct{}) (*Ticker, error) {
	schedule, err := parser.Parse(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	t := Ticker{
		C:    c,
		stop: stop,

		schedule: schedule,
		parser:   parser,
	}

	go t.runTimer()
	return &t, nil
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a concurrent goroutine
// reading from the channel from seeing an erroneous "tick".
func (t *Ticker) Stop() {
	t.stop <- struct{}{}
}

// Reset stops a ticker and resets its crontab schedule to the specified spec.
// An error is returned if the crontab specification cant be pased.
func (t *Ticker) Reset(spec string) error {
	t.Stop()

	var err error
	t, err = newTicker(spec, t.parser, t.C, t.stop)
	if err != nil {
		return err
	}

	return nil
}

// runTimer handles the ticker logic.
func (t *Ticker) runTimer() {
	next := t.schedule.Next(time.Now())
	timer := time.NewTimer(time.Until(next))

	for {
		select {
		case tick := <-timer.C:
			t.C <- tick
			next := t.schedule.Next(tick)
			timer.Reset(time.Until(next))
		case <-t.stop:
			timer.Stop()
			return
		}
	}
}
