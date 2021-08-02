<template>
  <div class="create-dialog-content-batch">
    <el-form class="control-panel" inline>
      <el-form-item>
        <Button type="primary" @click="onAdd">
          <font-awesome-icon :icon="['fa', 'plus']"/>
          Add
        </Button>
      </el-form-item>
      <el-form-item label="Edit All">
        <Switch v-model="editAll"/>
      </el-form-item>
    </el-form>
    <FormTable
        :data="data"
        :fields="fields"
        @add="onAdd"
        @clone="onClone"
        @delete="onDelete"
        @field-change="onFieldChange"
        @field-register="onFieldRegister"
    />
  </div>
</template>

<script lang="ts">
import {defineComponent, inject, PropType, Ref, ref} from 'vue';
import FormTable from '@/components/form/FormTable.vue';
import {emptyArrayFunc} from '@/utils/func';
import Switch from '@/components/switch/Switch.vue';
import Button from '@/components/button/Button.vue';

export default defineComponent({
  name: 'CreateDialogContentBatch',
  components: {
    Button,
    Switch,
    FormTable,
  },
  props: {
    data: {
      type: Array as PropType<TableData>,
      required: true,
      default: emptyArrayFunc,
    },
    fields: {
      type: Array as PropType<FormTableField[]>,
      required: true,
      default: emptyArrayFunc,
    }
  },
  setup(props: CreateDialogContentBatchProps) {
    const editAll = ref<boolean>(false);

    const actionFunctions = inject('action-functions') as CreateEditDialogActionFunctions;

    const onAdd = (rowIndex: number) => {
      actionFunctions?.onAdd?.(rowIndex);
    };

    const onClone = (rowIndex: number) => {
      actionFunctions?.onClone?.(rowIndex);
    };

    const onDelete = (rowIndex: number) => {
      actionFunctions?.onDelete?.(rowIndex);
    };

    const onFieldChange = (rowIndex: number, prop: string, value: any) => {
      if (editAll.value) {
        // edit all rows
        rowIndex = -1;
      }

      actionFunctions?.onFieldChange?.(rowIndex, prop, value);
    };

    const onFieldRegister = (rowIndex: number, prop: string, formRef: Ref) => {
      actionFunctions?.onFieldRegister(rowIndex, prop, formRef);
    };

    return {
      editAll,
      onAdd,
      onClone,
      onDelete,
      onFieldChange,
      onFieldRegister,
    };
  },
});
</script>

<style lang="scss" scoped>
.create-dialog-content-batch {
  .control-panel {
    margin-bottom: 10px;

    .el-form-item {
      margin: 0;
    }
  }
}

</style>
