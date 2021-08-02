import {GetterTree, Module, MutationTree} from 'vuex';

declare global {
  interface LayoutStoreModule extends Module<LayoutStoreState, RootStoreState> {
    getters: LayoutStoreGetters;
    mutations: LayoutStoreMutations;
  }

  interface LayoutStoreState {
    // sidebar
    sidebarCollapsed: boolean;
    menuItems: MenuItem[];

    // tabs view
    maxTabId: number;
    tabs: Tab[];
    activeTabId?: number;
    draggingTab?: Tab;
    targetTab?: Tab;
    isTabsDragging: boolean;
  }

  interface LayoutStoreGetters extends GetterTree<LayoutStoreState, RootStoreState> {
    tabs: StoreGetter<LayoutStoreState, Tab[]>;
    activeTab: StoreGetter<LayoutStoreState, Tab | undefined>;
  }

  interface LayoutStoreMutations extends MutationTree<LayoutStoreState> {
    setSideBarCollapsed: StoreMutation<LayoutStoreState, boolean>;
    setTabs: StoreMutation<LayoutStoreState, Tab[]>;
    setActiveTabId: StoreMutation<LayoutStoreState, number>;
    addTab: StoreMutation<LayoutStoreState, Tab>;
    updateTab: StoreMutation<LayoutStoreState, Tab>;
    removeTab: StoreMutation<LayoutStoreState, Tab>;
    removeAllTabs: StoreMutation<LayoutStoreState>;
    setDraggingTab: StoreMutation<LayoutStoreState, Tab>;
    resetDraggingTab: StoreMutation<LayoutStoreState>;
    setTargetTab: StoreMutation<LayoutStoreState, Tab>;
    resetTargetTab: StoreMutation<LayoutStoreState>;
    setIsTabsDragging: StoreMutation<LayoutStoreState, boolean>;
  }
}
