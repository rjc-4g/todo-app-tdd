# docker-go-server

Go 1.24.5とDockerを使用したシンプルなHTTPサーバーです。開発時のホットリロード機能付き。

## 特徴

- シンプルなHTTP APIサーバー
- "Hello World" JSONレスポンスを返す
- Airによるホットリロード機能
- Docker環境での開発
- ミニマルな構成

## 必要な環境変数

`.env`ファイルを作成して以下の環境変数を設定してください：

```bash
# AI API Configuration
OPENAI_API_KEY=your_openai_api_key_here
```

## OpenAI APIキーの取得

1. [OpenAI Platform](https://platform.openai.com/)にアクセス
2. アカウントを作成またはログイン
3. 設定（歯車） → API Keysセクションで新しいAPIキーを生成
4. 生成されたキーを.envの`OPENAI_API_KEY`に設定

## クイックスタート

```bash
# リポジトリをクローン
git clone https://github.com/your-username/docker-go-server.git
cd docker-go-server

# 環境変数を設定
cp .env.example .env
# .envファイルを編集してAPIキーを設定

# 依存関係をインストール
go mod tidy

# 開発サーバー起動
docker compose up --build
```

## APIエンドポイント

- `GET /` - Hello World JSONを返す

レスポンス例：
```json
{
  "message": "Hello World!"
}
```

### `GET /api/themes/random` - ランダムテーマ生成
技術関連のランダムなテーマを3つ返します。

レスポンス例：
```json
{
  "themes": [
    "最近学んだ新しいプログラミング言語の特徴",
    "開発環境の効率化テクニック",
    "コードレビューのベストプラクティス"
  ],
  "count": 3
}
```

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
