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
import {
  NODE_STATUS_OFFLINE,
  NODE_STATUS_ONLINE,
  NODE_STATUS_REGISTERED,
  NODE_STATUS_UNREGISTERED
} from '@/constants/node';

export default defineComponent({
  name: 'NodeStatus',
  components: {
    Tag,
  },
  props: {
    status: {
      type: String as PropType<TaskStatus>,
      required: false,
    },
    size: {
      type: String as PropType<BasicSize>,
      required: false,
      default: 'mini',
    },
  },
  emits: ['click'],
  setup(props: NodeStatusProps, {emit}) {
    const data = computed<TagData>(() => {
      const {status} = props;
      switch (status) {
        case NODE_STATUS_UNREGISTERED:
          return {
            label: 'Unregistered',
            tooltip: 'Node is waiting to be registered',
            type: 'danger',
            icon: ['fa', 'exclamation'],
          };
        case NODE_STATUS_REGISTERED:
          return {
            label: 'Registered',
            tooltip: 'Node is registered and wait to be online',
            type: 'warning',
            icon: ['far', 'check-square'],
          };
        case NODE_STATUS_ONLINE:
          return {
            label: 'Online',
            tooltip: 'Node is currently online',
            type: 'success',
            icon: ['fa', 'check'],
          };
        case NODE_STATUS_OFFLINE:
          return {
            label: 'Offline',
            tooltip: 'Node is currently offline',
            type: 'info',
            icon: ['fa', 'times'],
          };
        default:
          return {
            label: 'Unknown',
            tooltip: 'Unknown node status',
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
.task-status {
  width: 80px;
  cursor: default;

  .icon {
    width: 10px;
    margin-right: 5px;
  }
}
</style>
