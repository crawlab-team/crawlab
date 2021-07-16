<template>
  <Tag
      :key="data"
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
import {computed, defineComponent, PropType} from 'vue';
import Tag from '@/components/tag/Tag.vue';

export default defineComponent({
  name: 'NodeRunners',
  components: {
    Tag,
  },
  props: {
    available: {
      type: Number as PropType<number | undefined>,
      required: false,
    },
    max: {
      type: Number as PropType<number | undefined>,
      required: false,
    },
    size: {
      type: String as PropType<BasicSize>,
      required: false,
      default: 'mini',
    }
  },
  emits: ['click'],
  setup(props: NodeRunnersProps, {emit}) {
    const running = computed<number>(() => {
      const {available, max} = props;
      if (available === undefined ||
          max === undefined ||
          isNaN(available) ||
          isNaN(max)
      ) {
        return 0;
      }
      return (max - available) as number;
    });

    const label = computed<string>(() => {
      const {max} = props;
      return `${running.value} / ${max}`;
    });

    const data = computed<TagData>(() => {
      const max = props.max as number;
      if (running.value === max) {
        return {
          label: label.value,
          tooltip: 'No runners available at this moment',
          type: 'danger',
          icon: ['fa', 'ban'],
        };
      } else if (running.value > 0) {
        return {
          label: label.value,
          tooltip: `${running.value} out of ${max} runners are running`,
          type: 'warning',
          icon: ['far', 'check-square'],
        };
      } else {
        return {
          label: label.value,
          tooltip: `All runners available`,
          type: 'success',
          icon: ['far', 'check-square'],
        };
      }
    });

    return {
      label,
      data,
    };
  },
});
</script>

<style lang="scss" scoped>
.node-runners {
  cursor: pointer;

  .icon {
    margin-right: 5px;
  }
}
</style>
