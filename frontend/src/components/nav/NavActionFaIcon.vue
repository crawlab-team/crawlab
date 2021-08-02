<template>
  <NavActionItem class="nav-action-fa-icon" is-label>
    <el-tooltip :content="tooltip" :disabled="!tooltip">
      <font-awesome-icon :icon="icon" class="title"/>
    </el-tooltip>
  </NavActionItem>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import NavActionItem from '@/components/nav/NavActionItem.vue';

export default defineComponent({
  name: 'NavActionFaIcon',
  components: {
    NavActionItem,
  },
  props: {
    tooltip: {
      type: String,
    },
    icon: {
      type: [String, Array] as PropType<Icon>,
    },
  },
  setup(props, {emit}) {
    return {};
  },
});
</script>

<style lang="scss" scoped>
</style>
