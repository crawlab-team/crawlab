<template>
  <el-dialog
      :modal-class="modalClass"
      :before-close="onClose"
      :model-value="visible"
      :title="title"
      :top="top"
      :width="width"
      :z-index="zIndex"
  >
    <slot/>
    <template #footer>
      <slot name="prefix"/>
      <Button plain size="mini" type="info" @click="onClose">Cancel</Button>
      <Button
          :disabled="confirmDisabled"
          :loading="confirmLoading"
          size="mini"
          type="primary"
          @click="onConfirm"
      >
        Confirm
      </Button>
      <slot name="suffix"/>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import Button from '@/components/button/Button.vue';

export default defineComponent({
  name: 'Dialog',
  components: {Button},
  props: {
    visible: {
      type: Boolean,
      required: false,
      default: false,
    },
    modalClass: {
      type: String,
    },
    title: {
      type: String,
      required: false,
    },
    top: {
      type: String,
      required: false,
    },
    width: {
      type: String,
      required: false,
    },
    zIndex: {
      type: Number,
      required: false,
    },
    confirmDisabled: {
      type: Boolean,
      default: false,
    },
    confirmLoading: {
      type: Boolean,
      default: false,
    },
  },
  emits: [
    'close',
    'confirm',
  ],
  setup(props: DialogProps, {emit}) {
    const onClose = () => {
      emit('close');
    };

    const onConfirm = () => {
      emit('confirm');
    };

    return {
      onClose,
      onConfirm,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
