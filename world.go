package ecsbit

//  defaultOptions : Worldのデフォルトオプション
var defaultOptions = WorldOptions{
	OnCreateCallbacksCapacity: 256,
	OnRemoveCallbacksCapacity: 256,
}

// WorldOptions : Worldのオプションを提供する構造体
type WorldOptions struct {
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
		onCreateCallbacks: make([]func(w *World, e Entity), option.OnCreateCallbacksCapacity),
		onRemoveCallbacks: make([]func(w *World, e Entity), option.OnRemoveCallbacksCapacity),
	}
}

// World : ECSの仕組みを提供する構造体
type World struct {
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

// RegisterArchetype : Archetypeを登録します
// Archetypeは、特定のComponentの組み合わせを持つパターンのことです
func (w *World) RegisterArchetype(components ...ComponentType) {
	panic("not implemented")
}
