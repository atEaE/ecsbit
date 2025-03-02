package ecsbit

import (
	"testing"

	"github.com/atEaE/ecsbit/internal/bits"
)

func TestArchetype_Remove(t *testing.T) {
	t.Run("remove entity swap false", func(t *testing.T) {
		// arrange
		a := newArchetype(0, newArchetypeData(256))
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
		a := newArchetype(0, newArchetypeData(256))
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
		a := newArchetype(0, newArchetypeData(256))
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

func TestConvertToComponentIDs(t *testing.T) {
	// arrange
	m := bits.Mask256{}
	index := []uint32{1, 10, 124}

	for _, i := range index {
		m.Set(i, true)
	}

	// act
	ids := convertToComponentIDs(&m)

	// assert
	_ = ids
	for i := range ids {
		if uint32(ids[i]) != index[i] {
			t.Errorf("unexpected component id: %d", ids[i])
		}
	}
}
