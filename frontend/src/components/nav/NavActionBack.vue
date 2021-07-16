<template>
  <NavActionGroup>
    <NavActionButton
        :button-type="buttonType"
        :disabled="disabled"
        :icon="icon"
        :label="label"
        :tooltip="tooltip"
        type="primary"
        @click="() => $emit('click')"
    />
  </NavActionGroup>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionButton from '@/components/nav/NavActionButton.vue';

export default defineComponent({
  name: 'NavActionBack',
  components: {NavActionButton, NavActionGroup},
  props: {
    buttonType: {
      type: String as PropType<ButtonType>,
      default: 'label',
    },
    label: {
      type: String,
      default: 'Back'
    },
    tooltip: {
      type: String,
    },
    icon: {
      type: [String, Array] as PropType<Icon>,
      default: () => {
        return ['fa', 'undo'];
      }
    },
    type: {
      type: String as PropType<BasicType>,
      default: 'primary'
    },
    size: {
      type: String as PropType<BasicSize>,
      default: 'mini'
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
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
