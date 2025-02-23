package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12, // TLS1.2以上を許可
		CipherSuites: []uint16{ // 使用する暗号スイート
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true, // サーバーの暗号スイートを優先
		GetConfigForClient: func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
			// どの暗号スイートが選択されたかを確認
			log.Printf("Client connected with TLS version: %x, Cipher Suite: %x",
				clientHello.SupportedVersions, clientHello.CipherSuites)
			return nil, nil
		},
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, TLS World!\n")
}

func main() {
	// サーバーの設定
	server := &http.Server{
		Addr:      ":8443",
		Handler:   http.HandlerFunc(handler),
		TLSConfig: getTLSConfig(),
	}

	log.Printf("Starting server on :8443")
	log.Fatal(server.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}
