<template>
  <div class="nav-tabs">
    <el-menu
        :default-active="activeKey"
        mode="horizontal"
        @select="onSelect"
    >
      <el-menu-item
          v-for="item in items"
          :key="item.id"
          :class="item.emphasis ? 'emphasis' : ''"
          :index="item.id"
          :style="item.style"
      >
        <el-tooltip :content="item.tooltip" :disabled="!item.tooltip">
          <template v-if="!!item.icon">
            <font-awesome-icon :icon="item.icon"/>
          </template>
          <template v-else>
            {{ item.title }}
          </template>
        </el-tooltip>
      </el-menu-item>
      <div class="extra">
        <slot name="extra">
        </slot>
      </div>
    </el-menu>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'NavTabs',
  props: {
    items: Array,
    activeKey: String,
  },
  setup(props, {emit}) {
    const onSelect = (index: string) => {
      emit('select', index);
    };

    return {
      onSelect,
    };
  },
});
</script>
<style lang="scss" scoped>
@import "../../styles/variables.scss";

.nav-tabs {
  .el-menu {
    height: calc(#{$navTabsHeight} + 1px);

    .el-menu-item {
      height: $navTabsHeight;
      line-height: $navTabsHeight;

      &.emphasis {
        color: $infoColor;
        border-bottom: none;
      }
    }

    .extra {
      float: right;
      height: $navTabsHeight;
      line-height: $navTabsHeight;
    }
  }
}
</style>
