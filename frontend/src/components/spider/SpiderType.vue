<template>
  <Tag
      :icon="data.icon"
      :label="data.label"
      :tooltip="data.tooltip"
      :type="data.type"
      width="100px"
      @click="$emit('click')"
  />
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import {SPIDER_TYPE_CONFIGURABLE, SPIDER_TYPE_CUSTOMIZED} from '@/constants/spider';
import Tag from '@/components/tag/Tag.vue';

export default defineComponent({
  name: 'SpiderType',
  components: {
    Tag,
  },
  props: {
    type: {
      type: String as PropType<SpiderType>,
      default: SPIDER_TYPE_CUSTOMIZED,
    }
  },
  emits: ['click'],
  setup(props: SpiderTypeProps, {emit}) {
    const data = computed<TagData>(() => {
      const {type} = props;
      switch (type) {
        case SPIDER_TYPE_CUSTOMIZED:
          return {
            tooltip: 'Customized Spider',
            label: 'Customized',
            type: 'success',
            icon: ['fa', 'code'],
          };
        case SPIDER_TYPE_CONFIGURABLE:
          return {
            label: 'Configurable',
            tooltip: 'Configurable Spider',
            type: 'primary',
            icon: ['fa', 'cog'],
          };
        default:
          return {
            label: 'Unknown',
            tooltip: 'Unknown Type',
            type: 'info',
            icon: ['fa', 'question'],
          };
      }
    });

    return {
      data,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
