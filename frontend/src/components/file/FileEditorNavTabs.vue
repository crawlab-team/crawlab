<template>
  <div
      ref="navTabs"
      :style="{
        backgroundColor: style.backgroundColorGutters,
        color: style.color,
      }"
      class="file-editor-nav-tabs"
  >
    <slot name="prefix"></slot>
    <DraggableList
        item-key="path"
        :items="tabs"
        @d-end="onDragEnd"
    >
      <template v-slot="{item}">
        <FileEditorNavTabsContextMenu
            :clicking="contextMenuClicking"
            :visible="isShowContextMenu(item)"
            @close="onClose(item)"
            @hide="onContextMenuHide"
            @close-others="onCloseOthers(item)"
            @close-all="onCloseAll"
        >
          <div
              :class="activeTab && activeTab.path === item.path ? 'active' : ''"
              :style="{
                backgroundColor: style.backgroundColor,
              }"
              class="file-editor-nav-tab"
              @click="onClick(item)"
              @contextmenu.prevent="onContextMenuShow(item)"
          >
            <span class="icon">
              <atom-material-icon :is-dir="item.is_dir" :name="item.name"/>
            </span>
            <el-tooltip :content="getTitle(item)" :show-after="500">
              <span class="title">
                {{ getTitle(item) }}
              </span>
            </el-tooltip>
            <span class="close-btn" @click.stop="onClose(item)">
              <i class="el-icon-close"></i>
            </span>
            <div class="background"/>
          </div>
        </FileEditorNavTabsContextMenu>
      </template>
    </DraggableList>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, ref, watch} from 'vue';
import DraggableList from '@/components/drag/DraggableList.vue';
import AtomMaterialIcon from '@/components/icon/AtomMaterialIcon.vue';
import FileEditorNavTabsContextMenu from '@/components/file/FileEditorNavTabsContextMenu.vue';

export default defineComponent({
  name: 'FileEditorNavTabs',
  components: {FileEditorNavTabsContextMenu, AtomMaterialIcon, DraggableList},
  props: {
    activeTab: {
      type: Object,
      required: false,
    },
    tabs: {
      type: Array,
      required: true,
      default: () => {
        return [];
      },
    },
    style: {
      type: Object,
      required: false,
      default: () => {
        return {};
      },
    },
  },
  emits: [
    'tab-click',
    'tab-close',
    'tab-close-others',
    'tab-close-all',
    'tab-dragend',
    'show-more',
  ],
  setup(props, {emit}) {
    const activeContextMenuItem = ref<FileNavItem>();

    const navTabs = ref<HTMLDivElement>();

    const navTabsWidth = ref<number>();

    const navTabsOverflowWidth = ref<number>();

    const showMoreVisible = computed<boolean>(() => {
      if (navTabsWidth.value === undefined || navTabsOverflowWidth.value === undefined) return false;
      return navTabsOverflowWidth.value > navTabsWidth.value;
    });

    const contextMenuClicking = ref<boolean>(false);

    const tabs = computed<FileNavItem[]>(() => {
      const {tabs} = props as FileEditorNavTabsProps;
      return tabs;
    });

    const getTitle = (item: FileNavItem) => {
      return item.name;
    };

    const onClick = (item: FileNavItem) => {
      emit('tab-click', item);
    };

    const onClose = (item: FileNavItem) => {
      emit('tab-close', item);
    };

    const onCloseOthers = (item: FileNavItem) => {
      emit('tab-close-others', item);
    };

    const onCloseAll = () => {
      emit('tab-close-all');
    };

    const onDragEnd = (items: FileNavItem[]) => {
      emit('tab-dragend', items);
    };

    const onContextMenuShow = (item: FileNavItem) => {
      contextMenuClicking.value = true;
      activeContextMenuItem.value = item;

      setTimeout(() => {
        contextMenuClicking.value = false;
      }, 500);
    };

    const onContextMenuHide = () => {
      activeContextMenuItem.value = undefined;
    };

    const isShowContextMenu = (item: FileNavItem) => {
      return activeContextMenuItem.value?.path === item.path;
    };

    const updateWidths = () => {
      if (!navTabs.value) return;

      // width
      navTabsWidth.value = Number(getComputedStyle(navTabs.value).width.replace('px', ''));

      // overflow width
      const el = navTabs.value.querySelector('.draggable-list');
      if (el) {
        navTabsOverflowWidth.value = Number(getComputedStyle(el).width.replace('px', ''));
      }
    };

    watch(tabs.value, () => {
      setTimeout(updateWidths, 100);
    });

    onMounted(() => {
      // update tabs widths
      updateWidths();
    });

    return {
      activeContextMenuItem,
      navTabs,
      navTabsWidth,
      navTabsOverflowWidth,
      showMoreVisible,
      contextMenuClicking,
      getTitle,
      onClick,
      onClose,
      onCloseOthers,
      onCloseAll,
      onDragEnd,
      onContextMenuShow,
      onContextMenuHide,
      isShowContextMenu,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.file-editor-nav-tabs {
  position: relative;
  display: flex;
  align-items: center;
  overflow: auto;
  height: $fileEditorNavTabsHeight;

  .file-editor-nav-tab {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: left;
    height: $fileEditorNavTabsHeight;
    max-width: $fileEditorNavTabsItemMaxWidth;
    white-space: nowrap;
    text-overflow: ellipsis;
    padding: 0 10px;
    font-size: 14px;
    cursor: pointer;
    box-sizing: border-box;
    z-index: 1;

    &:hover {
      .background {
        background-color: $fileEditorMaskBg;
      }
    }

    &.active {
      border-bottom: 2px solid $primaryColor;

      .title {
        color: $fileEditorNavTabsItemActiveColor;
      }
    }

    .background {
      position: absolute;
      left: 0;
      top: 0;
      width: 100%;
      height: 100%;
      z-index: 0;
    }

    .icon {
      margin-right: 5px;
      z-index: 1;
    }

    .title {
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      color: $fileEditorNavTabsItemColor;
      z-index: 1;
    }

    .close-btn {
      margin-left: 5px;
      z-index: 1;
    }
  }

  .suffix {
    position: absolute;
    right: 0;
    top: 0;
    height: 100%;
    display: flex;
    align-items: center;
  }
}
</style>
