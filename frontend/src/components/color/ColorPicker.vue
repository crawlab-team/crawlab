<template>
  <div class="color-picker">
    <el-color-picker
        v-model="internalValue"
        :disabled="disabled"
        :predefine="predefine"
        :show-alpha="showAlpha"
        @change="onChange"
    />
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, ref, watch} from 'vue';

export default defineComponent({
  name: 'ColorPicker',
  props: {
    modelValue: {
      type: String,
    },
    disabled: {
      type: Boolean,
    },
    predefine: {
      type: Array as PropType<string[]>,
    },
    showAlpha: {
      type: Boolean,
      default: true,
    },
  },
  emits: [
    'update:model-value',
    'change',
  ],
  setup(props: ColorPickerProps, {emit}) {
    const internalValue = ref<string>();

    watch(() => props.modelValue, () => {
      internalValue.value = props.modelValue;
    });

    const onChange = (value: string) => {
      emit('update:model-value', value);
      emit('change', value);
    };

    onMounted(() => {
      const {modelValue} = props;
      internalValue.value = modelValue;
    });

    return {
      internalValue,
      onChange,
    };
  },
});
</script>

<style scoped>
.color-picker >>> .el-color-picker__trigger {
  border: none;
  padding: 0;
}


.color-picker >>> .el-color-picker__mask {
  background: transparent;
  border-radius: 0;
  left: 0;
}
</style>
