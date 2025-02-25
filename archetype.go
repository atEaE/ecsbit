package ecsbit

import (
	"sync/atomic"

	"github.com/atEaE/ecsbit/internal/primitive"
)

// archetypeIDCounter : ArchetypeIDを生成するためのカウンタ
var archetypeIDCounter uint32 = 0

// nextArchetypeID : ArchetypeIDを生成する
func nextArchetypeID() primitive.ArchetypeID {
	return primitive.ArchetypeID(atomic.AddUint32(&archetypeIDCounter, 1) - 1)
}

// Archetype : Entityの構成要素を表す構造体
type Archetype struct {
	id       primitive.ArchetypeID       // Archetypeを一意に識別するID
	layout   []primitive.ComponentTypeID // Archetypeの構成要素
	entities []Entity                    // Archetypeに属するEntity
}

// ID : archetypeを一意に識別するIDを取得する
func (a *Archetype) ID() primitive.ArchetypeID {
	return a.id
}

// Layout : Archetypeの構成要素を取得する
func (a *Archetype) Layout() []primitive.ComponentTypeID {
	return a.layout
}

// Entities : Archetypeに属するEntityを取得する
func (a *Archetype) Entities() []Entity {
	return a.entities
}

// Count : Archetypeに属するEntityの数を取得する
func (a *Archetype) Count() int {
	return len(a.entities)
}

// ArchetypeLayoutIndex : Archetypeの構成要素をインデックスとして管理する
type ArchetypeLayoutIndex struct {
	layouts [][]primitive.ComponentTypeID
}

// NewArchetypeIndex : ArchetypeIndexを生成する
func NewArchetypeIndex() *ArchetypeLayoutIndex {
	return &ArchetypeLayoutIndex{
		layouts: make([][]primitive.ComponentTypeID, 256, 512),
	}
}

// Set : Archetypeの構成要素を登録する
func (i *ArchetypeLayoutIndex) Set(id primitive.ArchetypeID, layout []primitive.ComponentTypeID) {
	if len(i.layouts) > int(id) {
		i.layouts[id] = layout
		return
	}
	// IDが大きい場合は、スライスを拡張する
	newLayouts := make([][]primitive.ComponentTypeID, id+1)
	copy(newLayouts, i.layouts)
	newLayouts[id] = layout
	i.layouts = newLayouts
}
