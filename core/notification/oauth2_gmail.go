package notification

import (
	"context"
	"github.com/apex/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/smtp"
	"time"
)

// 获取服务账户的OAuth2配置
func getGmailOAuth2Token(oauth2Json string) (token *oauth2.Token, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 读取服务账户 JSON 密钥
	b := []byte(oauth2Json)

	// 使用服务账户 JSON 密钥文件创建 JWT 配置
	config, err := google.JWTConfigFromJSON(b, "https://mail.google.com/")
	if err != nil {
		log.Errorf("Unable to parse service account key file to config: %v", err)
		return nil, err
	}

	// 使用服务账户的电子邮件和访问令牌
	token, err = config.TokenSource(ctx).Token()
	if err != nil {
		log.Errorf("Unable to generate token: %v", err)
		return nil, err
	}
	return token, nil
}

// GmailOAuth2Auth 自定义OAuth2认证
type GmailOAuth2Auth struct {
	username, accessToken string
}

func (a *GmailOAuth2Auth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "XOAUTH2", []byte("user=" + a.username + "\x01auth=Bearer " + a.accessToken + "\x01\x01"), nil
}

func (a *GmailOAuth2Auth) Next(_ []byte, _ bool) ([]byte, error) {
	return nil, nil
}
