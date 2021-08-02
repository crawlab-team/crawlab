<template>
  <el-form ref="formRef" :model="form" :rules="computedFormRules" inline-message>
    <el-form-item ref="formItemRef" :prop="prop" :required="isRequired">
      <el-input
          v-if="fieldType === FORM_FIELD_TYPE_INPUT"
          v-model="internalValue"
          :placeholder="placeholder"
          size="mini"
          :disabled="disabled"
          @input="onInputChange"
      />
      <el-input
          v-else-if="fieldType === FORM_FIELD_TYPE_INPUT_PASSWORD"
          v-model="internalValue"
          :disabled="disabled"
          :placeholder="placeholder"
          size="mini"
          type="password"
          @input="onInputChange"
      />
      <el-input
          v-else-if="fieldType === FORM_FIELD_TYPE_INPUT_TEXTAREA"
          v-model="internalValue"
          :placeholder="placeholder"
          size="mini"
          type="textarea"
          :disabled="disabled"
          @input="onInputChange"
      />
      <el-select
          v-else-if="fieldType === FORM_FIELD_TYPE_SELECT"
          v-model="internalValue"
          :placeholder="placeholder"
          size="mini"
          :disabled="disabled"
          @change="onInputChange"
      >
        <el-option
            v-for="op in options"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
      <InputWithButton
          v-else-if="fieldType === FORM_FIELD_TYPE_INPUT_WITH_BUTTON"
          v-model="internalValue"
          :placeholder="placeholder"
          button-label="Edit"
          size="mini"
          :disabled="disabled"
          @input="onInputChange"
      />
      <TagInput
          v-else-if="fieldType === FORM_FIELD_TYPE_TAG_INPUT"
          v-model="internalValue"
          :disabled="disabled"
          @change="onInputChange"
      />
      <Switch
          v-else-if="fieldType === FORM_FIELD_TYPE_SWITCH"
          v-model="internalValue"
          :disabled="disabled"
          @change="onInputChange"
      />
      <!-- TODO: implement more field types -->
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  onBeforeMount,
  onMounted,
  PropType,
  Ref,
  ref,
  SetupContext,
  watch
} from 'vue';
import {
  FORM_FIELD_TYPE_CHECK_TAG_GROUP,
  FORM_FIELD_TYPE_INPUT,
  FORM_FIELD_TYPE_INPUT_PASSWORD,
  FORM_FIELD_TYPE_INPUT_TEXTAREA,
  FORM_FIELD_TYPE_INPUT_WITH_BUTTON,
  FORM_FIELD_TYPE_SELECT,
  FORM_FIELD_TYPE_SWITCH,
  FORM_FIELD_TYPE_TAG_INPUT,
  FORM_FIELD_TYPE_TAG_SELECT,
} from '@/constants/form';
import TagInput from '@/components/input/TagInput.vue';
import {emptyArrayFunc, voidFunc} from '@/utils/func';
import {ElForm, ElFormItem} from 'element-plus';
import InputWithButton from '@/components/input/InputWithButton.vue';
import Switch from '@/components/switch/Switch.vue';

export default defineComponent({
  name: 'FormTableField',
  components: {
    Switch,
    InputWithButton,
    TagInput,
  },
  props: {
    form: {
      type: Object as PropType<any>,
      required: true,
    },
    formRules: {
      type: Object as PropType<FormRuleItem[]>,
      required: false,
    },
    prop: {
      type: String,
      required: true,
    },
    fieldType: {
      type: String as PropType<FormFieldType>,
      required: true,
    },
    options: {
      type: Array as PropType<SelectOption[]>,
      default: emptyArrayFunc,
    },
    required: {
      type: Boolean,
      default: false,
    },
    placeholder: {
      type: String,
      default: 'Please enter value',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    onChange: {
      type: Function as PropType<(value: any) => void>,
      default: voidFunc,
    },
    onRegister: {
      type: Function as PropType<(formRef: Ref) => void>,
      default: voidFunc,
    }
  },
  setup(props: FormTableFieldProps, {emit}: SetupContext) {
    // form ref
    const formRef = ref<typeof ElForm>();

    // form item ref
    const formItemRef = ref<typeof ElFormItem>();

    // internal value
    const internalValue = ref<any>();

    // computed field value
    const fieldValue = computed(() => {
      const {form, prop} = props;
      return form[prop];
    });
    watch(() => fieldValue.value, () => {
      if (internalValue.value !== fieldValue.value) {
        internalValue.value = fieldValue.value;
      }
    });

    const onInputChange = (value: any) => {
      const {onChange} = props;
      onChange?.(value);
    };

    const isEmptyForm = inject('fn:isEmptyForm') as (d: any) => boolean;

    const isRequired = computed<boolean>(() => {
      const {form, required} = props;
      if (isEmptyForm(form)) return false;
      return required || false;
    });

    const isErrorMessageVisible = computed<boolean>(() => !!formItemRef.value?.validateMessage);

    const computedFormRules = computed<FormRuleItem[]>(() => {
      const {form, formRules} = props;
      if (isEmptyForm(form)) {
        return [];
      } else {
        return formRules || [];
      }
    });

    onBeforeMount(() => {
      const {form, prop} = props;

      // initialize internal value
      internalValue.value = form[prop];
    });

    onMounted(() => {
      const {onRegister} = props;

      // register form ref
      onRegister?.(formRef);
    });

    return {
      FORM_FIELD_TYPE_INPUT,
      FORM_FIELD_TYPE_INPUT_PASSWORD,
      FORM_FIELD_TYPE_INPUT_TEXTAREA,
      FORM_FIELD_TYPE_INPUT_WITH_BUTTON,
      FORM_FIELD_TYPE_SELECT,
      FORM_FIELD_TYPE_TAG_INPUT,
      FORM_FIELD_TYPE_TAG_SELECT,
      FORM_FIELD_TYPE_CHECK_TAG_GROUP,
      FORM_FIELD_TYPE_SWITCH,
      formRef,
      formItemRef,
      internalValue,
      onInputChange,
      isRequired,
      isErrorMessageVisible,
      computedFormRules,
    };
  },
});
</script>

<style lang="scss" scoped>
.el-form {
  margin: 0;

  .el-form-item {
    margin: 0;
  }
}
</style>
