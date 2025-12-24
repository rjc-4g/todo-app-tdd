FROM golang:1.23-alpine

WORKDIR /app

# Air のバージョンを固定してインストール
RUN go install github.com/air-verse/air@v1.52.3

COPY go.mod ./
# go.sum があれば追加
# COPY go.mod go.sum ./ 
RUN go mod download

CMD ["air", "-c", ".air.toml"]