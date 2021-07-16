<template>
  <div ref="formItem" :style="style" class="form-item">
    <el-form-item
        :prop="prop"
        :label="label"
        :required="isRequired"
        :rules="rules"
        :size="size || formContext?.size"
    >
      <template #label>
        <el-tooltip :content="labelTooltip" :disabled="!labelTooltip">
          <span class="form-item-label">
            <span :class="showRequiredAsterisk ? 'required' : ''" class="form-item-label-text">
              {{ label }}
            </span>
            <el-tooltip v-if="isSelectiveForm" :content="editableTooltip">
              <el-checkbox
                  v-model="internalEditable"
                  :disabled="notEditable"
                  class="editable-checkbox"
                  @change="onEditableChange"
              />
            </el-tooltip>
          </span>
        </el-tooltip>
      </template>
      <div class="form-item-content">
        <slot></slot>
      </div>
    </el-form-item>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, inject, onMounted, PropType, ref, watch} from 'vue';
import {RuleItem} from 'async-validator';
import {cloneArray} from '@/utils/object';

export default defineComponent({
  name: 'FormItem',
  props: {
    prop: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    labelTooltip: {
      type: String,
      required: false,
    },
    labelWidth: {
      type: String,
      required: false,
    },
    size: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      required: false,
      default: false,
    },
    span: {
      type: Number,
      required: false,
      default: 1,
    },
    offset: {
      type: Number,
      required: false,
      default: 0,
    },
    rules: {
      type: [Object, Array] as PropType<RuleItem | RuleItem[]>,
    },
    notEditable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props: FormItemProps, {emit}) {
    const formItem = ref<HTMLDivElement>();

    // form context
    const formContext = inject<FormContext>('form-context', {} as FormContext);

    // store context
    const storeContext = inject<ListStoreContext<BaseModel>>('store-context');
    const ns = storeContext?.namespace;
    const store = storeContext?.store;
    const state = storeContext?.state;
    const isSelectiveForm = computed<boolean | undefined>(() => state?.isSelectiveForm);
    const selectedFormFields = computed<string[] | undefined>(() => state?.selectedFormFields);
    const isBatchForm = computed<boolean>(() => store?.getters[`${ns}/isBatchForm`]);
    const activeDialogKey = computed<DialogKey | undefined>(() => state?.activeDialogKey);

    const style = computed<Partial<CSSStyleDeclaration>>(() => {
      const {span, offset} = props;
      return {
        flexBasis: `calc(100% / ${formContext.grid} * ${span})`,
        marginRight: `calc(100% / ${formContext.grid} * ${offset})`,
      };
    });

    const internalEditable = ref<boolean>(false);
    watch(() => state?.activeDialogKey, () => internalEditable.value = false);

    const editableTooltip = computed<string>(() => {
      const {notEditable} = props;
      if (notEditable) return 'Unable to edit';
      return internalEditable.value ? 'Uncheck to disable editing' : 'Check to enable editing';
    });

    const onEditableChange = (value: boolean) => {
      const {prop} = props;
      if (!selectedFormFields.value || !prop) return;
      const fields = cloneArray<string>(selectedFormFields.value);
      if (value) {
        if (!fields.includes(prop)) {
          fields.push(prop);
        }
      } else {
        if (fields.includes(prop)) {
          const idx = fields.findIndex(field => field === prop);
          fields.splice(idx, 1);
        }
      }
      store?.commit(`${ns}/setSelectedFormFields`, fields);
    };

    const isRequired = computed<boolean>(() => {
      switch (activeDialogKey.value) {
        case 'edit':
          if (isBatchForm.value) {
            if (!internalEditable.value) return false;
          }
          break;
      }

      const {required} = props;
      return required;
    });

    const showRequiredAsterisk = computed<boolean>(() => {
      if (isSelectiveForm.value) return false;
      const {required} = props;
      return required;
    });

    onMounted(() => {
      if (formItem.value) {
        const {labelWidth} = formContext;
        const el = formItem.value?.querySelector('.el-form-item__content') as HTMLDivElement;
        if (labelWidth && el.style) {
          el.style.width = `calc(100% - ${labelWidth})`;
        }
      }
    });

    return {
      formItem,
      formContext,
      style,
      isSelectiveForm,
      internalEditable,
      editableTooltip,
      onEditableChange,
      isRequired,
      showRequiredAsterisk,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables";

.form-item {
  .el-form-item {
    width: 100%;
    margin-right: 0;

    .form-item-label {
      .editable-checkbox {
        margin-left: 10px;
      }

      .form-item-label-text.required {
        &:before {
          content: "*";
          color: $red;
          margin-right: 4px;
        }
      }
    }
  }
}
</style>
<style scoped>
.form-item >>> .form-item-content > .el-select,
.form-item >>> .form-item-content > .el-autocomplete,
.form-item >>> .form-item-content > .el-input {
  width: 100%;
}
</style>
