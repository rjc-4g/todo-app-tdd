# 開発用Dockerfile
FROM golang:1.24.5-alpine

# 作業ディレクトリを設定
WORKDIR /app

# 必要なパッケージを追加（git: go modules用, wget: ヘルスチェック用）
RUN apk add --no-cache git wget

# ホットリロード用のairをインストール
RUN go install github.com/air-verse/air@latest

# go.modとgo.sumを先にコピー（依存関係のキャッシュ最適化）
COPY go.mod go.sum* ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# ポート8080を公開
EXPOSE 8080

# airを使ってホットリロードで起動
CMD ["air", "-c", ".air.toml"]
