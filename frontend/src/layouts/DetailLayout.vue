<template>
  <div :class="(noSidebar || sidebarCollapsed) ? 'collapsed' : ''" class="detail-layout">
    <div class="sidebar">
      <NavSidebar
          v-if="!noSidebar"
          ref="navSidebar"
          :active-key="activeId"
          :collapsed="noSidebar || sidebarCollapsed"
          :items="computedNavItems"
          @select="onNavSidebarSelect"
          @toggle="onNavSidebarToggle"
      />
    </div>
    <div class="content">
      <NavTabs :active-key="activeTabName" :items="tabs" class="nav-tabs" @select="onNavTabsSelect">
        <template v-slot:extra>
          <el-tooltip
              v-model="showActionsToggleTooltip"
              :content="actionsCollapsed ? 'Expand actions bar' : 'Collapse actions bar'"
          >
            <div :class="actionsCollapsed ? 'collapsed' : ''" class="actions-toggle" @click="onActionsToggle">
              <font-awesome-icon :icon="['fa', 'angle-up']" class="icon"/>
            </div>
          </el-tooltip>
        </template>
      </NavTabs>
      <NavActions ref="navActions" :collapsed="actionsCollapsed" class="nav-actions">
        <NavActionGroupDetailCommon
            @back="onBack"
            @save="onSave"
        />
        <slot name="actions"/>
      </NavActions>
      <div :style="contentContainerStyle" class="content-container">
        <router-view/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import useDetail from '@/layouts/detail';
import NavSidebar from '@/components/nav/NavSidebar.vue';
import NavTabs from '@/components/nav/NavTabs.vue';
import NavActions from '@/components/nav/NavActions.vue';
import NavActionGroupDetailCommon from '@/components/nav/NavActionGroupDetailCommon.vue';
import {useStore} from 'vuex';

export default defineComponent({
  name: 'DetailLayout',
  components: {
    NavActionGroupDetailCommon,
    NavSidebar,
    NavTabs,
    NavActions,
  },
  props: {
    storeNamespace: {
      type: String as PropType<ListStoreNamespace>,
      required: true,
    },
    noSidebar: {
      type: Boolean,
      default: false,
    },
    navItemNameKey: {
      type: String,
      default: 'name',
    }
  },
  setup(props: DetailLayoutProps, {emit}) {
    const ns = computed(() => props.storeNamespace);
    const store = useStore();
    const state = store.state[ns.value] as BaseStoreState;

    const computedNavItems = computed<NavItem[]>(() => state.allList.map((d: BaseModel) => {
      const {navItemNameKey} = props;
      return {
        id: d._id,
        title: d[navItemNameKey],
      } as NavItem;
    }));

    return {
      ...useDetail(ns.value),
      computedNavItems,
    };
  },
});
</script>

<style lang="scss" scoped>

@import "../styles/variables.scss";

.detail-layout {
  display: flex;
  height: 100%;

  &.collapsed {
    .sidebar {
      flex-basis: 0;
      width: 0;
    }

    .content {
      flex: 1;
      max-width: 100%;
    }
  }

  .sidebar {
    flex-basis: $navSidebarWidth;
    transition: all $navSidebarCollapseTransitionDuration;
  }

  .content {
    //margin: 10px;
    flex: 1;
    background-color: $containerWhiteBg;
    display: flex;
    flex-direction: column;
    max-width: calc(100% - #{$navSidebarWidth});

    .nav-tabs {
      height: calc(#{$navTabsHeight} + 1px);
    }

    .nav-actions {
      height: fit-content;
    }

    .content-container {
      flex: 1;
    }
  }

  .actions-toggle {
    height: $navTabsHeight;
    color: $infoColor;
    cursor: pointer;
    padding: 0 20px;
    border-left: 1px solid $containerBorderColor;

    .icon {
      transition: all $defaultTransitionDuration;
    }

    &.collapsed {
      .icon {
        transform: rotateZ(180deg);
      }
    }
  }
}
</style>
