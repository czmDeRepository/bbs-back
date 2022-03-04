package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
)

var RsaPrivateKey = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAoKjkQK3Km4iVXuzehVUr7dLjRAQ7vQdurVOBa4LrlYL1L4bm\ndmm8WPsaI5b540E0i/MPBvQbmjXxwlreECIPuMs2ni2yESSeLsfKUwEUipVv1r0X\n/LhRSFh/G0qVe6uheCIGdgiMXu0g+HlEj5FOjoZwKYUwo7CpQprP4Mi/aRK7hQnZ\nqV7eWTPRj187BqNpecx+0uD4ju6bbePr+r/CQ7s879l+jZk4J40sY1wq6+Xmq2yZ\nd9xT+sZtW8DhbIb8+vgdA6Qn/lydIwW/BphYVdi6/89Ah/Yx975tjJYFZppT3O0B\nnzka4LYZaHvypqlPwvGpDfT6cQkEIycWqUn0NwIDAQABAoIBAGP7VzkNLsaGRHbz\nsDOH4rO5hyef/tFPm8gP7L3MlvHPsuhl6myMtuMlOYomfdK4lIv3skVgiwD4S0Dp\nrcKaf/A/vvjHdUaH5E02lqn7RR5Ni3E28oOa7TK5qDiCvO5ezDjn699uyHSW2+rm\nBQ3XnuzBq3Gbar6jxWIm1/GK3Nv1CQOZbHg85082XAxwTuobdaUmPzjj8q8IMzLC\nc0ve3tDlA8ujc8oCzbuetuo/YdvvWQisSCkwKDBBcCOI3qwRFv1gEoPzusUGxIB1\nW5ymbbBUubjE0Zu7GzXPTL9dxGs42Yy/XHr58w+olL3wNA1arGlIMKqp/vUrIjhR\nMTUjOUECgYEAydEqZtTc53dvHXn7QGEzGp2HHmd34Qrp4AHO2FUxI2hbM6x3PajT\nrRcUx61wfIrHeGs5p5XNcKr07FbHWmdMZ10/N6gDvV98MC7oT8DBoK5fhtcgb/UJ\nTp6B/YDKisNMRJokDs3JYcPsZFawGrCKxJviSINoc8QRTpIyeEWLSxcCgYEAy8r+\nb9AmMz0y1bjn2lFqNWKiG6Zfte+5HsDKScHTrHxB3yqyDeQ/OI8nCUUrFSdryupR\nTo6rLAuCluAfm4osX5UUJIivtWQ/cNTLZwmeYxnHKVdCfktAjVm1vW3ENju0URHO\nuMSamYx+DnxdSkI9cIUzyHE5XezvIFtv0rRv0+ECgYBlz/WaJuzCgMg3kKAmHGMR\nnELcHcqmZ0ERVxgonuHJQQ4xhWIqYw9WlPxQt7i1u7VhlIZjevlHS5d/20961f3/\nb1VDGKm6UX9vN1rPUSjdjNp4RfMBSBbH6MMfRmfnlRrWyQRDy6E6hwKso+b3r/Rx\nt0py1ohNTq6EetCnSD47RQKBgFouXQOLv6vC9CDhbzAMAQzYtdW4fPgcufWi6KFU\n9V+JqPihgyNkkplrt6GBizwUMr4bjJlPxu15tnMfgL0qmtI9PSmhlueVEgHTGKNi\n/UTrXler9o++qzUhsqu9zCsXpoaNc8YNskAqjInKfNnkkB2fxDd56yHmPDc8XzKF\niErhAoGAPnxvndy6iLkQOFvjwvHQ+FzSNu1Q58Oi7JK5D74gDzFPJ3kQqSFQltqW\n7Lmh/Bog3bswYXPXwrES63dW3HpWHBnoPTw0W1W3wqjPl2WTdQXjcJ7gxoHecGwr\ngSnlixOgMAddasHsjj0X9VERiIrmDhCRezMWIQKa7JvmqOBUSK4=\n-----END RSA PRIVATE KEY-----")

// DecryptRSA RSA解密
func DecryptRSA(cipherText string) ([]byte, error) {
	cipher, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	//pem解码
	block, _ := pem.Decode(RsaPrivateKey)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipher)
	if err != nil {
		return nil, err
	}
	//返回明文
	return plainText, nil
}

// DecryptPassword RSA解密密码并MD5加密
// 返回长度32的16进制字符串
func DecryptPassword(password string) (string, error) {
	decryptRSA, err := DecryptRSA(password)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	hash.Write(decryptRSA)
	return hex.EncodeToString(hash.Sum(nil)), err
}
