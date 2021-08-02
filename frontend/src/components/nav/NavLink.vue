<template>
  <div class="nav-link" @click="onClick">
    <Icon :icon="icon" class="icon"/>
    <span class="title">{{ label }}</span>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import Icon from '@/components/icon/Icon.vue';
import {useRouter} from 'vue-router';

export default defineComponent({
  name: 'NavLink',
  components: {Icon},
  props: {
    path: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    icon: {
      type: [String, Array] as PropType<Icon>,
      default: '',
    },
  },
  emits: [
    'click',
  ],
  setup(props: NavLinkProps, {emit}) {
    const router = useRouter();

    const onClick = () => {
      const {path} = props;
      if (path) {
        router.push(path);
      }
      emit('click');
    };

    return {
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/color";

.nav-link {
  cursor: pointer;
  color: $blue;

  &:hover {
    text-decoration: underline;
  }

  .icon {
    margin-right: 3px;
  }
}
</style>
