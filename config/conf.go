package config

import (
	"github.com/atEaE/ecsbit/internal/config"
)

// defaultConfig : Worldのデフォルトオプション
var defaultConfig = config.WorldConfig{
	ArchetypeDefaultCapacity:         256,
	EntityPoolDefaultCapacity:        1024,
	OnCreateCallbacksDefaultCapacity: 256,
	OnRemoveCallbacksDefaultCapacity: 256,
}

// Default : Worldのデフォルトオプションを取得する
func Default() config.WorldConfig {
	return defaultConfig
}

// config.WorldConfigOption : config.WorldConfigのオプションを提供する関数
type WorldConfigOption func(*config.WorldConfig)

// WithArchetypeDefaultCapacity : Archetypeのキャパシティを設定する
func WithArchetypeDefaultCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.ArchetypeDefaultCapacity = capacity
	}
}

// WithEntityPoolDefaultCapacity : Entity Poolのキャパシティを設定する
func WithEntityPoolDefaultCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.EntityPoolDefaultCapacity = capacity
	}
}

// WithOnCreateCallbacksDefaultCapacity : Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティを設定する
func WithOnCreateCallbacksDefaultCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.OnCreateCallbacksDefaultCapacity = capacity
	}
}

// WithOnRemoveCallbacksDefaultCapacity : Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティを設定する
func WithOnRemoveCallbacksDefaultCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.OnRemoveCallbacksDefaultCapacity = capacity
	}
}
