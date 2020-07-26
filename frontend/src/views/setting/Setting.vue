<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="setting"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <!--新增全局变量-->
    <el-dialog
      :title="$t('Add Global Variable')"
      :visible.sync="addDialogVisible"
    >
      <el-form ref="globalVariableForm" label-width="80px">
        <el-form-item :label="$t('Key')">
          <el-input v-model="globalVariableForm.key" size="small" />
        </el-form-item>
        <el-form-item :label="$t('Value')">
          <el-input v-model="globalVariableForm.value" size="small" />
        </el-form-item>
        <el-form-item :label="$t('Remark')">
          <el-input v-model="globalVariableForm.remark" size="small" />
        </el-form-item>
        <el-form-item>
          <div style="text-align: right">
            <el-button type="danger" size="small" @click="addDialogVisible = false">{{ $t('Cancel') }}</el-button>
            <el-button type="success" size="small" @click="addGlobalVariableHandle(false)">{{ $t('Save') }}</el-button>
          </div>
        </el-form-item>
      </el-form>
    </el-dialog>
    <!--./新增全局变量-->

    <el-tabs v-model="activeName" type="border-card" @tab-click="tabActiveHandle">
      <!--通用-->
      <el-tab-pane :label="$t('General')" name="general">
        <el-form
          ref="setting-form"
          :model="userInfo"
          class="setting-form"
          label-width="200px"
          :rules="rulesNotification"
          inline-message
        >
          <el-form-item prop="username" :label="$t('Username')">
            <el-input v-model="userInfo.username" disabled />
          </el-form-item>
          <el-form-item :label="$t('Allow Sending Statistics')">
            <el-switch
              v-model="isAllowSendingStatistics"
              active-color="#67C23A"
              inactive-color="#909399"
              @change="onAllowSendingStatisticsChange"
            />
          </el-form-item>
          <el-form-item :label="$t('Enable Tutorial')">
            <el-switch
              v-model="isEnableTutorial"
              active-color="#67C23A"
              inactive-color="#909399"
              @change="onEnableTutorialChange"
            />
          </el-form-item>
          <el-form-item>
            <div style="text-align: right">
              <el-button type="success" size="small" @click="saveUserInfo">
                {{ $t('Save') }}
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <!--./通用-->

      <!--更改密码-->
      <el-tab-pane :label="$t('Change Password')" name="change-password">
        <el-form
          ref="change-password-form"
          :model="userInfo"
          class="change-password-form"
          label-width="200px"
          inline-message
        >
          <el-form-item prop="password" :label="$t('Password')" required>
            <el-input v-model="userInfo.password" type="password" :placeholder="$t('Password')" />
          </el-form-item>
          <el-form-item prop="confirm_password" :label="$t('Confirm Password')" required>
            <el-input v-model="userInfo.confirm_password" type="password" :placeholder="$t('Confirm Password')" />
          </el-form-item>
          <el-form-item>
            <div style="text-align: right">
              <el-button type="success" size="small" @click="changePassword">
                {{ $t('Save') }}
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <!--./更改密码-->

      <!--消息通知-->
      <el-tab-pane :label="$t('Notifications')" name="notify">
        <el-form
          ref="setting-form"
          :model="userInfo"
          class="setting-form"
          label-width="200px"
          :rules="rulesNotification"
          inline-message
        >
          <el-form-item :label="$t('Notification Trigger Timing')">
            <el-radio-group v-model="userInfo.setting.notification_trigger">
              <el-radio label="notification_trigger_on_task_end">
                {{ $t('On Task End') }}
              </el-radio>
              <el-radio label="notification_trigger_on_task_error">
                {{ $t('On Task Error') }}
              </el-radio>
              <el-radio label="notification_trigger_never">
                {{ $t('Never') }}
              </el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item prop="enabledNotifications" :label="$t('消息通知方式')">
            <el-checkbox-group v-model="userInfo.setting.enabled_notifications">
              <el-checkbox label="notification_type_mail">{{ $t('邮件') }}</el-checkbox>
              <el-checkbox label="notification_type_ding_talk">{{ $t('钉钉') }}</el-checkbox>
              <el-checkbox label="notification_type_wechat">{{ $t('企业微信') }}</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <el-form-item prop="email" :label="$t('Email')">
            <el-input v-model="userInfo.email" :placeholder="$t('Email')" />
          </el-form-item>
          <el-form-item prop="setting.ding_talk_robot_webhook" :label="$t('DingTalk Robot Webhook')">
            <el-input
              v-model="userInfo.setting.ding_talk_robot_webhook"
              :placeholder="$t('DingTalk Robot Webhook')"
            />
          </el-form-item>
          <el-form-item prop="setting.wechat_robot_webhook" :label="$t('Wechat Robot Webhook')">
            <el-input
              v-model="userInfo.setting.wechat_robot_webhook"
              :placeholder="$t('Wechat Robot Webhook')"
            />
          </el-form-item>
          <el-form-item>
            <div style="text-align: right">
              <el-button type="success" size="small" @click="saveUserInfo">
                {{ $t('Save') }}
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <!--./消息通知-->

      <!--日志-->
      <el-tab-pane :label="$t('Log')" name="log">
        <el-form
          ref="log-form"
          :model="userInfo"
          class="setting-form"
          label-width="200px"
          :rules="rulesLog"
          inline-message
        >
          <el-form-item :label="$t('Error Regex Pattern')" prop="setting.error_regex_pattern">
            <el-input
              v-model="userInfo.setting.error_regex_pattern"
              :placeholder="$t('By default: ') + $utils.log.errorRegex.source"
              clearable
            />
          </el-form-item>
          <el-form-item :label="$t('Max Error Logs Display')" prop="setting.max_error_log">
            <el-select
              v-model="userInfo.setting.max_error_log"
              clearable
            >
              <el-option :value="100" label="100" />
              <el-option :value="500" label="500" />
              <el-option :value="1000" label="1000" />
              <el-option :value="5000" label="5000" />
              <el-option :value="10000" label="10000" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('Log Expire Duration')" prop="setting.log_expire_duration">
            <el-select
              v-model="userInfo.setting.log_expire_duration"
              clearable
            >
              <el-option :value="0" :label="$t('No Expire')" />
              <el-option :value="3600" :label="'1 ' + $t('Hour')" />
              <el-option :value="3600 * 6" :label="'6 ' + $t('Hours')" />
              <el-option :value="3600 * 12" :label="'12 ' + $t('Hours')" />
              <el-option :value="3600 * 24" :label="'1 ' + $t('Day')" />
              <el-option :value="3600 * 24 * 7" :label="'7 ' + $t('Days')" />
              <el-option :value="3600 * 24 * 14" :label="'14 ' + $t('Days')" />
              <el-option :value="3600 * 24 * 30" :label="'30 ' + $t('Days')" />
              <el-option :value="3600 * 24 * 30 * 3" :label="'90 ' + $t('Days')" />
              <el-option :value="3600 * 24 * 30 * 6" :label="'180 ' + $t('Days')" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <div style="text-align: right">
              <el-button type="success" size="small" @click="saveUserInfo">
                {{ $t('Save') }}
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <!--./日志-->

      <!--API Token-->
      <el-tab-pane label="API Token" name="api-token">
        <input id="clipboard">
        <el-alert
          type="primary"
        />
        <div class="actions">
          <el-button
            size="small"
            type="primary"
            @click="onAddApiToken"
          >
            {{ $t('Add') }}
          </el-button>
        </div>
        <el-table
          :data="apiTokens"
          border
        >
          <el-table-column
            label="Token"
          >
            <template slot-scope="scope">
              {{ scope.row.visible ? scope.row.token : getMaskValue(scope.row.token) }}
            </template>
          </el-table-column>
          <el-table-column
            :label="$t('Action')"
            width="200px"
          >
            <template slot-scope="scope">
              <el-button
                type="warning"
                size="mini"
                icon="el-icon-view"
                circle
                @click="toggleTokenVisible(scope.row)"
              />
              <el-button
                type="primary"
                size="mini"
                icon="el-icon-document-copy"
                circle
                @click="copyToken(scope.row.token)"
              />
              <el-button
                type="danger"
                size="mini"
                icon="el-icon-delete"
                circle
                @click="onDeleteToken(scope.row)"
              />
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <!--./API Token-->

      <!--全局变量-->
      <el-tab-pane :label="$t('Global Variable')" name="global-variable">
        <div style="text-align: right;margin-bottom: 10px">
          <el-button
            size="small"
            icon="el-icon-plus"
            type="primary"
            @click="addGlobalVariableHandle(true)"
          >
            {{ $t('Add') }}
          </el-button>
          <el-button size="small" type="success" @click="saveUserInfo">{{ $t('Save') }}</el-button>
        </div>
        <el-table :data="globalVariableList" border>
          <el-table-column prop="key" :label="$t('Key')" />
          <el-table-column prop="value" :label="$t('Value')" />
          <el-table-column prop="remark" :label="$t('Remark')" />
          <el-table-column prop="" :label="$t('Action')" width="80">
            <template slot-scope="scope">
              <el-button
                icon="el-icon-delete"
                type="danger"
                size="mini"
                @click="deleteGlobalVariableHandle(scope.row._id)"
              />
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <!--./全局变量-->
    </el-tabs>
  </div>
</template>

<script>
  import { mapState } from 'vuex'

  export default {
    name: 'Setting',
    data() {
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
        userInfo: { setting: { enabled_notifications: [] }},
        rulesNotification: {
          password: [{ trigger: 'blur', validator: validatePass }],
          email: [{ trigger: 'blur', validator: validateEmail }],
          'setting.ding_talk_robot_webhook': [{ trigger: 'blur', validator: validateDingTalkRobotWebhook }],
          'setting.wechat_robot_webhook': [{ trigger: 'blur', validator: validateWechatRobotWebhook }]
        },
        rulesLog: {},
        isShowDingTalkAppSecret: false,
        activeName: 'general',
        addDialogVisible: false,
        tourSteps: [
          {
            target: '#tab-general',
            content: this.$t('Here you can set your general settings.'),
            params: {
              placement: 'right'
            }
          },
          {
            target: '#tab-notify',
            content: this.$t('In this tab you can configure your notification settings.')
          },
          {
            target: '#tab-global-variable',
            content: this.$t('Here you can add/edit/delete global environment variables which will be passed into your spider programs.')
          }
        ],
        tourCallbacks: {
          onStop: () => {
            this.$utils.tour.finishTour('setting')
          },
          onPreviousStep: (currentStep) => {
            if (currentStep === 1) {
              this.activeName = 'password'
            } else if (currentStep === 2) {
              this.activeName = 'notify'
            }
            this.$utils.tour.prevStep('setting', currentStep)
          },
          onNextStep: (currentStep) => {
            if (currentStep === 0) {
              this.activeName = 'notify'
            } else if (currentStep === 1) {
              this.activeName = 'global-variable'
            }
            this.$utils.tour.nextStep('setting', currentStep)
          }
        },
        isAllowSendingStatistics: localStorage.getItem('useStats') === '1',
        isEnableTutorial: localStorage.getItem('enableTutorial') === '1',
        apiTokens: []
      }
    },
    computed: {
      ...mapState('user', [
        'globalVariableList',
        'globalVariableForm'
      ])
    },
    watch: {
      async userInfoStr() {
        await this.saveUserInfo()
        await this.$store.dispatch('user/getInfo')
      }
    },
    async created() {
      await this.$store.dispatch('user/getInfo')
      await this.$store.dispatch('user/getGlobalVariable')
      this.getUserInfo()
      await this.getApiTokens()
    },
    mounted() {
      if (!this.$utils.tour.isFinishedTour('setting')) {
        this.$utils.tour.startTour(this, 'setting')
      }
    },
    methods: {
      deleteGlobalVariableHandle(id) {
        this.$confirm(this.$t('Are you sure to delete this global variable'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('user/deleteGlobalVariable', id).then(() => {
            this.$store.dispatch('user/getGlobalVariable')
          })
        }).catch(() => {
        })
      },
      addGlobalVariableHandle(isShow) {
        if (isShow) {
          this.addDialogVisible = true
          return
        }
        this.$store.dispatch('user/addGlobalVariable')
          .then(() => {
            this.addDialogVisible = false
            this.$st.sendEv('设置', '添加全局变量')
          })
          .then(() => {
            this.$store.dispatch('user/getGlobalVariable')
          })
      },
      getUserInfo() {
        const data = localStorage.getItem('user_info')
        if (!data) {
          return {}
        }
        this.userInfo = JSON.parse(data)
        if (!this.userInfo.setting) {
          this.userInfo.setting = {}
        }
        if (!this.userInfo.setting.enabled_notifications) {
          this.userInfo.setting.enabled_notifications = []
        }
      },
      saveUserInfo() {
        this.$refs['setting-form'].validate(async valid => {
          if (!valid) return
          const res = await this.$store.dispatch('user/postInfo', this.userInfo)
          if (!res || res.error) {
            this.$message.error(res.error)
          } else {
            this.$message.success(this.$t('Saved successfully'))
          }
        })
        this.$st.sendEv('设置', '保存')
      },
      tabActiveHandle() {
      },
      onAllowSendingStatisticsChange(value) {
        if (value) {
          this.$st.sendPv('/allow_stats')
          this.$st.sendEv('全局', '允许/禁止统计', '允许')
        } else {
          this.$st.sendPv('/disallow_stats')
          this.$st.sendEv('全局', '允许/禁止统计', '禁止')
        }
        this.$message.success(this.$t('Saved successfully'))
        localStorage.setItem('useStats', value ? '1' : '0')
      },
      onEnableTutorialChange(value) {
        this.$message.success(this.$t('Saved successfully'))
        localStorage.setItem('enableTutorial', value ? '1' : '0')
      },
      onAddApiToken() {
        this.$confirm(this.$t('Are you sure to add an API token?'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(async() => {
          const res = await this.$request.put('/tokens')
          if (!res.data.error) {
            this.$message.success(this.$t('Added API token successfully'))
            await this.getApiTokens()
          }
        })
      },
      onDeleteToken(row) {
        this.$confirm(this.$t('Are you sure to delete this API token?'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(async() => {
          const res = await this.$request.delete(`/tokens/${row._id}`)
          if (!res.data.error) {
            this.$message.success(this.$t('Deleted API token successfully'))
            await this.getApiTokens()
          }
        })
      },
      async addApiToken() {
        await this.$request.put('/tokens')
      },
      async getApiTokens() {
        const res = await this.$request.get('/tokens')
        this.apiTokens = res.data.data
      },
      toggleTokenVisible(row) {
        this.$set(row, 'visible', !row.visible)
      },
      getMaskValue(str) {
        let s = ''
        for (let i = 0; i < str.length; i++) {
          s += '*'
        }
        return s
      },
      copyToken(str) {
        const input = document.getElementById('clipboard')
        input.value = str
        input.select()
        document.execCommand('copy')
        this.$message.success(this.$t('Token copied'))
      },
      changePassword() {
        this.$refs['change-password-form'].validate(async valid => {
          if (!valid) return
          if (this.userInfo.password !== this.userInfo.confirm_password) {
            this.$message.error(this.$t('Two passwords do not match'))
            return
          }
          if (this.userInfo.password.length < 5) {
            this.$message.error(this.$t('Password length should be no shorter than 5'))
            return
          }
          const res = await this.$request.post(`/me/change-password`, {
            password: this.userInfo.password
          })
          if (!res.data.error) {
            this.$message.success(this.$t('Changed password successfully'))
          }
        })
      }
    }
  }
</script>

<style scoped>
  .setting-form,
  .change-password-form {
    width: 600px;
  }

  .setting-form .buttons,
  .change-password-form .buttons {
    text-align: right;
  }

  .setting-form .icon,
  .change-password-form .icon {
    top: calc(50% - 14px / 2);
    right: 14px;
    position: absolute;
    color: #DCDFE6;
  }

  .setting-form >>> .el-form-item__label,
  .change-password-form >>> .el-form-item__label {
    height: 40px;
  }

  .actions {
    margin-bottom: 10px;
    text-align: right;
  }

  #clipboard {
    position: fixed;
    z-index: -99999;
    top: 9999px;
    right: 9999px;
  }
</style>
