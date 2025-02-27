package ecsbit

import (
	"github.com/atEaE/ecsbit/internal/primitive"
)

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

// Remove : Archetypeに属するEntityを削除する
func (a *Archetype) Remove(index uint32) {
	// 対象Indexを削除して、スライスを詰める
	a.entities = append(a.entities[:index], a.entities[index+1:]...)
}
