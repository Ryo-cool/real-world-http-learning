# TLS 学習プロジェクト

このプロジェクトは、TLS（Transport Layer Security）の学習を目的としたサンプル実装です。
クライアント認証を含む TLS 通信の基本的な実装を提供します。

## プロジェクト構成

```
.
real-world-http-learning/
│── cmd/                 # 各機能のエントリポイント
│   ├── server/          # HTTPSサーバー関連
│   │   ├── main.go      # HTTPSサーバーのメインコード
│   │   ├── server.go    # サーバーロジック
│   │   ├── config.go    # TLS設定管理
│   │   ├── handlers.go  # HTTPハンドラ
│   │   ├── middleware.go # 認証やロギング
│   │   ├── self_signed/ # 自己署名証明書関連
│   │   ├── mtls/        # mTLS（双方向認証）
│   │   ├── h2_h3/       # HTTP/2・HTTP/3比較
│   │   ├── cert_renew/  # 証明書自動更新
│   ├── client/          # クライアント関連
│   │   ├── main.go      # HTTPSクライアントのメインコード
│   │   ├── client.go    # HTTPリクエスト処理
│   │   ├── verify.go    # 証明書検証ロジック
│   │   ├── pinning.go   # TLSピンニング
│   │   ├── handshake/   # ハンドシェイクログ
│── pkg/                 # 共有ライブラリ
│   ├── tlsutil/         # TLS設定や証明書管理
│   ├── logger/          # ログ出力
│── certs/               # 証明書（.gitignore推奨）
│   ├── server.crt       # サーバー証明書
│   ├── server.key       # サーバー秘密鍵
│   ├── ca.crt           # クライアント用CA証明書
│   ├── client.crt       # クライアント証明書
│   ├── client.key       # クライアント秘密鍵
│── scripts/             # 証明書生成やテスト用スクリプト
│   ├── generate_certs.sh # OpenSSLで自己署名証明書を作成
│   ├── renew_cert.sh     # Let's Encryptの証明書自動更新
│── docs/                # ドキュメント
│   ├── README.md        # プロジェクト概要
│   ├── tls_handshake.md # TLSハンドシェイクの解説
│   ├── http2_vs_http3.md # HTTP/2とHTTP/3の違い
│── go.mod               # Goモジュールファイル
│── go.sum               # 依存関係管理
```

## セットアップ

1. 証明書の生成:

```bash
./scripts/generate_certs.sh
```

2. サーバーの起動:

```bash
go run cmd/server/main.go
```

3. クライアントの実行:

```bash
go run cmd/client/main.go
```

## 実装の特徴

- TLS 1.2 以上を使用
- 相互 TLS 認証（クライアント認証）
- 強力な暗号スイートの使用
- X.509 証明書による認証

## セキュリティ注意事項

このプロジェクトは学習目的で作成されています。本番環境での使用には、
追加のセキュリティ対策が必要です。
