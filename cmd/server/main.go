package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13, // TLS1.3以上を許可
		CipherSuites: []uint16{ // 使用する暗号スイート
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		NextProtos:               []string{"h2", "http/1.1"},
		PreferServerCipherSuites: true, // サーバーの暗号スイートを優先
		GetConfigForClient: func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
			// ハンドシェイクの時間を計測
			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("🚀 ハンドシェイク時間: %s", elapsed)
			}()

			// どの暗号スイートが選択されたかを確認
			log.Printf("🚀 クライアント接続: TLS version: %x, Cipher Suite: %x",
				clientHello.SupportedVersions, clientHello.CipherSuites)
			return nil, nil
		},
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Alt-Svcヘッダーを追加してHTTP/3の利用を通知
	if r.ProtoMajor == 2 {
		w.Header().Add("Alt-Svc", `h3=":8444"; ma=2592000`)
	}

	log.Printf("🔍 リクエスト: プロトコル=%s, リモートアドレス=%s",
		r.Proto, r.RemoteAddr)
	fmt.Fprintf(w, "HTTP/2 and HTTP/3 World!\n")
}

func main() {
	// サーバーの設定
	server := &http.Server{
		Addr:      ":8443",
		Handler:   http.HandlerFunc(handler),
		TLSConfig: getTLSConfig(),
	}

	// HTTP/3のサーバーを起動
	http3Server := &http3.Server{
		Addr:      ":8444",
		Handler:   http.HandlerFunc(handler),
		TLSConfig: getTLSConfig(),
		QuicConfig: &quic.Config{
			EnableDatagrams: true,
			MaxIdleTimeout:  30 * time.Second,
		},
	}

	go func() {
		log.Printf("Starting HTTP/3 server on :8444")
		err := http3Server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
		if err != nil {
			log.Fatal("🚨 サーバーの起動に失敗しました。", err)
		}
	}()

	log.Printf("Starting HTTP/2 server on :8443")
	err := server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal("🚨 サーバーの起動に失敗しました。", err)
	}
}
