# ベースイメージを指定
FROM golang:1.16 as build

# 必要なパッケージをインストール
RUN apt-get update && apt-get -y install git

# コンテナ内の作業ディレクトリを指定
WORKDIR /app

# go.modとgo.sumをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy


# ソースコードをコピー
COPY . .

# ビルドを実行
RUN go build -o app .

# 軽量化のため、最終的なイメージを別のベースイメージにする
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/app .

# コンテナが起動した際に実行されるコマンドを指定
CMD ["./app"]
