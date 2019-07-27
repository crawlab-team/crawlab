<template>
  <el-scrollbar wrap-class="scrollbar-wrapper">
    <el-menu
      :show-timeout="200"
      :default-active="routeLevel1"
      :collapse="isCollapse"
      :background-color="variables.menuBg"
      :text-color="variables.menuText"
      :active-text-color="variables.menuActiveText"
      mode="vertical"
    >
      <sidebar-item v-for="route in routes" :key="route.path" :item="route" :base-path="route.path"/>
    </el-menu>
  </el-scrollbar>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import variables from '@/styles/variables.scss'
import SidebarItem from './SidebarItem'

export default {
  components: { SidebarItem },
  computed: {
    ...mapState('user', [
      'adminPaths'
    ]),
    ...mapGetters([
      'sidebar'
    ]),
    routeLevel1 () {
      let pathArray = this.$route.path.split('/')
      return `/${pathArray[1]}`
    },
    routes () {
      return this.$router.options.routes.filter(d => {
        const role = this.$store.getters['user/userInfo'].role
        if (role === 'admin') return true
        console.log(role)
        return !this.adminPaths.includes(d.path)
      })
    },
    variables () {
      return variables
    },
    isCollapse () {
      return !this.sidebar.opened
    }
  }
}
</script>
