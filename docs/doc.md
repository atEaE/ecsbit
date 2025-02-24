# Entity-Component-Sysmte(ECS)理解のためのドキュメント

## 概要

Entity-Component-System（以下、ECS）を事前知識のない状態で、作っていくのはほぼ不可能なのでそのためのドキュメントです.

## 前提

ECSは、ソフトウェアアーキテクチャのパータンの１つです.  
一般的にオブジェクト指向の設計とは異なり、 **コンポーネント指向**, **データ指向の設計** の２つをベースに構築されることが多いアーキテクチャとなっています.  
そのため、以下の２点が理解できていない場合は、まずそちらを理解してからECSの詳細を理解することをおすすめします.  

- コンポーネント指向
- データ指向
- データ指向が高速な理由

## データ指向型のアーキテクチャはなぜ高速なのか？

### キャッシュメモリ


## 参考

- [C++でEntity Component Systemを実装してみる](https://zenn.dev/kd_gamegikenblg/articles/4ca7b1ec032329)
- [ECS (Entity Component System) ざっくり概念解説](https://qiita.com/aobat/items/262293651fbbd696c171)
- [ECS (Entity Component System)について](https://zenn.dev/suuta/articles/0aa567690ec52a)
- [【Rustのまほう2】#2 Entity Component Systemの基本](https://qiita.com/hiruberuto/items/9cb625a0a8f253764bd8)
- [ECSの仕組みを理解し、使いどころを把握する](https://edom18.hateblo.jp/entry/2024/04/07/172558)