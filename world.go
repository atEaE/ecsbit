package ecsbit

import (
	"github.com/atEaE/ecsbit/internal/primitive"
)

// defaultOptions : Worldのデフォルトオプション
var defaultOptions = WorldOptions{
	EntityPoolCapacity:        1024,
	OnCreateCallbacksCapacity: 256,
	OnRemoveCallbacksCapacity: 256,
}

// WorldOptions : Worldのオプションを提供する構造体
type WorldOptions struct {
	EntityPoolCapacity        int // Entity Poolのキャパシティ
	OnCreateCallbacksCapacity int // Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティ
	OnRemoveCallbacksCapacity int // Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティ
}

// NewWorld : Worldを生成します
func NewWorld(opts ...WorldOptions) *World {
	option := defaultOptions
	if len(opts) > 0 {
		// optional args にしたかっただけなので、先頭要素以外は不要
		option = opts[0]
	}

	return &World{
		entities:          newEntityPool(uint32(option.EntityPoolCapacity)),
		onCreateCallbacks: make([]func(w *World, e Entity), 0, option.OnCreateCallbacksCapacity),
		onRemoveCallbacks: make([]func(w *World, e Entity), 0, option.OnRemoveCallbacksCapacity),
	}
}

// World : ECSの仕組みを提供する構造体
type World struct {
	entities entityPool

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
func (w *World) NewEntity(components []primitive.ComponentType) Entity {
	if len(components) == 0 {
		return w.createEntity(nil)
	}
	panic("not implemented")
}

func (w *World) createEntity(archetype *Archetype) Entity {

	panic("not implemented")
}

// duplicateComponents
func (w *World) duplicateComponents(c []primitive.ComponentType) bool {
	for i := 0; i < len(c); i++ {
		for j := i + 1; j < len(c); j++ {
			if c[i] == c[j] {
				return true
			}
		}
	}
	return false
}
