package ecsbit

import (
	"github.com/atEaE/ecsbit/config"
	"github.com/atEaE/ecsbit/internal/bits"
	internalconfig "github.com/atEaE/ecsbit/internal/config"
	"github.com/atEaE/ecsbit/internal/primitive"
	"github.com/atEaE/ecsbit/stats"
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
		componentStorage:  newComponentStorage(registeredComponentMaxSize),
		archetypeData:     make([]archetypeData, 0, conf.ArchetypeDefaultCapacity),
		archetypes:        make([]archetype, 0, conf.ArchetypeDefaultCapacity),
		entityIndices:     make([]EntityIndex, 0, conf.EntityPoolDefaultCapacity),
		entityPool:        newEntityPool(conf.EntityPoolDefaultCapacity),
		onCreateCallbacks: make([]func(w *World, e Entity), 0, conf.OnCreateCallbacksDefaultCapacity),
		onRemoveCallbacks: make([]func(w *World, e Entity), 0, conf.OnRemoveCallbacksDefaultCapacity),
		config:            conf,
	}
	// entitiesに先頭sentinelを追加
	// entity側もEntityID = 0がsentinelに該当するため、ID = Indexとして扱うこの仕様に合わせてsentinelを設定している
	world.entityIndices = append(world.entityIndices, EntityIndex{index: 0, archetype: nil})
	// LayoutなしのArchetypeをあらかじめ生成しておく
	world.createArchetype(bits.Mask256{})

	return world
}

// World : ECSの仕組みを提供する構造体
type World struct {
	componentStorage componentStorage            // Componentを管理するStorage
	archetypeData    []archetypeData             // Archetypeから生成されたEntityのデータを保持するSlice
	archetypeLayouts map[bits.Mask256]*archetype // LayoutMaskからArchetypeを取得するためのMap
	archetypes       []archetype                 // Achetypeを管理するSlice
	entityIndices    []EntityIndex               // Archetype内に置けるEntityIndexとArchetypeの関連性を管理する（EntityIDでIndexにアクセスする）
	entityPool       entityPool                  // Entityを管理するPool（生成とリサイクルを管理する）

	onCreateCallbacks []func(w *World, e Entity) // Entity生成時に呼び出すコールバック
	onRemoveCallbacks []func(w *World, e Entity) // Entity削除時に呼び出すコールバック

	config internalconfig.WorldConfig // Worldの設定（内部関数で使う場合があるので予め保持しておく）
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

// CreateEntity : 新しいEntityを生成します
func (w *World) CreateEntity(components ...ComponentID) Entity {
	return w.createEntity(w.findOrCreateArchetype(components))
}

// createEntity : Entityを生成します
func (w *World) createEntity(archetype *archetype) Entity {
	entity := w.entityPool.Get()
	index := archetype.Add(entity)
	w.entityIndices = append(w.entityIndices, EntityIndex{index: index, archetype: archetype})

	for i := range w.onCreateCallbacks {
		w.onCreateCallbacks[i](w, entity)
	}
	return entity
}

// findOrCreateArchetype : 指定されたComponentIDからArchetypeを取得します
// 存在しない場合は新しいArchetypeを生成します
func (w *World) findOrCreateArchetype(components []ComponentID) *archetype {
	if len(components) == 0 {
		return &w.archetypes[noLayoutArchetypeIndex]
	}

	layout := createLayoutMask(components)
	if archetype, ok := w.archetypeLayouts[layout]; ok {
		return archetype
	}
	return w.createArchetype(layout)
}

// RemoveEntity : Entityを削除します
func (w *World) RemoveEntity(e Entity) {
	// 死んでいるEntityをリサイクルするとpoolが破損するのでエラーを返す
	if !w.entityPool.Alive(e) {
		panic(ErrDeadEntityOperation)
	}

	// archetype周りの処理
	index := &w.entityIndices[e.ID()]
	oldArchetype := index.archetype

	swapped := oldArchetype.Remove(index.index)
	w.entityPool.Recycle(e)
	if swapped {
		// Swapが発生した場合、削除指定したIndexの位置にSwapして移動させてEntityがいるので、それを取得してEntityIndexを更新する
		swappedEntity := oldArchetype.GetEntity(index.index)
		w.entityIndices[swappedEntity.ID()].index = index.index
	}
	index.Clear()

	// panic("not implemented")
}

// RegisterComponent : Componentを登録します
func (w *World) RegisterComponent(c component) ComponentID {
	id := w.componentStorage.ComponentID(c)
	return id
}

// createArchetype : Archetypeを生成します
func (w *World) createArchetype(layoutMask bits.Mask256) *archetype {
	idx := primitive.ArchetypeID(len(w.archetypes))
	w.archetypeData = append(w.archetypeData, *newArchetypeData(w.config.EntityPoolDefaultCapacity, layoutMask))
	w.archetypes = append(w.archetypes, *newArchetype(idx, &w.archetypeData[idx]))
	return &w.archetypes[idx]
}

// Stats : Worldの統計情報を取得します
func (w *World) Stats() *stats.World {
	stats := &stats.World{
		Entities: stats.Entities{
			Used:     w.entityPool.Used(),
			Total:    w.entityPool.Total(),
			Capacity: w.entityPool.Cap(),
			Recycled: w.entityPool.Available(),
		},
	}
	return stats
}

// createLayoutMask : 引数に指定されたComponentIDからLayoutMaskを生成します
func createLayoutMask(components []ComponentID) bits.Mask256 {
	mask := bits.Mask256{}
	for _, c := range components {
		mask.Set(uint32(c), true)
	}
	return mask
}
