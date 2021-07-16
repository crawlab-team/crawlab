<template>
  <div ref="fileEditor" class="file-editor">
    <div :class="navMenuCollapsed ? 'collapsed' : ''" class="nav-menu">
      <div
          :style="{
            backgroundColor: style.backgroundColorGutters,
            color: style.color,
          }"
          class="nav-menu-top-bar"
      >
        <div class="left">
          <el-input
              v-model="fileSearchString"
              :style="{
                color: style.color,
              }"
              class="search"
              clearable
              placeholder="Search files..."
              prefix-icon="fa fa-search"
              size="mini"
          />
        </div>
        <div class="right">
          <el-tooltip content="Settings">
            <span class="action-icon" @click="showSettings = true">
              <div class="background"/>
              <font-awesome-icon :icon="['fa', 'cog']"/>
            </span>
          </el-tooltip>
          <el-tooltip content="Hide files">
            <span class="action-icon" @click="onToggleNavMenu">
              <div class="background"/>
              <font-awesome-icon :icon="['fa', 'minus']"/>
            </span>
          </el-tooltip>
        </div>
      </div>
      <FileEditorNavMenu
          :active-item="activeFileItem"
          :default-expand-all="!!fileSearchString"
          :items="files"
          :style="style"
          @node-click="onNavItemClick"
          @node-db-click="onNavItemDbClick"
          @node-drop="onNavItemDrop"
          @ctx-menu-new-file="onContextMenuNewFile"
          @ctx-menu-new-directory="onContextMenuNewDirectory"
          @ctx-menu-rename="onContextMenuRename"
          @ctx-menu-clone="onContextMenuClone"
          @ctx-menu-delete="onContextMenuDelete"
          @drop-files="onDropFiles"
      />
    </div>
    <div class="file-editor-content">
      <FileEditorNavTabs
          ref="navTabs"
          :active-tab="activeFileItem"
          :tabs="tabs"
          :style="style"
          @tab-click="onTabClick"
          @tab-close="onTabClose"
          @tab-close-others="onTabCloseOthers"
          @tab-close-all="onTabCloseAll"
          @tab-dragend="onTabDragEnd"
      >
        <template v-if="navMenuCollapsed" #prefix>
          <el-tooltip content="Show files">
            <span class="action-icon expand-files" @click="onToggleNavMenu">
              <div class="background"/>
              <font-awesome-icon :icon="['fa', 'bars']"/>
            </span>
          </el-tooltip>
        </template>
      </FileEditorNavTabs>
      <div
          ref="codeMirrorEditor"
          :class="showCodeMirrorEditor ? '' : 'hidden'"
          :style="{
            scrollbar: style.backgroundColorGutters,
          }"
          class="code-mirror-editor"
      />
      <div
          v-show="!showCodeMirrorEditor"
          :style="{
            backgroundColor: style.backgroundColor,
            color: style.color,
          }"
          class="empty-content"
      >
        You can edit or view a file by double-clicking one of the files on the left.
      </div>
      <template v-if="navTabs && navTabs.showMoreVisible">
        <FileEditorNavTabsShowMoreContextMenu
            :tabs="tabs"
            :visible="showMoreContextMenuVisible"
            @hide="onShowMoreHide"
            @tab-click="onClickShowMoreContextMenuItem"
        >
          <div
              :style="{
                background: style.backgroundColor,
                color: style.color,
              }"
              class="nav-tabs-suffix"
          >
            <el-tooltip content="Show more">
              <span class="action-icon" @click.prevent="onShowMoreShow">
                <div class="background"/>
                <font-awesome-icon :icon="['fa', 'angle-down']"/>
              </span>
            </el-tooltip>
          </div>
        </FileEditorNavTabsShowMoreContextMenu>
      </template>
    </div>
  </div>
  <div ref="codeMirrorTemplate" class="code-mirror-template"/>
  <div ref="styleRef" v-html="extraStyle"/>
  <FileEditorSettingsDialog/>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, onUnmounted, PropType, ref, watch} from 'vue';
import CodeMirror, {Editor, EditorConfiguration, KeyMap} from 'codemirror';
import {MimeType} from 'codemirror/mode/meta';
import {useStore} from 'vuex';
import {getCodemirrorEditor, getCodeMirrorTemplate, initTheme} from '@/utils/codemirror';
import variables from '@/styles/variables.scss';
import {FILE_ROOT} from '@/constants/file';

// codemirror css
import 'codemirror/lib/codemirror.css';

// codemirror mode
import 'codemirror/mode/meta';

// codemirror utils
import '@/utils/codemirror';

// components
import FileEditorNavMenu from '@/components/file/FileEditorNavMenu.vue';
import FileEditorNavTabs from '@/components/file/FileEditorNavTabs.vue';
import FileEditorSettingsDialog from '@/components/file/FileEditorSettingsDialog.vue';
import FileEditorNavTabsShowMoreContextMenu from '@/components/file/FileEditorNavTabsShowMoreContextMenu.vue';

// codemirror mode import cache
const codeMirrorModeCache = new Set<string>();

// codemirror tab content cache
const codeMirrorTabContentCache = new Map<string, string>();

export default defineComponent({
  name: 'FileEditor',
  components: {
    FileEditorSettingsDialog,
    FileEditorNavTabs,
    FileEditorNavMenu,
    FileEditorNavTabsShowMoreContextMenu,
  },
  props: {
    content: {
      type: String,
      required: true,
      default: '',
    },
    activeNavItem: {
      type: Object as PropType<FileNavItem>,
      required: false,
    },
    navItems: {
      type: Array,
      required: true,
      default: () => {
        return [];
      },
    },
  },
  emits: [
    'content-change',
    'tab-click',
    'node-click',
    'node-db-click',
    'node-drop',
    'save-file',
    'ctx-menu-new-file',
    'ctx-menu-new-directory',
    'ctx-menu-rename',
    'ctx-menu-clone',
    'ctx-menu-delete',
    'drop-files',
  ],
  setup(props, {emit}) {
    const ns = 'spider';
    const store = useStore();
    const {file} = store.state as RootStoreState;

    const fileEditor = ref<HTMLDivElement>();

    const codeMirrorEditor = ref<HTMLDivElement>();

    const tabs = ref<FileNavItem[]>([]);

    const activeFileItem = computed<FileNavItem | undefined>(() => props.activeNavItem);

    const style = ref<FileEditorStyle>({});

    const fileSearchString = ref<string>('');

    const navMenuCollapsed = ref<boolean>(false);

    const showSettings = ref<boolean>(false);

    const styleRef = ref<HTMLDivElement>();

    let editor: Editor | null = null;

    let codeMirrorTemplateEditor: Editor | null = null;

    const codeMirrorTemplate = ref<HTMLDivElement>();

    let codeMirrorEditorSearchLabel: HTMLSpanElement | undefined;

    let codeMirrorEditorSearchInput: HTMLInputElement | undefined;

    const navTabs = ref<typeof FileEditorNavTabs>();

    const showMoreContextMenuVisible = ref<boolean>(false);

    const showCodeMirrorEditor = computed<boolean>(() => {
      return !!activeFileItem.value;
    });

    const language = computed<MimeType | undefined>(() => {
      const fileName = activeFileItem.value?.name;
      if (!fileName) return;
      return CodeMirror.findModeByFileName(fileName);
    });

    const languageMime = computed<string | undefined>(() => language.value?.mime);

    const options = computed<FileEditorConfiguration>(() => {
      const {editorOptions} = file as FileStoreState;
      return {
        mode: languageMime.value || 'text',
        ...editorOptions,
      };
    });

    const content = computed<string>(() => {
      const {content} = props as FileEditorProps;
      return content || '';
    });

    const extraStyle = computed<string>(() => {
      return `<style>
.file-editor .file-editor-nav-menu::-webkit-scrollbar {
  background-color: ${style.value.backgroundColor};
  width: 8px;
  height: 8px;
}
.file-editor .file-editor-nav-menu::-webkit-scrollbar-thumb {
  background-color: ${variables.primaryColor};
  border-radius: 4px;
}
.file-editor .file-editor-content .code-mirror-editor .CodeMirror-vscrollbar::-webkit-scrollbar {
  background-color: ${style.value.backgroundColor};
  width: 8px;
}
.file-editor .file-editor-content .code-mirror-editor .CodeMirror-hscrollbar::-webkit-scrollbar {
  background-color: ${style.value.backgroundColor};
  height: 8px;
}
.file-editor .file-editor-content .code-mirror-editor .CodeMirror-vscrollbar::-webkit-scrollbar-thumb,
.file-editor .file-editor-content .code-mirror-editor .CodeMirror-hscrollbar::-webkit-scrollbar-thumb {
  background-color: ${variables.primaryColor};
  border-radius: 4px;
}
.file-editor .file-editor-nav-tabs::-webkit-scrollbar {
  display: none;
}
</style>`;
    });

    const codeMirrorTemplateContent = computed<string>(() => {
      return getCodeMirrorTemplate();
    });

    const updateEditorOptions = () => {
      for (const k in options.value) {
        const key = k as keyof EditorConfiguration;
        const value = options.value[key];
        editor?.setOption(key, value);
      }
    };

    const updateEditorContent = () => {
      editor?.setValue(content.value || '');
    };

    const updateStyle = () => {
      // codemirror style: background color / color / height
      const el = codeMirrorEditor.value as HTMLElement;
      const cm = el.querySelector('.CodeMirror');
      if (!cm) return;
      const computedStyle = window.getComputedStyle(cm);
      style.value = {
        backgroundColor: computedStyle.backgroundColor,
        color: computedStyle.color,
        height: computedStyle.height,
      };

      // gutter
      const cmGutters = el.querySelector('.CodeMirror-gutters');
      if (!cmGutters) return;
      const computedStyleGutters = window.getComputedStyle(cmGutters);
      style.value.backgroundColorGutters = computedStyleGutters.backgroundColor;
    };

    const updateTheme = async () => {
      await initTheme(options.value.theme);
    };

    const updateMode = async () => {
      const mode = language.value?.mode;
      if (!mode || codeMirrorModeCache.has(mode)) return;
      // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
      // @ts-ignore
      await import(`codemirror/mode/${mode}/${mode}.js`);
      codeMirrorModeCache.add(mode);
    };

    const updateSearchInput = () => {
      codeMirrorEditorSearchLabel = codeMirrorEditor.value?.querySelector<HTMLSpanElement>('.CodeMirror-search-label') || undefined;
      codeMirrorEditorSearchInput = codeMirrorEditor.value?.querySelector<HTMLInputElement>('.CodeMirror-search-field') || undefined;
      if (!codeMirrorEditorSearchInput) return;
      codeMirrorEditorSearchInput.onblur = () => {
        if (codeMirrorEditorSearchLabel?.textContent?.includes('Search')) {
          setTimeout(() => {
            codeMirrorEditorSearchInput?.parentElement?.remove();
            editor?.focus();
          }, 10);
        }
      };
    };

    const getContentCache = (tab: FileNavItem) => {
      if (!tab.path) return;
      const key = tab.path;
      const content = codeMirrorTabContentCache.get(key);
      emit('content-change', content as string);
      setTimeout(updateEditorContent, 0);
    };

    const updateContentCache = (tab: FileNavItem, content: string) => {
      if (!tab.path) return;
      const key = tab.path as string;
      codeMirrorTabContentCache.set(key, content as string);
    };

    const deleteContentCache = (tab: FileNavItem) => {
      if (!tab.path) return;
      const key = tab.path;
      codeMirrorTabContentCache.delete(key);
    };

    const deleteOtherContentCache = (tab: FileNavItem) => {
      if (!tab.path) return;
      const key = tab.path;
      const content = codeMirrorTabContentCache.get(key);
      codeMirrorTabContentCache.clear();
      codeMirrorTabContentCache.set(key, content as string);
    };

    const clearContentCache = () => {
      codeMirrorTabContentCache.clear();
    };

    const getFilteredFiles = (items: FileNavItem[]): FileNavItem[] => {
      return items
          .filter(d => {
            if (!d.is_dir) {
              return d.name?.toLowerCase().includes(fileSearchString.value.toLowerCase());
            }
            if (d.children) {
              const children = getFilteredFiles(d.children);
              if (children.length > 0) {
                return true;
              }
            }
            return false;
          })
          .map(d => {
            if (!d.is_dir) return d;
            d.children = getFilteredFiles(d.children || []);
            return d;
          });
    };

    const files = computed<FileNavItem[]>(() => {
      const {navItems} = props as FileEditorProps;
      const root: FileNavItem = {
        path: FILE_ROOT,
        name: FILE_ROOT,
        is_dir: true,
        children: fileSearchString.value ? getFilteredFiles(navItems) : navItems,
      };
      return [root];
    });

    const updateTabs = (item?: FileNavItem) => {
      // add tab
      if (item && !tabs.value.find(t => t.path === item.path)) {
        if (tabs.value.length === 0) {
          store.commit(`${ns}/setActiveFileNavItem`, item);
          editor?.focus();
        }
        tabs.value.push(item);
        getContentCache(item);
      }
    };

    const onNavItemClick = (item: FileNavItem) => {
      emit('node-click', item);
    };

    const onNavItemDbClick = (item: FileNavItem) => {
      store.commit(`${ns}/setActiveFileNavItem`, item);
      emit('node-db-click', item);

      // update tabs
      updateTabs(item);
    };

    const onNavItemDrop = (draggingItem: FileNavItem, dropItem: FileNavItem) => {
      emit('node-drop', draggingItem, dropItem);
    };

    const onContextMenuNewFile = (item: FileNavItem, name: string) => {
      emit('ctx-menu-new-file', item, name);
    };

    const onContextMenuNewDirectory = (item: FileNavItem, name: string) => {
      emit('ctx-menu-new-directory', item, name);
    };

    const onContextMenuRename = (item: FileNavItem, name: string) => {
      emit('ctx-menu-rename', item, name);
    };

    const onContextMenuClone = (item: FileNavItem, name: string) => {
      emit('ctx-menu-clone', item, name);
    };

    const onContextMenuDelete = (item: FileNavItem) => {
      emit('ctx-menu-delete', item);
    };

    const onContentChange = (cm: Editor) => {
      const content = cm.getValue();
      if (!activeFileItem.value) return;
      emit('content-change', content);

      // update in cache
      updateContentCache(activeFileItem.value, content);
    };

    const onTabClick = (tab: FileNavItem) => {
      store.commit(`${ns}/setActiveFileNavItem`, tab);
      emit('tab-click', tab);

      // get from cache and update content
      getContentCache(tab);
    };

    const closeTab = (tab: FileNavItem) => {
      const idx = tabs.value.findIndex(t => t.path === tab.path);
      if (idx !== -1) {
        tabs.value.splice(idx, 1);
      }
      if (activeFileItem.value) {
        if (activeFileItem.value.path === tab.path) {
          if (idx === 0) {
            store.commit(`${ns}/setActiveFileNavItem`, tabs.value[0]);
          } else {
            store.commit(`${ns}/setActiveFileNavItem`, tabs.value[idx - 1]);
          }
        }

        // get from cache
        setTimeout(() => {
          if (activeFileItem.value) {
            getContentCache(activeFileItem.value);
          }
        }, 0);
      }

      // delete in cache
      deleteContentCache(tab);
    };

    const onTabClose = (tab: FileNavItem) => {
      closeTab(tab);
    };

    const onTabCloseOthers = (tab: FileNavItem) => {
      tabs.value = [tab];
      store.commit(`${ns}/setActiveFileNavItem`, tab);

      // clear cache and update current tab content
      deleteOtherContentCache(tab);
    };

    const onTabCloseAll = () => {
      tabs.value = [];
      store.commit(`${ns}/resetActiveFileNavItem`);

      // clear cache
      clearContentCache();
    };

    const onTabDragEnd = (newTabs: FileNavItem[]) => {
      tabs.value = newTabs;
    };

    const onShowMoreShow = () => {
      showMoreContextMenuVisible.value = true;
    };

    const onShowMoreHide = () => {
      showMoreContextMenuVisible.value = false;
    };

    const onClickShowMoreContextMenuItem = (tab: FileNavItem) => {
      store.commit(`${ns}/setActiveFileNavItem`, tab);
      emit('tab-click', tab);
    };

    const keyMapSave = () => {
      if (!activeFileItem.value) return;
      emit('save-file', activeFileItem.value);
    };

    const keyMapClose = () => {
      if (!activeFileItem.value) return;
      closeTab(activeFileItem.value);
    };

    const addSaveKeyMap = (cm: Editor) => {
      const map = {
        'Cmd-S': keyMapSave,
        'Ctrl-S': keyMapSave,
        // 'Cmd-W': keyMapClose,
        'Ctrl-W': keyMapClose,
      } as KeyMap;
      cm.addKeyMap(map);
    };

    const onToggleNavMenu = () => {
      navMenuCollapsed.value = !navMenuCollapsed.value;
    };

    const listenToKeyboardEvents = () => {
      editor?.on('blur', () => {
        updateSearchInput();
      });
    };

    const unlistenToKeyboardEvents = () => {
      document.onkeydown = null;
    };

    watch(options, async () => {
      await Promise.all([
        updateMode(),
        updateTheme(),
      ]);
      updateEditorOptions();
      updateStyle();
    });

    const onDropFiles = (files: InputFile[]) => {
      emit('drop-files', files);
    };

    onMounted(async () => {
      // init codemirror editor
      const el = codeMirrorEditor.value as HTMLElement;
      editor = getCodemirrorEditor(el, options.value);

      // add save key map
      addSaveKeyMap(editor);

      // on editor change
      editor.on('change', onContentChange);

      // update editor options
      updateEditorOptions();

      // update editor content
      updateEditorContent();

      // update editor theme
      await updateTheme();

      // update styles
      updateStyle();

      // listen to keyboard events key
      listenToKeyboardEvents();

      // init codemirror template
      const elTemplate = codeMirrorTemplate.value as HTMLElement;
      codeMirrorTemplateEditor = getCodemirrorEditor(elTemplate, options.value);
      codeMirrorTemplateEditor.setValue(codeMirrorTemplateContent.value);
      codeMirrorTemplateEditor.setOption('mode', 'text/x-python');
    });

    onUnmounted(() => {
      // turnoff listening to keyboard events
      unlistenToKeyboardEvents();
    });

    return {
      fileEditor,
      codeMirrorEditor,
      tabs,
      activeFileItem,
      fileSearchString,
      navMenuCollapsed,
      styleRef,
      codeMirrorTemplate,
      showSettings,
      showCodeMirrorEditor,
      navTabs,
      showMoreContextMenuVisible,
      languageMime,
      options,
      style,
      files,
      extraStyle,
      variables,
      onNavItemClick,
      onNavItemDbClick,
      onNavItemDrop,
      onContextMenuNewFile,
      onContextMenuNewDirectory,
      onContextMenuRename,
      onContextMenuClone,
      onContextMenuDelete,
      onContentChange,
      onTabClick,
      onTabClose,
      onTabCloseOthers,
      onTabCloseAll,
      onTabDragEnd,
      onToggleNavMenu,
      onShowMoreShow,
      onShowMoreHide,
      onClickShowMoreContextMenuItem,
      updateTabs,
      updateEditorContent,
      updateContentCache,
      onDropFiles,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.file-editor {
  height: 100%;
  display: flex;

  .nav-menu {
    flex-basis: $fileEditorNavMenuWidth;
    min-width: $fileEditorNavMenuWidth;
    display: flex;
    flex-direction: column;
    transition: all $fileEditorNavMenuCollapseTransitionDuration;

    &.collapsed {
      min-width: 0;
      flex-basis: 0;
      overflow: hidden;
    }

    .nav-menu-top-bar {
      flex-basis: $fileEditorNavMenuTopBarHeight;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-size: 12px;
      padding: 0 10px 0 0;

      .left,
      .right {
        display: flex;
      }
    }
  }

  .file-editor-content {
    position: relative;
    flex: 1;
    display: flex;
    min-width: calc(100% - #{$fileEditorNavMenuWidth});
    flex-direction: column;

    .code-mirror-editor {
      flex: 1;

      &.hidden {
        position: fixed;
        top: -100vh;
        left: 0;
        height: 100vh;
      }
    }

    .empty-content {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .nav-tabs-suffix {
      width: 30px;
      position: absolute;
      top: 0;
      right: 0;
      z-index: 5;
      display: flex;
      align-items: center;
      justify-content: center;
      height: $fileEditorNavTabsHeight;
    }
  }

  .action-icon {
    position: relative;
    height: 16px;
    width: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-size: 12px;

    &:hover {
      .background {
        background-color: $fileEditorMaskBg;
        border-radius: 8px;
      }
    }

    &.expand-files {
      width: 29px;
      text-align: center;
    }

    .background {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
    }
  }
}

.code-mirror-template {
  position: fixed;
  top: -100vh;
  left: 0;
  height: 100vh;
}
</style>
<style scoped>
.file-editor .nav-menu .nav-menu-top-bar >>> .search.el-input .el-input__inner {
  border: none;
  background: transparent;
  color: inherit;
}

.file-editor .file-editor-content .code-mirror-editor >>> .CodeMirror {
  position: relative;
  min-height: 100%;
}

.file-editor .file-editor-content .code-mirror-editor >>> .CodeMirror.dialog-opened {
  position: relative;
}

.file-editor .file-editor-content .code-mirror-editor >>> .CodeMirror-dialog {
  font-size: 14px;
}

.file-editor .file-editor-content .code-mirror-editor >>> .CodeMirror-dialog.CodeMirror-dialog-top {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 2;
}

.file-editor .file-editor-content .code-mirror-editor >>> .CodeMirror-dialog.CodeMirror-dialog-bottom {
  position: absolute;
  bottom: 10px;
  right: 10px;
  z-index: 2;
}

.file-editor .file-editor-content .code-mirror-editor >>> .CodeMirror-search-field {
  background-color: transparent;
  border: 1px solid;
  color: inherit;
  outline: none;
}
</style>
