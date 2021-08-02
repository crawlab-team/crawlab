import useDetail from '@/layouts/detail';
import {computed, onBeforeMount} from 'vue';
import {useStore} from 'vuex';
import useSpiderService from '@/services/spider/spiderService';
import {useRoute} from 'vue-router';

const useSpiderDetail = () => {
  const ns = 'spider';
  const store = useStore();
  const {
    spider: state,
  } = store.state as RootStoreState;

  const route = useRoute();

  const id = computed(() => route.params.id as string);

  const activeNavItem = computed<FileNavItem | undefined>(() => state.activeNavItem);

  const fileContent = computed<string>(() => state.fileContent);

  const {
    saveFile: save,
  } = useSpiderService(store);

  const saveFile = async () => {
    if (!id.value || !activeNavItem.value?.path) return;
    await save(id.value, activeNavItem.value?.path, fileContent.value);
  };

  onBeforeMount(async () => {
    await Promise.all([
      store.dispatch(`project/getAllList`),
    ]);

    store.commit(`${ns}/setAfterSave`, [
      saveFile,
    ]);
  });

  return {
    ...useDetail('spider'),
  };
};

export default useSpiderDetail;
