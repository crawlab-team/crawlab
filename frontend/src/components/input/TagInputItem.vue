<template>
  <div
      :class="[
          isFocus ? 'is-focus' : '',
          isNew ? 'is-new' : '',
      ]"
      class="tag-input-item"
  >
    <!-- Input -->
    <div class="input-wrapper">
      <el-autocomplete
          ref="inputRef"
          v-model="internalValue.name"
          :disabled="disabled"
          :fetch-suggestions="fetchSuggestions"
          :placeholder="placeholder"
          :size="size"
          popper-class="tag-input-item-popper"
          class="input"
          value-key="name"
          @blur="onBlur"
          @focus="onFocus"
          @select="onSelect"
          @keyup.enter="onCheck"
      />
      <div class="actions">
        <font-awesome-icon
            :class="[isDisabled('check') ? 'disabled' : '']"
            :icon="['fa', 'check']"
            class="action-btn check"
            @click="onCheck"
        />
        <font-awesome-icon
            :class="[isDisabled('close') ? 'disabled' : '']"
            :icon="['fa', 'times']"
            class="action-btn close"
            @click="onClose"
        />
        <font-awesome-icon
            :class="[isDisabled('delete') ? 'disabled' : '']"
            :icon="['fa', 'trash']"
            class="action-btn delete"
            @click="onDelete"
        />
      </div>
    </div>
    <!-- ./Input -->

    <!-- Color Picker -->
    <ColorPicker
        v-model="internalValue.color"
        :disabled="!isNew"
        :predefine="predefinedColors"
        class="color-picker"
        show-alpha
    />
    <!--    <el-color-picker-->
    <!--        v-model="internalValue.color"-->
    <!--        :disabled="!isNew"-->
    <!--        :predefine="predefinedColors"-->
    <!--        class="color-picker"-->
    <!--        show-alpha-->
    <!--    />-->
    <!-- ./Color Picker -->
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, inject, onMounted, PropType, readonly, ref, watch} from 'vue';
import {ElInput} from 'element-plus';
import {plainClone} from '@/utils/object';
import useTagService from '@/services/tag/tagService';
import {useStore} from 'vuex';
import {FILTER_OP_CONTAINS, FILTER_OP_EQUAL} from '@/constants/filter';
import {getNewTag} from '@/components/tag/tag';
import {getPredefinedColors} from '@/utils/color';
import ColorPicker from '@/components/color/ColorPicker.vue';

export default defineComponent({
  name: 'TagInputItem',
  components: {
    ColorPicker,
  },
  props: {
    modelValue: {
      type: Object as PropType<Tag>,
    },
    placeholder: {
      type: String,
    },
    size: {
      type: String as PropType<BasicSize>,
      default: 'mini',
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
    'close',
    'check',
    'delete',
  ],
  setup(props: TagInputItemProps, {emit}) {
    const store = useStore();

    const internalValue = ref<Tag>(getNewTag());

    const isFocus = ref<boolean>(false);

    const inputRef = ref<typeof ElInput>();

    const isNew = computed<boolean>(() => !internalValue.value._id);

    // predefined colors
    const predefinedColors = readonly<string[]>(getPredefinedColors());

    watch(() => props.modelValue, () => {
      if (!props.modelValue) {
        internalValue.value = getNewTag();
      } else {
        internalValue.value = plainClone(props.modelValue);
      }
    });

    const isDisabled = (key: string) => {
      switch (key) {
        case 'check':
          return !internalValue.value.name;
        case 'close':
          return false;
        case 'delete':
          return false;
        default:
          return false;
      }
    };

    const onInput = (name: string) => {
      const value = {...props.modelValue, name};
      emit('input', value);
    };

    const onClick = () => {
      emit('click');
    };

    const onBlur = () => {
      isFocus.value = false;
      emit('blur');
    };

    const onFocus = () => {
      isFocus.value = true;
      emit('focus');
    };

    const focus = () => {
      inputRef.value?.focus();
    };

    const onSelect = (value: Tag) => {
      internalValue.value = value;
    };

    const onCheck = () => {
      if (isDisabled('check')) return;
      emit('update:model-value', internalValue.value);
      emit('check', internalValue.value);
    };

    const onClose = () => {
      if (isDisabled('close')) return;
      emit('close');
    };

    const onDelete = () => {
      if (isDisabled('delete')) return;
      emit('delete');
    };

    const ctx = inject<ListStoreContext<BaseModel>>('store-context');

    const fetchSuggestions = async (queryString: string, callback: (data: Tag[]) => void) => {
      const {
        getList,
      } = useTagService(store);
      const params = {
        page: 1,
        size: 50,
        conditions: [
          {key: 'col', op: FILTER_OP_EQUAL, value: `${ctx?.namespace}s`}
        ]
      } as ListRequestParams;
      if (queryString) {
        const conditions = params.conditions as FilterConditionData[];
        conditions.push({key: 'name', op: FILTER_OP_CONTAINS, value: queryString});
      }
      try {
        const res = await getList(params);
        return callback(res.data || []);
      } catch (e) {
        console.error(e);
        callback([]);
      }
    };

    onMounted(() => {
      if (!props.modelValue) {
        internalValue.value = getNewTag();
      } else {
        internalValue.value = plainClone(props.modelValue);
      }
    });

    return {
      predefinedColors,
      internalValue,
      isFocus,
      inputRef,
      isNew,
      onClick,
      onInput,
      onBlur,
      onFocus,
      focus,
      onCheck,
      onClose,
      onDelete,
      onSelect,
      isDisabled,
      fetchSuggestions,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.tag-input-item {
  display: flex;
  align-items: center;
  //height: 28px;

  .input-wrapper {
    display: inherit;
    border: none;
    position: relative;
    height: 28px;

    .actions {
      position: absolute;
      top: 0;
      right: 5px;

      .action-btn {
        width: 14px;
        height: 14px;
        padding: 3px;
        color: $infoMediumColor;
        cursor: pointer;

        &:hover:not(.disabled) {
          &.check {
            color: $successColor;
          }

          &.close {
            color: $infoColor;
          }

          &.delete {
            color: $dangerColor;
          }
        }

        &.disabled {
          color: $infoMediumLightColor;
          cursor: not-allowed;
        }
      }
    }
  }
}
</style>

<style scoped>
.tag-input-item >>> .input,
.tag-input-item >>> .actions,
.tag-input-item >>> .color-picker,
.tag-input-item >>> .color-picker .el-color-picker {
  margin: 0;
  padding: 0;
  height: 28px;
  line-height: 28px;
}

.tag-input-item >>> .input {
  display: inherit;
}

.tag-input-item >>> .input .el-input__inner {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
  transition: none;
}

.tag-input-item >>> .color-picker .el-color-picker__trigger {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  border-left: none;
  border-top: 1px solid #DCDFE6;
  border-right: 1px solid #DCDFE6;
  border-bottom: 1px solid #DCDFE6;
  padding: 0;
}

.tag-input-item.is-focus >>> .color-picker .el-color-picker__trigger {
  border-color: #409eff;
}

.tag-input-item >>> .color-picker .el-color-picker__color {
  border: none;
}

.tag-input-item >>> .color-picker .el-color-picker__mask {
  background: transparent;
  border-radius: 0;
  left: 0;
  height: 28px;
  width: 28px;
}

.tag-input-item >>> .el-autocomplete-suggestion__list > li {
  height: 28px;
}

</style>

<style>
.tag-input-item-popper >>> .el-autocomplete-suggestion__list > li {
  height: 28px;
}
</style>
