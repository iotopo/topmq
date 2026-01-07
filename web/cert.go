package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/iotopo/topmq/internal/utils"
	"math/big"
	rd "math/rand"
	"os"
	"sync"
	"time"
)

// createRootCert 生成 root.crt, client.key
func createRootCert() (err error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()), //证书序列号
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"GlobalSign nv-sa"},
			OrganizationalUnit: []string{"Root CA"},
			//Province:           []string{"DL"},
			CommonName: "GlobalSign Root CA",
			//Locality:           []string{"DaLian"},
		},
		NotBefore:             time.Now(),                                                                 //证书有效期开始时间
		NotAfter:              time.Now().AddDate(20, 0, 0),                                               //证书有效期结束时间
		BasicConstraintsValid: true,                                                                       //基本的有效性约束
		IsCA:                  false,                                                                      //是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating private key error: %w", err)
	}
	var crtBuf []byte
	var keyBuf []byte

	// 生成自签名证书
	crtBuf, err = x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}
	keyBuf = x509.MarshalPKCS1PrivateKey(privateKey)

	certFileOut, err := os.Create(certFile)
	defer certFileOut.Close()

	err = pem.Encode(certFileOut, &pem.Block{Type: "CERTIFICATE", Bytes: crtBuf})
	if err != nil {
		return
	}

	keyFileOut, err := os.Create(keyFile)
	defer keyFileOut.Close()
	err = pem.Encode(keyFileOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBuf})

	return
}

var certFile = "./tls/cert.cer"
var keyFile = "./tls/cert.key"
var locker sync.Mutex

func initCert() error {
	locker.Lock()
	defer locker.Unlock()
	if !utils.PathExists(certFile) { // 自动创建自签名证书
		err := createRootCert()
		if err != nil {
			return fmt.Errorf("failed to create tls cert: %w", err)
		}
	}
	return nil
}
