package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
)

//生成RSA私钥和公钥，保存到文件中
// bits 证书大小
func GenerateRSAKey(bits int) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.OpenFile("./rsa_private_key.pem", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: X509PrivateKey}
	//将数据保存到文件
	pem.Encode(privateFile, &privateBlock)

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.OpenFile("./rsa_public_key.pem", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA PUBLIC KEY", Bytes: X509PublicKey}
	//保存到文件
	pem.Encode(publicFile, &publicBlock)
}

func TestCreateRSA(t *testing.T) {
	GenerateRSAKey(2048)
}

//RSA加密
// plainText 要加密的数据
// path 公钥匙文件地址
func RSA_Encrypt(plainText []byte, path string) []byte {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText
}

//RSA解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func RSA_Decrypt(cipherText []byte, path string) []byte {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	tt := string(buf)
	fmt.Println(tt)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		panic(err)
	}
	//返回明文
	return plainText
}

func TestRSA(t *testing.T) {
	////加密
	//data := []byte("123456")
	//encrypt := RSA_Encrypt(data, "./rsa_public_key.pem")
	//t.Log(len(encrypt), string(encrypt))
	//encrypt2 := RSA_Encrypt(data, "./rsa_public_key.pem")
	//t.Log(len(encrypt2), string(encrypt2))
	//t.Log(string(encrypt) == string(encrypt2))
	//// 解密
	//decrypt := RSA_Decrypt(encrypt, "./rsa_private_key.pem")
	//t.Log(string(decrypt))

	pass := "J0z2Z1yWczFH8GfbxCjJWGbKsMNfFwENBnivC76GlXeeIx70GZ6QOWehX6xP82tO3H8KVY4YOtpnw+2UwKXTC9BmucqI8D3d0x1AcdJGEfQ/2E98MBmNvGUnqm37okGzjsVGHfaDOWobuVlIqrPE3m/vkP/cMFKgg77/lbApfiFAlL4YeE24CO7BaZ3pet2K9ikdecoNr9auMwvYtUDzmgTM8Bw1zLZpSKAF4svn+BaxOIwRzBW298XrF0nbPTVVtfH8tr2yJXmdtfUwBYkx0DzabTZYLIg7Wd0Xsj96zqPdgGcnSNRcn1N1cIj+WSfOar7sP5+h244VG5bl+UFqrw=="
	buf, _ := base64.StdEncoding.DecodeString(pass)
	t.Log(string(buf))
	decrypt := RSA_Decrypt(buf, "./rsa_private_key.pem")
	t.Log(string(decrypt))
}
