<template>
  <el-tooltip content="New Tab">
    <Tab
        :icon="icon"
        :show-close="false"
        :show-title="false"
        class="action-tab"
        @click="onClick"
    />
  </el-tooltip>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import Icon from '@/components/icon/Icon.vue';
import Tab from '@/components/tab/Tab.vue';

export default defineComponent({
  name: 'ActionTab',
  components: {Tab},
  props: {
    icon: {
      type: [String, Array] as PropType<Icon>,
    },
    title: {
      type: String,
    }
  },
  emits: [
    'click',
  ],
  setup(props: ActionTabProps, {emit}) {
    const onClick = () => {
      emit('click');
    };

    return {
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>
.action-tab {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>

<style scoped>
.action-tab >>> .icon {
  color: #909399;
}
</style>
