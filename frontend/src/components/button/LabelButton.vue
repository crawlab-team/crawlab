<template>
  <Button
      :circle="circle"
      :disabled="disabled"
      :plain="plain"
      :round="round"
      :size="size"
      :tooltip="tooltip"
      :type="type"
      class="label-button"
      @click="() => $emit('click')"
  >
    <font-awesome-icon v-if="icon" :icon="icon" class="icon"/>
    {{ label }}
  </Button>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import Button, {buttonProps} from '@/components/button/Button.vue';

export default defineComponent({
  name: 'LabelButton',
  components: {Button},
  props: {
    icon: {
      type: [Array, String] as PropType<Icon>,
    },
    label: {
      type: String,
      required: true,
    },
    ...buttonProps
  },
  emits: [
    'click',
  ],
  setup(props: LabelButtonProps, {emit}) {
    return {};
  },
});
</script>

<style lang="scss" scoped>
.label-button {
  .icon {
    margin-right: 3px;
  }
}
</style>
