package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

//生成md5
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 调整key长度
func resetKey(key []byte, size uint) []byte {
	if len(key) < int(size) {
		key_tmp := make([]byte, size)
		copy(key_tmp, key)
		return key_tmp
	}
	return key[:size]
}

//加密字符串
func AESEncrypt(data []byte, key []byte) ([]byte, error) {
	// 通过key解析
	var iv = resetKey(key, aes.BlockSize)
	aesBlockEncrypter, err := aes.NewCipher(iv)
	if err != nil {
		return nil, err
	}
	// 解析key
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	encrypted := make([]byte, len(data))
	aesEncrypter.XORKeyStream(encrypted, data)
	return encrypted, nil
}

//解密字符串
func AESDecrypt(src []byte, key []byte) ([]byte, error) {
	// 创建解析器
	var iv = resetKey(key, aes.BlockSize) // 截取key长度
	aesBlockDecrypter, err := aes.NewCipher(iv)
	if err != nil {
		return nil, err
	}
	// 解析参数
	decrypted := make([]byte, len(src))
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	if e := recover(); e != nil {
		return nil, e.(error)
	}
	return decrypted, nil
}

func AESEncrpyAndBase64(data []byte, key []byte) (string, error) {
	//aes加密
	e, err_e := AESEncrypt(data, []byte(key))
	if err_e != nil {
		return "", err_e
	}
	// base64加密
	outStr := base64.StdEncoding.EncodeToString(e)
	return outStr, nil
}

func AESDecryptAndBase64(data []byte, key []byte) ([]byte, error) {
	// base64解密
	d_base, err_bd := base64.StdEncoding.DecodeString(string(data))
	if err_bd != nil {
		return nil, err_bd
	}
	//aes解密
	d, err_d := AESDecrypt(d_base, key)
	if err_d != nil {
		return nil, err_d
	}
	return d, nil
}

// base64 解析
func Base64Decode(str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, err
}

// base64 加密
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
