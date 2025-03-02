package bits

import (
	"strconv"
	"strings"
)

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

// Get : 指定したIndexのビットを取得する
func (m *Mask256) Get(index uint32) bool {
	word := index / 64
	bit := index % 64
	return m.bits[word]&(1<<bit) != 0
}

// Set : 指定したIndexのビットを設定する
// trueの場合は1、falseの場合は0
func (m *Mask256) Set(index uint32, value bool) {
	word := index / 64
	bit := index % 64
	if value {
		m.bits[word] |= (1 << bit) // bitを立てる
	} else {
		m.bits[word] &^= (1 << bit) // bitを落とす
	}
}

// IsZero : ビットマスクが0かどうかを判定する
func (m *Mask256) IsZero() bool {
	return m.bits[0] == 0 && m.bits[1] == 0 && m.bits[2] == 0 && m.bits[3] == 0
}

// Reset : ビットマスクをリセットする
func (m *Mask256) Reset() {
	m.bits = [4]uint64{0, 0, 0, 0}
}

// Bits : ビットマスクのビットを取得する
func (m *Mask256) Bits() *[4]uint64 {
	return &m.bits
}

// String : ビットマスクを2進数表記で表示する
func (m *Mask256) String() string {
	var sb strings.Builder

	// 逆順から表示していく必要があるので注意
	for i := 3; i >= 0; i-- {
		// FormatUintに２を設定することで、2進数表記に変換できる
		// ただし、そのままだと64bit未満の場合に0埋めがされないので、0埋めを行っている
		bitStr := strconv.FormatUint(m.bits[i], 2)
		bitStr = strings.Repeat("0", 64-len(bitStr)) + bitStr

		// 0000 0000 0000 ...みたいな表記にしたいので
		for j := 0; j < 64; j += 4 {
			sb.WriteString(bitStr[j : j+4])
			sb.WriteString(" ")
		}
	}
	return sb.String()
}
