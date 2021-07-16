import {readonly} from 'vue';
import {Store} from 'vuex';
import useForm from '@/components/form/form';
import useNodeService from '@/services/node/nodeService';
import {getDefaultFormComponentData} from '@/utils/form';
import {FORM_FIELD_TYPE_INPUT, FORM_FIELD_TYPE_INPUT_TEXTAREA, FORM_FIELD_TYPE_SWITCH} from '@/constants/form';

type Node = CNode;

// get new node
export const getNewNode = (): Node => {
  return {
    tags: [],
    max_runners: 8,
    enabled: true,
  };
};

// form component data
const formComponentData = getDefaultFormComponentData<Node>(getNewNode);

const useNode = (store: Store<RootStoreState>) => {
  // store
  const ns = 'node';
  const {node: state} = store.state as RootStoreState;

  // batch form fields
  const batchFormFields: FormTableField[] = [
    {
      prop: 'name',
      label: 'Name',
      width: '150',
      fieldType: FORM_FIELD_TYPE_INPUT,
      required: true,
      placeholder: 'Name',
    },
    {
      prop: 'enabled',
      label: 'Enabled',
      width: '120',
      fieldType: FORM_FIELD_TYPE_SWITCH,
    },
    {
      prop: 'description',
      label: 'Description',
      width: '200',
      fieldType: FORM_FIELD_TYPE_INPUT_TEXTAREA,
    },
  ];

  // form rules
  const formRules = readonly<FormRules>({});

  return {
    ...useForm(ns, store, useNodeService(store), formComponentData),
    batchFormFields,
    formRules,
  };
};

export default useNode;
