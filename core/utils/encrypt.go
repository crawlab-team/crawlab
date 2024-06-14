package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"io"
)

func GetSecretKey() string {
	return constants.DefaultEncryptServerKey
}

func GetSecretKeyBytes() []byte {
	return []byte(GetSecretKey())
}

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func EncryptMd5(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

func unPadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

func EncryptAES(src string) (res string, err error) {
	srcBytes := []byte(src)
	key := GetSecretKeyBytes()
	block, err := aes.NewCipher(key)
	if err != nil {
		return res, err
	}
	srcBytes = padding(srcBytes, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(srcBytes, srcBytes)
	res = hex.EncodeToString(srcBytes)
	return res, nil
}

func DecryptAES(src string) (res string, err error) {
	srcBytes, err := hex.DecodeString(src)
	if err != nil {
		return res, err
	}
	key := GetSecretKeyBytes()
	block, err := aes.NewCipher(key)
	if err != nil {
		return res, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(srcBytes, srcBytes)
	res = string(unPadding(srcBytes))
	return res, nil
}
