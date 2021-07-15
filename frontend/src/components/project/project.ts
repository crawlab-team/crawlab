import {computed, readonly} from 'vue';
import {Store} from 'vuex';
import {isDuplicated} from '@/utils/array';
import useForm from '@/components/form/form';
import useProjectService from '@/services/project/projectService';
import {getDefaultFormComponentData} from '@/utils/form';
import {FORM_FIELD_TYPE_INPUT, FORM_FIELD_TYPE_INPUT_TEXTAREA, FORM_FIELD_TYPE_TAG_INPUT} from '@/constants/form';

// get new project
export const getNewProject = (): Project => {
  return {
    tags: [],
  };
};

// form component data
const formComponentData = getDefaultFormComponentData<Project>(getNewProject);

const useProject = (store: Store<RootStoreState>) => {
  // store
  const ns = 'project';
  const state = store.state[ns];

  // batch form fields
  const batchFormFields = computed<FormTableField[]>(() => [
    {
      prop: 'name',
      label: 'Name',
      width: '150',
      fieldType: FORM_FIELD_TYPE_INPUT,
      placeholder: 'Name',
      required: true,
    },
    {
      prop: 'tags',
      label: 'Tags',
      width: '200',
      placeholder: 'Tags',
      fieldType: FORM_FIELD_TYPE_TAG_INPUT,
    },
    {
      prop: 'description',
      label: 'Description',
      width: '800',
      placeholder: 'Description',
      fieldType: FORM_FIELD_TYPE_INPUT_TEXTAREA,
    },
  ]);

  // form rules
  const formRules = readonly<FormRules>({
    tags: {
      validator: ((_, value, callback) => {
        if (isDuplicated(value)) {
          return callback('Duplicated tags');
        }
        callback();
      }),
    },
  });

  // all project select options
  const allProjectSelectOptions = computed<SelectOption[]>(() => state.allList.map(d => {
    return {
      label: d.name,
      value: d._id,
    };
  }));

  return {
    ...useForm('project', store, useProjectService(store), formComponentData),
    batchFormFields,
    formRules,
    allProjectSelectOptions,
  };
};

export default useProject;
