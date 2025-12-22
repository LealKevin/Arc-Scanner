// Package keyboard provides global keyboard hook functionality.
package keyboard

import (
	"context"
	"log/slog"
	"sync"

	hook "github.com/robotn/gohook"
)

// Handler is a function called when a registered key is pressed.
type Handler func()

// Hook manages global keyboard event listening with context-aware lifecycle.
type Hook struct {
	ctx      context.Context
	handlers map[rune]Handler
	mu       sync.RWMutex
	started  bool
}

// New creates a new Hook that will respect the given context for shutdown.
func New(ctx context.Context) *Hook {
	return &Hook{
		ctx:      ctx,
		handlers: make(map[rune]Handler),
	}
}

// Register adds a handler for the specified key character.
// The handler will be called each time the key is pressed.
func (h *Hook) Register(key rune, handler Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[key] = handler
	slog.Debug("keyboard handler registered", "key", string(key))
}

// Start begins listening for keyboard events in a goroutine.
// It returns immediately. The hook will automatically stop when
// the context is cancelled.
func (h *Hook) Start() {
	h.mu.Lock()
	if h.started {
		h.mu.Unlock()
		return
	}
	h.started = true
	h.mu.Unlock()

	go h.run()
}

func (h *Hook) run() {
	evChan := hook.Start()
	defer hook.End()

	slog.Info("keyboard hook started")

	for {
		select {
		case <-h.ctx.Done():
			slog.Info("keyboard hook stopping (context cancelled)")
			return
		case ev, ok := <-evChan:
			if !ok {
				slog.Info("keyboard hook stopping (channel closed)")
				return
			}
			h.handleEvent(ev)
		}
	}
}

func (h *Hook) handleEvent(ev hook.Event) {
	h.mu.RLock()
	handler, ok := h.handlers[ev.Keychar]
	h.mu.RUnlock()

	if ok {
		slog.Debug("key pressed", "key", string(ev.Keychar))
		handler()
	}
}
