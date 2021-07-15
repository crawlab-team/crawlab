<template>
  <CheckTag
      v-for="op in options"
      :key="{v: op.value, c:checkedMap[op.value]}"
      v-model="checkedMap[op.value]"
      :disabled="disabled"
      :label="op.label"
      clickable
      style="margin-right: 10px"
      @change="onChange"
  />
</template>

<script lang="ts">
import {computed, defineComponent, PropType, reactive} from 'vue';
import CheckTag from '@/components/tag/CheckTag.vue';

export default defineComponent({
  name: 'CheckTagGroup',
  components: {
    CheckTag,
  },
  props: {
    modelValue: {
      type: Array as PropType<string[]>,
      default: () => {
        return [];
      }
    },
    options: {
      type: Array as PropType<SelectOption[]>,
      default: () => {
        return [];
      }
    },
    disabled: {
      type: Boolean,
      default: false,
    }
  },
  emits: [
    'update:model-value',
    'change',
  ],
  setup(props: CheckTagGroupProps, {emit}) {
    const checkedMap = reactive<{ [key: string]: boolean }>({});

    const checkedKeys = computed<string[]>(() => {
      return Object.keys(checkedMap).filter(k => checkedMap[k]);
    });

    const onChange = () => {
      emit('update:model-value', checkedKeys.value);
      emit('change', checkedKeys.value);
    };

    return {
      checkedMap,
      onChange,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
