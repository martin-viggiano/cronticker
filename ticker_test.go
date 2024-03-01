package crontick

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// defaultOptions is used for the test parsers.
	defaultOptions cron.ParseOption = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor

	// timeout is the time duration to wait for a tick to happen or not
	timeout = 2 * time.Second
)

func TestNewTicker(t *testing.T) {
	tt := []struct {
		name      string
		spec      string
		wantError error
	}{
		{name: "default", spec: "0 0 * * *", wantError: nil},
		{name: "invalid spec", spec: "invalid", wantError: fmt.Errorf("failed to parse spec: %w", errors.New("expected exactly 5 fields, found 1: [invalid]"))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTicker(tc.spec)

			assert.Equal(t, tc.wantError, err)
		})
	}
}

func TestNewTickerWithParser(t *testing.T) {
	tt := []struct {
		name      string
		spec      string
		wantError error
	}{
		{name: "default", spec: "0", wantError: nil},
		{name: "invalid spec", spec: "100", wantError: fmt.Errorf("failed to parse spec: %w", errors.New("end of range (100) above maximum (59): 100"))},
	}

	parser := cron.NewParser(cron.Minute)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTickerWithParser(tc.spec, parser)

			assert.Equal(t, tc.wantError, err)
		})
	}
}

func TestTickerStop(t *testing.T) {
	ticker, err := NewTickerWithParser("* * * * * *", cron.NewParser(defaultOptions))
	require.NoError(t, err)

	assert.NotPanics(t, ticker.Stop)

	timeout := time.NewTimer(2 * time.Second)

	select {
	case <-ticker.C:
		assert.Fail(t, "tick not expected")
	case <-timeout.C:
	}
}

func TestTickerReset(t *testing.T) {
	ticker, err := NewTickerWithParser("@yearly", cron.NewParser(defaultOptions))
	require.NoError(t, err)

	err = ticker.Reset("* * * * * *")
	assert.NoError(t, err)

	timeout := time.NewTimer(2 * time.Second)

	select {
	case <-ticker.C:
	case <-timeout.C:
		assert.Fail(t, "expected a tick")
	}
}

func TestTicker(t *testing.T) {
	ticker, err := NewTickerWithParser("* * * * * *", cron.NewParser(defaultOptions))
	require.NoError(t, err)

	timeout := time.NewTimer(5 * time.Second)

	i := 0

loop:
	select {
	case <-ticker.C:
		i++
		if i >= 2 {
			break loop
		}

	case <-timeout.C:
		assert.Fail(t, "expected a tick")
	}
}
