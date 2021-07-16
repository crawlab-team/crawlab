<template>
  <el-popover
      :placement="placement"
      :show-arrow="false"
      :visible="visible"
      popper-class="context-menu"
      trigger="manual"
  >
    <template #default>
      <slot name="default"></slot>
    </template>
    <template #reference>
      <div v-click-outside="onClickOutside">
        <slot name="reference"></slot>
      </div>
    </template>
  </el-popover>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {ClickOutside} from 'element-plus/lib/directives';

export const contextMenuDefaultProps = {
  visible: {
    type: Boolean,
    default: false,
  },
  placement: {
    type: String,
    default: 'right-start',
  },
  clicking: {
    type: Boolean,
    default: false,
  }
};

export const contextMenuDefaultEmits = [
  'hide',
];

export default defineComponent({
  name: 'ContextMenu',
  directives: {
    ClickOutside,
  },
  emits: contextMenuDefaultEmits,
  props: contextMenuDefaultProps,
  setup(props, {emit}) {
    const onClickOutside = () => {
      if (props.clicking) return;
      emit('hide');
    };

    return {
      onClickOutside,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
