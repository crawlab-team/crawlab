<script lang="ts">
import {defineComponent, inject} from 'vue';

export default defineComponent({
  name: 'DraggableItemContent',
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    return () => {
      const {item} = props;
      const content = inject<DraggableListContext>('list');
      if (!content || !content.ctx || !content.ctx.slots || !content.ctx.slots.default) return '';
      return content.ctx.slots.default({item});
    };
  },
});
</script>
