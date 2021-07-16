<template>
  <ContextMenu :clicking="clicking" :placement="placement" :visible="visible" @hide="$emit('hide')">
    <template #default>
      <ContextMenuList :items="items" @hide="$emit('hide')"/>
    </template>
    <template #reference>
      <slot></slot>
    </template>
  </ContextMenu>
</template>

<script lang="ts">
import {defineComponent, readonly} from 'vue';
import ContextMenu, {contextMenuDefaultProps} from '@/components/context-menu/ContextMenu.vue';
import ContextMenuList from '@/components/context-menu/ContextMenuList.vue';

export default defineComponent({
  name: 'FileEditorNavMenuContextMenu',
  components: {ContextMenuList, ContextMenu},
  props: contextMenuDefaultProps,
  emits: [
    'hide',
    'new-file',
    'new-directory',
    'rename',
    'clone',
    'delete',
  ],
  setup(props, {emit}) {
    const items = readonly<ContextMenuItem[]>([
      {title: 'New File', icon: ['fa', 'file-alt'], action: () => emit('new-file')},
      {title: 'New Directory', icon: ['fa', 'folder-plus'], action: () => emit('new-directory')},
      {title: 'Rename', icon: ['fa', 'edit'], action: () => emit('rename')},
      {title: 'Duplicate', icon: ['fa', 'clone'], action: () => emit('clone')},
      {title: 'Delete', icon: ['fa', 'trash'], action: () => emit('delete')},
    ]);

    return {
      items,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
