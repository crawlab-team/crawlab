<template>
  <div class="feedback app-container">
    <div class="content">
      <el-card
        class="form"
      >
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

export default {
  name: 'Feedback',
  data () {
    return {
      form: {
        email: '',
        wechat: '',
        content: ''
      },
      isLoading: false
    }
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
            v: sessionStorage.getItem('v')
          })
          if (res && res.data.error) {
            this.$message.error(res.data.error)
            return
          }
          this.form = {
            email: '',
            wechat: '',
            content: ''
          }
          this.$message.success(this.$t('Submitted successfully'))
        } catch (e) {
          this.$message.error(e.toString())
        } finally {
          this.isLoading = false
        }
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
</style>
