package config

// WorldConfig : Worldのオプションを提供する構造体
type WorldConfig struct {
	RegisterdComponentMaxSize uint32 // 登録可能なComponentの最大数
	ArchetypeCapacity         uint32 // Archetypeのキャパシティ
	EntityPoolCapacity        uint32 // Entity Poolのキャパシティ
	OnCreateCallbacksCapacity uint32 // Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティ
	OnRemoveCallbacksCapacity uint32 // Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティ
}
