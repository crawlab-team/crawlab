<template>
  <div
      :class="classes"
      :draggable="true"
      class="tab"
      @click="onClick"
  >
    <span :key="item?.icon || icon" class="icon">
      <MenuItemIcon v-if="!icon" :item="item" size="10px"/>
      <Icon v-else :icon="icon" size="10px"/>
    </span>
    <span v-if="showTitle" class="title">
      {{ title }}
    </span>
    <span v-if="showClose" class="close-btn" @click.stop="onClose">
      <i class="el-icon-close"></i>
    </span>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import MenuItemIcon from '@/components/icon/MenuItemIcon.vue';
import {useStore} from 'vuex';
import {getPrimaryPath} from '@/utils/path';
// import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import Icon from '@/components/icon/Icon.vue';

export default defineComponent({
  name: 'Tab',
  components: {
    Icon,
    MenuItemIcon,
  },
  props: {
    tab: {
      type: Object as PropType<Tab>,
    },
    icon: {
      type: [String, Array] as PropType<Icon>
    },
    showTitle: {
      type: Boolean,
      default: true,
    },
    showClose: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    }
  },
  emits: [
    'click',
  ],
  setup(props: TabProps, {emit}) {
    // const {tm} = useI18n();
    const router = useRouter();
    const storeNamespace = 'layout';
    const store = useStore();
    const {layout: state} = store.state as RootStoreState;

    const item = computed(() => {
      const {tab} = props as TabProps;
      if (!tab) return;
      const {menuItems} = state;
      for (const _item of menuItems) {
        const primaryPath = getPrimaryPath(tab.path);
        if (primaryPath === _item.path) {
          return _item;
        }
      }
    });

    const title = computed(() => {
      // TODO: detailed title
      // return item.value?.title || tm('No Title');
      return item.value?.title;
    });

    const active = computed(() => {
      const {tab} = props as TabProps;
      const {activeTabId} = state;
      return tab?.id === activeTabId;
    });

    const dragging = computed<boolean>(() => {
      const {tab} = props as TabProps;
      return !!tab?.dragging;
    });

    const isTabsDragging = computed<boolean>(() => state.isTabsDragging);

    const classes = computed(() => {
      const cls = [];
      if (active.value) cls.push('active');
      if (dragging.value) cls.push('dragging');
      if (isTabsDragging.value) cls.push('is-tabs-dragging');
      return cls;
    });

    const onClick = () => {
      emit('click');
      const {tab} = props as TabProps;
      if (!tab) return;
      store.commit(`${storeNamespace}/setActiveTabId`, tab.id);
      router.push(tab.path);
    };

    const onClose = () => {
      // current tab
      const {tab} = props as TabProps;
      if (!tab) return;

      // tabs
      const {tabs} = state;

      // index of current tab (to be removed)
      const idx = tabs.findIndex(d => d.id === tab.id);

      // remove tab
      store.commit(`${storeNamespace}/removeTab`, tab);

      // after-remove actions
      if (active.value) {
        if (tabs.length === 0) {
          const newTab: Tab = {path: '/'};
          store.commit(`${storeNamespace}/addTab`, newTab);
          store.commit(`${storeNamespace}/setActiveTabId`, newTab.id);
          router.push(newTab.path);
        } else if (idx === 0) {
          router.push(tabs[0].path);
          store.commit(`${storeNamespace}/setActiveTabId`, tabs[0].id);
        } else {
          router.push(tabs[idx - 1].path);
          store.commit(`${storeNamespace}/setActiveTabId`, tabs[idx - 1].id);
        }
      }
    };

    return {
      item,
      title,
      active,
      dragging,
      isTabsDragging,
      classes,
      onClick,
      onClose,
    };
  },
});
</script>
<style lang="scss" scoped>
@import "../../styles/variables";

.tab {
  display: flex;
  align-items: center;
  padding: 3px 5px;
  max-width: $tabsViewTabMaxWidth;
  border: 1px solid $tabsViewTabBorderColor;
  cursor: pointer;
  background-color: $tabsViewTabBg;
  user-select: none;
  color: $tabsViewTabColor;

  &.disabled {
    cursor: not-allowed;
    background-color: $disabledBgColor;
    border-color: $disabledBorderColor;
    color: $disabledColor;
  }

  &:focus:not(.disabled) {
    color: inherit;
  }

  &:hover:not(.disabled) {
    &.dragging {
      .title,
      .icon {
        color: inherit;
      }
    }

    .title,
    .icon {
      color: $tabsViewActiveTabColor;
    }
  }

  &.active {
    color: $tabsViewActiveTabColor;
    border-color: $tabsViewActiveTabColor;
    background-color: $tabsViewActiveTabPlainColor;
  }

  .close-btn,
  .icon {
    font-weight: 100;
    display: flex;
    align-items: center;
    font-size: 12px;
  }

  .title {
    display: flex;
    align-items: center;
    margin: 0 3px;
    font-size: 12px;
    height: $tabsViewTabHeight;
  }
}
</style>
