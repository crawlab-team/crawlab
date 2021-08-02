import {RuleItem} from 'async-validator';
import {Ref} from 'vue';
import {
  FORM_FIELD_TYPE_CHECK_TAG_GROUP,
  FORM_FIELD_TYPE_INPUT,
  FORM_FIELD_TYPE_INPUT_TEXTAREA,
  FORM_FIELD_TYPE_INPUT_WITH_BUTTON,
  FORM_FIELD_TYPE_TAG_INPUT,
  FORM_FIELD_TYPE_TAG_SELECT,
} from '@/constants/form';

declare global {
  interface FormProps {
    inline: boolean;
    labelWidth?: string;
    size?: string;
    grid: number;
    rules?: FormRuleItem | FormRuleItem[];
  }

  interface FormContext {
    labelWidth?: string;
    size?: string;
    grid: number;
  }

  interface FormModel {
    [key: string]: any;
  }

  interface FormRuleItem extends RuleItem {
    trigger?: string;
  }

  interface FormRules {
    [key: string]: FormRuleItem | FormRuleItem[];
  }

  type FormValidateCallback = (valid: boolean) => void;

  interface FormComponentData<T> {
    form: Ref<T>;
    formRef: Ref;
    formList: Ref<T[]>;
    formTableFieldRefsMap: Ref<FormTableFieldRefsMap>;
  }

  type FormTableFieldRefsMapKey = [number, string];
  type FormTableFieldRefsMap = Map<FormTableFieldRefsMapKey, Ref>;

  type DefaultFormFunc<T> = () => T;

  type FormFieldType = FORM_FIELD_TYPE_INPUT |
    FORM_FIELD_TYPE_INPUT_TEXTAREA |
    FORM_FIELD_TYPE_INPUT_WITH_BUTTON |
    FORM_FIELD_TYPE_TAG_INPUT |
    FORM_FIELD_TYPE_TAG_SELECT |
    FORM_FIELD_TYPE_CHECK_TAG_GROUP;
}
