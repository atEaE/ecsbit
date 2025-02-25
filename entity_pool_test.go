package ecsbit

import (
	"errors"
	"testing"
)

func TestEntityPool_RecycleAndGet(t *testing.T) {
	getPool := func() entityPool {
		pool := newEntityPool(10)
		for i := 1; i < 10; i++ {
			pool.entities = append(pool.entities, NewEntity(EntityID(i)))
		}
		return pool
	}

	t.Run("linked list recycle", func(t *testing.T) {
		// setup
		pool := getPool()

		expectedFirstRecycleStackID := EntityID(0)
		// recycle marked phase1
		{
			// arrange
			e := NewEntity(2)
			expectedFirstRecycleStackID = e.ID()

			// act
			pool.Recycle(e)

			// assert
			if pool.Alive(e) {
				t.Errorf("expected entity is dead, but alive")
			}
			if !pool.IsRecycleWait(e.ID()) {
				t.Errorf("expected entity is recycle wait, but not")
			}
			if pool.Available() != 1 {
				t.Errorf("unexpected available count: %d", pool.Available())
			}
		}

		expectedSecondRecycleStackID := EntityID(0)
		// recycle marked phase2
		{
			// arrange
			e := NewEntity(8)
			expectedSecondRecycleStackID = e.ID()

			// act
			pool.Recycle(e)

			// assert
			if pool.Alive(e) {
				t.Errorf("expected entity is dead, but alive")
			}
			if !pool.IsRecycleWait(e.ID()) {
				t.Errorf("expected entity is recycle wait, but not")
			}
			if pool.Available() != 2 {
				t.Errorf("unexpected available count: %d", pool.Available())
			}
		}

		// recycle get phase1
		{
			// arrange
			e := pool.Get()

			// act & assert
			// stackなので、最後にRecycleしたものが最初に取り出される
			if e.ID() != expectedSecondRecycleStackID {
				t.Errorf("unexpected entity id: %d", e.ID())
			}
			if e.Version() != 1 {
				t.Errorf("unexpected entity version: %d", e.Version())
			}
			if pool.IsRecycleWait(e.ID()) {
				t.Errorf("expected entity is not recycle wait, but recycle wait")
			}
			if pool.Available() != 1 {
				t.Errorf("unexpected available count: %d", pool.Available())
			}
		}

		// recycle get phase2
		{
			// arrange
			e := pool.Get()

			// act & assert
			if e.ID() != expectedFirstRecycleStackID {
				t.Errorf("unexpected entity id: %d", e.ID())
			}
			if e.Version() != 1 {
				t.Errorf("unexpected entity version: %d", e.Version())
			}
			if pool.IsRecycleWait(e.ID()) {
				t.Errorf("expected entity is not recycle wait, but recycle wait")
			}
			if pool.Available() != 0 {
				t.Errorf("unexpected available count: %d", pool.Available())
			}
		}
	})

	t.Run("sentinel recycle error", func(t *testing.T) {
		// setup
		pool := newEntityPool(10)

		// arrange
		e := zeroEntity

		// act & assert
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected panic, but not occurred")
			}

			if !errors.Is(err.(error), ErrRecycleSentinel) {
				t.Errorf("unexpected error: %v", err)
			}
		}()
		pool.Recycle(e)
	})
}
