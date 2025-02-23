mkdir tls-learning
cd tls-learning
go mod init github.com/yourusername/tls-learning

# メインディレクトリの作成
mkdir -p cmd/{server,client}
mkdir -p cmd/server/{self_signed,mtls,h2_h3,cert_renew}
mkdir -p pkg/{tlsutil,logger}
mkdir -p certs
mkdir -p scripts
mkdir -p docs

# 証明書を.gitignoreに追加するための設定
echo "certs/*.crt" >> .gitignore
echo "certs/*.key" >> .gitignore

go get -u golang.org/x/crypto/...
go get -u golang.org/x/net/... 