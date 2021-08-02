<template>
  <div class="input-with-button">
    <!-- Input -->
    <el-input
        v-model="internalValue"
        :placeholder="placeholder"
        :size="size"
        class="input"
        :disabled="disabled"
        @input="onInput"
        @blur="onBlur"
        @focus="onFocus"
        @keyup.enter="onBlur"
    />
    <!-- ./Input -->

    <!-- Button -->
    <Button
        v-if="buttonLabel"
        :disabled="disabled"
        :size="size"
        :type="buttonType"
        class="button"
        no-margin
        @click="onClick"
    >
      <Icon v-if="buttonIcon" :icon="buttonIcon" />
      {{ buttonLabel }}
    </Button>
    <template v-else-if="buttonIcon">
      <FaIconButton
          v-if="isFaIcon"
          :disabled="disabled"
          :icon="buttonIcon"
          :size="size"
          :type="buttonType"
          class="button"
          @click="onClick"
      />
      <IconButton
          v-else
          :disabled="disabled"
          :icon="buttonIcon"
          :size="size"
          :type="buttonType"
          class="button"
          @click="onClick"
      />
    </template>
    <!-- ./Button -->
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, ref, watch} from 'vue';
import Button from '@/components/button/Button.vue';
import Icon from '@/components/icon/Icon.vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import useIcon from '@/components/icon/icon';
import IconButton from '@/components/button/IconButton.vue';

export default defineComponent({
  name: 'InputWithButton',
  components: {
    IconButton,
    FaIconButton,
    Icon,
    Button,
  },
  props: {
    modelValue: {
      type: String,
    },
    placeholder: {
      type: String,
    },
    size: {
      type: String,
      default: 'mini',
    },
    buttonType: {
      type: String as PropType<BasicType>,
      default: 'primary',
    },
    buttonLabel: {
      type: String,
      default: 'Click',
    },
    buttonIcon: {
      type: [String, Array] as PropType<string | string[]>,
    },
    disabled: {
      type: Boolean,
      default: false,
    }
  },
  emits: [
    'update:model-value',
    'input',
    'click',
    'blur',
    'focus',
    'keyup.enter',
  ],
  setup(props: InputWithButtonProps, {emit}) {
    const internalValue = ref<string>();

    const {
      isFaIcon: _isFaIcon,
    } = useIcon();

    const isFaIcon = () => {
      const {buttonIcon} = props;
      if (!buttonIcon) return false;
      return _isFaIcon(buttonIcon);
    };

    watch(() => props.modelValue, () => {
      internalValue.value = props.modelValue;
    });

    const onInput = (value: string) => {
      emit('update:model-value', value);
      emit('input', value);
    };

    const onClick = () => {
      emit('click');
    };

    const onBlur = () => {
      emit('blur');
    };

    const onFocus = () => {
      emit('focus');
    };

    const onKeyUpEnter = () => {
      emit('keyup.enter');
    };

    onMounted(() => {
      const {modelValue} = props;
      internalValue.value = modelValue;
    });

    return {
      internalValue,
      isFaIcon,
      onClick,
      onInput,
      onBlur,
      onFocus,
      onKeyUpEnter,
    };
  },
});
</script>

<style lang="scss" scoped>
.input-with-button {
  display: inline-table;
  vertical-align: middle;
  //align-items: start;

  .input {
    display: table-cell;
  }

  .button {
    display: table-cell;
  }
}
</style>

<style scoped>
.input-with-button >>> .input.el-input .el-input__inner {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.input-with-button >>> .button .el-button {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  height: 28px;
}
</style>
