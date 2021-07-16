<template>
  <el-tag :type="type" class="node-active" size="mini">
    <font-awesome-icon :icon="icon" class="icon"/>
    <span>{{ label }}</span>
  </el-tag>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';

export default defineComponent({
  name: 'NodeActive',
  props: {
    active: {
      type: Boolean,
      default: false,
    }
  },
  setup(props: NodeActiveProps, {emit}) {
    const type = computed<string>(() => {
      const {active} = props;
      return active ? 'success' : 'info';
    });

    const label = computed<string>(() => {
      const {active} = props;
      return active ? 'Online' : 'Offline';
    });

    const icon = computed<string[]>(() => {
      const {active} = props;
      return active ? ['fa', 'check-circle'] : ['fa', 'times-circle'];
    });

    return {
      type,
      label,
      icon,
    };
  },
});
</script>

<style lang="scss" scoped>
.node-active {
  .icon {
    margin-right: 5px;
  }
}
</style>
