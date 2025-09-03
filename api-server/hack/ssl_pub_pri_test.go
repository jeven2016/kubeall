package hack

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

// loadCertificate 从 server.crt 加载公钥
func loadCertificate(certPath string) (*rsa.PublicKey, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取证书文件: %v", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("无效的证书格式")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %v", err)
	}

	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥不是 RSA 格式")
	}

	return publicKey, nil
}

// loadPrivateKey 从 server.key 加载私钥
func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
		return nil, fmt.Errorf("无效的私钥格式")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}

	return privateKey, nil
}

// publicKeyEncrypt 使用公钥加密数据
func publicKeyEncrypt(publicKey *rsa.PublicKey, data []byte) (string, error) {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return "", fmt.Errorf("公钥加密失败: %v", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// privateKeyDecrypt 使用私钥解密数据
func privateKeyDecrypt(privateKey *rsa.PrivateKey, encryptedBase64 string) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, fmt.Errorf("Base64 解码失败: %v", err)
	}

	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("私钥解密失败: %v", err)
	}
	return decrypted, nil
}

// privateKeyEncrypt 使用私钥加密数据（签名场景）
func privateKeyEncrypt(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return "", fmt.Errorf("私钥加密（签名）失败: %v", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// publicKeyDecrypt 使用公钥解密数据（验证签名）
func publicKeyDecrypt(publicKey *rsa.PublicKey, data []byte, signatureBase64 string) error {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("Base64 解码失败: %v", err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, nil)
	if err != nil {
		return fmt.Errorf("公钥解密（验证签名）失败: %v", err)
	}
	return nil
}

func TestPubPri(t *testing.T) {
	// 文件路径
	certPath := "/home/jeven/Desktop/workspace/project/k8s-env/ca_TLS/work/certs/server.crt"
	keyPath := "/home/jeven/Desktop/workspace/project/k8s-env/ca_TLS/work/certs/server.key"

	// 加载公钥和私钥
	publicKey, err := loadCertificate(certPath)
	if err != nil {
		log.Fatalf("加载公钥失败: %v", err)
	}
	privateKey, err := loadPrivateKey(keyPath)
	if err != nil {
		log.Fatalf("加载私钥失败: %v", err)
	}

	// 示例数据
	data := []byte("Hello, IoT Device!")

	// 1. 公钥加密 / 私钥解密
	fmt.Println("=== 公钥加密 / 私钥解密 ===")
	encrypted, err := publicKeyEncrypt(publicKey, data)
	if err != nil {
		log.Fatalf("公钥加密失败: %v", err)
	}
	fmt.Printf("加密结果 (Base64): %s\n", encrypted)

	decrypted, err := privateKeyDecrypt(privateKey, encrypted)
	if err != nil {
		log.Fatalf("私钥解密失败: %v", err)
	}
	fmt.Printf("解密结果: %s\n", string(decrypted))

	// 2. 私钥加密 / 公钥解密（签名场景）
	fmt.Println("\n=== 私钥加密 / 公钥解密（签名） ===")
	signature, err := privateKeyEncrypt(privateKey, data)
	if err != nil {
		log.Fatalf("私钥加密（签名）失败: %v", err)
	}
	fmt.Printf("签名 (Base64): %s\n", signature)

	err = publicKeyDecrypt(publicKey, data, signature)
	if err != nil {
		log.Fatalf("公钥解密（验证签名）失败: %v", err)
	}
	fmt.Println("签名验证通过")
}
