<template>
  <Tag
      :clickable="computedClickable"
      :label="label"
      :tooltip="tooltip"
      :type="computedType"
      :effect="computedEffect"
      :icon="computedIcon"
      :spinning="spinning"
      :width="width"
      class="checked-tag"
      @click="onClick"
      @mouseenter="onMouseEnter"
      @mouseleave="onMouseLeave"
  />
</template>

<script lang="ts">
import {computed, defineComponent, ref} from 'vue';
import Tag, {tagProps} from '@/components/tag/Tag.vue';

const checkTagProps = {
  modelValue: {
    type: Boolean,
    default: false,
  },
  ...tagProps,
};

export default defineComponent({
  name: 'CheckTag',
  components: {
    Tag,
  },
  props: checkTagProps,
  emits: [
    'update:model-value',
    'change',
  ],
  setup(props: CheckTagProps, {emit}) {
    const isHover = ref<boolean>(false);

    const computedType = computed<BasicType | undefined>(() => {
      const {modelValue, type, disabled} = props;
      if (modelValue) {
        return 'primary';
      }
      return disabled ? 'info' : type;
    });

    const computedIcon = computed<Icon>(() => {
      const {modelValue} = props;
      return modelValue ? ['far', 'check-square'] : ['far', 'square'];
    });

    const computedClickable = computed<boolean>(() => {
      const {clickable, disabled} = props;
      if (disabled) {
        return false;
      }
      if (clickable === undefined) {
        return true;
      }
      return clickable;
    });

    const computedEffect = computed<BasicEffect>(() => {
      const {modelValue} = props;
      if (modelValue) {
        return 'dark';
      }
      if (!computedClickable.value) {
        return 'plain';
      }
      return isHover.value ? 'light' : 'plain';
    });

    const onClick = () => {
      const {modelValue} = props;
      const newValue = !modelValue;
      emit('update:model-value', newValue);
      emit('change', newValue);
    };

    const onMouseEnter = () => {
      isHover.value = true;
    };

    const onMouseLeave = () => {
      isHover.value = false;
    };

    return {
      isHover,
      computedType,
      computedIcon,
      computedEffect,
      computedClickable,
      onClick,
      onMouseEnter,
      onMouseLeave,
    };
  },
});
</script>

<style lang="scss" scoped>
</style>
