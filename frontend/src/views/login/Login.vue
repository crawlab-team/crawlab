<template>
  <div class="login-container">
    <canvas id="canvas"/>
    <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        auto-complete="on"
        class="login-form"
        label-position="left"
    >
      <h3 class="title">
        <img :src="logo" alt="logo" className="logo-img"/>
        <span class="logo-title">Crawlab</span>
        <span class="logo-sub-title">
          <div class="logo-sub-title-block">
            Community
          </div>
          <div class="logo-sub-title-block">
            v0.6.0
          </div>
        </span>
      </h3>
      <el-form-item prop="username" style="margin-bottom: 28px;">
        <el-input
            v-model="loginForm.username"
            :placeholder="$t('Username')"
            auto-complete="on"
            name="username"
            type="text"
            @keyup.enter="onLogin"
        />
      </el-form-item>
      <el-form-item prop="password" style="margin-bottom: 28px;">
        <el-input
            v-model="loginForm.password"
            :placeholder="$t('Password')"
            auto-complete="on"
            name="password"
            type="password"
            @keyup.enter="onLogin"
        />
      </el-form-item>
      <el-form-item v-if="isSignup" prop="confirmPassword" style="margin-bottom: 28px;">
        <el-input
            v-model="loginForm.confirmPassword"
            :placeholder="$t('Confirm Password')"
            auto-complete="on"
            name="password"
        />
      </el-form-item>
      <el-form-item v-if="isSignup" prop="email" style="margin-bottom: 28px;">
        <el-input
            v-model="loginForm.email"
            :placeholder="$t('Email')"
            name="email"
        />
      </el-form-item>
      <el-form-item style="border: none">
        <el-button
            v-if="isSignup"
            :loading="loading"
            style="width:100%;"
            type="primary"
        >
          {{ $t('Sign up') }}
        </el-button>
        <el-button
            v-if="!isSignup"
            :loading="loading"
            style="width:100%;"
            type="primary"
            @click="onLogin"
        >
          {{ $t('Sign in') }}
        </el-button>
      </el-form-item>
      <div class="alternatives">
        <div class="left">
          <el-tooltip content="Please follow the Reset Password section in the documentation." trigger="click">
            <span class="forgot-password">{{ $t('Forgot Password') }}</span>
          </el-tooltip>
        </div>
        <div class="right" v-if="allowRegister">
          <span v-if="isSignup">{{ $t('Has Account') }}, </span>
          <span v-if="isSignup" class="sign-in" @click="$router.push('/login')">{{ $t('Sign-in') }} ></span>
          <span v-if="!isSignup">{{ $t('New to Crawlab') }}, </span>
          <span v-if="!isSignup" class="sign-up" @click="$router.push('/signup')">{{ $t('Sign-up') }} ></span>
        </div>
      </div>
      <div class="tips">
        <span>{{ $t('Initial Username/Password') }}: admin/admin</span>
        <a href="https://github.com/crawlab-team/crawlab" style="float:right" target="_blank">
          <img alt="github-stars" src="https://img.shields.io/github/stars/crawlab-team/crawlab?logo=github">
        </a>
      </div>
      <!-- TODO: implement -->
      <div v-if="false" class="lang">
        <span :class="lang==='zh'?'active':''" @click="setLang('zh')">中文</span>
        |
        <span :class="lang==='en'?'active':''" @click="setLang('en')">English</span>
      </div>
      <div v-if="false" class="documentation">
        <a href="https://docs.crawlab.cn" target="_blank">{{ $t('Documentation') }}</a>
      </div>
      <div class="mobile-warning" v-if="isShowMobileWarning">
        <el-alert :closable="false" type="error">
          {{
            $t('You are running on a mobile device, which is not optimized yet. Please try with a laptop or desktop.')
          }}
        </el-alert>
      </div>
    </el-form>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, onUnmounted, ref} from 'vue';
import {isValidUsername} from '@/utils/validate';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import logo from '@/assets/logo.svg';
import {ElForm, ElMessage} from 'element-plus';
import useRequest from '@/services/request';

const {
  post,
} = useRequest();

export default defineComponent({
  name: 'Login',
  setup() {
    const route = useRoute();

    const router = useRouter();

    const loading = ref<boolean>(false);

    const isSignup = computed(() => route.path === '/signup');

    const loginForm = ref<LoginForm>({
      username: '',
      password: '',
      confirmPassword: '',
      email: '',
    });

    const loginFormRef = ref<typeof ElForm>();

    const validateUsername = (rule: any, value: any, callback: any) => {
      if (!isValidUsername(value)) {
        callback(new Error('Please enter the correct username'));
      } else {
        callback();
      }
    };

    const validatePass = (rule: any, value: any, callback: any) => {
      if (value.length < 5) {
        callback(new Error('Password length should be no shorter than 5'));
      } else {
        callback();
      }
    };

    const validateConfirmPass = (rule: any, value: any, callback: any) => {
      if (!isSignup.value) return callback();
      if (value !== loginForm.value.password) {
        callback(new Error('Two passwords must be the same'));
      } else {
        callback();
      }
    };

    const loginRules: LoginRules = {
      username: [{required: true, trigger: 'blur', validator: validateUsername}],
      password: [{required: true, trigger: 'blur', validator: validatePass}],
      confirmPassword: [{required: true, trigger: 'blur', validator: validateConfirmPass}]
    };

    const isShowMobileWarning = ref<boolean>(false);

    const allowRegister = ref<boolean>(false);

    const lang = computed<string | null>(() => localStorage.getItem('lang'));

    const setLang = (lang: string) => {
      localStorage.setItem('lang', lang);
    };

    const onLogin = async () => {
      if (!loginFormRef.value) return;
      await loginFormRef.value.validate();
      const {username, password} = loginForm.value;
      loading.value = true;
      try {
        const res = await post<LoginForm, ResponseWithData>('/login', {
          username,
          password,
        });
        if (!res.data) {
          ElMessage.error('No token returned');
          return;
        }
        localStorage.setItem('token', res.data);
        await router.push('/');
      } catch (e) {
        if (e.toString().includes('401')) {
          ElMessage.error('Unauthorized. Please check username and password.');
        } else {
          ElMessage.error(e.toString());
        }
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      if (window.innerWidth >= 1024) {
        if (!window.initCanvas) {
          import('../../assets/js/loginCanvas.js');
        } else {
          window.initCanvas();
        }
      } else {
        isShowMobileWarning.value = true;
      }
    });
    onUnmounted(() => {
      if (window.resetCanvas) {
        window.resetCanvas();
      }
    });

    return {
      loginForm,
      loginFormRef,
      loginRules,
      isShowMobileWarning,
      allowRegister,
      isSignup,
      loading,
      lang,
      logo,
      setLang,
      onLogin,
    };
  }
});
</script>

<style lang="scss" rel="stylesheet/scss" scoped>
@import "../../styles/variables.scss";

$bg: white;
$dark_gray: #889aa4;
$light_gray: #aaa;
.login-container {
  position: fixed;
  height: 100%;
  width: 100%;
  background-color: $bg;

  .login-form {
    background: transparent;
    position: absolute;
    left: 0;
    right: 0;
    width: 480px;
    max-width: 100%;
    padding: 35px 35px 15px 35px;
    margin: 120px auto;
  }

  .tips {
    font-size: 14px;
    color: #666;
    margin-bottom: 10px;
    background: transparent;

    span {
      &:first-of-type {
        margin-right: 22px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title {
    font-family: "Verdana", serif;
    /*font-style: italic;*/
    font-weight: 600;
    font-size: 24px;
    color: #409EFF;
    margin: 0px auto 20px auto;
    text-align: center;
    cursor: default;

    display: flex;
    align-items: center;
    height: 128px;

    .logo-img {
      height: 80px;
    }

    .logo-title {
      font-family: BlinkMacSystemFont, -apple-system, segoe ui, roboto, oxygen, ubuntu, cantarell, fira sans, droid sans, helvetica neue, helvetica, arial, sans-serif;
      font-size: 56px;
      font-weight: 600;
      margin-left: 24px;
      color: #409eff;
    }

    .logo-sub-title {
      font-family: BlinkMacSystemFont, -apple-system, segoe ui, roboto, oxygen, ubuntu, cantarell, fira sans, droid sans, helvetica neue, helvetica, arial, sans-serif;
      font-size: 20px;
      height: 48px;
      line-height: 48px;
      margin-left: 20px;
      font-weight: 500;
      color: $infoMediumColor;
      opacity: 0.8;

      .logo-sub-title-block {
        display: flex;
        align-items: center;
        height: 24px;
        line-height: 24px;
      }
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }

  .alternatives {
    border-bottom: 1px solid #ccc;
    display: flex;
    justify-content: space-between;
    font-size: 14px;
    color: #666;
    font-weight: 400;
    margin-bottom: 10px;
    padding-bottom: 10px;

    .forgot-password {
      cursor: pointer;
    }

    .sign-in,
    .sign-up {
      cursor: pointer;
      color: #409EFF;
      font-weight: 600;
    }
  }

  .lang {
    margin-top: 20px;
    text-align: center;
    color: #666;

    span {
      cursor: pointer;
      margin: 10px;
      font-size: 14px;
    }

    span.active {
      font-weight: 600;
      text-decoration: underline;
    }

    span:hover {
      text-decoration: underline;
    }
  }

  .documentation {
    margin-top: 20px;
    text-align: center;
    font-size: 14px;
    color: #409eff;
    font-weight: bolder;

    &:hover {
      text-decoration: underline;
    }
  }

  .mobile-warning {
    margin-top: 20px;
  }

}
</style>
<style scoped>
.mobile-warning >>> .el-alert .el-alert__description {
  font-size: 1.2rem;
}

#canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
}
</style>
