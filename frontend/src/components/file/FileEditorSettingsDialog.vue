<template>
  <div class="file-editor-settings-dialog">
    <el-dialog
        :model-value="visible"
        title="File Editor Settings"
        @close="onClose"
    >
      <el-menu :default-active="activeTabName" class="nav-menu" mode="horizontal" @select="onTabChange">
        <el-menu-item v-for="tab in tabs" :key="tab.name" :index="tab.name">
          {{ tab.title }}
        </el-menu-item>
      </el-menu>
      <el-form
          :label-width="variables.fileEditorSettingsDialogLabelWidth"
          class="form"
          size="small"
      >
        <el-form-item
            v-for="name in optionNames[activeTabName]"
            :key="name"
        >
          <template #label>
            <el-tooltip :content="getDefinitionDescription(name)" popper-class="help-tooltip" trigger="click">
              <font-awesome-icon :icon="['far', 'question-circle']" class="icon" size="sm"/>
            </el-tooltip>
            {{ getDefinitionTitle(name) }}
          </template>
          <FileEditorSettingsFormItem v-model="options[name]" :name="name"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button plain size="small" type="info" @click="onClose">Cancel</el-button>
        <el-button size="small" type="primary" @click="onConfirm">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onBeforeMount, readonly, ref} from 'vue';
import {useStore} from 'vuex';
import {plainClone} from '@/utils/object';
import variables from '@/styles/variables.scss';
import {getOptionDefinition, getThemes} from '@/utils/codemirror';
import FileEditorSettingsFormItem from '@/components/file/FileEditorSettingsFormItem.vue';
import {onBeforeRouteLeave} from 'vue-router';

export default defineComponent({
  name: 'FileEditorSettingsDialog',
  components: {FileEditorSettingsFormItem},
  setup() {
    const storeNamespace = 'file';
    const store = useStore();
    const {file} = store.state as RootStoreState;

    const options = ref<FileEditorConfiguration>({});

    const tabs = readonly([
      {name: 'general', title: 'General'},
      {name: 'edit', title: 'Edit'},
      {name: 'indentation', title: 'Indentation'},
      {name: 'cursor', title: 'Cursor'},
    ]);

    const optionNames = readonly({
      general: [
        'theme',
        'keyMap',
        'lineWrapping',
        'lineNumbers',
        'maxHighlightLength',
        'spellcheck',
        'autocorrect',
        'autocapitalize',
      ],
      edit: [
        'lineWiseCopyCut',
        'pasteLinesPerSelection',
        'undoDepth',
      ],
      indentation: [
        'indentUnit',
        'smartIndent',
        'tabSize',
        'indentWithTabs',
        'electricChars',
      ],
      cursor: [
        'showCursorWhenSelecting',
        'cursorBlinkRate',
        'cursorScrollMargin',
        'cursorHeight',
      ],
    });

    const activeTabName = ref<string>(tabs[0].name);

    const visible = computed<boolean>(() => {
      const {editorSettingsDialogVisible} = file;
      return editorSettingsDialogVisible;
    });

    const themes = computed<string[]>(() => {
      return getThemes();
    });

    const resetOptions = () => {
      const {editorOptions} = file;
      options.value = plainClone(editorOptions);
    };

    const onClose = () => {
      store.commit(`${storeNamespace}/setEditorSettingsDialogVisible`, false);
      resetOptions();
    };

    const onConfirm = () => {
      store.commit(`${storeNamespace}/setEditorOptions`, options.value);
      store.commit(`${storeNamespace}/setEditorSettingsDialogVisible`, false);
      resetOptions();
    };

    const onTabChange = (tabName: string) => {
      activeTabName.value = tabName;
    };

    const getDefinitionDescription = (name: string) => {
      return getOptionDefinition(name)?.description;
    };

    const getDefinitionTitle = (name: string) => {
      return getOptionDefinition(name)?.title;
    };

    onBeforeMount(() => {
      resetOptions();
    });

    onBeforeRouteLeave(() => {
      store.commit(`${storeNamespace}/setEditorSettingsDialogVisible`, false);
    });

    return {
      variables,
      options,
      activeTabName,
      tabs,
      optionNames,
      visible,
      themes,
      onClose,
      onConfirm,
      onTabChange,
      getDefinitionDescription,
      getDefinitionTitle,
    };
  },
});
</script>

<style lang="scss" scoped>
.file-editor-settings-dialog {
  .nav-menu {
    .el-menu-item {
      height: 40px;
      line-height: 40px;
    }
  }

  .form {
    margin: 20px;
  }
}
</style>

<style scoped>
.file-editor-settings-dialog >>> .el-dialog .el-dialog__body {
  padding: 10px 20px;
}

.file-editor-settings-dialog >>> .el-form-item > .el-form-item__label .icon {
  cursor: pointer;
}

.file-editor-settings-dialog >>> .el-form-item > .el-form-item__content {
  width: 240px;
}

.file-editor-settings-dialog >>> .el-form-item > .el-form-item__content .el-input,
.file-editor-settings-dialog >>> .el-form-item > .el-form-item__content .el-select {
  width: 100%;
}
</style>
<style>
.help-tooltip {
  max-width: 240px;
}
</style>
