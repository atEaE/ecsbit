package ecsbit

import "fmt"

var (
	ErrRecycleSentinel = fmt.Errorf("can't recycle reserved entity")
)

// newEntityPool : Entity Poolを生成します
func newEntityPool(capacity uint32) entityPool {
	entities := make([]Entity, 1, capacity)
	entities[0] = zeroEntity // リンクリストの用に扱いたいのでsentinelを入れておく
	return entityPool{
		entities:  entities,
		next:      0,
		available: 0,
	}
}

// entityPool : Entity用のプール
// entityのGeneratorでもあるので、これは公開しない
// refs: https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	entities  []Entity // 使用中の生きているEntityと死んでいるEntityが一緒に入っている点に注意してください.
	next      EntityID // 次に利用するEntityID (RecycleされたEntityIDを再利用するために利用する)
	available uint32   // 利用可能なEntityの数
}

// Get : Entity PoolからEntityを取得します
func (p *entityPool) Get() Entity {
	// Entityが0の場合や、すべてが利用中のなどPoolから取得可能なEntityが存在しない場合は新たに作り出す必要がある
	if p.available == 0 {
		return p.new()
	}

	// 以下は、RecycleされたEntityを再利用するための処理
	// 本来はp.available < 0の場合のエラー処理も考慮するべきだが、Entityの生成、Recycleはworldが行うので、性能重視で無視する
	recycledID := p.next
	p.next, p.entities[p.next] = p.entities[p.next].ID(), switchID(p.next, p.entities[p.next])
	p.available--
	return p.entities[recycledID]
}

func (p *entityPool) new() Entity {
	e := NewEntity(EntityID(len(p.entities)))
	p.entities = append(p.entities, e)
	return e
}

// Recycle : 指定したEntityをリサイクル可能な状態にする
// この関数に渡したEntityは、その時点で無効な状態になります.
// この関数を呼び出した後に、Alive関数を呼び出すとfalseが返ります.
func (p *entityPool) Recycle(e Entity) {
	// sentinelはリサイクルしない.
	// 本来であれば、out of rangeのエラーも考慮するべきだが、Entityの生成、Recycleはworldが行うので、性能重視で無視する
	if e.ID() == 0 {
		panic(ErrRecycleSentinel)
	}

	// versionを上げることで、現在のEntityを無効な状態にする
	// 次に再利用されるEntityを記録し、合わせてEntityのID部分に次に再利用されるEntityIDを設定する
	// これによってリンクリストのようにRecycle待機しているEntityを再利用していく.
	p.entities[e.ID()] = p.entities[e.ID()].IncrementVersion()
	p.next, p.entities[e.ID()] = e.ID(), switchID(p.next, p.entities[e.ID()])
	p.available++
}

// switchID : 指定されたEntityのIDを引数に指定のEntityIDに切り替える
func switchID(next EntityID, target Entity) Entity {
	return (Entity(uint64(next) << 32)) | (target & versionMask)
}

// IsRecycleWait : 指定したEntityがリサイクル待ちかどうかを返します
func (p *entityPool) IsRecycleWait(eid EntityID) bool {
	return p.entities[eid].ID() != eid
}

// Alive : 該当のEntityが生存しているかどうかを返します
func (p *entityPool) Alive(e Entity) bool {
	// NOTE: versionが異なる場合は、リサイクル済みでEntityとしてはすでに死んでいるためfalseを返す
	return e.Version() == p.entities[e.ID()].Version()
}

// Len : Entity Poolに含まれるEntityの数を返します(index = 0は,sentinelのためカウントされない)
func (p *entityPool) Len() int {
	return len(p.entities) - 1
}

// Available : 利用可能なEntityの数を返します
func (p *entityPool) Available() int {
	return int(p.available)
}
