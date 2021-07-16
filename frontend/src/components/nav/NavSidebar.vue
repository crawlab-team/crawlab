<template>
  <div class="nav-sidebar" :class="classes">
    <div class="search">
      <el-input
          v-model="searchString"
          class="search-input"
          :prefix-icon="collapsed ? '' : 'fa fa-search'"
          placeholder="Search..."
          :clearable="true"
      />
      <div v-if="!collapsed" class="search-suffix" @click.stop="onToggle">
        <el-tooltip
            v-model="toggleTooltipValue"
            content="Collapse sidebar"
        >
          <font-awesome-icon :icon="['fa', 'outdent']" class="icon"/>
        </el-tooltip>
      </div>
    </div>
    <el-menu
        ref="navMenu"
        v-if="filteredItems && filteredItems.length > 0"
        class="nav-menu"
        :default-active="activeKey"
        @select="onSelect"
    >
      <el-menu-item v-for="item in filteredItems" :key="item.id" :index="item.id" class="nav-menu-item">
        <span class="title">{{ item.title }}</span>
        <!-- TODO: implement -->
        <span v-if="false" class="actions">
          <font-awesome-icon class="icon" :icon="['far', 'star']" @click="onStar(item.id)"/>
          <font-awesome-icon class="icon" :icon="['far', 'clone']" @click="onDuplicate(item.id)"/>
          <font-awesome-icon class="icon" :icon="['fa', 'trash-alt']" @click="onRemove(item.id)"/>
        </span>
      </el-menu-item>
    </el-menu>
    <Empty
        v-else
        description="No Items"
    />
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, ref} from 'vue';
import {ElMenu} from 'element-plus';
import variables from '@/styles/variables.scss';
import Empty from '@/components/empty/Empty.vue';

export default defineComponent({
  name: 'NavSidebar',
  components: {Empty},
  props: {
    items: Array,
    activeKey: String,
    collapsed: Boolean,
    showActions: Boolean,
  },
  setup(props, {emit}) {
    const toggling = ref(false);
    const searchString = ref('');
    const navMenu = ref<typeof ElMenu | null>(null);
    const toggleTooltipValue = ref(false);

    const filteredItems = computed<NavItem[]>(() => {
      const items = props.items as NavItem[];
      if (!searchString.value) return items;
      return items.filter(d => d.title?.toLocaleLowerCase().includes(searchString.value.toLocaleLowerCase()));
    });

    const classes = computed(() => {
      const {collapsed} = props as NavSidebarProps;
      const cls = [];
      if (collapsed) cls.push('collapsed');
      // if (toggling.value) cls.push('toggling');
      return cls;
    });

    const onSelect = (index: string) => {
      emit('select', index);
    };

    const onStar = (index: string) => {
      emit('star', index);
    };

    const onDuplicate = (index: string) => {
      emit('duplicate', index);
    };

    const onRemove = (index: string) => {
      emit('remove', index);
    };

    const onToggle = () => {
      const {collapsed} = props as NavSidebarProps;
      emit('toggle', !collapsed);
      toggleTooltipValue.value = false;
    };

    const scroll = (id: string) => {
      const idx = filteredItems.value.findIndex(d => d.id === id);
      if (idx === -1) return;
      const {navSidebarItemHeight} = variables;
      const navSidebarItemHeightNumber = Number(navSidebarItemHeight.replace('px', ''));
      if (!navMenu.value) return;
      const $el = navMenu.value.$el as HTMLDivElement;
      $el.scrollTo({
        top: navSidebarItemHeightNumber * idx,
      });
    };

    return {
      toggling,
      searchString,
      navMenu,
      toggleTooltipValue,
      filteredItems,
      classes,
      onSelect,
      onStar,
      onDuplicate,
      onRemove,
      onToggle,
      scroll,
    };
  },
});
</script>
<style scoped lang="scss">
@import "../../styles/variables.scss";

.nav-sidebar {
  position: relative;
  //margin: 10px;
  width: $navSidebarWidth;
  border-right: 1px solid $navSidebarBorderColor;
  background-color: $navSidebarBg;
  height: calc(100vh - #{$headerHeight} - #{$tabsViewHeight} - 1px);
  transition: width $navSidebarCollapseTransitionDuration;

  &.collapsed {
    margin: 10px 0;
    width: 0;
    border: none;

    .search {
      position: relative;
    }
  }

  .search {
    position: relative;
    height: $navSidebarSearchHeight;
    box-sizing: content-box;
    border-bottom: 1px solid $navSidebarBorderColor;

    .search-input {
      width: 100%;
      border: none;
      padding: 0;
      margin: 0;
    }

    .search-suffix {
      position: absolute;
      top: 0;
      right: 0;
      display: inline-flex;
      align-items: center;
      height: 40px;
      width: 25px;
      color: $navSidebarItemActionColor;
      cursor: pointer;
      //transition: right $navSidebarCollapseTransitionDuration;
    }
  }

  .nav-menu {
    list-style: none;
    padding: 0;
    margin: 0;
    border: none;
    max-height: calc(100% - #{$navSidebarSearchHeight});
    overflow-y: auto;
    color: $navSidebarColor;

    &.empty {
      height: $navSidebarItemHeight;
      display: flex;
      align-items: center;
      padding-left: 20px;
      font-size: 14px;
    }

    .nav-menu-item {
      position: relative;
      height: $navSidebarItemHeight;
      line-height: $navSidebarItemHeight;

      &:hover {
        .actions {
          display: inherit;
        }
      }

      .title {
        font-size: 14px;
        margin-bottom: 3px;
      }

      .subtitle {
        font-size: 12px;
      }

      .actions {
        display: none;
        position: absolute;
        top: 0;
        right: 10px;

        .icon {
          color: $navSidebarItemActionColor;
          margin-left: 3px;

          &:hover {
            color: $primaryColor;
          }
        }
      }
    }
  }

  .toggle-expand {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    display: flex;
    align-items: center;
    z-index: 100;
    cursor: pointer;

    &:hover {
      opacity: 0.7;
    }

    .wrapper {
      height: 24px;
      width: 24px;
      background-color: $infoPlainColor;
      border: 1px solid $infoColor;
      border-bottom-right-radius: 5px;
      border-top-right-radius: 5px;
      border-left: none;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>
<style scoped>
.nav-sidebar > .search >>> .el-input__inner {
  border: none;
}

.nav-sidebar.collapsed > .search >>> .el-input__inner {
  padding: 0;
  width: 0;
}
</style>
