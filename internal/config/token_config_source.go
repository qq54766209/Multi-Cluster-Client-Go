package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"k8s.io/client-go/rest"
)

// TokenConfigSource 是 K8sConfigSource 的一个实现，可以从文件读取或直接传递 Token
type TokenConfigSource struct {
	APIServerURL string
	Token        string // 直接传递的 Token
	TokenPath    string // Token 文件路径
	Insecure     bool   // 是否忽略 TLS 证书验证
	CACertPath   string // CA 证书文件路径
}

// GetConfig 实现了 K8sConfigSource 接口，返回 rest.Config
func (t *TokenConfigSource) GetConfig() (*rest.Config, error) {
	var token string
	var err error

	// 如果 Token 已经通过字符串直接传递，则使用该 Token
	if t.Token != "" {
		token = t.Token
		log.Printf("Using provided token for API server %s", t.APIServerURL)
	} else if t.TokenPath != "" {
		// 否则尝试从文件中读取 Token
		log.Printf("Attempting to read token from file: %s", t.TokenPath)
		token, err = readTokenFromFile(t.TokenPath)
		if err != nil {
			log.Printf("Error reading token from %s: %v", t.TokenPath, err)
			return nil, fmt.Errorf("failed to read token from %s: %w", t.TokenPath, err)
		}
		log.Printf("Successfully read token from file: %s", t.TokenPath)
	} else {
		err := fmt.Errorf("no token provided for API server %s", t.APIServerURL)
		log.Printf("Error: %v", err)
		return nil, err
	}

	tlsConfig := rest.TLSClientConfig{
		Insecure: t.Insecure, // 是否忽略 TLS 证书验证
	}

	if !t.Insecure {
		if t.CACertPath != "" {
			// 如果 Insecure = false，且指定了 CA 证书路径，则设置 CA 证书
			if _, err := os.Stat(t.CACertPath); os.IsNotExist(err) {
				log.Printf("CA certificate file does not exist at path: %s", t.CACertPath)
				return nil, fmt.Errorf("CA certificate file not found: %s", t.CACertPath)
			}
			tlsConfig.CAFile = t.CACertPath
			log.Printf("Using CA certificate from: %s", t.CACertPath)
		} else {
			err := fmt.Errorf("insecure is set to false, but no CA certificate path is provided for API server %s", t.APIServerURL)
			log.Printf("Error: %v", err)
			return nil, err
		}
	} else {
		log.Printf("Insecure TLS is enabled, certificate validation is skipped for API server %s", t.APIServerURL)
	}

	config := &rest.Config{
		Host:            t.APIServerURL,
		BearerToken:     token,
		TLSClientConfig: tlsConfig,
		Timeout:         10 * time.Second, // 设置超时时间
	}

	log.Printf("Successfully created config for API server %s", t.APIServerURL)
	return config, nil
}

// 读取 Token 文件内容
func readTokenFromFile(tokenPath string) (string, error) {
	data, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Printf("Failed to read token file: %v", err)
		return "", err
	}
	return string(data), nil
}

// NewTokenConfigSourceWithToken 直接使用 Token 字符串创建 TokenConfigSource
func NewTokenConfigSourceWithToken(apiServerURL, token string, insecure bool, caCertPath string) *TokenConfigSource {
	log.Printf("Creating TokenConfigSource with provided token for API server %s", apiServerURL)
	return &TokenConfigSource{
		APIServerURL: apiServerURL,
		Token:        token,
		Insecure:     insecure,
		CACertPath:   caCertPath,
	}
}

// NewTokenConfigSourceWithFile 从 Token 文件路径创建 TokenConfigSource
func NewTokenConfigSourceWithFile(apiServerURL, tokenPath string, insecure bool, caCertPath string) *TokenConfigSource {
	log.Printf("Creating TokenConfigSource with token file %s for API server %s", tokenPath, apiServerURL)
	return &TokenConfigSource{
		APIServerURL: apiServerURL,
		TokenPath:    tokenPath,
		Insecure:     insecure,
		CACertPath:   caCertPath,
	}
}
