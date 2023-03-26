# ベースイメージを指定
FROM golang:1.19 as build

# 必要なパッケージをインストール
RUN apt-get update && apt-get -y install git

WORKDIR /app

COPY . .
RUN go mod download
RUN go mod tidy
RUN go build -o app .

# 軽量化のため、最終的なイメージを別のベースイメージにする
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/app .

# コンテナが起動した際に実行されるコマンドを指定
CMD ["./app"]
