package ecsbit

import "fmt"

// golangでは、0xを先頭につけることで16進数を表すことができる。
// この場合、1桁が4bitになるため、0(0000), F(1111)というbitで表現できるため、IDのみを取り出すための32bitマスクとしてbitを表すと...
// 1111 1111 1111 1111 1111 1111 1111 1111 0000 0000 0000 0000 0000 0000 0000 0000
// となる。
const (
	// idMask : EntityIDを取り出すためのマスク
	idMask      Entity = 0xFFFFFFFF00000000
	versionMask Entity = 0xFFFFFFFF
)

const (
	zeroEntity Entity = 0x00000000FFFFFFFF
)

// Entity : uint64で表されるEntityの識別子
// * 最初の32bitがEntityIDを表す
// * 残りの32bitがversionを表す
type Entity uint64

// EntityID : EntityのID部分(先頭32bit)
type EntityID uint32

// NewEntity : Entityを作成します
func NewEntity(id EntityID) Entity {
	// uint64に返還後、左ビットシフトのオペレーターを使って、IDを32bit左にシフトする
	// さらに idMaskのbit maskと AND演算子を使ってbit単位のAND演算を行い、Version部分を0クリアして、ID部分を取り出している
	return Entity(uint64(id)<<32) & idMask
}

// ID : EntityのID部分を取得します
func (e Entity) ID() EntityID {
	// EntityIDを左ビットシフトしてEntityを作っているので、逆に右ビットシフトすることでIDを取り出している
	return EntityID(e >> 32)
}

// Version : EntityのVersion部分を取得します
func (e Entity) Version() uint32 {
	// versionMaskのbit maskと AND演算子を使ってbit単位のAND演算を行い、ID部分を0クリアして、version部分を取り出している
	return uint32(e & versionMask)
}

// IncrementVersion : EntityのVersionをインクリメントします
func (e Entity) IncrementVersion() Entity {
	// idMask済みのbit列と、versionMask済みのbit列をOR演算することで、ID部分は変更せずにVersion部分をインクリメントしている
	// オーバーフローが発生した場合は、0に戻る
	return e&idMask | (e+1)&versionMask
}

// String : Entityを文字列に変換します
func (e Entity) String() string {
	return fmt.Sprintf("Entity: {id: %d, version: %d}", e.ID(), e.Version())
}
