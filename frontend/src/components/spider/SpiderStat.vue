<template>
  <div class="spider-stat">
    <Tag
        :icon="['fa', 'tasks']"
        :label="labels.tasks"
        :tooltip="tooltips.tasks"
        type="primary"
    />
    <Tag
        :icon="['fa', 'database']"
        :label="labels.results"
        :tooltip="tooltips.results"
        type="success"
    />
    <Tag
        :icon="['fa', 'stopwatch']"
        :label="labels.duration"
        type="warning"
    >
      <template #tooltip>
        <div v-html="tooltips.duration"/>
      </template>
    </Tag>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import Tag from '@/components/tag/Tag.vue';
import colors from '@/styles/color.scss';
import humanizeDuration from 'humanize-duration';

export default defineComponent({
  name: 'SpiderStats',
  components: {
    Tag,
  },
  props: {
    stat: {
      type: Object as PropType<SpiderStat>,
      required: false,
    },
  },
  setup(props: SpiderStatProps, {emit}) {
    const labels = computed<SpiderStatLabels>(() => {
      const {stat} = props;
      const {
        tasks,
        results,
        // average_wait_duration,
        // average_runtime_duration,
        average_total_duration,
      } = stat as SpiderStat;

      return {
        tasks: `${tasks}`,
        results: `${results}`,
        duration: `${average_total_duration}`
      };
    });

    const tooltips = computed<SpiderStatTooltips>(() => {
      const {stat} = props;
      const {
        tasks,
        results,
        average_wait_duration,
        average_runtime_duration,
        average_total_duration,
      } = stat as SpiderStat;

      return {
        tasks: `Total Tasks: ${tasks}`,
        results: `Total Results: ${results}`,
        duration: `
<span class="label">Average Wait Duration:</span>
<span class="value" style="color: ${colors.blue}">${humanizeDuration(average_wait_duration * 1000, {spacer: ' '})}</span><br>
<span class="label">Average Runtime Duration:</span>
<span class="value" style="color: ${colors.orange}">${humanizeDuration(average_runtime_duration * 1000, {spacer: ' '})}</span><br>
<span class="label">Average Total Duration:</span>
<span class="value" style="color: ${colors.white}">${humanizeDuration(average_total_duration * 1000, {spacer: ' '})}</span><br>
`,
      };
    });

    return {
      tooltips,
      labels,
    };
  },
});
</script>

<style lang="scss" scoped>
.spider-stat {
}
</style>

<style scoped>
</style>
