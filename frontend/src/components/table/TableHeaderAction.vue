<template>
  <span :class="classes" class="action" @click="onClick">
    <el-tooltip :content="tooltip">
      <template v-if="isHtml" #content>
        <div v-html="tooltip"/>
      </template>
      <font-awesome-icon :icon="icon"/>
    </el-tooltip>
  </span>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';

export default defineComponent({
  name: 'TableHeaderAction',
  props: {
    tooltip: {
      type: [String, Object],
      required: false,
    },
    isHtml: {
      type: Boolean,
      required: false,
      default: false,
    },
    icon: {
      type: [Array, String],
      required: true,
    },
    status: {
      type: Object,
      required: false,
      default: () => {
        return {active: false, focused: false};
      }
    }
  },
  emits: [
    'click',
  ],
  setup(props, {emit}) {
    const classes = computed<string[]>(() => {
      const {status} = props as TableHeaderActionProps;
      if (!status) return [];
      const {active, focused} = status;
      const cls = [];
      if (active) cls.push('active');
      if (focused) cls.push('focused');
      return cls;
    });

    const onClick = () => {
      emit('click');
    };

    return {
      classes,
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.action {
  margin-left: 3px;
  font-size: 10px;

  &:hover {
    color: $primaryColor;
  }

  &.focused {
    display: inline !important;
    color: $primaryColor;
  }

  &.active {
    display: inline !important;
    color: $warningColor;
  }
}
</style>
