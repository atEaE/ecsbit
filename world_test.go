package ecsbit

import "testing"

func TestWorldOptions(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// arrange
		w := NewWorld()

		// assert
		if got := cap(w.onCreateCallbacks); got != defaultOptions.OnCreateCallbacksCapacity {
			t.Errorf("unexpected result: got %v, want %v", got, defaultOptions.OnCreateCallbacksCapacity)
		}
		if got := cap(w.onRemoveCallbacks); got != defaultOptions.OnRemoveCallbacksCapacity {
			t.Errorf("unexpected result: got %v, want %v", got, defaultOptions.OnRemoveCallbacksCapacity)
		}
	})

	t.Run("custom", func(t *testing.T) {
		// arrange
		opts := WorldOptions{
			OnCreateCallbacksCapacity: 512,
			OnRemoveCallbacksCapacity: 128,
		}
		w := NewWorld(opts)

		// assert
		if got := cap(w.onCreateCallbacks); got != opts.OnCreateCallbacksCapacity {
			t.Errorf("unexpected result: got %v, want %v", got, opts.OnCreateCallbacksCapacity)
		}
		if got := cap(w.onRemoveCallbacks); got != opts.OnRemoveCallbacksCapacity {
			t.Errorf("unexpected result: got %v, want %v", got, opts.OnRemoveCallbacksCapacity)
		}
	})
}
