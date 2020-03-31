<template>
  <div class="feedback app-container">
    <div class="content">
      <el-card
        class="form"
      >
        <el-alert
          type="info"
          effect="light"
          class="notice"
          :closable="false"
        >
          <template v-if="lang === 'zh'">
            <strong>您的反馈意见对我们优化产品非常重要！</strong><br>
            您可以在这里畅所欲言，提供您的建议和我们需要完善提升的地方。<br>
            您可以选择留下您的联系方式，方便我们进一步了解您的使用情况。
          </template>
          <template v-else>
            <strong>Your feedback is very important for us to improve the product!</strong><br>
            You can comment anything here and provide any suggestions and what we should enhance about.<br>
            You can leave your contact info here for us to get better understanding about how you are using our product.
          </template>
        </el-alert>
        <el-form
          ref="form"
          v-model="form"
          :model="form"
          label-width="80px"
        >
          <el-form-item
            :label="$t('Email')"
            prop="email"
          >
            <el-input
              v-model="form.email"
              :placeholder="$t('Please enter your email')"
            />
          </el-form-item>
          <el-form-item
            :label="$t('Wechat')"
            prop="wechat"
          >
            <el-input
              v-model="form.wechat"
              :placeholder="$t('Please enter your Wechat account')"
            />
          </el-form-item>
          <el-form-item
            :label="$t('Content')"
            prop="content"
            required
          >
            <el-input
              type="textarea"
              rows="5"
              v-model="form.content"
              :placeholder="$t('Please enter your feedback content')"
            />
          </el-form-item>
          <el-form-item
            class="rating"
            :label="$t('Rating')"
            prop="rating"
            required
          >
            <el-rate
              v-model="form.rating"
            />
          </el-form-item>
          <el-form-item>
            <div class="actions">
              <el-button
                type="primary"
                size="small"
                :icon="isLoading ? 'el-icon-loading' : ''"
                :disabled="isLoading"
                @click="submit"
              >
                {{$t('Submit')}}
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import {
  mapState
} from 'vuex'

export default {
  name: 'Feedback',
  data () {
    return {
      form: {
        email: '',
        wechat: '',
        content: '',
        rating: 0
      },
      isLoading: false
    }
  },
  computed: {
    ...mapState('lang', [
      'lang'
    ])
  },
  methods: {
    submit () {
      this.$refs['form'].validate(async valid => {
        if (!valid) return
        this.isLoading = true
        try {
          const res = await axios.put(process.env.VUE_APP_CRAWLAB_BASE_URL + '/feedback', {
            uid: localStorage.getItem('uid'),
            sid: sessionStorage.getItem('sid'),
            email: this.form.email,
            wechat: this.form.wechat,
            content: this.form.content,
            rating: this.form.rating,
            v: sessionStorage.getItem('v')
          })
          if (res && res.data.error) {
            this.$message.error(res.data.error)
            return
          }
          this.form = {
            email: '',
            wechat: '',
            content: '',
            rating: 0
          }
          this.$message.success(this.$t('Submitted successfully'))
        } catch (e) {
          this.$message.error(e.toString())
        } finally {
          this.isLoading = false
        }
        this.$st.sendEv('反馈', '提交反馈')
      })
    }
  }
}
</script>

<style scoped>
  .content {
    width: 900px;
    margin-left: calc(50% - 450px);
  }

  .actions {
    text-align: right;
  }

  .rating >>> .el-form-item__content {
    display: flex;
    align-items: center;
    height: 40px;
  }

  .notice {
    margin-bottom: 20px;
  }

  .notice >>> .el-alert__description {
    line-height: 24px;
    font-size: 16px;
  }
</style>
