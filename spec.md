# ecsbit 仕様書

## 1. 概要

ecsbitは、Go言語のゲームエンジン **ebitengine** 向けのEntity Component System（ECS）ライブラリである。
高速なデータ指向設計（Data-Oriented Design）を採用し、Archetype方式によるキャッシュフレンドリーなデータ格納を実現する。

### 1.1 設計原則

- **データ指向設計**: CPUキャッシュの効率的な利用を最優先し、Structure of Arrays（SoA）形式でコンポーネントデータを格納する。
- **ebitengine互換**: ebitengineの`Update`/`Draw`ライフサイクルに自然に統合できるAPIを提供する。ebitengineを直接の依存に含める。
- **型安全**: Go 1.24のGenericsを活用し、コンポーネントやクエリの操作を型安全に行える。
- **シンプルさ**: シングルスレッド前提の設計とし、過度な抽象化を避ける。

### 1.2 モジュール情報

```
module github.com/atEaE/ecsbit
go 1.24.0
```

### 1.3 依存関係

- `github.com/hajimehoshi/ebiten/v2` — Draw系Systemの引数に`*ebiten.Image`を使用するため。

---

## 2. Entity

### 2.1 概要

Entityはゲーム内のオブジェクトを表す一意な識別子である。
Entity自体はデータを持たず、Componentの組み合わせによって振る舞いが決まる。

### 2.2 データ構造

Entityは64bit符号なし整数として表現する。

```
Entity (uint64)
┌─────────────────────────────┬─────────────────────────────┐
│       EntityID (上位32bit)    │       Version (下位32bit)     │
└─────────────────────────────┴─────────────────────────────┘
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| EntityID | `uint32` | Entityの一意な識別子。Entity Poolによって管理される。 |
| Version | `uint32` | 世代番号。Entityがリサイクルされるたびにインクリメントされる。同一IDの旧参照を無効化するために使用する。 |

### 2.3 型定義

```go
type Entity uint64
type EntityID uint32
```

### 2.4 操作

| メソッド | シグネチャ | 説明 |
|---------|-----------|------|
| NewEntity | `NewEntity(id EntityID) Entity` | 指定IDでVersion=0のEntityを生成する |
| ID | `(e Entity) ID() EntityID` | 上位32bitからEntityIDを取得する |
| Version | `(e Entity) Version() uint32` | 下位32bitからVersionを取得する |
| IncrementVersion | `(e Entity) IncrementVersion() Entity` | Versionをインクリメントした新しいEntityを返す。オーバーフロー時は0に戻る |
| String | `(e Entity) String() string` | デバッグ用の文字列表現を返す |

### 2.5 Entity Pool

EntityのIDの割り当てとリサイクルを管理するプール機構。

#### データ構造

```go
type entityPool struct {
    entities  []Entity   // Entity格納用スライス
    next      EntityID   // リサイクルリストの先頭を指すID
    available uint32     // リサイクル可能なEntityの数
}
```

#### リサイクル機構

- リサイクル済みのEntityIDはリンクリスト（LIFO/スタック）形式で管理する。
- `next`フィールドがリサイクルリストの先頭を指す。
- Entity取得時（`Get`）、リサイクル可能なEntityがあればそれを再利用する。なければ新規IDを割り当てる。
- Entity解放時（`Recycle`）、Versionをインクリメントし、リサイクルリストに追加する。

#### 制約

- EntityID `0` はセンチネル（番兵）として予約されており、リサイクル不可。リサイクルを試みた場合はpanicする。

#### 操作

| メソッド | シグネチャ | 説明 |
|---------|-----------|------|
| Get | `(p *entityPool) Get() Entity` | 新規または再利用EntityIDからEntityを取得する |
| Recycle | `(p *entityPool) Recycle(e Entity)` | EntityをリサイクルリストへFreightする |
| Alive | `(p *entityPool) Alive(e Entity) bool` | Entityがまだ生存しているかを返す |
| Used | `(p *entityPool) Used() int` | 使用中のEntity数を返す |
| Available | `(p *entityPool) Available() int` | リサイクル可能なEntity数を返す |
| Total | `(p *entityPool) Total() int` | 生成されたEntityの総数を返す |
| Cap | `(p *entityPool) Cap() int` | プールの容量を返す |

### 2.6 EntityIndex

EntityがArchetype内のどの位置に存在するかを管理する索引。

```go
type EntityIndex struct {
    index     uint32      // Archetype内のインデックス
    archetype *archetype  // 所属するArchetypeへの参照（削除済みの場合はnil）
}
```

---

## 3. Component

### 3.1 概要

Componentは純粋なデータの入れ物であり、振る舞い（ロジック）を持たない。
各ComponentはGoの構造体として定義され、Worldに型として登録される。

### 3.2 型定義

```go
type ComponentID uint32

type component struct {
    name string       // コンポーネント名
    typ  reflect.Type // 型情報
}
```

### 3.3 コンポーネント登録

Genericsを使用してコンポーネントを生成し、Worldに登録する。

```go
// コンポーネントの生成
posComponent := ecsbit.NewComponent[Position]()

// Worldへの登録（ComponentIDが返る）
posID := world.RegisterComponent(posComponent)
```

### 3.4 Component Storage

コンポーネントの型とIDの対応を管理するストレージ。

```go
type componentStorage struct {
    Components map[component]ComponentID  // component → ID の逆引き
    Types      []component                // ComponentID → component の正引き
    IDs        []ComponentID              // 登録順のID一覧
    maxSize    int                        // 最大登録数
}
```

#### 制約

- 1つのWorldに登録できるコンポーネント種類の最大数は **256** である。これは256bitマスクの制約に由来する。

### 3.5 Tag

Tagはデータを持たない特殊なComponentである。Entityのグルーピングやフラグとして使用する。

```go
type Tag string

func NewTag(t Tag) component
```

#### 使用例

```go
playerTag := ecsbit.NewTag("Player")
playerTagID := world.RegisterComponent(playerTag)

// Entityにタグを付ける
entity := world.CreateEntity(posID, velID, playerTagID)
```

---

## 4. Archetype

### 4.1 概要

Archetypeは、同一のComponent構成を持つEntityの集合である。
データ指向設計の核となる仕組みであり、同じ種類のデータを連続したメモリ領域に配置することでキャッシュ効率を最大化する。

### 4.2 レイアウトとマスク

ArchetypeのComponent構成は256bitのビットマスク（`Mask256`）で表現する。
各bitがComponentIDに対応し、そのComponentを含むかどうかを示す。

```go
type Mask256 struct {
    bits [4]uint64  // 4 × 64bit = 256bit
}
```

#### Mask256の操作

| メソッド | シグネチャ | 説明 |
|---------|-----------|------|
| Get | `(m *Mask256) Get(index uint32) bool` | 指定bitが立っているか確認 |
| Set | `(m *Mask256) Set(index uint32, value bool)` | 指定bitを設定 |
| Equal | `(m *Mask256) Equal(other *Mask256) bool` | マスクが等しいか比較 |
| IsZero | `(m *Mask256) IsZero() bool` | 全bitが0か判定 |
| Reset | `(m *Mask256) Reset()` | 全bitをクリア |

### 4.3 データ構造

```go
type archetype struct {
    id         primitive.ArchetypeID  // Archetype識別子
    layoutMask bits.Mask256           // Component構成を表すビットマスク
    *archetypeData
}

type archetypeData struct {
    entities   []Entity               // このArchetypeに属するEntity一覧
    columns    []componentColumn      // SoA形式のコンポーネントデータ列
    layoutMask bits.Mask256           // Component構成マスク
}
```

### 4.4 SoA（Structure of Arrays）形式のデータ格納

Archetype内のコンポーネントデータは、SoA形式で格納する。
各コンポーネント型ごとに独立した配列（カラム）を持ち、Entity間で同種のデータが連続して配置される。

```
Archetype (Position, Velocity を持つEntity群)
┌──────────────────────────────────────────┐
│ entities:   [Entity0, Entity1, Entity2]  │
│ column[0]:  [Pos0,    Pos1,    Pos2   ]  │  ← Position配列
│ column[1]:  [Vel0,    Vel1,    Vel2   ]  │  ← Velocity配列
└──────────────────────────────────────────┘
```

#### componentColumn

各カラムはComponentIDに対応するデータ配列を保持する。

```go
type componentColumn struct {
    componentID ComponentID
    data        reflect.Value  // reflect.SliceOf(componentType) で生成されたスライス
}
```

**注**: Generics + `unsafe`パッケージの利用、または`reflect`によるスライス操作によって型安全なアクセスを提供する。パフォーマンスクリティカルなパスでは`unsafe.Pointer`による直接アクセスも検討する。

### 4.5 操作

| メソッド | 説明 |
|---------|------|
| `Add(e Entity) uint32` | Entityを追加し、追加されたインデックスを返す |
| `Remove(index uint32) bool` | 指定インデックスのEntityを削除する。最後のEntityとスワップすることでO(1)で実行する。スワップが発生した場合はtrueを返す |
| `Count() int` | Archetype内のEntity数を返す |
| `GetEntity(index uint32) Entity` | 指定インデックスのEntityを取得する |

### 4.6 特殊なArchetype

- **インデックス0**: コンポーネントを持たないEntity用のArchetype（No-Layout Archetype）。Entity生成時にComponentを指定しなかった場合に使用される。

---

## 5. World

### 5.1 概要

WorldはECSの中心的な管理構造体であり、全てのEntity、Component、Archetype、System、Resourceを統括する。
アプリケーション内で1つ以上のWorldを持つことができる。

### 5.2 データ構造

```go
type World struct {
    // コンポーネント管理
    componentStorage  componentStorage

    // Archetype管理
    archetypeData     []archetypeData
    archetypeLayouts  map[bits.Mask256]*archetype
    archetypes        []archetype

    // Entity管理
    entityIndices     []EntityIndex
    entityPool        entityPool

    // System管理
    updateSystems     []systemEntry       // Update系System（登録順）
    drawSystems       []drawSystemEntry   // Draw系System（登録順）

    // Resource管理
    resources         map[reflect.Type]any

    // Event管理
    eventBus          eventBus

    // Callback
    onCreateCallbacks []func(*World, Entity)
    onRemoveCallbacks []func(*World, Entity)

    // 設定
    config            internalconfig.WorldConfig
}
```

### 5.3 生成と設定

Functional Optionsパターンを使用してWorldを生成する。

```go
world := ecsbit.NewWorld(
    config.WithArchetypeDefaultCapacity(512),
    config.WithEntityPoolDefaultCapacity(2048),
)
```

#### 設定項目

| オプション | デフォルト値 | 説明 |
|-----------|------------|------|
| ArchetypeDefaultCapacity | 256 | Archetype内のEntity配列の初期容量 |
| EntityPoolDefaultCapacity | 1024 | Entity Poolの初期容量 |
| OnCreateCallbacksDefaultCapacity | 256 | OnCreateコールバック配列の初期容量 |
| OnRemoveCallbacksDefaultCapacity | 256 | OnRemoveコールバック配列の初期容量 |

### 5.4 Entity操作

| メソッド | シグネチャ | 説明 |
|---------|-----------|------|
| CreateEntity | `(w *World) CreateEntity(components ...ComponentID) Entity` | 指定したComponent構成でEntityを生成する。適切なArchetypeに配置される |
| RemoveEntity | `(w *World) RemoveEntity(e Entity)` | Entityを削除する。Archetype上のデータはスワップ削除される |
| Alive | `(w *World) Alive(e Entity) bool` | Entityが生存中か確認する |

#### Entity生成フロー

1. `CreateEntity(componentIDs...)` が呼ばれる
2. ComponentIDからビットマスクを生成する
3. そのマスクに対応するArchetypeを検索する。なければ新規作成する
4. Entity Poolから新規EntityIDを取得（またはリサイクル）する
5. Archetypeに追加し、EntityIndexを記録する
6. 各コンポーネントカラムにゼロ値のデータを追加する
7. 登録されたOnCreateコールバックを全て実行する
8. 生成されたEntityを返す

#### Entity削除フロー

1. `RemoveEntity(entity)` が呼ばれる
2. EntityのAliveチェックを行う（死亡済みの場合はエラー）
3. EntityIndexからArchetypeとインデックスを取得する
4. 登録されたOnRemoveコールバックを全て実行する
5. Archetype上でスワップ削除を実行する（Entity配列 + 全コンポーネントカラム）
6. スワップが発生した場合、スワップ先EntityのEntityIndexを更新する
7. 削除されたEntityのEntityIndex.archetypeをnilに設定する
8. Entity PoolでEntityをリサイクルする（Versionインクリメント）

### 5.5 Component動的追加・削除（Entity Migration）

既存のEntityに対してComponentを動的に追加・削除できる。
この操作はEntityを元のArchetypeから新しいArchetypeへ移動（Migration）させる。

```go
// Componentの追加
ecsbit.AddComponent[Velocity](world, entity, Velocity{X: 1.0, Y: 0.5})

// Componentの削除
ecsbit.RemoveComponent[Velocity](world, entity)
```

#### Migration フロー（Component追加の場合）

1. EntityのEntityIndexから現在のArchetypeを取得する
2. 現在のArchetypeのマスクに新しいComponentIDのbitを追加した新マスクを計算する
3. 新マスクに対応するArchetypeを検索する。なければ新規作成する
4. 元のArchetypeから全コンポーネントデータをコピーし、新Archetypeに追加する
5. 新しいComponentのデータを新Archetypeのカラムに設定する
6. 元のArchetypeからEntityをスワップ削除する
7. EntityIndexを新しいArchetypeの位置で更新する

### 5.6 Callback

Entity生成・削除時のフック機構。

```go
world.PushOnCreateCallback(func(w *World, e Entity) {
    // Entityが生成された直後に呼ばれる
})

world.PushOnRemoveCallback(func(w *World, e Entity) {
    // Entityが削除される直前に呼ばれる
})
```

### 5.7 Update / Draw

WorldはebitengineのGame loopに統合するための`Update`メソッドと`Draw`メソッドを公開する。
ただし、`ebiten.Game`インターフェースは実装しない。ユーザーが自身のGame構造体内でWorldのメソッドを呼び出す形とする。

```go
// Update は登録された全てのUpdate系Systemを登録順に実行する
func (w *World) Update()

// Draw は登録された全てのDraw系Systemを登録順に実行する
func (w *World) Draw(screen *ebiten.Image)
```

#### 使用例

```go
type MyGame struct {
    world *ecsbit.World
}

func (g *MyGame) Update() error {
    g.world.Update()
    return nil
}

func (g *MyGame) Draw(screen *ebiten.Image) {
    g.world.Draw(screen)
}

func (g *MyGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 320, 240
}

func main() {
    world := ecsbit.NewWorld()
    // ... System登録、Entity生成等 ...
    ebiten.RunGame(&MyGame{world: world})
}
```

### 5.8 Statistics

Worldの統計情報をJSON形式で取得できる。

```go
type stats.World struct {
    Entities stats.Entities `json:"entities"`
}

type stats.Entities struct {
    Used     int `json:"used"`
    Total    int `json:"total"`
    Recycled int `json:"recycled"`
    Capacity int `json:"capacity"`
}
```

---

## 6. System

### 6.1 概要

SystemはECSのロジックを担う部分である。各Systemは特定のComponent構成を持つEntity群に対して処理を行う。
Systemは **Update系** と **Draw系** の2種類に分類され、それぞれWorldの`Update()`と`Draw()`で実行される。

### 6.2 System インターフェース

SystemはインターフェースによるStructベースと、関数ベースの両方で登録できる。

#### インターフェース定義

```go
// Update系System
type System interface {
    Update(w *World)
}

// Draw系System
type DrawSystem interface {
    Draw(w *World, screen *ebiten.Image)
}
```

#### 関数型

```go
// Update系Systemの関数型
type UpdateFunc func(w *World)

// Draw系Systemの関数型
type DrawFunc func(w *World, screen *ebiten.Image)
```

### 6.3 System登録

```go
// インターフェースによる登録
world.AddSystem(&MovementSystem{})
world.AddDrawSystem(&RenderSystem{})

// 関数による登録
world.AddUpdateFunc(func(w *World) {
    // 更新処理
})

world.AddDrawFunc(func(w *World, screen *ebiten.Image) {
    // 描画処理
})
```

### 6.4 実行順序

- Update系System・Draw系Systemともに **登録順** に実行される。
- Update系SystemはDraw系Systemよりも先に実行される（`World.Update()`と`World.Draw()`が別メソッドであるため、呼び出し順はユーザーが制御する）。
- System内からEntityの生成・削除を行った場合、その変更は即座に反映される（遅延なし）。

### 6.5 使用例

```go
// 構造体ベースのSystem
type MovementSystem struct{}

func (s *MovementSystem) Update(w *ecsbit.World) {
    q := ecsbit.Query2[Position, Velocity](w)
    q.Each(func(e ecsbit.Entity, pos *Position, vel *Velocity) {
        pos.X += vel.X
        pos.Y += vel.Y
    })
}

// 関数ベースのSystem
world.AddUpdateFunc(func(w *ecsbit.World) {
    q := ecsbit.Query1[Health](w)
    q.Each(func(e ecsbit.Entity, hp *Health) {
        if hp.Value <= 0 {
            w.RemoveEntity(e)
        }
    })
})
```

---

## 7. Query

### 7.1 概要

Queryは特定のComponent構成を持つEntity群を検索・イテレートする仕組みである。
Go GenericsのType Parameterを使用し、型安全にコンポーネントデータへアクセスできる。

### 7.2 Query関数

コンポーネント数に応じたトップレベル関数を提供する（最大4コンポーネント）。

```go
func Query1[A any](w *World, filters ...Filter) *QueryResult1[A]
func Query2[A, B any](w *World, filters ...Filter) *QueryResult2[A, B]
func Query3[A, B, C any](w *World, filters ...Filter) *QueryResult3[A, B, C]
func Query4[A, B, C, D any](w *World, filters ...Filter) *QueryResult4[A, B, C, D]
```

### 7.3 QueryResult

各QueryResult型はイテレーションメソッドを持つ。

#### Query1の場合

```go
type QueryResult1[A any] struct { /* ... */ }

// 全マッチEntityに対して関数を実行する
func (q *QueryResult1[A]) Each(fn func(e Entity, a *A))

// マッチするEntity数を返す
func (q *QueryResult1[A]) Count() int
```

#### Query2の場合

```go
type QueryResult2[A, B any] struct { /* ... */ }

func (q *QueryResult2[A, B]) Each(fn func(e Entity, a *A, b *B))
func (q *QueryResult2[A, B]) Count() int
```

#### Query3, Query4も同様のパターン

```go
type QueryResult3[A, B, C any] struct { /* ... */ }
func (q *QueryResult3[A, B, C]) Each(fn func(e Entity, a *A, b *B, c *C))
func (q *QueryResult3[A, B, C]) Count() int

type QueryResult4[A, B, C, D any] struct { /* ... */ }
func (q *QueryResult4[A, B, C, D]) Each(fn func(e Entity, a *A, b *B, c *C, d *D))
func (q *QueryResult4[A, B, C, D]) Count() int
```

### 7.4 Filter

Queryに追加条件を指定できるFilter機構を提供する。

```go
type Filter interface {
    apply(mask *bits.Mask256, storage *componentStorage) filterResult
}
```

#### 提供するFilter

| Filter | 説明 |
|--------|------|
| `With[T any]()` | 指定したComponentを **追加で** 持つEntityのみにマッチする。Queryの型パラメータとは別に、追加条件として使用する |
| `Without[T any]()` | 指定したComponentを **持たない** Entityのみにマッチする |

#### 使用例

```go
// Position と Velocity を持ち、かつ Health も持つが、Dead タグは持たないEntity
q := ecsbit.Query2[Position, Velocity](world,
    ecsbit.With[Health](),
    ecsbit.Without[Dead](),
)
q.Each(func(e ecsbit.Entity, pos *Position, vel *Velocity) {
    pos.X += vel.X
    pos.Y += vel.Y
})
```

### 7.5 Queryの内部動作

1. Type Parameterの各コンポーネント型からComponentIDを解決する
2. Filter条件も含めた必須/除外マスクを生成する
3. 全Archetypeを走査し、マスク条件にマッチするArchetypeを収集する
4. `Each`呼び出し時、マッチした各ArchetypeのSoAカラムに対して直接イテレーションを行う

#### パフォーマンス考慮事項

- Queryの生成（Archetype走査）はSystem実行のたびに行われる。頻繁に使用するQueryはキャッシュする最適化を将来的に検討する。
- `Each`のイテレーションはSoAカラムに対する連続メモリアクセスとなるため、キャッシュ効率が高い。
- コンポーネントデータはポインタで渡されるため、コールバック内で直接書き換えが可能である。

---

## 8. Resource

### 8.1 概要

Resourceは、World全体で1つだけ存在するシングルトンデータである。
Entityに紐づかないグローバルな状態（スクリーンサイズ、入力状態、設定値など）を管理する。

### 8.2 API

Genericsを使用した型安全なトップレベル関数で操作する。

```go
// Resourceの追加（型ごとに1つのみ。同じ型を二重登録した場合は上書き）
func AddResource[T any](w *World, value T)

// Resourceの取得（未登録の場合はpanicする）
func GetResource[T any](w *World) *T

// Resourceの存在確認
func HasResource[T any](w *World) bool

// Resourceの削除
func RemoveResource[T any](w *World)
```

### 8.3 内部実装

ResourceはWorldの`resources map[reflect.Type]any`に格納される。型をキーとして1つの値のみを保持する。

### 8.4 使用例

```go
// Resource型の定義
type GameConfig struct {
    ScreenWidth  int
    ScreenHeight int
    Debug        bool
}

type InputState struct {
    MouseX, MouseY int
    LeftClick      bool
}

// Resourceの登録
ecsbit.AddResource[GameConfig](world, GameConfig{
    ScreenWidth:  320,
    ScreenHeight: 240,
    Debug:        false,
})

// System内でのResourceの利用
func (s *RenderSystem) Draw(w *ecsbit.World, screen *ebiten.Image) {
    cfg := ecsbit.GetResource[GameConfig](w)
    // cfg.ScreenWidth, cfg.ScreenHeight を利用
}
```

---

## 9. Event

### 9.1 概要

EventはSystem間のメッセージングを実現する仕組みである。
あるSystemから発火されたイベントを、別のSystemが受信・処理できる。

### 9.2 設計方針

- イベントはフレーム単位でバッファリングされる。
- `World.Update()`の開始時に前フレームのイベントがクリアされる。
- イベントの型はGenericsで指定し、型安全に送受信する。

### 9.3 API

```go
// イベントの送信（同一フレーム内で複数回呼び出し可能）
func EmitEvent[T any](w *World, event T)

// イベントの読み取り（現在フレームに蓄積されたイベントを全件返す）
func ReadEvents[T any](w *World) []T

// イベントの存在確認
func HasEvents[T any](w *World) bool
```

### 9.4 Event Bus 内部構造

```go
type eventBus struct {
    events map[reflect.Type]any  // reflect.Type → []T のスライスを保持
}
```

### 9.5 ライフサイクル

```
フレーム N:
  World.Update() 開始
    → 前フレームのイベントをクリア
    → System A 実行: EmitEvent[CollisionEvent](w, event)
    → System B 実行: events := ReadEvents[CollisionEvent](w) ← System Aのイベントを受信
  World.Update() 終了

フレーム N+1:
  World.Update() 開始
    → フレーム N のイベントをクリア
    → ...
```

### 9.6 使用例

```go
// イベント型の定義
type CollisionEvent struct {
    EntityA ecsbit.Entity
    EntityB ecsbit.Entity
}

// 衝突判定System
type CollisionSystem struct{}

func (s *CollisionSystem) Update(w *ecsbit.World) {
    q := ecsbit.Query1[Collider](w)
    // 衝突判定ロジック...
    ecsbit.EmitEvent[CollisionEvent](w, CollisionEvent{
        EntityA: entityA,
        EntityB: entityB,
    })
}

// ダメージ処理System（CollisionSystemより後に登録されている前提）
type DamageSystem struct{}

func (s *DamageSystem) Update(w *ecsbit.World) {
    events := ecsbit.ReadEvents[CollisionEvent](w)
    for _, event := range events {
        // ダメージ計算処理...
    }
}
```

### 9.7 注意事項

- イベントの受信はSystem登録順に依存する。イベントを発火するSystemよりも後に登録されたSystemのみがそのフレーム内でイベントを受信できる。
- イベントは値コピーで格納される。大きなデータ構造をイベントとして送る場合はポインタ型を使用することを推奨する。

---

## 10. コンポーネントデータアクセス

### 10.1 概要

Query以外にも、特定のEntityのコンポーネントデータに直接アクセスするための関数を提供する。

### 10.2 API

```go
// 指定EntityのComponentデータを取得する（ポインタで返す）
func GetComponent[T any](w *World, e Entity) *T

// 指定EntityのComponentデータを設定する
func SetComponent[T any](w *World, e Entity, value T)

// 指定EntityがComponentを持つか確認する
func HasComponent[T any](w *World, e Entity) bool

// 指定EntityにComponentを追加する（Entity Migration発生）
func AddComponent[T any](w *World, e Entity, value T)

// 指定EntityからComponentを削除する（Entity Migration発生）
func RemoveComponent[T any](w *World, e Entity)
```

### 10.3 注意事項

- `GetComponent` / `SetComponent` はMigrationを発生させない。既にEntityが持っているComponentに対してのみ使用できる。
- `AddComponent` / `RemoveComponent` はArchetype間のMigrationを発生させるため、`GetComponent`/`SetComponent`よりコストが高い。頻繁な呼び出しは避けるべきである。
- 死亡済みEntityに対する操作はpanicする。

---

## 11. エラーハンドリング

### 11.1 エラー定義

| エラー | 説明 |
|--------|------|
| `ErrDeadEntityOperation` | 削除済みEntityに対する操作を行った場合 |
| `ErrDuplicateComponent` | 同一Entity作成時に重複するComponentを指定した場合 |
| `ErrRecycleSentinel` | センチネルEntity（ID=0）のリサイクルを試みた場合 |

### 11.2 方針

- Entity関連の操作で不正な状態が検出された場合はpanicする。これはプログラミングエラーであり、リカバリすべきでないためである。
- Componentの型解決に失敗した場合（未登録の型が指定された場合）もpanicする。

---

## 12. パッケージ構成

```
ecsbit/
├── ecsbit (root)          # メインパッケージ。公開API全体を提供する
│   ├── entity.go          # Entity, EntityID, EntityIndex
│   ├── component.go       # Component, ComponentID, componentStorage
│   ├── archetype.go       # archetype, archetypeData, componentColumn
│   ├── world.go           # World
│   ├── system.go          # System, DrawSystem インターフェース、登録関数
│   ├── query.go           # Query1〜Query4, QueryResult, Filter
│   ├── resource.go        # AddResource, GetResource, HasResource, RemoveResource
│   ├── event.go           # EmitEvent, ReadEvents, HasEvents, eventBus
│   ├── tag.go             # Tag
│   ├── entity_pool.go     # entityPool
│   ├── error.go           # エラー定義
│   └── migration.go       # AddComponent, RemoveComponent（Entity Migration）
│
├── config/                # World設定オプション
│   └── conf.go
│
├── stats/                 # 統計情報
│   └── stats.go
│
└── internal/              # 内部パッケージ（外部非公開）
    ├── bits/
    │   └── mask_256.go    # 256bitマスク
    ├── config/
    │   └── config.go      # 内部設定構造体
    └── primitive/
        └── archetype.go   # ArchetypeID型
```

---

## 13. ebitengineとの統合パターン

### 13.1 基本パターン

```go
package main

import (
    "github.com/atEaE/ecsbit"
    "github.com/atEaE/ecsbit/config"
    "github.com/hajimehoshi/ebiten/v2"
)

// コンポーネント定義
type Position struct{ X, Y float64 }
type Velocity struct{ X, Y float64 }
type Sprite struct{ Image *ebiten.Image }

// System定義
type MovementSystem struct{}
func (s *MovementSystem) Update(w *ecsbit.World) {
    q := ecsbit.Query2[Position, Velocity](w)
    q.Each(func(e ecsbit.Entity, pos *Position, vel *Velocity) {
        pos.X += vel.X
        pos.Y += vel.Y
    })
}

type RenderSystem struct{}
func (s *RenderSystem) Draw(w *ecsbit.World, screen *ebiten.Image) {
    q := ecsbit.Query2[Position, Sprite](w)
    q.Each(func(e ecsbit.Entity, pos *Position, spr *Sprite) {
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(pos.X, pos.Y)
        screen.DrawImage(spr.Image, op)
    })
}

// Game構造体
type Game struct {
    world *ecsbit.World
}

func (g *Game) Update() error {
    g.world.Update()
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.world.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 320, 240
}

func main() {
    world := ecsbit.NewWorld(
        config.WithEntityPoolDefaultCapacity(4096),
    )

    // コンポーネント登録
    posID := world.RegisterComponent(ecsbit.NewComponent[Position]())
    velID := world.RegisterComponent(ecsbit.NewComponent[Velocity]())
    sprID := world.RegisterComponent(ecsbit.NewComponent[Sprite]())

    // System登録（実行順序 = 登録順）
    world.AddSystem(&MovementSystem{})
    world.AddDrawSystem(&RenderSystem{})

    // Entity生成
    entity := world.CreateEntity(posID, velID, sprID)
    ecsbit.SetComponent[Position](world, entity, Position{X: 100, Y: 100})
    ecsbit.SetComponent[Velocity](world, entity, Velocity{X: 1.0, Y: 0.5})

    // ゲーム起動
    ebiten.RunGame(&Game{world: world})
}
```

---

## 14. 制約と制限

| 項目 | 制約 |
|------|------|
| 最大コンポーネント種類数 | 256（256bitマスクに由来） |
| 最大Entity数 | 2^32 - 1（EntityIDがuint32）。ただしメモリ量に実質的に制限される |
| Query型パラメータ数 | 最大4 |
| スレッド安全性 | なし。シングルスレッド前提。ebitengineのUpdate/Drawと同一goroutineで使用すること |
| EntityID 0 | センチネルとして予約。通常のEntityには使用されない |
| Versionオーバーフロー | uint32の最大値を超えると0に戻る。極端に長時間のEntityリサイクルでは同一ID+Versionの衝突がありうるが、実用上問題にならない |

---

## 15. 用語集

| 用語 | 説明 |
|------|------|
| Entity | ゲーム内オブジェクトの一意な識別子。データを持たない |
| Component | 純粋なデータの入れ物。Goの構造体で表現される |
| System | ロジックの実行単位。ComponentのQueryを通じてEntityを処理する |
| Archetype | 同一Component構成を持つEntityの集合。SoA形式でデータを格納する |
| World | ECSの中央管理構造体。Entity, Component, System, Resource, Eventを統括する |
| Resource | World全体で1つだけ存在するシングルトンデータ |
| Event | System間のフレーム単位メッセージング機構 |
| Tag | データを持たない特殊なComponent。フラグやグルーピングに使用する |
| Migration | Entityが異なるArchetype間を移動すること。Componentの追加/削除時に発生する |
| SoA | Structure of Arrays。同種のデータを連続配列に格納する方式 |
| Filter | Queryに追加条件を指定する機構（With / Without） |
