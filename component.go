package ecsbit

import (
	"reflect"
	"sync/atomic"

	"github.com/atEaE/ecsbit/internal/primitive"
)

// componentIDCounter : ComponentID生成のためのカウンタ
var componentIDCounter uint32 = 0

// nextComponentID : ComponentIDを生成する
func nextComponentID() primitive.ComponentTypeID {
	return primitive.ComponentTypeID(atomic.AddUint32(&componentIDCounter, 1) - 1)
}

// NewComponent : Generics指定の型をComponentとして生成する
func NewComponent[T any]() *Component[T] {
	var t T
	typ := reflect.TypeOf(t)
	return &Component[T]{
		id:   nextComponentID(),
		name: typ.Name(),
		typ:  typ,
	}
}

// Component : Componentを表す構造体
type Component[T any] struct {
	id   primitive.ComponentTypeID
	name string
	typ  reflect.Type
}

// ID : componentTypeを一意に識別するIDを取得する
func (c *Component[T]) ID() primitive.ComponentTypeID {
	return c.id
}

// Name : componentの名前（基本的にベースになっているものの型名）
func (c *Component[T]) Name() string {
	return c.name
}

// OriginType : componentの元になっている型情報を取得する
func (c *Component[T]) OriginType() reflect.Type {
	return c.typ
}

// SetName : componentの名前を設定する.
// 自分で設定しない場合は、基本的にベースになっているものの型名が名前になる
func (c *Component[T]) SetName(n string) {
	c.name = n
}

type ComponentPool[T any] struct {
	data []T
}
