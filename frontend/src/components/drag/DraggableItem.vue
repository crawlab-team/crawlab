<template>
  <div
      :class="classes"
      :draggable="true"
      class="draggable-item"
      @dragstart="$emit('d-start', item)"
      @dragend="$emit('d-end', item)"
      @dragenter="$emit('d-enter', item)"
      @dragleave="$emit('d-leave', item)"
  >
    <DraggableItemContent :item="item"/>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import DraggableItemContent from '@/components/drag/DraggableItemContent.vue';

export default defineComponent({
  name: 'DraggableItem',
  components: {DraggableItemContent},
  props: {
    item: {
      type: Object,
      required: true,
    },
    dragging: {
      type: Boolean,
      default: false,
    }
  },
  emits: [
    'd-start',
    'd-end',
    'd-enter',
    'd-leave',
  ],
  setup(props) {
    const dragging = computed(() => {
      const {item} = props as DraggableItemProps;
      return item.dragging;
    });

    const classes = computed(() => {
      const cls = [];
      if (dragging.value) cls.push('dragging');
      return cls;
    });

    return {
      classes,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.draggable-item {
  position: relative;

  &.dragging {
    visibility: hidden;
  }

  &.dragging * {
    pointer-events: none;
  }
}
</style>
