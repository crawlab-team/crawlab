package notification

import (
	"context"
	"encoding/base64"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/trace"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"strings"
)

func sendMailGmail(ch *models.NotificationChannelV2, smtpConfig smtpAuthentication, options sendOptions, htmlBody, txtBody string) error {
	// 读取服务账户 JSON 密钥
	b := []byte(ch.GoogleOAuth2Json)

	// 使用服务账户 JSON 密钥文件创建 JWT 配置
	config, err := google.JWTConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Errorf("Unable to parse service account key file to config: %v", err)
		return trace.TraceError(err)
	}

	// 使用服务账户的电子邮件地址来模拟用户
	config.Subject = ch.SMTPUsername

	// 创建 Gmail 服务
	client := config.Client(context.Background())
	srv, err := gmail.New(client)
	if err != nil {
		log.Errorf("Unable to create Gmail client: %v", err)
		return trace.TraceError(err)
	}

	// 创建 MIME 邮件
	m, err := getMailMessage(smtpConfig, options, htmlBody, txtBody)
	if err != nil {
		return err
	}

	var buf strings.Builder
	if _, err := m.WriteTo(&buf); err != nil {
		log.Errorf("Unable to write message: %v", err)
		return trace.TraceError(err)
	}

	// 将邮件内容进行 base64 编码
	gmsg := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(buf.String())),
	}

	// 发送邮件
	_, err = srv.Users.Messages.Send("me", gmsg).Do()
	if err != nil {
		log.Errorf("Unable to send email: %v", err)
		return trace.TraceError(err)
	}

	return nil
}
