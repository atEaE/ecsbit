package config

// defaultConfig : Worldのデフォルトオプション
var defaultConfig = worldConfig{
	RegisterdComponentMaxSize: 256,
	EntityPoolCapacity:        1024,
	OnCreateCallbacksCapacity: 256,
	OnRemoveCallbacksCapacity: 256,
}

// Default : Worldのデフォルトオプションを取得する
func Default() worldConfig {
	return defaultConfig
}

// WorldConfig : Worldのオプションを提供する構造体
type worldConfig struct {
	RegisterdComponentMaxSize uint32 // 登録可能なComponentの最大数
	EntityPoolCapacity        uint32 // Entity Poolのキャパシティ
	OnCreateCallbacksCapacity uint32 // Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティ
	OnRemoveCallbacksCapacity uint32 // Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティ
}

// WorldConfigOption : WorldConfigのオプションを提供する関数
type WorldConfigOption func(*worldConfig)

// WithRegisterdComponentMaxSize : 登録可能なComponentの最大数を設定する
func WithRegisterdComponentMaxSize(size uint32) WorldConfigOption {
	if size == 0 {
		panic("RegisterdComponentMaxSize must be greater than 0")
	}
	return func(c *worldConfig) {
		c.RegisterdComponentMaxSize = size
	}
}

// WithEntityPoolCapacity : Entity Poolのキャパシティを設定する
func WithEntityPoolCapacity(capacity uint32) WorldConfigOption {
	return func(c *worldConfig) {
		c.EntityPoolCapacity = capacity
	}
}

// WithOnCreateCallbacksCapacity : Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティを設定する
func WithOnCreateCallbacksCapacity(capacity uint32) WorldConfigOption {
	return func(c *worldConfig) {
		c.OnCreateCallbacksCapacity = capacity
	}
}

// WithOnRemoveCallbacksCapacity : Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティを設定する
func WithOnRemoveCallbacksCapacity(capacity uint32) WorldConfigOption {
	return func(c *worldConfig) {
		c.OnRemoveCallbacksCapacity = capacity
	}
}
