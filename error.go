package ecsbit

import "fmt"

var (
	// ErrDeadEntityOperation : DeadなEntityに対して操作しようとした場合に発生するエラー
	ErrDeadEntityOperation = fmt.Errorf("can't operate a dead entity")
	// ErrDuplicateComponent : 重複したComponentを一緒にEntityに対して追加しようとした場合に発生するエラー
	ErrDuplicateComponent = fmt.Errorf("duplicate components")
)
