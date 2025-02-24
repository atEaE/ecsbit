package ecsbit

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestEntityString(t *testing.T) {
	// setup
	bytes := []byte{0, 0, 0, 1, 255, 255, 255, 255}
	base := binary.BigEndian.Uint64(bytes)

	testcases := []struct {
		title  string
		entity Entity
		opts   func(e Entity) Entity
		want   string
	}{
		{
			title:  "id only",
			entity: NewEntity(2),
			opts:   nil,
			want:   "Entity: {id: 2, version: 0}",
		},
		{
			title:  "version incremented",
			entity: NewEntity(5),
			opts: func(e Entity) Entity {
				return e.IncrementVersion()
			},
			want: "Entity: {id: 5, version: 1}",
		},
		{
			title:  "version overflow",
			entity: Entity(base),
			opts: func(e Entity) Entity {
				return e.IncrementVersion()
			},
			want: "Entity: {id: 1, version: 0}",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			// arrange
			if tc.opts != nil {
				tc.entity = tc.opts(tc.entity)
			}

			// act & assert
			if got := tc.entity.String(); got != tc.want {
				t.Errorf("unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEntityVersion(t *testing.T) {
	t.Run("version 1", func(t *testing.T) {
		// act & assert
		e := NewEntity(EntityID(1))
		if got := e.Version(); got != 0 {
			t.Errorf("unexpected version: got %v, want %v", got, 1)
		}
	})

	t.Run("version limit", func(t *testing.T) {
		// arrange
		bytes := []byte{0, 0, 0, 1, 255, 255, 255, 255}
		base := binary.BigEndian.Uint64(bytes)
		e := Entity(base)

		// act & assert
		if got := e.Version(); got != math.MaxUint32 {
			t.Errorf("unexpected version: got %v, want %v", got, math.MaxUint32)
		}
	})
}

func TestEntityIncrementVersion(t *testing.T) {
	t.Run("increment version", func(t *testing.T) {
		// arrange
		e := NewEntity(EntityID(1))
		if check := e.ID(); check != 1 {
			t.Fatalf("unexpected id: got %v, want %v", check, 1)
		}
		if check := e.Version(); check != 0 {
			t.Fatalf("unexpected version: got %v, want %v", check, 0)
		}

		// act & assert
		e = e.IncrementVersion()
		if got := e.ID(); got != 1 {
			t.Errorf("unexpected id: got %v, want %v", got, 1)
		}
		if got := e.Version(); got != 1 {
			t.Errorf("unexpected version: got %v, want %v", got, 1)
		}
	})

	t.Run("overflow version", func(t *testing.T) {
		// arrange
		bytes := []byte{0, 0, 0, 23, 255, 255, 255, 255}
		base := binary.BigEndian.Uint64(bytes)
		e := Entity(base)
		if check := e.ID(); check != 23 {
			t.Fatalf("unexpected id: got %v, want %v", check, 23)
		}
		if check := e.Version(); check != math.MaxUint32 {
			t.Fatalf("unexpected version: got %v, want %v", check, math.MaxUint32)
		}

		// act & assert
		e = e.IncrementVersion()
		if got := e.ID(); got != 23 {
			t.Errorf("unexpected id: got %v, want %v", got, 23)
		}
		if got := e.Version(); got != 0 {
			t.Errorf("unexpected version: got %v, want %v", got, 0)
		}
	})
}
