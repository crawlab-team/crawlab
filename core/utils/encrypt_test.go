package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncryptAesPassword(t *testing.T) {
	plainText := "crawlab"
	encryptedText, err := EncryptAES(plainText)
	require.Nil(t, err)
	decryptedText, err := DecryptAES(encryptedText)
	require.Nil(t, err)
	fmt.Println(fmt.Sprintf("plainText: %s", plainText))
	fmt.Println(fmt.Sprintf("encryptedText: %s", encryptedText))
	fmt.Println(fmt.Sprintf("decryptedText: %s", decryptedText))
	require.Equal(t, decryptedText, plainText)
	require.NotEqual(t, decryptedText, encryptedText)
}
