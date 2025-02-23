#!/bin/bash

# 証明書を格納するディレクトリの作成
mkdir -p certs
cd certs

# CA証明書の生成
openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout ca.key -out ca.crt \
  -subj "/C=JP/ST=Tokyo/L=Tokyo/O=TLS Learning/CN=Learning CA"

# サーバー証明書の生成
openssl req -newkey rsa:4096 -nodes \
  -keyout server.key -out server.csr \
  -subj "/C=JP/ST=Tokyo/L=Tokyo/O=TLS Learning/CN=localhost"

openssl x509 -req -days 365 \
  -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt

# クライアント証明書の生成
openssl req -newkey rsa:4096 -nodes \
  -keyout client.key -out client.csr \
  -subj "/C=JP/ST=Tokyo/L=Tokyo/O=TLS Learning/CN=client"

openssl x509 -req -days 365 \
  -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out client.crt

# CSRファイルの削除
rm *.csr

echo "証明書の生成が完了しました。" 