package ecsbit

import (
	"testing"

	"github.com/atEaE/ecsbit/config"
)

func TestWorld_NewWorldWithOptions(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// arrange
		w := NewWorld()

		// assert
		if got := cap(w.onCreateCallbacks); got != int(config.Default().OnCreateCallbacksDefaultCapacity) {
			t.Errorf("unexpected result: got %v, want %v", got, int(config.Default().OnCreateCallbacksDefaultCapacity))
		}
		if got := cap(w.onRemoveCallbacks); got != int(config.Default().OnRemoveCallbacksDefaultCapacity) {
			t.Errorf("unexpected result: got %v, want %v", got, int(config.Default().OnRemoveCallbacksDefaultCapacity))
		}
	})

	t.Run("custom", func(t *testing.T) {
		// arrange
		expectedOnCreateCap := uint32(512)
		expectedOnRemoveCap := uint32(128)

		w := NewWorld(
			config.WithOnCreateCallbacksDefaultCapacity(expectedOnCreateCap),
			config.WithOnRemoveCallbacksDefaultCapacity(expectedOnRemoveCap),
		)

		// assert
		if got := cap(w.onCreateCallbacks); got != int(expectedOnCreateCap) {
			t.Errorf("unexpected result: got %v, want %v", got, int(expectedOnCreateCap))
		}
		if got := cap(w.onRemoveCallbacks); got != int(expectedOnRemoveCap) {
			t.Errorf("unexpected result: got %v, want %v", got, int(expectedOnRemoveCap))
		}
	})
}

func TestWorld_getArchetype(t *testing.T) {
	t.Run("no layout archetype", func(t *testing.T) {
		// arrange
		w := NewWorld()

		// act
		a := w.getArchetype(nil)

		// assert
		if got := a.ID(); got != 0 {
			t.Errorf("unexpected result: got %v, want %v", got, 0)
		}
		if got := a.Count(); got != 0 {
			t.Errorf("unexpected result: got %v, want %v", got, 0)
		}
	})
}
