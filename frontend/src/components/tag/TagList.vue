<template>
  <Tag
      v-for="tag in tags"
      :key="tag._id"
      :clickable="clickable"
      :tag="tag"
      @click="onClick"
  />
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import {emptyArrayFunc} from '@/utils/func';
import Tag from '@/components/tag/Tag.vue';

export default defineComponent({
  name: 'TagList',
  components: {
    Tag,
  },
  props: {
    tags: {
      type: Array as PropType<Tag[]>,
      default: emptyArrayFunc,
    },
    clickable: {
      type: Boolean,
      default: true,
    },
  },
  emits: [
    'click'
  ],
  setup(props: TagListProps, {emit}) {
    const onClick = (tag: Tag) => {
      emit('click', tag);
    };

    return {
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
