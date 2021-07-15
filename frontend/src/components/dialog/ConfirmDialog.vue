<template>
  <Dialog
      :confirm-loading="confirmLoading"
      :title="title"
      :visible="visible"
      @cancel="onCancel"
      @confirm="onConfirm"
  >
    <template v-if="content">
      {{ content }}
    </template>
    <template v-else>
      <slot></slot>
    </template>
  </Dialog>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue';
import Dialog from '@/components/dialog/Dialog.vue';
import {voidFunc} from '@/utils/func';

export default defineComponent({
  name: 'ConfirmDialog',
  components: {Dialog},
  props: {
    confirmFunc: {
      type: Function,
      default: voidFunc,
    },
    title: {
      type: String,
      required: true,
    },
    content: {
      type: String,
    }
  },
  emits: [
    'confirm',
    'cancel',
  ],
  setup(props: ConfirmDialogProps, {emit}) {
    const visible = ref<boolean>(false);

    const confirmLoading = ref<boolean>(false);

    const onCancel = () => {
      visible.value = false;
      emit('cancel');
    };

    const onConfirm = async () => {
      const {confirmFunc} = props;
      confirmLoading.value = true;
      await confirmFunc();
      confirmLoading.value = false;
      visible.value = false;
      emit('confirm');
    };

    return {
      visible,
      confirmLoading,
      onCancel,
      onConfirm,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
