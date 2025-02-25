package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
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

func measureRequest(url string, protocol string) {
	start := time.Now()
	var resp *http.Response
	var err error

	switch protocol {
	case "HTTP/3":
		// HTTP/3クライアントを作成
		tlsConfig := getTLSConfig(tls.VersionTLS13)
		roundTripper := &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
			QuicConfig: &quic.Config{
				EnableDatagrams: true,
				MaxIdleTimeout:  30 * time.Second,
			},
		}
		defer roundTripper.Close()

		client := &http.Client{
			Transport: roundTripper,
			Timeout:   30 * time.Second,
		}

		// HTTP/3サーバーにリクエスト（ポート8444）
		http3URL := "https://localhost:8444"
		fmt.Println("🚀 HTTP/3リクエストを送信中...")
		resp, err = client.Get(http3URL)

	case "HTTP/2":
		// HTTP/2クライアントを作成
		tlsConfig := getTLSConfig(tls.VersionTLS13)
		tr := &http.Transport{
			TLSClientConfig:   tlsConfig,
			ForceAttemptHTTP2: true, // HTTP/2を強制
		}

		client := &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		}

		fmt.Println("🚀 HTTP/2リクエストを送信中...")
		resp, err = client.Get(url)

	case "HTTP/1.1":
		// HTTP/1.1クライアントを作成
		tlsConfig := getTLSConfig(tls.VersionTLS13)
		tr := &http.Transport{
			TLSClientConfig:   tlsConfig,
			ForceAttemptHTTP2: false,                                                                     // HTTP/2を無効化
			TLSNextProto:      make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0), // HTTP/2を無効化
		}

		client := &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		}

		fmt.Println("🚀 HTTP/1.1リクエストを送信中...")
		resp, err = client.Get(url)
	}

	if err != nil {
		log.Printf("🚨 %sリクエストエラー: %v", protocol, err)
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("🚨 レスポンスの読み込みに失敗しました: %v", err)
		return
	}

	// レスポンスの詳細を表示
	fmt.Printf("📡 プロトコル: %s (要求: %s)\n", resp.Proto, protocol)
	fmt.Printf("📊 ステータス: %d\n", resp.StatusCode)
	fmt.Printf("📝 レスポンス: %s\n", string(body))
	fmt.Printf("⏱️ 接続時間: %s\n\n", elapsed)
}

func main() {
	// コマンドライン引数でプロトコルを選択できるようにする
	measureAll := flag.Bool("all", false, "すべてのプロトコルを計測")
	protocol := flag.String("protocol", "", "計測するプロトコル (HTTP/1.1, HTTP/2, HTTP/3)")
	flag.Parse()

	url := "https://localhost:8443"

	if *measureAll {
		fmt.Println("🔹 HTTP/1.1 の計測開始...")
		measureRequest(url, "HTTP/1.1")

		fmt.Println("🔹 HTTP/2 の計測開始...")
		measureRequest(url, "HTTP/2")

		fmt.Println("🔹 HTTP/3 の計測開始...")
		measureRequest(url, "HTTP/3")
	} else if *protocol != "" {
		fmt.Printf("🔹 %s の計測開始...\n", *protocol)
		measureRequest(url, *protocol)
	} else {
		// デフォルトはHTTP/2
		fmt.Println("🔹 HTTP/2 の計測開始...")
		measureRequest(url, "HTTP/2")
	}
}
