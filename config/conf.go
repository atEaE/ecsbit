package config

import (
	"github.com/atEaE/ecsbit/internal/config"
)

// defaultConfig : Worldのデフォルトオプション
var defaultConfig = config.WorldConfig{
	RegisterdComponentMaxSize: 256,
	ArchetypeCapacity:         256,
	EntityPoolCapacity:        1024,
	OnCreateCallbacksCapacity: 256,
	OnRemoveCallbacksCapacity: 256,
}

// Default : Worldのデフォルトオプションを取得する
func Default() config.WorldConfig {
	return defaultConfig
}

// config.WorldConfigOption : config.WorldConfigのオプションを提供する関数
type WorldConfigOption func(*config.WorldConfig)

// WithRegisterdComponentMaxSize : 登録可能なComponentの最大数を設定する
func WithRegisterdComponentMaxSize(size uint32) WorldConfigOption {
	if size == 0 {
		panic("RegisterdComponentMaxSize must be greater than 0")
	}
	return func(c *config.WorldConfig) {
		c.RegisterdComponentMaxSize = size
	}
}

// WithEntityPoolCapacity : Entity Poolのキャパシティを設定する
func WithEntityPoolCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.EntityPoolCapacity = capacity
	}
}

// WithOnCreateCallbacksCapacity : Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティを設定する
func WithOnCreateCallbacksCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.OnCreateCallbacksCapacity = capacity
	}
}

// WithOnRemoveCallbacksCapacity : Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティを設定する
func WithOnRemoveCallbacksCapacity(capacity uint32) WorldConfigOption {
	return func(c *config.WorldConfig) {
		c.OnRemoveCallbacksCapacity = capacity
	}
}
