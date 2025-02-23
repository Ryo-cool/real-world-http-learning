# TLS 学習プロジェクト

このプロジェクトは TLS、HTTPS、証明書管理について学ぶための実装例を提供します。

## 機能

- 基本的な HTTPS サーバー/クライアント
- 自己署名証明書の生成と使用
- 相互 TLS 認証（mTLS）
- HTTP/2 と HTTP/3 の実装比較
- 証明書の自動更新
- TLS ピンニング

## セットアップ

1. リポジトリをクローン
2. 証明書の生成: `./scripts/generate_certs.sh`
3. サーバーの起動: `go run cmd/server/main.go`
4. クライアントの実行: `go run cmd/client/main.go`
