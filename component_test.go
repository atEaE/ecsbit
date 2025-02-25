package ecsbit

import "testing"

func TestComponentStorage_ComponentID(t *testing.T) {
	maxSize := uint32(256)

	type Vector2 struct {
		X, Y float64
	}

	type Rotation struct {
		F float64
	}

	var (
		vector2Comp  = NewComponent[Vector2]()
		rotationComp = NewComponent[Rotation]()
	)

	t.Run("new component", func(t *testing.T) {
		// arrange
		cs := newComponentStorage(maxSize)

		// act & assert
		id := cs.ComponentID(vector2Comp)
		if id != 0 {
			t.Errorf("want %d, got %d", 0, id)
		}
	})

	t.Run("existing component", func(t *testing.T) {
		// setup
		cs := newComponentStorage(maxSize)
		expectedID := cs.ComponentID(vector2Comp)

		// act & assert
		id := cs.ComponentID(vector2Comp)
		if id != expectedID {
			t.Errorf("want %d, got %d", expectedID, id)
		}

		id = cs.ComponentID(rotationComp)
		if id != 1 {
			t.Errorf("want %d, got %d", 1, id)
		}
	})
}
