package primitive

import "reflect"

// ComponentTypeID : componentを一意に表す型
type ComponentTypeID uint32

// ComponentType : componentの型情報を提供するインターフェース
type ComponentType interface {
	ID() ComponentTypeID      // componentTypeを一意に識別するIDを取得する
	Name() string             // componentの名前（基本的にベースになっているものの型名）
	OriginType() reflect.Type // componentの元になっている型情報を取得する
}
