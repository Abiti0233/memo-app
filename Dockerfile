# Dockerfile
FROM golang:1.23.4-alpine

# 必要なパッケージをインストール
RUN apk update && apk add --no-cache git

# 作業ディレクトリの設定
WORKDIR /app

# モジュールファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN go build -o main .

# ポートを開放
EXPOSE 8081

# 実行
CMD ["./main"]
