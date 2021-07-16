<template>
  <div
      :style="{
        backgroundColor: style.backgroundColorGutters,
        borderRight: `1px solid ${style.backgroundColor}`
      }"
      ref="fileEditorNavMenu"
      class="file-editor-nav-menu"
  >
    <el-tree
        ref="tree"
        :render-after-expand="defaultExpandAll"
        :data="items"
        :expand-on-click-node="false"
        :highlight-current="false"
        :allow-drop="allowDrop"
        empty-text="No files available"
        icon-class="fa fa-angle-right"
        :style="{
          backgroundColor: style.backgroundColorGutters,
          color: style.color,
        }"
        node-key="path"
        :default-expanded-keys="defaultExpandedKeys"
        draggable
        @node-drag-enter="onNodeDragEnter"
        @node-drag-leave="onNodeDragLeave"
        @node-drag-end="onNodeDragEnd"
        @node-drop="onNodeDrop"
        @node-click="onNodeClick"
        @node-contextmenu="onNodeContextMenuShow"
        @node-expand="onNodeExpand"
        @node-collapse="onNodeCollapse"
    >
      <template #default="{ data }">
        <FileEditorNavMenuContextMenu
            :clicking="contextMenuClicking"
            :visible="isShowContextMenu(data)"
            @hide="onNodeContextMenuHide"
            @clone="onNodeContextMenuClone(data)"
            @delete="onNodeContextMenuDelete(data)"
            @rename="onNodeContextMenuRename(data)"
            @new-file="onNodeContextMenuNewFile(data)"
            @new-directory="onNodeContextMenuNewDirectory(data)"
        >
          <div
              v-bind="getBindDir(data)"
              :class="getItemClass(data)"
              class="nav-item-wrapper"
          >
            <div class="background"/>
            <div class="nav-item">
              <span class="icon">
                <atom-material-icon :is-dir="data.is_dir" :name="data.name"/>
              </span>
              <span class="title">
                {{ data.name }}
              </span>
            </div>
          </div>
        </FileEditorNavMenuContextMenu>
      </template>
    </el-tree>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, onUnmounted, reactive, ref} from 'vue';
import {ClickOutside} from 'element-plus/lib/directives';
import Node from 'element-plus/lib/el-tree/src/model/node';
import {DropType} from 'element-plus/lib/el-tree/src/tree.type';
import AtomMaterialIcon from '@/components/icon/AtomMaterialIcon.vue';
import {KEY_CONTROL, KEY_META} from '@/constants/keyboard';
import FileEditorNavMenuContextMenu from '@/components/file/FileEditorNavMenuContextMenu.vue';
import {ElMessageBox, ElTree} from 'element-plus';
import {useDropzone} from 'vue3-dropzone';

export default defineComponent({
  name: 'FileEditorNavMenu',
  components: {
    FileEditorNavMenuContextMenu,
    AtomMaterialIcon
  },
  directives: {
    ClickOutside,
  },
  props: {
    activeItem: {
      type: Object,
      required: false,
    },
    items: {
      type: Array,
      required: true,
      default: () => {
        return [];
      },
    },
    defaultExpandAll: {
      type: Boolean,
      required: true,
      default: false,
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
    'node-click',
    'node-db-click',
    'node-drop',
    'ctx-menu-new-file',
    'ctx-menu-new-directory',
    'ctx-menu-rename',
    'ctx-menu-clone',
    'ctx-menu-delete',
    'drop-files',
  ],
  setup(props, ctx) {
    const {emit} = ctx;

    const tree = ref<typeof ElTree>();

    const fileEditorNavMenu = ref<HTMLDivElement>();

    const clickStatus = reactive<FileEditorNavMenuClickStatus>({
      clicked: false,
      item: undefined,
    });

    const selectedCache = reactive<FileEditorNavMenuCache<boolean>>({});

    const dragCache = reactive<FileEditorNavMenuCache<boolean>>({});

    const isCtrlKeyPressed = ref<boolean>(false);

    const activeContextMenuItem = ref<FileNavItem>();

    const contextMenuClicking = ref<boolean>(false);

    const expandedKeys = ref<string[]>([]);

    const defaultExpandedKeys = computed<string[]>(() => {
      return ['~'].concat(expandedKeys.value);
    });

    const addDefaultExpandedKey = (key: string) => {
      if (!expandedKeys.value.includes(key)) expandedKeys.value.push(key);
    };

    const removeDefaultExpandedKey = (key: string) => {
      if (!expandedKeys.value.includes(key)) return;
      const idx = expandedKeys.value.indexOf(key);
      expandedKeys.value.splice(idx, 1);
    };

    const resetDefaultExpandedKeys = () => {
      expandedKeys.value = [];
    };

    const resetClickStatus = () => {
      clickStatus.clicked = false;
      clickStatus.item = undefined;
      activeContextMenuItem.value = undefined;
    };

    const updateSelectedMap = (item: FileNavItem) => {
      const key = item.path;
      if (!key) {
        console.warn('No path specified for FileNavItem');
        return;
      }
      if (!selectedCache[key]) {
        selectedCache[key] = false;
      }
      selectedCache[key] = !selectedCache[key];

      // if Ctrl key is not pressed, clear other selection
      if (!isCtrlKeyPressed.value) {
        Object.keys(selectedCache).filter(k => k !== key).forEach(k => {
          selectedCache[k] = false;
        });
      }
    };

    const onNodeClick = (item: FileNavItem) => {
      if (clickStatus.clicked && clickStatus.item?.path === item.path) {
        emit('node-db-click', item);
        updateSelectedMap(item);
        resetClickStatus();
        return;
      }

      clickStatus.item = item;
      clickStatus.clicked = true;
      setTimeout(() => {
        if (clickStatus.clicked) {
          emit('node-click', item);
          updateSelectedMap(item);
        }
        resetClickStatus();
      }, 200);
    };

    const onNodeContextMenuShow = (ev: Event, item: FileNavItem) => {
      contextMenuClicking.value = true;
      activeContextMenuItem.value = item;
      setTimeout(() => {
        contextMenuClicking.value = false;
      }, 500);
    };

    const onNodeContextMenuHide = () => {
      activeContextMenuItem.value = undefined;
    };

    const onNodeContextMenuNewFile = async (item: FileNavItem) => {
      const res = await ElMessageBox.prompt('Please enter the name of the new file', 'New File');
      emit('ctx-menu-new-file', item, res.value);
    };

    const onNodeContextMenuNewDirectory = async (item: FileNavItem) => {
      const res = await ElMessageBox.prompt('Please enter the name of the new directory', 'New Directory');
      emit('ctx-menu-new-directory', item, res.value);
    };

    const onNodeContextMenuRename = async (item: FileNavItem) => {
      const res = await ElMessageBox.prompt('Please enter the new name', 'Rename');
      emit('ctx-menu-rename', item, res.value);
    };

    const onNodeContextMenuClone = async (item: FileNavItem) => {
      const res = await ElMessageBox.prompt('Please enter the new name', 'Clone');
      emit('ctx-menu-clone', item, res.value);
    };

    const onNodeContextMenuDelete = async (item: FileNavItem) => {
      await ElMessageBox.confirm('Are you sure to delete?', 'Delete');
      emit('ctx-menu-delete', item);
    };

    const onNodeDragEnter = (draggingNode: Node, dropNode: Node) => {
      const item = dropNode.data as FileNavItem;
      if (!item.path) return;
      dragCache[item.path] = true;
    };

    const onNodeDragLeave = (draggingNode: Node, dropNode: Node) => {
      const item = dropNode.data as FileNavItem;
      if (!item.path) return;
      dragCache[item.path] = false;
    };

    const onNodeDragEnd = () => {
      for (const key in dragCache) {
        dragCache[key] = false;
      }
    };

    const onNodeDrop = (draggingNode: Node, dropNode: Node) => {
      const draggingItem = draggingNode.data as FileNavItem;
      const dropItem = dropNode.data as FileNavItem;
      emit('node-drop', draggingItem, dropItem);
    };

    const onNodeExpand = (data: FileNavItem) => {
      addDefaultExpandedKey(data.path as string);
    };

    const onNodeCollapse = (data: FileNavItem) => {
      removeDefaultExpandedKey(data.path as string);
    };

    const isSelected = (item: FileNavItem): boolean => {
      if (!item.path) return false;
      return selectedCache[item.path] || false;
    };

    const isDroppable = (item: FileNavItem): boolean => {
      if (!item.path) return false;
      return dragCache[item.path] || false;
    };

    const isShowContextMenu = (item: FileNavItem) => {
      return activeContextMenuItem.value?.path === item.path;
    };

    const allowDrop = (draggingNode: Node, dropNode: Node, type: DropType) => {
      if (type !== 'inner') return false;
      if (draggingNode.data?.path === dropNode.data?.path) return false;
      if (draggingNode.parent?.data?.path === dropNode.data?.path) return false;
      const item = dropNode.data as FileNavItem;
      return item.is_dir;
    };

    const getItemClass = (item: FileNavItem): string[] => {
      const cls = [];
      if (isSelected(item)) cls.push('selected');
      if (isDroppable(item)) cls.push('droppable');
      return cls;
    };

    const {
      getRootProps,
    } = useDropzone({
      onDrop: (files: InputFile[]) => {
        emit('drop-files', files);
      },
    });

    const getBindDir = (item: FileNavItem) => getRootProps({
      onDragEnter: (ev: DragEvent) => {
        ev.stopPropagation();
        if (!item.is_dir || !item.path) return;
        dragCache[item.path] = true;
      },
      onDragLeave: (ev: DragEvent) => {
        ev.stopPropagation();
        if (!item.is_dir || !item.path) return;
        dragCache[item.path] = false;
      },
      onDrop: () => {
        for (const key in dragCache) {
          dragCache[key] = false;
        }
      },
    });

    onMounted(() => {
      // listen to keyboard events
      document.onkeydown = (ev: KeyboardEvent) => {
        if (!ev) return;
        if (ev.key === KEY_CONTROL || ev.key === KEY_META) {
          isCtrlKeyPressed.value = true;
        }
      };
      document.onkeyup = (ev: KeyboardEvent) => {
        if (!ev) return;
        if (ev.key === KEY_CONTROL || ev.key === KEY_META) {
          isCtrlKeyPressed.value = false;
        }
      };
    });

    onUnmounted(() => {
      // turnoff listening to keyboard events
      document.onkeydown = null;
      document.onkeyup = null;
    });

    return {
      tree,
      activeContextMenuItem,
      fileEditorNavMenu,
      contextMenuClicking,
      defaultExpandedKeys,
      onNodeClick,
      onNodeContextMenuShow,
      onNodeContextMenuHide,
      onNodeContextMenuNewFile,
      onNodeContextMenuNewDirectory,
      onNodeContextMenuRename,
      onNodeContextMenuClone,
      onNodeContextMenuDelete,
      onNodeDragEnter,
      onNodeDragLeave,
      onNodeDragEnd,
      onNodeDrop,
      onNodeExpand,
      onNodeCollapse,
      isSelected,
      isDroppable,
      isShowContextMenu,
      allowDrop,
      getItemClass,
      resetDefaultExpandedKeys,
      addDefaultExpandedKey,
      removeDefaultExpandedKey,
      getBindDir,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.file-editor-nav-menu {
  flex: 1;
  max-height: 100%;
  overflow: auto;

  .el-tree {
    height: 100%;
    min-width: 100%;
    max-width: fit-content;

    .el-tree-node {
      .nav-item-wrapper {
        z-index: 2;

        & * {
          pointer-events: none;
        }

        &.selected {
          .background {
            background-color: $fileEditorMaskBg;
          }

          .nav-item {
            color: $fileEditorNavMenuItemSelectedColor;
          }
        }

        &.droppable {
          & * {
            pointer-events: none;
          }

          .nav-item {
            border: 1px dashed $fileEditorNavMenuItemDragTargetBorderColor;
          }
        }

        .nav-item:hover,
        .background:hover + .nav-item {
          color: $fileEditorNavMenuItemSelectedColor;
        }

        .background {
          position: absolute;
          width: 100%;
          height: 100%;
          left: 0;
          top: 0;
          z-index: -1;
        }

        .nav-item {
          display: flex;
          align-items: center;
          font-size: 14px;
          user-select: none;

          .icon {
            margin-right: 5px;
          }
        }
      }
    }
  }
}
</style>
<style scoped>
.file-editor-nav-menu >>> .el-tree .el-tree-node > .el-tree-node__content {
  background-color: inherit;
  position: relative;
  z-index: 0;
}

.file-editor-nav-menu >>> .el-tree .el-tree-node > .el-tree-node__content .el-tree-node__expand-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  font-size: 16px;
  padding: 0;
  margin: 0;
}

.file-editor-nav-menu >>> .el-tree .el-tree-node * {
  transition: none;
}
</style>
