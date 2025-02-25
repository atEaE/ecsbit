package ecsbit

import (
	"fmt"

	"github.com/atEaE/ecsbit/config"
)

var (
	ErrRemoveDeadEntity = fmt.Errorf("can't remove a dead entity")
)

// NewWorld : Worldを生成します
func NewWorld(opts ...config.WorldConfigOption) *World {
	conf := config.Default()
	for _, opt := range opts {
		opt(&conf)
	}

	return &World{
		entityPool:        newEntityPool(uint32(conf.EntityPoolCapacity)),
		onCreateCallbacks: make([]func(w *World, e Entity), 0, conf.OnCreateCallbacksCapacity),
		onRemoveCallbacks: make([]func(w *World, e Entity), 0, conf.OnRemoveCallbacksCapacity),
	}
}

// World : ECSの仕組みを提供する構造体
type World struct {
	entityPool entityPool // Entityを管理するPool（生成とリサイクルを管理する）

	onCreateCallbacks []func(w *World, e Entity) // Entity生成時に呼び出すコールバック
	onRemoveCallbacks []func(w *World, e Entity) // Entity削除時に呼び出すコールバック
}

// PushOnCreateCallback : Entity生成時に呼び出すコールバックを追加します。
// 追加されたコールバックは、追加順に全てのEntity生成時に呼び出されます。
// 追加順と特定のEntityに対する制御を加えたい場合は、コールバック内で制御してください
func (w *World) PushOnCreateCallback(f func(w0 *World, e Entity)) {
	w.onCreateCallbacks = append(w.onCreateCallbacks, f)
}

// PushOnRemoveCallback : Entity削除時に呼び出すコールバックを追加します。
// 追加されたコールバックは、追加順に全てのEntity削除時に呼び出されます。
// 追加順と特定のEntityに対する制御を加えたい場合は、コールバック内で制御してください
func (w *World) PushOnRemoveCallback(f func(w *World, e Entity)) {
	w.onRemoveCallbacks = append(w.onRemoveCallbacks, f)
}

// NewEntity : 新しいEntityを生成します
func (w *World) NewEntity(components []Component) Entity {
	if len(components) == 0 {
		return w.createEntity(nil)
	}
	panic("not implemented")
}

// createEntity : Entityを生成します
func (w *World) createEntity(archetype *Archetype) Entity {
	entity := w.entityPool.Get()

	for i := range w.onCreateCallbacks {
		w.onCreateCallbacks[i](w, entity)
	}
	return entity
}

// RemoveEntity : Entityを削除します
func (w *World) RemoveEntity(e Entity) {
	// 死んでいるEntityをリサイクルすることはできない
	if !w.entityPool.Alive(e) {
		panic(ErrRemoveDeadEntity)
	}

	w.entityPool.Recycle(e)
}

// duplicateComponents
func (w *World) duplicateComponents(c []Component) bool {
	for i := 0; i < len(c); i++ {
		for j := i + 1; j < len(c); j++ {
			if c[i] == c[j] {
				return true
			}
		}
	}
	return false
}
