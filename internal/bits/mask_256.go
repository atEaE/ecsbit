package bits

const (
	// BitMask256Max : 256bitのビットマスクの最大値
	Mask256Max = 256
)

// BitMask256 : 256bitのビットマスク
type Mask256 struct {
	// uint64の配列を4つ用意することで、256bitのビットマスクを表現する
	bits [4]uint64
}

// Equal : 2つのビットマスクが同じかどうかを比較する
func (m *Mask256) Equal(other *Mask256) bool {
	return m == other
}
