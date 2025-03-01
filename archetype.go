package ecsbit

import (
	"github.com/atEaE/ecsbit/internal/primitive"
)

// NewArchetype : Archetypeを生成する
func newArchetype(
	id primitive.ArchetypeID,
	data *archetypeData,
) *archetype {
	return &archetype{
		id:            id,
		archetypeData: data,
	}
}

// archetype : Entityの構成要素を表す構造体
type archetype struct {
	id primitive.ArchetypeID // Archetypeを一意に識別するID

	*archetypeData // archetypeから生成されたEntityのデータを保持する構造体
}

// ID : archetypeを一意に識別するIDを取得する
func (a *archetype) ID() primitive.ArchetypeID {
	return a.id
}

// Count : Archetypeに属するEntityの数を取得する
func (a *archetype) Count() int {
	return len(a.entities)
}

// GetEntity : 指定したIndexのEntityを取得する
func (a *archetype) GetEntity(index uint32) Entity {
	return a.entities[index]
}

// Add : ArchetypeにEntityを追加する
func (a *archetype) Add(e Entity) uint32 {
	a.entities = append(a.entities, e)
	return uint32(len(a.entities) - 1)
}

// Remove : Archetypeに属するEntityを削除する
// 削除Entityと末尾のEntityを入れ替えることで、削除処理を高速化する
func (a *archetype) Remove(index uint32) bool {
	last := len(a.entities) - 1
	if index == uint32(last) {
		a.entities = a.entities[:last]
		return false
	}

	// 末尾のEntityを削除対象のEntityの位置に移動
	a.entities[index], a.entities[last] = a.entities[last], a.entities[index]
	a.entities = a.entities[:last]
	return true
}

// newArchetypeData : archetypeDataを生成する
func newArchetypeData(
	entityCapacity uint32,
) *archetypeData {
	return &archetypeData{
		entities: make([]Entity, 0, entityCapacity),
	}
}

// archetypeData : archetypeから生成されたEntityのデータを保持する構造体
// archetype : archetypeData は 1 : 1 の関係
type archetypeData struct {
	entities []Entity // Archetypeに属するEntity
}
