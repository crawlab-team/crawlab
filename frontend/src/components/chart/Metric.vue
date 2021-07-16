<template>
  <div :class="[clickable ? 'clickable' : '']" :style="style" class="metric" @click="onClick">
    <div class="background"/>
    <div class="icon">
      <font-awesome-icon :icon="icon"/>
    </div>
    <div class="info">
      <div class="title">
        {{ title }}
      </div>
      <div class="value">
        {{ value }}
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';

export default defineComponent({
  name: 'Metric',
  props: {
    title: {
      type: String,
    },
    value: {
      type: [Number, String],
    },
    icon: {
      type: [String, Array] as PropType<Icon>,
    },
    color: {
      type: String,
    },
    clickable: {
      type: Boolean,
    }
  },
  emits: [
    'click',
  ],
  setup(props: MetricProps, {emit}) {
    const style = computed<Partial<CSSStyleDeclaration>>(() => {
      const {color} = props;
      return {
        backgroundColor: color,
      };
    });

    const onClick = () => {
      const {clickable} = props;
      if (!clickable) return;
      emit('click');
    };

    return {
      style,
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables";

.metric {
  padding: 10px;
  margin: 20px;
  display: flex;
  border: 1px solid $infoLightColor;
  border-radius: 5px;
  position: relative;


  &.clickable {
    cursor: pointer;

    &:hover {
      opacity: 0.9;
    }
  }

  .background {
    position: absolute;
    left: calc(64px + 10px);
    top: 0;
    width: calc(100% - 64px - 10px);
    height: 100%;
    background-color: white;
    filter: alpha(0.3);
    z-index: 1;
  }

  .icon {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-basis: 64px;
    font-size: 32px;
    color: white;
    padding-right: 10px;
    z-index: 2;
  }

  .info {
    margin-left: 20px;
    height: 48px;
    color: white;
    z-index: 2;

    .title {
      height: 24px;
      line-height: 24px;
      font-weight: bolder;
      //padding: 5px;
    }

    .value {
      height: 24px;
      line-height: 24px;
      font-weight: bold;
      //padding: 5px;
    }
  }
}
</style>
