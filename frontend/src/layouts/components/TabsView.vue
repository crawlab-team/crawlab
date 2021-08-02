<template>
  <div class="tabs-view">
    <DraggableList
        class="tab-list"
        :items="tabs"
        item-key="id"
        @d-end="onDragDrop"
    >
      <template v-slot="{item}">
        <Tab :tab="item"/>
      </template>
    </DraggableList>

    <!-- Add Tab -->
    <ActionTab :icon="['fa', 'plus']" class="add-tab" @click="onAddTab"/>
    <!-- ./Add Tab -->
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, onMounted, watch} from 'vue';
import {useStore} from 'vuex';
import TabComp from '@/components/tab/Tab.vue';
import {useRoute, useRouter} from 'vue-router';
import DraggableList from '@/components/drag/DraggableList.vue';
import {plainClone} from '@/utils/object';
import ActionTab from '@/components/tab/ActionTab.vue';

export default defineComponent({
  name: 'TabsView',
  components: {
    ActionTab,
    DraggableList,
    Tab: TabComp,
  },
  setup() {
    // store
    const storeNameSpace = 'layout';
    const store = useStore<RootStoreState>();

    // route
    const route = useRoute();

    // router
    const router = useRouter();

    // current path
    const currentPath = computed(() => route.path);

    // tabs
    const tabs = computed<Tab[]>(() => store.getters[`${storeNameSpace}/tabs`]);

    const addTab = (tab: Tab) => {
      store.commit(`${storeNameSpace}/addTab`, tab);
    };

    const setActiveTab = (tab: Tab) => {
      store.commit(`${storeNameSpace}/setActiveTabId`, tab.id);
    };

    const onAddTab = () => {
      addTab({path: '/'});
      const newTab = tabs.value[tabs.value.length - 1];
      setActiveTab(newTab);
      router.push(newTab.path);
    };

    const onDragDrop = (tabs: Tab[]) => {
      store.commit(`${storeNameSpace}/setTabs`, tabs);
    };

    watch(currentPath, (path) => {
      // active tab
      const activeTab = store.getters[`${storeNameSpace}/activeTab`] as Tab | undefined;

      // skip if active tab is undefined
      if (!activeTab) return;

      // set path to active tab
      activeTab.path = path;

      // update path of active tab
      store.commit(`${storeNameSpace}/updateTab`, plainClone(activeTab));
    });

    onMounted(() => {
      // add current page to tabs if no tab exists
      if (tabs.value.length === 0) {
        // add tab
        addTab({path: currentPath.value});

        // new tab
        const newTab = tabs.value[0];
        if (!newTab) return;

        // set active tab id
        setActiveTab(newTab);
      }
    });

    return {
      tabs,
      onAddTab,
      currentPath,
      onDragDrop,
    };
  }
});
</script>
<style lang="scss" scoped>
@import "../../styles/variables";

.tabs-view {
  padding: 10px 0;
  border-bottom: 1px solid $tabsViewBorderColor;
  background-color: $tabsViewBg;
  display: flex;
}
</style>
<style scoped>
.tabs-view >>> .draggable-item {
  margin: 0 5px;
}

.tabs-view >>> .draggable-item:first-child {
  margin-left: 10px;
}
</style>
