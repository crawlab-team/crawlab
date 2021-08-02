<template>
  <el-tooltip :content="tooltip" :disabled="!tooltip && !$slots.tooltip">
    <el-tag
        ref="tagRef"
        :closable="closable"
        :class="cls"
        :size="size"
        :type="type"
        :color="backgroundColor"
        :effect="effect"
        class="tag"
        @click="onClick($event)"
        @close="onClose($event)"
        @mouseenter="$emit('mouseenter')"
        @mouseleave="$emit('mouseleave')"
    >
      <Icon :icon="icon" :spinning="spinning"/>
      <span>{{ label || tag?.name }}</span>
    </el-tag>
    <template #content>
      <slot name="tooltip"></slot>
    </template>
  </el-tooltip>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, PropType, ref, watch} from 'vue';
import Icon from '@/components/icon/Icon.vue';
import {ElTag} from 'element-plus';

export const tagProps = {
  label: {
    type: String,
  },
  tooltip: {
    type: String,
  },
  type: {
    type: String as PropType<BasicType>,
    default: 'plain',
  },
  color: {
    type: String as PropType<string>,
  },
  backgroundColor: {
    type: String as PropType<string>,
  },
  borderColor: {
    type: String as PropType<string>,
  },
  icon: {
    type: [String, Array] as PropType<string | string[]>,
  },
  size: {
    type: String,
    default: 'mini',
  },
  spinning: {
    type: Boolean,
    default: false,
  },
  width: {
    type: String,
  },
  effect: {
    type: String as PropType<BasicEffect>,
  },
  clickable: {
    type: Boolean,
    default: false,
  },
  closable: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  tag: {
    type: Object as PropType<Tag>,
  },
};

export default defineComponent({
  name: 'Tag',
  components: {Icon},
  props: tagProps,
  emits: [
    'click',
    'close',
    'mouseenter',
    'mouseleave',
  ],
  setup(props: TagProps, {emit}) {
    const tagRef = ref<typeof ElTag>();

    const onClick = (ev?: Event) => {
      ev?.stopPropagation();
      const {clickable} = props;
      if (clickable) {
        emit('click');
      }
    };

    const onClose = (ev?: Event) => {
      ev?.stopPropagation();
      const {closable} = props;
      if (closable) {
        emit('close');
      }
    };

    const cls = computed<string[]>(() => {
      const {clickable, disabled} = props;
      const cls = [] as string[];
      if (clickable) cls.push('clickable');
      if (disabled) cls.push('disabled');
      return cls;
    });

    const setStyle = () => {
      const {color, borderColor, width, tag} = props;

      // normalize colors
      const color_ = color ?? tag?.color;
      const borderColor_ = borderColor ?? color_;

      // set style of tag
      const elTag = tagRef.value?.$el;
      if (!elTag) return;
      const styleTagList = [];
      if (color_) styleTagList.push(`color: ${color_}`);
      if (borderColor_) styleTagList.push(`border-color: ${borderColor_}`);
      if (width) styleTagList.push(`width: ${width}`);
      const styleTag = styleTagList.join(';');
      elTag.setAttribute('style', styleTag);

      // set style of tag close
      const elTagClose = elTag.querySelector('.el-tag__close');
      if (!elTagClose) return;
      const styleTagCloseList = [];
      if (color_) {
        styleTagCloseList.push(`color: ${color_}`);
        styleTagCloseList.push(`background-color: transparent`);
      }
      const styleTagClose = styleTagCloseList.join(';');
      elTagClose.setAttribute('style', styleTagClose);
    };

    watch(() => props.color, setStyle);
    watch(() => props.backgroundColor, setStyle);
    watch(() => props.borderColor, setStyle);

    onMounted(() => {
      setStyle();
    });

    return {
      tagRef,
      onClick,
      onClose,
      cls,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables";

.tag {
  cursor: default;

  &.disabled {
    cursor: not-allowed;
    background-color: $disabledBgColor;
    border-color: $disabledBorderColor;
    color: $disabledColor;
  }

  &.clickable {
    &:not(.disabled) {
      cursor: pointer;
    }
  }
}
</style>
<style scoped>
.tag >>> .el-tag__close:hover {
  font-weight: bolder;
}

.tag >>> .icon {
  margin-right: 5px;
}
</style>
