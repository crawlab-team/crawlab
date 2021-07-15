<template>
  <Tag
      :key="data"
      :color="data.color"
      :icon="data.icon"
      :label="data.label"
      :size="size"
      :spinning="data.spinning"
      :tooltip="data.tooltip"
      :type="data.type"
      @click="$emit('click')"
  />
</template>

<script lang="ts">
import colors from '@/styles/color.scss';
import {computed, defineComponent, PropType} from 'vue';
import Tag from '@/components/tag/Tag.vue';
import {getPriorityLabel} from '@/utils/task';

export default defineComponent({
  name: 'TaskPriority',
  components: {
    Tag,
  },
  props: {
    priority: {
      type: Number,
      required: false,
      default: 5,
    },
    size: {
      type: String as PropType<BasicSize>,
      required: false,
      default: 'mini',
    },
  },
  emits: ['click'],
  setup(props: TaskPriorityProps, {emit}) {
    const data = computed<TagData>(() => {
      const priority = props.priority as number;

      if (priority <= 2) {
        return {
          label: getPriorityLabel(priority),
          color: colors.red,
        };
      } else if (priority <= 4) {
        return {
          label: getPriorityLabel(priority),
          color: colors.orange,
        };
      } else if (priority <= 6) {
        return {
          label: getPriorityLabel(priority),
          color: colors.limeGreen,
        };
      } else if (priority <= 8) {
        return {
          label: getPriorityLabel(priority),
          color: colors.cyan,
        };
      } else {
        return {
          label: getPriorityLabel(priority),
          color: colors.blue,
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
