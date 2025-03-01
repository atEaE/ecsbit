package stats

// World : Worldの統計情報
type World struct {
	Entities Entities `json:"entities"` // Entityの統計情報
}

// Entities : Entityの統計情報
type Entities struct {
	// Used : 使用中のEntity数
	Used int `json:"used"`
	// Total : Entityの総数
	Total int `json:"total"`
	// Recycled : 再利用可能なEntity数
	Recycled int `json:"recycled"`
	// Capacity : Entity Poolのキャパシティ
	Capacity int `json:"capacity"`
}
