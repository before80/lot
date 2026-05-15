package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/before80/lot/lg"
	"golang.org/x/crypto/acme/autocert"
)

func startProductionServer() {
	lg.InfoToFile("生产环境：配置Let's Encrypt自动证书...")

	// 1. 配置域名（替换为你的实际域名）
	domains := []string{
		"be-better.top",
		"www.be-better.top",
		// 可以添加更多子域名
	}

	// 2. 创建证书缓存目录
	cacheDir := "./autoCertCache"
	if err := createCertCacheDir(cacheDir); err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建证书缓存目录失败: %v", err))
		return
	}

	// 3. 配置证书管理器
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
		Cache:      autocert.DirCache(cacheDir),
		// 可选：配置证书更新前的邮件通知
		// Email: "your-email@example.com",
	}

	// 4. 配置TLS
	tlsConfig := &tls.Config{
		GetCertificate: certManager.GetCertificate,
		MinVersion:     tls.VersionTLS12,
		NextProtos:     []string{"h2", "http/1.1"},
		// 推荐启用更安全的加密套件
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// 5. 创建HTTP服务器
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   http.DefaultServeMux,
	}

	// 6. 启动HTTP挑战服务器（必须在80端口）
	// 重要：必须在单独的goroutine中启动
	go func() {
		// 创建专门用于验证的HTTP服务器
		challengeHandler := certManager.HTTPHandler(nil)
		challengeServer := &http.Server{
			Addr:    ":80",
			Handler: challengeHandler,
		}

		lg.InfoToFile("启动HTTP挑战服务器在 :80")
		if err := challengeServer.ListenAndServe(); err != nil {
			// 注意：这里不要用Fatal，避免影响主服务器
			lg.InfoToFile(fmt.Sprintf("HTTP挑战服务器错误: %v", err))
		}
	}()

	// 7. 启动主HTTPS服务器
	lg.InfoToFile("启动HTTPS服务器在 :443")
	lg.InfoToFile(fmt.Sprintf("已配置域名: %v\n", domains))
	lg.InfoToFile("请确保阿里云安全组已开放80和443端口！")

	// 重要：这里传递空字符串，证书由GetCertificate动态提供
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("HTTPS服务器启动失败: %v", err)
	}
}

// createCertCacheDir 创建并配置证书缓存目录
func createCertCacheDir(dirPath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 创建目录
		lg.InfoToFile(fmt.Sprintf("创建证书缓存目录: %s\n", dirPath))
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return fmt.Errorf("无法创建目录: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("检查目录失败: %w", err)
	} else {
		lg.InfoToFile(fmt.Sprintf("证书缓存目录已存在: %s\n", dirPath))
	}

	// 验证目录可写性
	testFile := filepath.Join(dirPath, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("目录不可写: %w", err)
	}
	defer os.Remove(testFile) // 清理测试文件

	return nil
}

// 可选：健康检查端点
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok", "service": "autocert-demo"}`))
}
