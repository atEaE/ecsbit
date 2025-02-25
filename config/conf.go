package config

// defaultConfig : Worldのデフォルトオプション
var defaultConfig = worldConfig{
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
	EntityPoolCapacity        uint32 // Entity Poolのキャパシティ
	OnCreateCallbacksCapacity uint32 // Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティ
	OnRemoveCallbacksCapacity uint32 // Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティ
}

// WorldConfigOption : WorldConfigのオプションを提供する関数
type WorldConfigOption func(*worldConfig)

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
