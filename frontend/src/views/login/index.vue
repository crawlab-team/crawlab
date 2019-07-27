<template>
  <div class="login-container">
    <canvas id="canvas"></canvas>
    <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" auto-complete="on"
             label-position="left">
      <h3 class="title">
        Crawlab
      </h3>
      <el-form-item prop="username" style="margin-bottom: 28px;">
        <el-input
          v-model="loginForm.username"
          name="username"
          type="text"
          auto-complete="on"
          :placeholder="$t('Username')"
        />
      </el-form-item>
      <el-form-item prop="password" style="margin-bottom: 28px;">
        <el-input
          :type="pwdType"
          v-model="loginForm.password"
          name="password"
          auto-complete="on"
          :placeholder="$t('Password')"
          @keyup.enter.native="onKeyEnter"/>
      </el-form-item>
      <el-form-item v-if="isSignUp" prop="confirmPassword" style="margin-bottom: 28px;">
        <el-input
          :type="pwdType"
          v-model="loginForm.confirmPassword"
          name="password"
          auto-complete="on"
          :placeholder="$t('Confirm Password')"
          @keyup.enter.native="onKeyEnter"
        />
      </el-form-item>
      <el-form-item style="border: none">
        <el-button v-if="isSignUp" :loading="loading" type="primary" style="width:100%;"
                   @click.native.prevent="handleSignup">
          {{$t('Sign up')}}
        </el-button>
        <el-button v-if="!isSignUp" :loading="loading" type="primary" style="width:100%;"
                   @click.native.prevent="handleLogin">
          {{$t('Sign in')}}
        </el-button>
      </el-form-item>
      <div class="alternatives">
        <div class="left">
          <span v-if="!isSignUp" class="forgot-password">{{$t('Forgot Password')}}</span>
        </div>
        <div class="right">
          <span v-if="isSignUp">{{$t('Has Account')}}, </span>
          <span v-if="isSignUp" class="sign-in" @click="$router.push('/login')">{{$t('Sign-in')}} ></span>
          <span v-if="!isSignUp">{{$t('New to Crawlab')}}, </span>
          <span v-if="!isSignUp" class="sign-up" @click="$router.push('/signup')">{{$t('Sign-up')}} ></span>
        </div>
      </div>
      <div class="tips">
        <span>{{$t('Initial Username/Password')}}: admin/admin</span>
        <a href="https://github.com/tikazyq/crawlab" target="_blank" style="float:right">
          <img src="https://img.shields.io/badge/github-crawlab-blue">
        </a>
      </div>
    </el-form>
  </div>
</template>

<script>
import { isValidUsername } from '../../utils/validate'

export default {
  name: 'Login',
  data () {
    const validateUsername = (rule, value, callback) => {
      if (!isValidUsername(value)) {
        callback(new Error(this.$t('Please enter the correct username')))
      } else {
        callback()
      }
    }
    const validatePass = (rule, value, callback) => {
      if (value.length < 5) {
        callback(new Error(this.$t('Password length should be no shorter than 5')))
      } else {
        callback()
      }
    }
    const validateConfirmPass = (rule, value, callback) => {
      if (!this.isSignUp) return callback()
      if (value !== this.loginForm.password) {
        callback(new Error(this.$t('Two passwords must be the same')))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: '',
        password: '',
        confirmPassword: ''
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur', validator: validateUsername }],
        password: [{ required: true, trigger: 'blur', validator: validatePass }],
        confirmPassword: [{ required: true, trigger: 'blur', validator: validateConfirmPass }]
      },
      loading: false,
      pwdType: 'password'
    }
  },
  computed: {
    isSignUp () {
      return this.$route.path === '/signup'
    },
    redirect () {
      return this.$route.query.redirect
    }
  },
  methods: {
    handleLogin () {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          this.$store.dispatch('user/login', this.loginForm).then(() => {
            this.loading = false
            this.$router.push({ path: this.redirect || '/' })
            this.$store.dispatch('user/getInfo')
          }).catch(() => {
            this.$message.error(this.$t('Error when logging in (Please check username and password)'))
            this.loading = false
          })
        }
      })
    },
    handleSignup () {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          this.$store.dispatch('user/register', this.loginForm).then(() => {
            this.handleLogin()
            this.loading = false
          }).catch(err => {
            this.$message.error(this.$t(err))
            this.loading = false
          })
        }
      })
    },
    onKeyEnter () {
      const func = this.isSignUp ? this.handleSignup : this.handleLogin
      func()
    }
  },
  mounted () {
    initCanvas()
  }
}

const initCanvas = () => {
  var canvas = document.getElementById('canvas')
  var ctx = canvas.getContext('2d')

  resize()
  window.onresize = resize

  function resize () {
    canvas.width = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth
    canvas.height = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight
  }

  var RAF = (function () {
    return window.requestAnimationFrame || window.webkitRequestAnimationFrame || window.mozRequestAnimationFrame || window.oRequestAnimationFrame || window.msRequestAnimationFrame || function (callback) {
      window.setTimeout(callback, 1000 / 60)
    }
  })()

  // 鼠标活动时，获取鼠标坐标
  var warea = { x: null, y: null, max: 20000 }
  // window.onmousemove = function (e) {
  //   e = e || window.event
  //
  //   warea.x = e.clientX
  //   warea.y = e.clientY
  // }
  // window.onmouseout = function (e) {
  //   warea.x = null
  //   warea.y = null
  // }

  // 添加粒子
  // x，y为粒子坐标，xa, ya为粒子xy轴加速度，max为连线的最大距离
  var dots = []
  for (var i = 0; i < 300; i++) {
    var x = Math.random() * canvas.width
    var y = Math.random() * canvas.height
    var xa = Math.random() * 2 - 1
    var ya = Math.random() * 2 - 1

    dots.push({
      x: x,
      y: y,
      xa: xa,
      ya: ya,
      max: 6000
    })
  }

  // 延迟100秒开始执行动画，如果立即执行有时位置计算会出错
  setTimeout(function () {
    animate()
  }, 100)

  // 每一帧循环的逻辑
  function animate () {
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // 将鼠标坐标添加进去，产生一个用于比对距离的点数组
    var ndots = [warea].concat(dots)

    dots.forEach(function (dot) {
      // 粒子位移
      dot.x += dot.xa
      dot.y += dot.ya

      // 遇到边界将加速度反向
      dot.xa *= (dot.x > canvas.width || dot.x < 0) ? -1 : 1
      dot.ya *= (dot.y > canvas.height || dot.y < 0) ? -1 : 1

      // 绘制点
      ctx.fillRect(dot.x - 0.5, dot.y - 0.5, 1, 1)

      // 循环比对粒子间的距离
      for (var i = 0; i < ndots.length; i++) {
        var d2 = ndots[i]

        if (dot === d2 || d2.x === null || d2.y === null) continue

        var xc = dot.x - d2.x
        var yc = dot.y - d2.y

        // 两个粒子之间的距离
        var dis = xc * xc + yc * yc

        // 距离比
        var ratio

        // 如果两个粒子之间的距离小于粒子对象的max值，则在两个粒子间画线
        if (dis < d2.max) {
          // 如果是鼠标，则让粒子向鼠标的位置移动
          if (d2 === warea && dis > (d2.max / 2)) {
            dot.x -= xc * 0.03
            dot.y -= yc * 0.03
          }

          // 计算距离比
          ratio = (d2.max - dis) / d2.max

          // 画线
          ctx.beginPath()
          ctx.lineWidth = ratio / 2
          // 线条颜色
          ctx.strokeStyle = 'rgba(64,158,255,' + (ratio + 0.1) + ')'
          ctx.moveTo(dot.x, dot.y)
          ctx.lineTo(d2.x, d2.y)
          ctx.stroke()
        }
      }

      // 将已经计算过的粒子从数组中删除
      ndots.splice(ndots.indexOf(dot), 1)
    })

    RAF(animate)
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
  $bg: #2d3a4b;
  $light_gray: #eee;

  /* reset element-ui css */
  .login-container {
    .el-input {
      display: inline-block;
      width: calc(100% - 44px);
      margin-left: 22px;

      input {
        background: transparent;
        border: 0;
        -webkit-appearance: none;
        border-radius: 0;
        padding: 12px 5px 12px 15px;
        color: #666;
        height: 44px;
        line-height: 44px;
      }
    }

    .el-form-item {
      border: 1px solid #ddd;
      background: #fff;
      border-radius: 22px;
      color: #454545;
      height: 44px;
      /*margin-bottom: 28px;*/

      .el-form-item__content {
        line-height: 44px;
      }
    }

    .el-button {
      height: 44px;
      border-radius: 22px;
    }

    #canvas {
      position: fixed;
      top: 0;
      left: 0;
    }
  }
</style>

<style rel="stylesheet/scss" lang="scss" scoped>
  $bg: transparent;
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
      font-size: 32px;
      color: #666;
      margin: 0px auto 20px auto;
      text-align: center;
      font-weight: bolder;
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
  }
</style>
