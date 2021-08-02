<script lang="ts">
import {defineComponent, h} from 'vue';
import FaIconButton from '@/components/button/FaIconButton.vue';

export default defineComponent({
  name: 'TableCell',
  props: {
    column: {
      type: Object,
      required: true,
    },
    row: {
      type: Object,
      required: true,
    },
    rowIndex: {
      type: Number,
      required: true,
    }
  },
  emits: [
    'click',
  ],
  setup: function (props) {
    const getChildren = () => {
      const {row, column, rowIndex} = props as TableCellProps;
      const {value, buttons} = column;

      // value
      if (value !== undefined) {
        if (typeof value === 'function') {
          return [value(row, rowIndex, column)];
        } else {
          return value;
        }
      }

      // buttons
      if (buttons) {
        let _buttons: TableColumnButton[] = [];
        if (typeof buttons === 'function') {
          _buttons = buttons(row);
        } else if (Array.isArray(buttons) && buttons.length > 0) {
          _buttons = buttons;
        }

        return _buttons.map(btn => {
          const {tooltip, type, size, icon, disabled, onClick} = btn;
          const props = {
            key: JSON.stringify({tooltip, type, size, icon}),
            tooltip: typeof tooltip === 'function' ? tooltip(row) : tooltip,
            type,
            size,
            icon,
            disabled: disabled?.(row),
            onClick: () => {
              onClick?.(row, rowIndex, column);
            },
          };
          // FIXME: use "as any" to fix type errors temporarily
          return h(FaIconButton, props as any);
        });
      }

      // plain text
      return [row[column.key]];
    };

    return () => h('div', getChildren());
  },
});
</script>

<style lang="scss" scoped>

</style>
