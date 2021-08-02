<template>
  <div :class="sidebarCollapsed ? 'collapsed' : ''" class="header-container">
    <el-header :height="headerHeight" class="header">
      <div class="left">
      </div>
      <div class="right">
        <el-link :underline="false" href="javascript:;" @click="onLogout">
          Logout
        </el-link>
      </div>
    </el-header>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {useStore} from 'vuex';
import variables from '../../styles/variables.scss';
import {useRouter} from 'vue-router';

export default defineComponent({
  name: 'Header',
  setup() {
    // router
    const router = useRouter();

    // store
    const store = useStore();
    const {layout} = store.state as RootStoreState;

    const sidebarCollapsed = computed(() => {
      return layout.sidebarCollapsed;
    });

    const onLogout = () => {
      localStorage.removeItem('token');
      router.push('/login');
    };

    return {
      sidebarCollapsed,
      onLogout,
      ...variables,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.header-container {
  height: $headerHeight;
  width: calc(100vw - #{$sidebarWidth});
  background-color: $headerBg;
  transition: width $sidebarCollapseTransitionDuration;

  &.collapsed {
    width: calc(100vw - #{$sidebarWidthCollapsed});
  }

  .header {
    height: 100%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-left: none;
    border-bottom: 1px solid $headerBorderColor;

    .left {
      display: flex;
      align-items: center;
    }

    .right {
      display: flex;
      align-items: center;
    }
  }
}
</style>
