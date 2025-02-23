package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// サーバーの証明書を読み込む
	cert, err := os.ReadFile("certs/server.crt")
	if err != nil {
		log.Fatal("🚨 サーバーの証明書を読み込めませんでした。", err)
	}

	// CA証明書を信頼するよう設定
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(cert) {
		log.Fatal("🚨 CA証明書を信頼できませんでした。")
	}

	// TLS1.3のみを使用する。
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    caCertPool, // CA証明書を信頼するよう設定
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, // タイムアウト設定
	}

	// サーバーにGETリクエストを送信
	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		log.Fatal("🚨 サーバーに接続できませんでした。", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("🚨 レスポンスの読み込みに失敗しました。", err)
	}
	// レスポンスのステータスコードを確認
	log.Printf("Response status: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(body))
}
