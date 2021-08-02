<template>
  <el-tooltip :content="tooltip ? tooltip : undefined">
    <span>
      <el-button
          :circle="circle"
          :disabled="disabled"
          :icon="icon"
          :plain="plain"
          :round="round"
          :size="size"
          :title="tooltip"
          :type="type"
          class="button"
          style="padding: 7px"
          @click="() => $emit('click')"
      />
    </span>
  </el-tooltip>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {buttonProps} from '@/components/button/Button.vue';

export const iconButtonProps = {
  icon: {
    type: String,
    required: true,
  },
  ...buttonProps,
};

export default defineComponent({
  name: 'IconButton',
  props: iconButtonProps,
  emits: [
    'click',
  ],
  setup() {
    return {};
  },
});
</script>

<style lang="scss" scoped>

</style>
