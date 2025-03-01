package config

// WorldConfig : Worldのオプションを提供する構造体
type WorldConfig struct {
	ArchetypeDefaultCapacity         uint32 // Archetypeのキャパシティ
	EntityPoolDefaultCapacity        uint32 // Entity Poolのキャパシティ
	OnCreateCallbacksDefaultCapacity uint32 // Entity生成時に呼び出すコールバック群を保持するsliceのキャパシティ
	OnRemoveCallbacksDefaultCapacity uint32 // Entity削除時に呼び出すコールバック群を保持するsliceのキャパシティ
}
