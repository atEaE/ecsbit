package ecsbit

import (
	"testing"
)

func TestArchetype_Remove(t *testing.T) {
	t.Run("remove entity swap false", func(t *testing.T) {
		// arrange
		a := Archetype{
			entities: make([]Entity, 0, 256),
		}
		a.entities = append(a.entities, NewEntity(0), NewEntity(1), NewEntity(2), NewEntity(3))

		// act
		targetIndex := len(a.entities) - 1
		swapped := a.Remove(uint32(targetIndex))

		// assert
		if swapped {
			t.Errorf("unexpected swapped: %v", swapped)
		}
		if len(a.entities) != 3 {
			t.Errorf("unexpected entity count: %d", len(a.entities))
		}
		if cap(a.entities) != 256 {
			t.Errorf("unexpected entity capacity: %d", cap(a.entities))
		}
	})

	t.Run("remove entity swap true(top)", func(t *testing.T) {
		// arrange
		a := Archetype{
			entities: make([]Entity, 0, 256),
		}
		a.entities = append(a.entities, NewEntity(0), NewEntity(1), NewEntity(2), NewEntity(3))

		// act
		targetIndex := 0
		swapped := a.Remove(uint32(targetIndex))

		// assert
		if !swapped {
			t.Errorf("unexpected swapped: %v", swapped)
		}
		if len(a.entities) != 3 {
			t.Errorf("unexpected entity count: %d", len(a.entities))
		}
		if cap(a.entities) != 256 {
			t.Errorf("unexpected entity capacity: %d", cap(a.entities))
		}
	})

	t.Run("remove entity swap true(middle)", func(t *testing.T) {
		// arrange
		a := Archetype{
			entities: make([]Entity, 0, 256),
		}
		a.entities = append(a.entities, NewEntity(0), NewEntity(1), NewEntity(2), NewEntity(3))

		// act
		targetIndex := 2
		swapped := a.Remove(uint32(targetIndex))

		// assert
		if !swapped {
			t.Errorf("unexpected swapped: %v", swapped)
		}
		if len(a.entities) != 3 {
			t.Errorf("unexpected entity count: %d", len(a.entities))
		}
		if cap(a.entities) != 256 {
			t.Errorf("unexpected entity capacity: %d", cap(a.entities))
		}
	})
}
