package graceful

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitShutdown(t *testing.T) {
	tests := []struct {
		callbacks    []ShutdownFunc
		execOnErr    func(error)
		wantErr      bool
		emitBySignal bool
	}{
		{
			callbacks: []ShutdownFunc{
				ShutdownFunc(func() error {
					return nil
				}),
				ShutdownFunc(func() error {
					return errors.New("some error")
				}),
			},
			execOnErr: func(err error) {
				assert.Error(t, err)
			},
			wantErr:      false,
			emitBySignal: true,
		},
		{
			callbacks: []ShutdownFunc{
				ShutdownFunc(func() error {
					// NOTE(a.petrukhin): some very long shutdown...
					time.Sleep(*shutdownTimeout + time.Second)
					return nil
				}),
			},
			execOnErr: func(err error) {
				assert.NoError(t, err)
			},
			wantErr:      true,
			emitBySignal: true,
		},
		{
			callbacks: []ShutdownFunc{
				ShutdownFunc(func() error {
					return nil
				}),
			},
			execOnErr:    nil,
			wantErr:      false,
			emitBySignal: false,
		},
	}

	for _, tt := range tests {
		setupHandler()

		for _, cb := range tt.callbacks {
			AddCallback(cb)
		}

		if tt.emitBySignal {
			handler.stop <- syscall.SIGINT
		} else {
			ShutdownNow()
		}

		ExecOnError(tt.execOnErr)

		err := WaitShutdown()
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		isShuttingDown := IsShuttingDown()
		assert.True(t, isShuttingDown)
	}
}

func TestAddCallback(t *testing.T) {
	setupHandler()
	AddCallback(func() error {
		return nil
	})

	assert.Len(t, handler.callbacks, 1)
}

func TestHandlerShutdown(t *testing.T) {
	h := newHandler(make(chan os.Signal), make(chan struct{}))
	assert.False(t, h.isShuttingDown())
	h.markAsShutdown()
	assert.True(t, h.isShuttingDown())
}
