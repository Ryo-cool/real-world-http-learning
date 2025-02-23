#!/bin/bash

# 証明書を格納するディレクトリの作成
mkdir -p certs
cd certs

# SANs設定ファイルの作成
cat > san.cnf << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = JP
ST = Tokyo
L = Tokyo
O = TLS Learning
CN = localhost

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

# CA証明書の生成
openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout ca.key -out ca.crt \
  -subj "/C=JP/ST=Tokyo/L=Tokyo/O=TLS Learning/CN=Learning CA"

# サーバー証明書の生成
openssl req -newkey rsa:4096 -nodes \
  -keyout server.key -out server.csr \
  -config san.cnf

openssl x509 -req -days 365 \
  -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt \
  -extfile san.cnf -extensions v3_req

# クライアント証明書の生成
openssl req -newkey rsa:4096 -nodes \
  -keyout client.key -out client.csr \
  -subj "/C=JP/ST=Tokyo/L=Tokyo/O=TLS Learning/CN=client"

openssl x509 -req -days 365 \
  -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out client.crt

# 不要なファイルの削除
rm *.csr san.cnf

echo "証明書の生成が完了しました。" 