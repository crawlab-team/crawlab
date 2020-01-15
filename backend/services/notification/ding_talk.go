package notification

import (
	"crawlab/model"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/imroc/req"
	"github.com/royeo/dingrobot"
	"runtime/debug"
)

func SendDingTalkNotification(t model.Task, s model.Spider) error {
	// 获取用户
	user, _ := model.GetUser(t.UserId)

	// 如果AppKey或AppSecret未设置，则返回错误
	if user.Setting.DingTalkAppKey == "" || user.Setting.DingTalkAppSecret == "" {
		return errors.New("ding_talk_app_key or ding_talk_app_secret is empty")
	}

	// 获取access_token
	accessToken, err := GetDingTalkAccessToken(user)
	if err != nil {
		return err
	}

	// 时间戳
	//timestamp := time.Now().Unix()

	// 计算sign
	//signRawString := fmt.Sprintf("%d\n%s", timestamp, user.Setting.DingTalkAppSecret)
	//sign := utils.ComputeHmacSha256(signRawString, user.Setting.DingTalkAppSecret)

	// 请求数据
	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", accessToken)
	robot := dingrobot.NewRobot(url)

	text := "it works"
	if err := robot.SendText(text, []string{}, false); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}
	//header := req.Header{
	//	"Content-Type": "application/json; charset=utf-8",
	//	"timestamp":    fmt.Sprintf("%d000", timestamp),
	//	"sign":         sign,
	//}
	//data := req.Param{
	//	"msgtype": "text",
	//	"text": req.Param{
	//		"text": "it works",
	//	},
	//	"at": req.Param{
	//		"atMobiles": []string{},
	//		"isAtAll": false,
	//	},
	//}
	//res, err := req.Post(url, header, req.BodyJSON(&data))
	//if err != nil {
	//	log.Errorf("dingtalk notification error: " + err.Error())
	//	debug.PrintStack()
	//	return err
	//}
	//log.Infof(fmt.Sprintf("%+v", res))
	return nil
}

func GetDingTalkAccessToken(u model.User) (string, error) {
	type ResBody struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
	}

	// 请求数据
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", u.Setting.DingTalkAppKey, u.Setting.DingTalkAppSecret)
	res, err := req.Get(url)
	if err != nil {
		log.Errorf("get dingtalk access_token error: " + err.Error())
		debug.PrintStack()
		return "", err
	}

	// 解析相应body
	var resBody ResBody
	if err := res.ToJSON(&resBody); err != nil {
		log.Errorf("get dingtalk access_token error: " + err.Error())
		debug.PrintStack()
		return "", err
	}

	return resBody.AccessToken, nil
}
