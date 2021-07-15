<template>
  <Button
      :circle="circle"
      :disabled="disabled"
      :plain="plain"
      :round="round"
      :size="size"
      :tooltip="tooltip"
      :type="type"
      is-icon
      class="fa-icon-button"
      @click="() => $emit('click')"
  >
    <font-awesome-icon :icon="icon"/>
    <div v-if="badgeIcon" class="badge-icon">
      <font-awesome-icon :icon="badgeIcon"/>
    </div>
  </Button>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import {buttonProps} from './Button.vue';
import Button from '@/components/button/Button.vue';

export const faIconButtonProps = {
  icon: {
    type: [Array, String] as PropType<Icon>,
    required: true,
  },
  badgeIcon: {
    type: [Array, String] as PropType<Icon>,
    required: false,
  },
  ...buttonProps,
};

export default defineComponent({
  name: 'FaIconButton',
  components: {Button},
  props: faIconButtonProps,
  emits: [
    'click',
  ],
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
@import "../../styles/variables";

.badge-icon {
  position: absolute;
  top: -2px;
  right: 2px;
  font-size: 8px;
  color: $white;
}
</style>

<style scoped>
.el-button,
.el-button--mini,
.fa-icon-button,
.fa-icon-button >>> .el-button,
.fa-icon-button >>> .el-button--mini,
.fa-icon-button >>> .button {
  padding: 7px;
}
</style>
