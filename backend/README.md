# docker-go-server

Go 1.24.5とDockerを使用したシンプルなHTTPサーバーです。開発時のホットリロード機能付き。

## 特徴

- シンプルなHTTP APIサーバー
- Airによるホットリロード機能
- Docker環境での開発
- ミニマルな構成
cdcd
## 目的
- テスト駆動開発

## APIエンドポイント

サーバー側で​実装する​こと。​
目的の通り、テスト駆動開発が目的のため、まずはテストケースを作成すること。

・新規​（Create）​
　・メソッド：POST
　・エンドポイント：/api/v1/tasks
　・​機能：リクエスト内容を​JSONファイルに​書き出し
・​一覧​表示​（Read）​
　・メソッド：GET
　・エンドポイント：/api/v1/tasks
　・​機能：JSONファイルの​全内容を​返却
・チェック​（Update）​
　・メソッド：PATCH
　・エンドポイント：/api/v1/tasks/{id}
　・​機能：status の​ 0/1 を​更新
・削除​（Delete）​
　・メソッド：DELETE
　・エンドポイント：/api/v1/tasks/{id}
　・​機能：指定された​idの​deletedを​trueに​更新

※TDDが​主目的、​かつリレーショナルな​データ構造不要の​ため、​DBでなく​JSONファイルで​管理。

## JSONファイルの構成
~~~
id int(20) not null PK  auto incliment
name varchar(40) not null
status int(1) not null default 0
created timestamp not null
updated timestamp not null
deleted bool not null deafult false
~~~

## 動作確認

```bash
curl http://localhost:8080
```

## 開発について

サーバーはポート8080で動作し、ホットリロードが有効です。Goファイルを変更すると自動的にサーバーが再起動されます。

## 使用技術

- Go 1.24.5
- Docker & Docker Compose
- Alpine Linux
- Air（ホットリロード）

## ファイル構成

```
docker-go-server/
├── main.go              # メインのGoファイル
├── go.mod               # Go modules設定
├── Dockerfile           # Docker設定
├── docker-compose.yml   # Docker Compose設定
├── .air.toml            # ホットリロード設定
└── README.md            # このファイル
```

## 便利なコマンド

```bash
# 従来通りの起動
docker compose up --build

# Compose Watch機能を使用（Docker Compose v2.22+）
docker compose up --watch

# バックグラウンド実行
docker compose up -d

# ログ確認
docker compose logs app

# ヘルスチェック状況確認
docker ps

# 停止
docker compose down
```
