<template>
  <Tag
      v-if="!iconOnly"
      :key="data"
      :icon="data.icon"
      :label="data.label"
      :size="size"
      :spinning="data.spinning"
      :type="data.type"
      class="schedule-cron"
      @click="$emit('click')"
  >
    <template #tooltip>
      <div v-html="data.tooltip"/>
    </template>
  </Tag>
  <div v-else :class="[isValid ? 'valid' : 'invalid']" class="schedule-cron">
    <div class="row">
      <span class="title">
        <el-tooltip content="Cron Description">
          <font-awesome-icon :icon="['fa', 'info-circle']" class="description"/>
        </el-tooltip>
      </span>
      <span class="value description">
        {{ isValid ? description : 'Invalid' }}
      </span>
    </div>
    <div class="row">
      <span class="title">
        <el-tooltip content="Next Run">
          <font-awesome-icon :icon="['fa', 'arrow-right']" class="next"/>
        </el-tooltip>
      </span>
      <span class="value next">
        {{ isValid ? next : 'Invalid' }}
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import Tag from '@/components/tag/Tag.vue';
import {CronExpression, parseExpression} from 'cron-parser';
// import cronstrue from 'cronstrue/i18n';
import cronstrue from 'cronstrue';
import dayjs from 'dayjs';
import localizedFormat from 'dayjs/plugin/localizedFormat';
// import 'dayjs/locale/zh-cn';
import colors from '@/styles/color.scss';

// TODO: internalization
dayjs.extend(localizedFormat);
// dayjs.locale('zh-cn');

export default defineComponent({
  name: 'ScheduleCron',
  components: {
    Tag,
  },
  props: {
    cron: {
      type: String,
      required: false,
    },
    size: {
      type: String as PropType<BasicSize>,
      required: false,
      default: 'mini',
    },
    iconOnly: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  setup(props: ScheduleCronProps, {emit}) {
    const interval = computed<CronExpression | undefined>(() => {
      const {cron} = props;
      if (!cron) return;
      try {
        return parseExpression(cron);
      } catch (e) {
        // do nothing
      }
    });

    const next = computed<string | undefined>(() => {
      if (!interval.value) return;
      return dayjs(interval.value.next().toDate()).format('llll');
    });

    const description = computed<string | undefined>(() => {
      const {cron} = props;
      if (!cron) return;
      // TODO: internalization
      return cronstrue.toString(cron);
    });

    const tooltip = computed<string>(() => `<span class="title">Cron Expression: </span><span style="color: ${colors.blue}">${props.cron}</span><br>
<span class="title">Description: </span><span style="color: ${colors.orange}">${description.value}</span><br>
<span class="title">Next: </span><span style="color: ${colors.green}">${next.value}</span>`);

    const isValid = computed<boolean>(() => !!interval.value);

    const data = computed<TagData>(() => {
      const {cron} = props;
      if (!cron) {
        return {
          label: 'Unknown',
          tooltip: 'Unknown',
          type: 'info',
        };
      }

      return {
        label: cron,
        tooltip: tooltip.value,
        type: 'primary',
      };
    });

    return {
      data,
      next,
      description,
      isValid,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.schedule-cron {
  .row {
    min-height: 20px;

    .title {
      display: inline-block;
      width: 18px;
      text-align: right;
      font-size: 14px;
      margin-right: 10px;
    }

    .value {
      font-size: 14px;
    }

    .description {
      color: $warningColor;
    }

    .next {
      color: $successColor;
    }
  }

  &.invalid {
    .description,
    .next {
      color: $infoMediumColor;
    }
  }
}
</style>
