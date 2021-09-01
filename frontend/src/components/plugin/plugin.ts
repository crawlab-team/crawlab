import {readonly} from 'vue';
import {Store} from 'vuex';
import useForm from '@/components/form/form';
import usePluginService from '@/services/plugin/pluginService';
import {getDefaultFormComponentData} from '@/utils/form';

type Plugin = CPlugin;

// get new plugin
export const getNewPlugin = (): Plugin => {
  return {};
};

// form component data
const formComponentData = getDefaultFormComponentData<Plugin>(getNewPlugin);

const usePlugin = (store: Store<RootStoreState>) => {
  // store
  const ns = 'plugin';

  // form rules
  const formRules = readonly<FormRules>({});

  return {
    ...useForm(ns, store, usePluginService(store), formComponentData),
    formRules,
  };
};

export default usePlugin;
