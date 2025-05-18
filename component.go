package ecsbit

import (
	"reflect"

	"github.com/atEaE/ecsbit/internal/bits"
)

// ComponentID : World単位でComponentを一意に表すID
type ComponentID uint32

// NewComponent : Generics指定の型をComponentとして生成する
func NewComponent[T any]() component {
	var t T
	typ := reflect.TypeOf(t)
	return component{
		name: typ.Name(),
		typ:  typ,
	}
}

// component : componentを表す構造体
// ユーザーには公開せず、生成経路を制限する
type component struct {
	name string
	typ  reflect.Type
}

// Name : componentの名前（基本的にベースになっているものの型名）
func (c *component) Name() string {
	return c.name
}

// OriginType : componentの元になっている型情報を取得する
func (c *component) Type() reflect.Type {
	return c.typ
}

// SetName : componentの名前を設定する.
// 自分で設定しない場合は、基本的にベースになっているものの型名が名前になる
func (c *component) SetName(n string) {
	c.name = n
}

const (
	// registerdComponentMaxSize : 登録可能なComponentの最大数
	// ArchetypeのLayoutを表すビットマスクの最大サイズに合わせて設定している。これ以上登録してもBitMaskで表現できないため。
	registeredComponentMaxSize = bits.Mask256Max
)

// newComponentStorage : componentStorageを生成する
func newComponentStorage(maxSize uint32) componentStorage {
	maxSizeInt := int(maxSize)
	return componentStorage{
		Components: make(map[component]ComponentID, maxSizeInt),
		Types:      make([]component, maxSize),
		IDs:        make([]ComponentID, 0, maxSize),

		maxSize: int(maxSize),
	}
}

// componentStorage : componentを保管するストレージ
type componentStorage struct {
	Components map[component]ComponentID
	Types      []component
	IDs        []ComponentID

	maxSize int
}

// ComponentID : ComponentIDを取得する. storageに存在しない場合は、登録後のIDを返す
func (s *componentStorage) ComponentID(c component) ComponentID {
	if id, ok := s.Components[c]; ok {
		return id
	}
	return s.register(c)
}

// register : componentを登録する
func (s *componentStorage) register(c component) ComponentID {
	idInt := len(s.Components)
	if idInt >= s.maxSize {
		panic("componentStorage is full")
	}
	newID := ComponentID(idInt)
	s.Components[c], s.Types[newID] = newID, c
	s.IDs = append(s.IDs, newID)
	return newID
}
