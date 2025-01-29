# ベースイメージの指定
FROM golang:1.23.4-alpine

# 必要なパッケージのインストール
RUN apk update && apk add --no-cache git bash

# 作業ディレクトリの作成
WORKDIR /app

# wait-for-it スクリプトのコピー
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Goモジュールのキャッシュを利用するために依存ファイルを先にコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド
RUN go build -o main .

# ポートの指定
EXPOSE 8081

# アプリケーションの実行（wait-for-it を使用して db:5432 が利用可能になるまで待機）
CMD ["/wait-for-it.sh", "db:5432", "--", "./main"]
