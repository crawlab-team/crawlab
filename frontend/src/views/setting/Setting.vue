<template>
  <div class="app-container">
    <el-form :model="userInfo" class="setting-form" ref="setting-form" label-width="150px" :rules="rules"
             inline-message>
      <el-form-item prop="username" :label="$t('Username')">
        <el-input v-model="userInfo.username" disabled></el-input>
      </el-form-item>
      <el-form-item prop="password" :label="$t('Password')">
        <el-input v-model="userInfo.password" type="password" :placeholder="$t('Password')"></el-input>
      </el-form-item>
      <el-form-item prop="email" :label="$t('Email')">
        <el-input v-model="userInfo.email" :placeholder="$t('Email')"></el-input>
      </el-form-item>
      <el-form-item :label="$t('Notification Trigger')">
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
      <el-form-item prop="ding_talk_app_key" :label="$t('DingTalk AppKey')">
        <el-input v-model="userInfo.setting.ding_talk_app_key" :placeholder="$t('DingTalk AppKey')"></el-input>
      </el-form-item>
      <el-form-item prop="ding_talk_app_secret" :label="$t('DingTalk AppSecret')">
        <template v-if="isShowDingTalkAppSecret">
          <el-input
            v-model="userInfo.setting.ding_talk_app_secret"
            :placeholder="$t('DingTalk AppSecret')"
          />

          <i class="icon el-icon-view" @click="toggleDingTalkAppSecret"></i>
        </template>
        <template v-else>
          <el-input
            v-model="userInfo.setting.ding_talk_app_secret"
            type="password"
            :placeholder="$t('DingTalk AppSecret')"
          />
          <i class="icon el-icon-view" @click="toggleDingTalkAppSecret"></i>
        </template>
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
    return {
      userInfo: { setting: {} },
      rules: {
        password: [{ trigger: 'blur', validator: validatePass }],
        email: [{ trigger: 'blur', validator: validateEmail }]
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
    },
    saveUserInfo () {
      this.$refs['setting-form'].validate(async valid => {
        if (!valid) return
        const res = await this.$store.dispatch('user/postInfo', {
          password: this.userInfo.password,
          email: this.userInfo.email,
          notification_trigger: this.userInfo.setting.notification_trigger,
          ding_talk_app_key: this.userInfo.setting.ding_talk_app_key,
          ding_talk_app_secret: this.userInfo.setting.ding_talk_app_secret
        })
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
</style>
