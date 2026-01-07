package mqttserver

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/iotopo/topmq/config"
	"github.com/sirupsen/logrus"
	"math/big"
	"net"
	"os"
	"time"
)

// 保存证书到文件
func saveCertificate(filename string, derBytes []byte) {
	certOut, err := os.Create(filename)
	if err != nil {
		logrus.Fatal("创建证书文件失败:", err)
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	if err != nil {
		logrus.Fatal("编码证书失败:", err)
	}
}

// 保存私钥到文件
func savePrivateKey(filename string, privateKey *ecdsa.PrivateKey) {
	keyOut, err := os.Create(filename)
	if err != nil {
		logrus.Fatal("创建私钥文件失败:", err)
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		logrus.Fatal("编码私钥失败:", err)
	}

	err = pem.Encode(keyOut, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})
	if err != nil {
		logrus.Fatal("保存私钥失败:", err)
	}
}

const caCertFile = "./mqtt/certs/ca.crt"
const caKeyFile = "./mqtt/certs/ca.key"
const serverCertFile = "./mqtt/certs/server.crt"
const serverKeyFile = "./mqtt/certs/server.key"
const clientCertFile = "./mqtt/certs/client.crt"
const clientKeyFile = "./mqtt/certs/client.key"

func generateServerCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey) {
	conf := config.Conf.MQTTServer
	// 生成服务器证书
	serverCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Province:           []string{"Dalian"},
			Organization:       []string{"IOTOPO"},
			OrganizationalUnit: []string{"Server Unit"},
			CommonName:         "MQTT Server",
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0), // 10年有效期
		SubjectKeyId: []byte{1, 2, 3, 4, 5},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{"localhost", "mqtt.example.com"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	if len(conf.DNSNames) > 0 {
		serverCert.DNSNames = conf.DNSNames
	}
	if len(conf.IPAddresses) > 0 {
		var ipAddresses []net.IP
		for _, ip := range conf.IPAddresses {
			ipAddresses = append(ipAddresses, net.ParseIP(ip))
		}
		serverCert.IPAddresses = ipAddresses
	}

	// 生成服务器私钥
	serverPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logrus.WithError(err).Fatal("生成服务器私钥失败")
	}

	// 创建服务器证书
	serverBytes, err := x509.CreateCertificate(rand.Reader, serverCert, caCert, &serverPrivKey.PublicKey, caKey)
	if err != nil {
		logrus.WithError(err).Fatal("创建服务器证书失败")
	}

	// 保存服务器证书
	saveCertificate(serverCertFile, serverBytes)
	savePrivateKey(serverKeyFile, serverPrivKey)
}

func generateClientCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey) {
	clientCert := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Province:           []string{"Daliang"},
			Organization:       []string{"IOTOPO"},
			OrganizationalUnit: []string{"Client Unit"},
			CommonName:         "MQTT Client",
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0), // 10年有效期
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// 生成客户端私钥
	clientPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logrus.WithError(err).Fatal("生成客户端私钥失败")
	}

	// 创建客户端证书
	clientBytes, err := x509.CreateCertificate(rand.Reader, clientCert, caCert, &clientPrivKey.PublicKey, caKey)
	if err != nil {
		logrus.WithError(err).Fatal("创建客户端证书失败")
	}

	// 保存客户端证书
	saveCertificate(clientCertFile, clientBytes)
	savePrivateKey(clientKeyFile, clientPrivKey)
}

func generateCA() (*x509.Certificate, *ecdsa.PrivateKey) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Province:     []string{"Dalian"},
			Organization: []string{"IOTOPO"},
			//OrganizationalUnit: []string{"CA Unit"},
			CommonName: "IOTOPO MQTT CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10年有效期
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// 生成 CA 私钥
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logrus.WithError(err).Fatal("生成 CA 私钥失败")
	}

	// 创建 CA 证书
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		logrus.WithError(err).Fatal("创建 CA 证书失败")
	}

	// 保存 CA 证书
	saveCertificate(caCertFile, caBytes)
	savePrivateKey(caKeyFile, caPrivKey)

	return ca, caPrivKey
}

func generateCerts() {
	// 创建证书目录
	err := os.MkdirAll("./mqtt/certs", 0755)
	if err != nil {
		logrus.Fatal("创建证书目录失败:", err)
	}

	// 生成 CA 证书
	caCert, caKey := generateCA()

	// 生成服务器证书
	generateServerCert(caCert, caKey)
	// 生成客户端证书
	generateClientCert(caCert, caKey)

	logrus.Info("MQTT 证书生成完成！")
}
