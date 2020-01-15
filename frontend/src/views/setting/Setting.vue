<template>
  <div class="app-container">
    <el-form :model="userInfo" class="setting-form" ref="setting-form" label-width="200px" :rules="rules"
             inline-message>
      <el-form-item prop="username" :label="$t('Username')">
        <el-input v-model="userInfo.username" disabled></el-input>
      </el-form-item>
      <el-form-item prop="password" :label="$t('Password')">
        <el-input v-model="userInfo.password" type="password" :placeholder="$t('Password')"></el-input>
      </el-form-item>
      <div style="border-bottom: 1px solid #DCDFE6"></div>
      <el-form-item :label="$t('Notification Trigger Timing')">
        <el-radio-group v-model="userInfo.setting.notification_trigger">
          <el-radio label="notification_trigger_on_task_end">
            {{$t('On Task End')}}
          </el-radio>
          <el-radio label="notification_trigger_on_task_error">
            {{$t('On Task Error')}}
          </el-radio>
          <el-radio label="notification_trigger_never">
            {{$t('Never')}}
          </el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item prop="enabledNotifications" :label="$t('消息通知方式')">
        <el-checkbox-group v-model="userInfo.setting.enabled_notifications">
          <el-checkbox label="notification_type_mail">{{$t('邮件')}}</el-checkbox>
          <el-checkbox label="notification_type_ding_talk">{{$t('钉钉')}}</el-checkbox>
          <el-checkbox label="notification_type_wechat">{{$t('企业微信')}}</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item prop="email" :label="$t('Email')">
        <el-input v-model="userInfo.email" :placeholder="$t('Email')"></el-input>
      </el-form-item>
      <el-form-item prop="setting.ding_talk_robot_webhook" :label="$t('DingTalk Robot Webhook')">
        <el-input v-model="userInfo.setting.ding_talk_robot_webhook"
                  :placeholder="$t('DingTalk Robot Webhook')"></el-input>
      </el-form-item>
      <el-form-item prop="setting.wechat_robot_webhook" :label="$t('Wechat Robot Webhook')">
        <el-input v-model="userInfo.setting.wechat_robot_webhook" :placeholder="$t('Wechat Robot Webhook')"></el-input>
      </el-form-item>
      <el-form-item>
        <div class="buttons">
          <el-button type="success" @click="saveUserInfo">{{$t('Save')}}</el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: 'Setting',
  data () {
    const validatePass = (rule, value, callback) => {
      if (!value) return callback()
      if (value.length < 5) {
        callback(new Error(this.$t('Password length should be no shorter than 5')))
      } else {
        callback()
      }
    }
    const validateEmail = (rule, value, callback) => {
      if (!value) return callback()
      if (!value.match(/.+@.+/i)) {
        callback(new Error(this.$t('Email format invalid')))
      } else {
        callback()
      }
    }
    const validateDingTalkRobotWebhook = (rule, value, callback) => {
      if (!value) return callback()
      if (!value.match(/^https:\/\/oapi.dingtalk.com\/robot\/send\?access_token=[a-f0-9]+/i)) {
        callback(new Error(this.$t('DingTalk Robot Webhook format invalid')))
      } else {
        callback()
      }
    }
    const validateWechatRobotWebhook = (rule, value, callback) => {
      if (!value) return callback()
      if (!value.match(/^https:\/\/qyapi.weixin.qq.com\/cgi-bin\/webhook\/send\?key=.+/i)) {
        callback(new Error(this.$t('DingTalk Robot Webhook format invalid')))
      } else {
        callback()
      }
    }
    return {
      userInfo: { setting: { enabled_notifications: [] } },
      rules: {
        password: [{ trigger: 'blur', validator: validatePass }],
        email: [{ trigger: 'blur', validator: validateEmail }],
        'setting.ding_talk_robot_webhook': [{ trigger: 'blur', validator: validateDingTalkRobotWebhook }],
        'setting.wechat_robot_webhook': [{ trigger: 'blur', validator: validateWechatRobotWebhook }]
      },
      isShowDingTalkAppSecret: false
    }
  },
  methods: {
    getUserInfo () {
      const data = localStorage.getItem('user_info')
      if (!data) return {}
      this.userInfo = JSON.parse(data)
      if (!this.userInfo.setting) this.userInfo.setting = {}
      if (!this.userInfo.setting.enabled_notifications) this.userInfo.setting.enabled_notifications = []
    },
    saveUserInfo () {
      this.$refs['setting-form'].validate(async valid => {
        if (!valid) return
        const res = await this.$store.dispatch('user/postInfo', this.userInfo)
        if (!res || res.error) {
          this.$message.error(res.error)
        } else {
          this.$message.success(this.$t('Saved successfully'))
        }
      })
    },
    toggleDingTalkAppSecret () {
      this.isShowDingTalkAppSecret = !this.isShowDingTalkAppSecret
    }
  },
  async created () {
    await this.$store.dispatch('user/getInfo')
    this.getUserInfo()
  }
}
</script>

<style scoped>
  .setting-form {
    width: 600px;
  }

  .setting-form .buttons {
    text-align: right;
  }

  .setting-form .icon {
    top: calc(50% - 14px / 2);
    right: 14px;
    position: absolute;
    color: #DCDFE6;
  }

  .setting-form >>> .el-form-item__label {
    height: 40px;
  }
</style>
