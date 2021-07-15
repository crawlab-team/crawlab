import {computed, readonly} from 'vue';
import {Store} from 'vuex';
import useForm from '@/components/form/form';
import useUserService from '@/services/user/userService';
import {getDefaultFormComponentData} from '@/utils/form';
import {FORM_FIELD_TYPE_INPUT, FORM_FIELD_TYPE_INPUT_PASSWORD, FORM_FIELD_TYPE_SELECT,} from '@/constants/form';
import {getModeOptions} from '@/utils/task';
import {ROLE_ADMIN, ROLE_NORMAL} from '@/constants/user';

// get new user
export const getNewUser = (): User => {
  return {
    role: ROLE_NORMAL,
  };
};

// form component data
const formComponentData = getDefaultFormComponentData<User>(getNewUser);

const useUser = (store: Store<RootStoreState>) => {
  // store
  const ns = 'user';
  const state = store.state[ns];

  // options for default mode
  const modeOptions = getModeOptions();

  // batch form fields
  const batchFormFields = computed<FormTableField[]>(() => [
    {
      prop: 'username',
      label: 'Username',
      width: '150',
      fieldType: FORM_FIELD_TYPE_INPUT,
      placeholder: 'Username',
      required: true,
    },
    {
      prop: 'password',
      label: 'Password',
      width: '150',
      placeholder: 'Password',
      fieldType: FORM_FIELD_TYPE_INPUT_PASSWORD,
      required: true,
    },
    {
      prop: 'email',
      label: 'Email',
      width: '150',
      fieldType: FORM_FIELD_TYPE_INPUT,
      placeholder: 'Email',
    },
    {
      prop: 'role',
      label: 'Role',
      width: '150',
      placeholder: 'Role',
      fieldType: FORM_FIELD_TYPE_SELECT,
      options: [
        {label: 'Admin', value: ROLE_ADMIN},
        {label: 'Normal', value: ROLE_NORMAL},
      ],
      required: true,
    },
  ]);

  // form rules
  const formRules = readonly<FormRules>({
    password: {
      trigger: 'blur',
      validator: ((_, value: string, callback) => {
        const invalidMessage = 'Invalid password. Length must be no less than 5.';
        if (0 < value.length && value.length < 5) return callback(invalidMessage);
        return callback();
      }),
    },
  });

  // all user select options
  const allUserSelectOptions = computed<SelectOption[]>(() => state.allList.map(d => {
    return {
      label: d.username,
      value: d._id,
    };
  }));

  return {
    ...useForm('user', store, useUserService(store), formComponentData),
    modeOptions,
    batchFormFields,
    formRules,
    allUserSelectOptions,
  };
};

export default useUser;
