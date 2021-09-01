import {useRoute, useRouter} from 'vue-router';
import {useStore} from 'vuex';
import {computed, onBeforeMount, onMounted, provide, ref} from 'vue';
import variables from '@/styles/variables.scss';
import {plainClone} from '@/utils/object';
import {getRoutePathByDepth} from '@/utils/route';
import {ElMessage} from 'element-plus';

const IGNORE_GET_ALL_NS = [
  'task',
];

const useDetail = <T = BaseModel>(ns: ListStoreNamespace) => {
  const router = useRouter();
  const route = useRoute();

  // store state
  const store = useStore();
  const state = store.state[ns] as BaseStoreState;
  const {
    form,
  } = state;

  const navSidebar = ref<NavSidebar | null>(null);

  const navActions = ref<NavActions | null>(null);

  const showActionsToggleTooltip = ref<boolean>(false);

  const navItems = computed<NavItem<T>[]>(() => state.allList.map((d: BaseModel) => {
    return {
      id: d._id,
      title: d.name,
    } as NavItem;
  }));

  const activeId = computed<string>(() => {
    const {id} = route.params;
    return id as string || form._id || '';
  });

  const activeTabName = computed<SpiderTabName>(() => store.getters[`${ns}/tabName`]);

  const sidebarCollapsed = computed<boolean>(() => state.sidebarCollapsed);

  const actionsCollapsed = computed<boolean>(() => state.actionsCollapsed);

  const tabs = computed(() => {
    const {infoBorderColor} = variables;
    const tabs = plainClone(state.tabs) as NavItem[];
    if (sidebarCollapsed.value) {
      tabs.splice(0, 0, {
        id: 'toggle',
        icon: ['fa', 'indent'],
        tooltip: 'Expand sidebar',
        emphasis: true,
        style: {
          'border-right': `1px solid ${infoBorderColor}`,
        }
      });
    }
    return tabs;
  });

  const contentContainerStyle = computed(() => {
    return {
      height: `calc(100% - ${variables.navTabsHeight} - 1px${navActions.value ? ' - ' + navActions.value.getHeight() : ''})`,
    };
  });

  const primaryRoutePath = computed<string>(() => getRoutePathByDepth(route.path));

  const afterSave = computed<Function[]>(() => state.afterSave);

  const getForm = async () => {
    if (!activeId.value) return;
    return await store.dispatch(`${ns}/getById`, activeId.value);
  };

  const onNavSidebarSelect = async (id: string) => {
    if (!id) {
      console.error(new Error('id is empty'));
      return;
    }
    await router.push(`${primaryRoutePath.value}/${id}`);
    await getForm();
  };

  const onNavSidebarToggle = (value: boolean) => {
    if (value) {
      store.commit(`${ns}/collapseSidebar`);
    } else {
      store.commit(`${ns}/expandSidebar`);
    }
  };

  const onActionsToggle = () => {
    showActionsToggleTooltip.value = false;
    if (actionsCollapsed.value) {
      store.commit(`${ns}/expandActions`);
    } else {
      store.commit(`${ns}/collapseActions`);
    }
  };

  const onNavTabsSelect = (tabName: string) => {
    if (tabName === 'toggle') {
      store.commit(`${ns}/expandSidebar`);
      return;
    }
    router.push(`${primaryRoutePath.value}/${activeId.value}/${tabName}`);
  };

  const onBack = () => {
    router.push(`${primaryRoutePath.value}`);
  };

  const onSave = async () => {
    if (!activeId.value) {
      console.error('Active id is empty');
      return;
    }
    await store.dispatch(`${ns}/updateById`, {id: activeId.value, form: state.form});
    ElMessage.success('Saved successfully');
    await Promise.all([
      store.dispatch(`${ns}/getAllList`),
      store.dispatch(`${ns}/getById`, activeId.value),
    ]);

    // after save
    await afterSave.value.map(fn => fn());
  };

  onBeforeMount(getForm);
  onBeforeMount(async () => {
    if (IGNORE_GET_ALL_NS.includes(ns)) return;
    await store.dispatch(`${ns}/getAllList`);
  });

  onMounted(() => {
    if (!navSidebar.value) return;
    navSidebar.value.scroll(activeId.value);
  });

  // store context
  provide<DetailStoreContext<T>>('store-context', {
    namespace: ns,
    store,
    state,
  });

  return {
    navItems,
    activeId,
    navSidebar,
    navActions,
    showActionsToggleTooltip,
    tabs,
    activeTabName,
    sidebarCollapsed,
    actionsCollapsed,
    contentContainerStyle,
    onNavSidebarSelect,
    onNavSidebarToggle,
    onActionsToggle,
    onNavTabsSelect,
    onBack,
    onSave,
  };
};

export default useDetail;
