<template>
  <span :class="sidebarCollapsed ? 'collapsed' : ''" class="sidebar-toggle" @click="toggleSidebar">
    <font-awesome-icon v-if="!sidebarCollapsed" :icon="['fas', 'outdent']" />
    <font-awesome-icon v-else :icon="['fas', 'indent']" />
  </span>
  <el-aside :class="sidebarCollapsed ? 'collapsed' : ''" class="sidebar" width="inherit">
    <div class="logo-container">
      <div class="logo">
        <img :src="logo" alt="logo" className="logo-img" />
        <span class="logo-title">Crawlab</span>
        <span class="logo-sub-title">
          <div class="logo-sub-title-block">
            Community
          </div>
          <div class="logo-sub-title-block">
            v0.6.0
          </div>
        </span>
      </div>
    </div>
    <div class="sidebar-menu">
      <el-menu
          :collapse="sidebarCollapsed"
          :active-text-color="menuActiveText"
          :background-color="menuBg"
          :default-active="activePath"
          :text-color="menuText"
      >
        <el-menu-item
            v-for="(item, $index) in menuItems"
            :key="$index"
            :index="item.path"
            @click="onMenuItemClick(item)"
        >
          <MenuItemIcon :item="item" size="normal" />
          <template #title>
            <span class="menu-item-title">{{ item.title }}</span>
          </template>
        </el-menu-item>
        <div class="plugin-anchor" />
      </el-menu>
    </div>
  </el-aside>
  <div class="script-anchor" />
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {useStore} from 'vuex';
import {useRoute, useRouter} from 'vue-router';
import variables from '@/styles/variables.scss';
import logo from '@/assets/logo.svg';
import MenuItemIcon from '@/components/icon/MenuItemIcon.vue';
import {getPrimaryPath} from '@/utils/path';

export default defineComponent({
  name: 'Sidebar',
  components: {
    MenuItemIcon,
  },
  setup() {
    const router = useRouter();
    const route = useRoute();
    const store = useStore();
    const {layout} = store.state as RootStoreState;
    const {menuItems} = layout;
    const storeNamespace = 'layout';

    const activePath = computed<string>(() => {
      return getPrimaryPath(route.path);
    });

    const sidebarCollapsed = computed<boolean>(() => layout.sidebarCollapsed);

    const toggleIcon = computed<string[]>(() => {
      if (sidebarCollapsed.value) {
        return ['fas', 'indent'];
      } else {
        return ['fas', 'outdent'];
      }
    });

    const onMenuItemClick = (item: MenuItem) => {
      router.push(item.path);
    };

    const toggleSidebar = () => {
      store.commit(`${storeNamespace}/setSideBarCollapsed`, !sidebarCollapsed.value);
    };

    return {
      sidebarCollapsed,
      toggleIcon,
      menuItems,
      logo,
      activePath,
      onMenuItemClick,
      toggleSidebar,
      ...variables,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables";

.sidebar {
  overflow-x: hidden;
  user-select: none;

  &.collapsed {
    .logo-container,
    .sidebar-menu {
      width: $sidebarWidthCollapsed;
    }
  }

  .logo-container {
    display: inline-block;
    height: $headerHeight;
    width: $sidebarWidth;
    padding-left: 12px;
    padding-right: 20px;
    border-right: none;
    background-color: $menuBg;
    transition: width $sidebarCollapseTransitionDuration;

    .logo {
      display: flex;
      align-items: center;
      height: 100%;

      .logo-img {
        height: 40px;
        width: 40px;
      }

      .logo-title {
        font-family: BlinkMacSystemFont, -apple-system, segoe ui, roboto, oxygen, ubuntu, cantarell, fira sans, droid sans, helvetica neue, helvetica, arial, sans-serif;
        font-size: 28px;
        font-weight: 600;
        margin-left: 12px;
        color: #409eff;
      }

      .logo-sub-title {
        font-family: BlinkMacSystemFont, -apple-system, segoe ui, roboto, oxygen, ubuntu, cantarell, fira sans, droid sans, helvetica neue, helvetica, arial, sans-serif;
        font-size: 10px;
        height: 24px;
        line-height: 24px;
        margin-left: 10px;
        font-weight: 500;
        color: $menuText;

        .logo-sub-title-block {
          display: flex;
          align-items: center;
          height: 12px;
          line-height: 12px;
        }
      }
    }
  }

  .sidebar-menu {
    width: $sidebarWidth;
    height: calc(100vh - #{$headerHeight});
    margin: 0;
    padding: 0;
    transition: width $sidebarCollapseTransitionDuration;

    .el-menu {
      border-right: none;
      width: inherit !important;
      height: calc(100vh - #{$headerHeight});
      transition: none !important;

      .el-menu-item {
        &.is-active {
          background-color: $menuHover !important;
        }

        .menu-item-title {
          margin-left: 6px;
        }
      }
    }
  }
}

.sidebar-toggle {
  position: fixed;
  top: 0;
  left: $sidebarWidth;
  display: inline-flex;
  align-items: center;
  width: 18px;
  height: 64px;
  z-index: 5;
  color: $menuBg;
  font-size: 24px;
  margin-left: 10px;
  cursor: pointer;
  transition: left $sidebarCollapseTransitionDuration;

  &.collapsed {
    left: $sidebarWidthCollapsed;
  }
}
</style>
