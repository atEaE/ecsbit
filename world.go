package ecsbit

import (
	"fmt"

	"github.com/atEaE/ecsbit/config"
)

var (
	ErrRemoveDeadEntity = fmt.Errorf("can't remove a dead entity")
)

const (
	noLayoutArchetypeIndex = 0
)

// NewWorld : Worldを生成します
func NewWorld(opts ...config.WorldConfigOption) *World {
	conf := config.Default()
	for _, opt := range opts {
		opt(&conf)
	}

	world := &World{
		componentStorage:  newComponentStorage(conf.RegisterdComponentMaxSize),
		archetypes:        make([]archetype, 0, conf.ArchetypeCapacity),
		entities:          make([]EntityIndex, 0, conf.EntityPoolCapacity),
		entityPool:        newEntityPool(conf.EntityPoolCapacity),
		onCreateCallbacks: make([]func(w *World, e Entity), 0, conf.OnCreateCallbacksCapacity),
		onRemoveCallbacks: make([]func(w *World, e Entity), 0, conf.OnRemoveCallbacksCapacity),
	}
	// LayoutなしのArchetypeをあらかじめ生成しておく
	world.archetypes = append(world.archetypes, archetype{id: 0, entities: make([]Entity, 0, conf.EntityPoolCapacity)})

	return world
}

// World : ECSの仕組みを提供する構造体
type World struct {
	componentStorage componentStorage // Componentを管理するStorage
	archetypes       []archetype
	entities         []EntityIndex
	entityPool       entityPool // Entityを管理するPool（生成とリサイクルを管理する）

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
func (w *World) NewEntity(components ...ComponentID) Entity {
	if len(components) == 0 {
		return w.createEntity(w.getArchetype(nil))
	}
	if w.duplicateComponents(components) {
		panic("duplicate components")
	}

	panic("not implemented")
}

// createEntity : Entityを生成します
func (w *World) createEntity(archetype *archetype) Entity {
	entity := w.entityPool.Get()
	w.entities = append(w.entities, EntityIndex{index: uint32(entity.ID()), archetype: archetype})

	for i := range w.onCreateCallbacks {
		w.onCreateCallbacks[i](w, entity)
	}
	return entity
}

func (w *World) getArchetype(components []ComponentID) *archetype {
	if len(components) == 0 {
		return &w.archetypes[noLayoutArchetypeIndex]
	}
	panic("not implemented")
}

// RemoveEntity : Entityを削除します
func (w *World) RemoveEntity(e Entity) {
	// 死んでいるEntityをリサイクルするとpoolが破損するのでエラーを返す
	if !w.entityPool.Alive(e) {
		panic(ErrRemoveDeadEntity)
	}

	// archetype周りの処理
	index := w.entities[e.ID()]
	oldArchetype := index.archetype

	swapped := oldArchetype.Remove(index.index)
	w.entityPool.Recycle(e)
	if swapped {
		// Swapが発生した場合、削除指定したIndexの位置にSwapして移動させてEntityがいるので、それを取得してEntityIndexを更新する
		swappedEntity := oldArchetype.GetEntity(index.index)
		w.entities[swappedEntity.ID()].index = index.index
	}

	panic("not implemented")
}

// duplicateComponents
func (w *World) duplicateComponents(c []ComponentID) bool {
	for i := 0; i < len(c); i++ {
		for j := i + 1; j < len(c); j++ {
			if c[i] == c[j] {
				return true
			}
		}
	}
	return false
}

// RegisterComponent : Componentを登録します
func (w *World) RegisterComponent(c component) ComponentID {
	id := w.componentStorage.ComponentID(c)
	return id
}
