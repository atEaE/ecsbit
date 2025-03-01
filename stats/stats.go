package stats

import "encoding/json"

// ToJSON : Worldの統計情報をJSON形式に変換します
func ToJSON(stats World) []byte {
	b, err := json.Marshal(stats)
	// worldが生成したものを表示することになるので、基本的にエラーは発生しないはず
	// なので、使う側がエラー処理する手間を省くためにpanicで処理している
	if err != nil {
		panic(err)
	}
	return b
}

// ToJSONPretty : Worldの統計情報を整形されたJSON形式に変換します
func ToJSONPretty(stats World) []byte {
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}

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
