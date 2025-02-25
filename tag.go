package ecsbit

// Tag : Groupingなどに使用するタグ
type Tag string

// NewTag : Tag componentを生成する
func NewTag(t Tag) *Component[Tag] {
	c := NewComponent[Tag]()
	c.SetName(string(t))
	return c
}
