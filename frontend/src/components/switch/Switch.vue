<template>
  <el-switch
      v-model="internalValue"
      :active-color="activeColor"
      :active-icon-class="activeIconClass"
      :active-text="activeText"
      :disabled="disabled"
      :inactive-color="inactiveColor"
      :inactive-icon-class="inactiveIconClass"
      :inactive-text="inactiveText"
      :loading="loading"
      :width="width"
      @change="onChange"
  />
</template>

<script lang="ts">
import {defineComponent, onBeforeMount, ref, watch} from 'vue';
import variables from '@/styles/variables.scss';

export default defineComponent({
  name: 'Switch',
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    activeColor: {
      type: String,
      default: variables.successColor,
    },
    inactiveColor: {
      type: String,
      default: variables.infoMediumColor,
    },
    activeIconClass: {
      type: String,
      default: '',
    },
    inactiveIconClass: {
      type: String,
      default: '',
    },
    activeText: {
      type: String,
      default: '',
    },
    inactiveText: {
      type: String,
      default: '',
    },
    width: {
      type: Number,
      default: 40,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  emits: [
    'update:model-value',
    'change',
  ],
  setup(props: SwitchProps, {emit}) {
    const internalValue = ref<boolean>(false);
    watch(() => props.modelValue, () => {
      internalValue.value = props.modelValue;
    });

    const onChange = (value: boolean) => {
      internalValue.value = value;
      emit('update:model-value', value);
      emit('change', value);
    };

    onBeforeMount(() => {
      const {modelValue} = props;
      internalValue.value = modelValue;
    });

    return {
      variables,
      internalValue,
      onChange,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
