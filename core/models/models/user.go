package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password,omitempty" bson:"-"`
	Role     string             `json:"role" bson:"role"`
	Email    string             `json:"email" bson:"email"`
	//Setting  UserSetting        `json:"setting" bson:"setting"`
}

func (u *User) GetId() (id primitive.ObjectID) {
	return u.Id
}

func (u *User) SetId(id primitive.ObjectID) {
	u.Id = id
}

func (u *User) GetUsername() (name string) {
	return u.Username
}

func (u *User) GetPassword() (p string) {
	return u.Password
}

func (u *User) GetRole() (r string) {
	return u.Role
}

func (u *User) GetEmail() (email string) {
	return u.Email
}

//type UserSetting struct {
//	NotificationTrigger  string   `json:"notification_trigger" bson:"notification_trigger"`
//	DingTalkRobotWebhook string   `json:"ding_talk_robot_webhook" bson:"ding_talk_robot_webhook"`
//	WechatRobotWebhook   string   `json:"wechat_robot_webhook" bson:"wechat_robot_webhook"`
//	EnabledNotifications []string `json:"enabled_notifications" bson:"enabled_notifications"`
//	ErrorRegexPattern    string   `json:"error_regex_pattern" bson:"error_regex_pattern"`
//	MaxErrorLog          int      `json:"max_error_log" bson:"max_error_log"`
//	LogExpireDuration    int64    `json:"log_expire_duration" bson:"log_expire_duration"`
//}

type UserList []User

func (l *UserList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
