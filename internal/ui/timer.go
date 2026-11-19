package ui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
)

// Timer manages the monitoring timer functionality
type Timer struct {
	label      *widget.Label
	ticker     *time.Ticker
	stopChan   chan struct{}
	mainWindow fyne.Window
	mu         sync.Mutex
}

// NewTimer creates a new Timer instance
func NewTimer(mainWindow fyne.Window, label *widget.Label) *Timer {
	return &Timer{
		label:      label,
		mainWindow: mainWindow,
	}
}

// Start starts the timer for monitoring elapsed time
func (t *Timer) Start(startTime time.Time, timeLimitMinutes *int, onStop func() bool) {
	var limitReached bool

	t.mu.Lock()
	t.stopLocked()
	t.stopChan = make(chan struct{})
	t.ticker = time.NewTicker(10 * time.Millisecond)
	ticker, stopChan := t.ticker, t.stopChan
	t.mu.Unlock()

	go func() {
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)

				// Check if time limit is reached
				if timeLimitMinutes != nil && !limitReached {
					limitDuration := time.Duration(*timeLimitMinutes) * time.Minute
					if elapsed >= limitDuration {
						limitReached = true
						// Timer reached limit, stop it and notify on main thread
						fyne.Do(func() {
							if !t.stopIfCurrent(stopChan) {
								return
							}
							if onStop == nil || onStop() {
								dialog.ShowInformation(locale.T("dialog.info"), locale.T("dialog.time_limit_reached"), t.mainWindow)
							}
						})
						return
					}
				}

				// Update timer label on main thread
				finalElapsed := elapsed
				fyne.Do(func() {
					t.mu.Lock()
					defer t.mu.Unlock()
					if t.stopChan == stopChan && t.label != nil {
						t.label.SetText(t.formatDuration(finalElapsed))
					}
				})
			case <-stopChan:
				return
			}
		}
	}()
}

// Stop stops the running timer
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopLocked()
}

func (t *Timer) stopIfCurrent(stopChan chan struct{}) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopChan != stopChan {
		return false
	}
	t.stopLocked()
	return true
}

func (t *Timer) stopLocked() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
	if t.stopChan != nil {
		close(t.stopChan)
		t.stopChan = nil
	}
	// Do not clear timer label - keep showing the last elapsed time
	// Timer label is only cleared when switching to a different competition
}

// Reset clears the timer display
func (t *Timer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.label != nil {
		t.label.SetText("")
	}
}

// formatDuration formats a duration as HH:MM:SS.ss (hours:minutes:seconds.centiseconds)
func (t *Timer) formatDuration(d time.Duration) string {
	centisecond := 10 * time.Millisecond
	d = d.Round(centisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	cs := d / centisecond

	return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, cs)
}
