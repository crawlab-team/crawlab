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
  name: 'FileEditorNavTabsContextMenu',
  components: {ContextMenuList, ContextMenu},
  props: contextMenuDefaultProps,
  emits: [
    'hide',
    'close',
    'close-others',
    'close-all',
  ],
  setup(props, {emit}) {
    const items = readonly<ContextMenuItem[]>([
      {title: 'Close', icon: ['fa', 'times'], action: () => emit('close')},
      {title: 'Close Others', action: () => emit('close-others')},
      {title: 'Close All', action: () => emit('close-all')},
    ]);

    return {
      items,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
