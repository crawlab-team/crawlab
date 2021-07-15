<template>
  <template v-if="icon">
    <font-awesome-icon
        v-if="isFaIcon"
        :class="spinning ? 'fa-spin' : ''"
        :icon="icon"
        :style="{fontSize}"
        class="icon"
    />
    <i
        v-else
        :class="[spinning ? 'fa-spin' : '', icon, 'icon']"
        :style="{fontSize}"
    />
  </template>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import useIcon from '@/components/icon/icon';

export default defineComponent({
  name: 'Icon',
  props: {
    icon: {
      type: [String, Array] as PropType<Icon>,
    },
    spinning: {
      type: Boolean,
      default: false,
    },
    size: {
      type: String as PropType<IconSize>,
      default: 'mini',
    }
  },
  setup(props: IconProps, {emit}) {
    const {
      isFaIcon: _isFaIcon,
      getFontSize,
    } = useIcon();

    const fontSize = computed(() => {
      const {size} = props;
      return getFontSize(size);
    });

    const isFaIcon = computed<boolean>(() => {
      const {icon} = props;
      if (!icon) return false;
      return _isFaIcon(icon);
    });

    return {
      isFaIcon,
      fontSize,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
