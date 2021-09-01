<template>
  <el-container class="basic-layout">
    <Sidebar />
    <el-container :class="sidebarCollapsed ? 'collapsed' : ''" class="container">
      <Header />
      <TabsView />
      <div class="container-body">
        <router-view />
      </div>
    </el-container>
  </el-container>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted} from 'vue';
import Header from './components/Header.vue';
import Sidebar from './components/Sidebar.vue';
import {useStore} from 'vuex';
import TabsView from '@/layouts/components/TabsView.vue';
import {initPlugins} from '@/utils/plugin';

export default defineComponent({
  name: 'BasicLayout',
  components: {
    TabsView,
    Header,
    Sidebar,
  },
  setup() {
    const store = useStore();
    const {layout} = store.state as RootStoreState;

    const sidebarCollapsed = computed<boolean>(() => layout.sidebarCollapsed);

    return {
      sidebarCollapsed,
    };
  }
});
</script>

<style lang="scss" scoped>
@import "../styles/variables.scss";

.basic-layout {
  height: 100vh;

  .container {
    position: fixed;
    top: 0;
    left: $sidebarWidth;
    display: block;
    width: calc(100vw - #{$sidebarWidth});
    height: 100vh;
    transition: left $sidebarCollapseTransitionDuration;
    z-index: 2;

    &.collapsed {
      left: $sidebarWidthCollapsed;
      width: calc(100vw - #{$sidebarWidthCollapsed});
    }

    .container-body {
      background-color: $containerBg;
      height: calc(100vh - #{$headerHeight} - #{$tabsViewHeight});
      overflow: auto;
    }
  }
}
</style>
