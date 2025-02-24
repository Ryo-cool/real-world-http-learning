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

func getTLSConfig(tlsVersion uint16) *tls.Config {
	// CA証明書を読み込む
	cert, err := os.ReadFile("certs/server.crt")
	if err != nil {
		log.Fatal("🚨 サーバーの証明書を読み込めませんでした。", err)
	}

	// CA証明書を信頼するよう設定
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(cert) {
		log.Fatal("🚨 CA証明書を信頼できませんでした。")
	}

	return &tls.Config{
		MinVersion: tlsVersion,
		RootCAs:    caCertPool, // CA証明書を信頼するよう設定
	}
}

func main() {

	// TLSのバージョン指定(値をバージョン比較で指定)
	tlsConfig := getTLSConfig(tls.VersionTLS13)

	// HTTPトランスポートを作成
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, // タイムアウト設定
	}

	start := time.Now()
	// サーバーにGETリクエストを送信
	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		log.Fatal("🚨 サーバーに接続できませんでした。", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("🚨 レスポンスの読み込みに失敗しました。", err)
	}
	// レスポンスのステータスコードを確認
	log.Printf("Response status: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(body))
	// ハンドシェイク時間を表示
	log.Printf("🚀 ハンドシェイク時間: %s", elapsed)
}
