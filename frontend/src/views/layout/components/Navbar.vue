<template>
  <div class="navbar">
    <hamburger :toggle-click="toggleSideBar" :is-active="sidebar.opened" class="hamburger-container"/>
    <breadcrumb class="breadcrumb"/>
    <el-dropdown class="avatar-container" trigger="click">
      <span class="el-dropdown-link">
        {{$t('User')}}
        <i class="el-icon-arrow-down el-icon--right"></i>
      </span>
      <el-dropdown-menu slot="dropdown" class="user-dropdown">
        <el-dropdown-item>
          <span style="display:block;" @click="logout">{{$t('Logout')}}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    <el-dropdown class="lang-list" trigger="click">
      <span class="el-dropdown-link">
        {{$t($store.getters['lang/lang'])}}
        <i class="el-icon-arrow-down el-icon--right"></i>
      </span>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item @click.native="setLang('en')">
          <span>English</span>
        </el-dropdown-item>
        <el-dropdown-item @click.native="setLang('zh')">
          <span>中文</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'

export default {
  components: {
    Breadcrumb,
    Hamburger
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'avatar'
    ])
  },
  methods: {
    toggleSideBar () {
      this.$store.dispatch('ToggleSideBar')
    },
    logout () {
      this.$router.push('/login')
    },
    setLang (lang) {
      window.localStorage.setItem('lang', lang)
      this.$i18n.locale = lang
      this.$store.commit('lang/SET_LANG', lang)

      this.$st.sendEv('全局', '切换中英文', 'lang', lang)
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .navbar {
    height: 50px;
    line-height: 50px;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .12), 0 0 3px 0 rgba(0, 0, 0, .04);

    .hamburger-container {
      line-height: 58px;
      height: 50px;
      float: left;
      padding: 0 10px;
    }

    .lang-list {
      cursor: pointer;
      display: inline-block;
      float: right;
      margin-right: 35px;
      /*position: absolute;*/
      /*right: 80px;*/
    }

    .screenfull {
      position: absolute;
      right: 90px;
      top: 16px;
      color: red;
    }

    .avatar-container {
      cursor: pointer;
      height: 50px;
      display: inline-block;
      float: right;
      margin-right: 35px;
      /*position: absolute;*/
      /*right: 35px;*/
    }
  }
</style>
