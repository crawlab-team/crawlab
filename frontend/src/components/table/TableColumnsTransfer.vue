<template>
  <el-dialog
      :before-close="onClose"
      :model-value="visible"
      title="Table Columns Customization">
    <div class="table-columns-transfer-content">
      <Transfer
          :data="computedData"
          :titles="['Available', 'Selected']"
          :value="internalSelectedColumnKeys"
          @change="onChange"
      />
    </div>
    <template #footer>
      <Button plain size="mini" type="info" @click="onClose">Cancel</Button>
      <Button size="mini" @click="onApply">Apply</Button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import {computed, defineComponent, onBeforeMount, ref} from 'vue';
import Transfer from '@/components/transfer/Transfer.vue';
import Button from '@/components/button/Button.vue';
import {DataItem} from 'element-plus/lib/el-transfer/src/transfer';

export default defineComponent({
  name: 'TableColumnsTransfer',
  components: {Button, Transfer},
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    columns: {
      type: Array,
      default: () => {
        return [];
      },
    },
    selectedColumnKeys: {
      type: Array,
      default: () => {
        return [];
      }
    },
  },
  emits: [
    'close',
    'change',
    'sort',
    'apply',
  ],
  setup(props, {emit}) {
    const internalSelectedColumnKeys = ref<string[]>([]);

    const computedData = computed<DataItem[]>(() => {
      const {columns} = props as TableColumnsTransferProps;
      return columns.map(d => {
        const {key, label, disableTransfer} = d;
        return {
          key,
          label,
          disabled: disableTransfer || false,
        };
      });
    });

    const onClose = () => {
      emit('close');
    };

    const onChange = (value: string[]) => {
      internalSelectedColumnKeys.value = value;
    };

    const onApply = () => {
      emit('apply', internalSelectedColumnKeys.value);
      emit('close');
    };

    onBeforeMount(() => {
      const {selectedColumnKeys} = props as TableColumnsTransferProps;
      internalSelectedColumnKeys.value = selectedColumnKeys || [];
    });

    return {
      internalSelectedColumnKeys,
      computedData,
      onClose,
      onChange,
      onApply,
    };
  },
});
</script>

<style lang="scss" scoped>
.table-columns-transfer-content {
  display: flex;
  justify-content: center;
}
</style>
