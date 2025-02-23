package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TLS設定
	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	// HTTPハンドラーの設定
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, TLS World!\n")
	})

	// サーバーの設定
	server := &http.Server{
		Addr:      ":8443",
		Handler:   mux,
		TLSConfig: cfg,
	}

	log.Printf("Starting server on :8443")
	log.Fatal(server.ListenAndServeTLS("certs/server.crt", "certs/server.key"))
}
