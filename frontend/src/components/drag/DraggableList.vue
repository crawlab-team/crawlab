<template>
  <div class="draggable-list">
    <DraggableItem
        v-for="(item, $index) in orderedItems"
        :key="item[itemKey] === undefined ? $index : item[itemKey]"
        :item="item"
        @d-end="onTabDragEnd"
        @d-enter="onTabDragEnter"
        @d-leave="onTabDragLeave"
        @d-start="onTabDragStart"
    />
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, provide, reactive, ref} from 'vue';
import DraggableItem from '@/components/drag/DraggableItem.vue';
import {plainClone} from '@/utils/object';

export default defineComponent({
  name: 'DraggableList',
  components: {DraggableItem},
  props: {
    items: {
      type: Array,
      required: true,
    },
    itemKey: {
      type: String,
      default: 'key',
    },
  },
  emits: [
    'd-end',
  ],
  setup(props, ctx) {
    const {emit} = ctx;
    const internalItems = reactive<DraggableListInternalItems>({});
    const isDragging = ref(false);

    const orderedItems = computed(() => {
      const {items, itemKey} = props as DraggableListProps;
      const {draggingItem, targetItem} = internalItems;
      if (!draggingItem || !targetItem) return items;
      const orderedItems = plainClone(items) as DraggableItemData[];
      const draggingIdx = orderedItems.map(t => t[itemKey]).indexOf(draggingItem[itemKey]);
      const targetIdx = orderedItems.map(t => t[itemKey]).indexOf(targetItem[itemKey]);
      if (draggingIdx === -1 || targetIdx === -1) return items;
      orderedItems.splice(draggingIdx, 1);
      orderedItems.splice(targetIdx, 0, plainClone(draggingItem));
      return orderedItems;
    });

    const onTabDragStart = (item: DraggableItemData) => {
      internalItems.draggingItem = plainClone(item) as DraggableItemData;
      internalItems.draggingItem.dragging = true;
      isDragging.value = true;
    };

    const onTabDragEnd = () => {
      const items = orderedItems.value.map(d => {
        d.dragging = false;
        return d;
      });
      isDragging.value = false;
      internalItems.draggingItem = undefined;
      internalItems.targetItem = undefined;
      emit('d-end', items);
    };

    const onTabDragEnter = (item: DraggableItemData) => {
      const {itemKey} = props as DraggableListProps;
      const {draggingItem} = internalItems;
      if (!draggingItem || (draggingItem ? draggingItem[itemKey] : undefined) === item[itemKey]) return;
      const _item = plainClone(item) as DraggableItemData;
      _item.dragging = true;
      internalItems.targetItem = _item;
    };

    const onTabDragLeave = (item: DraggableItemData) => {
      const {itemKey} = props as DraggableListProps;
      const {draggingItem, targetItem} = internalItems;
      if (!!targetItem || !draggingItem || (draggingItem ? draggingItem[itemKey] : undefined) === item[itemKey]) return;
      internalItems.targetItem = undefined;
    };

    provide('list', {
      ctx,
      props,
    } as DraggableListContext);

    return {
      orderedItems,
      onTabDragStart,
      onTabDragEnd,
      onTabDragEnter,
      onTabDragLeave,
    };
  },
});
</script>

<style lang="scss" scoped>
.draggable-list {
  list-style: none;
  display: flex;
  align-items: center;
  margin: 0;
  padding: 0;
}
</style>
