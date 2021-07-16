import {Ref} from 'vue';

declare global {
  interface FormTableFieldProps {
    form: any;
    formRules?: FormRuleItem[];
    prop: string;
    fieldType: FormFieldType;
    options?: SelectOption[];
    required?: boolean;
    placeholder?: string;
    disabled?: boolean;
    onChange?: (value: any) => void;
    onRegister?: (formRef: Ref) => void;
  }
}
