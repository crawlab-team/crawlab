<template>
  <Tag
      :clickable="clickable"
      :color="color"
      :effect="effect"
      :icon="icon"
      :label="label"
      :spinning="spinning"
      :tooltip="tooltip"
      :type="type"
      :width="width"
      class="link-tag"
      @click="onClick"
  />
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import Tag, {tagProps} from '@/components/tag/Tag.vue';
import {useRouter} from 'vue-router';

const linkTagProps = {
  ...tagProps,
  path: {
    type: String,
    default: '',
  },
};

export default defineComponent({
  name: 'LinkTag',
  components: {Tag},
  props: linkTagProps,
  setup(props: LinkTagProps, {emit}) {
    const router = useRouter();

    const onClick = () => {
      const {path} = props;
      if (path) {
        router.push(path);
      }
    };

    return {
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
