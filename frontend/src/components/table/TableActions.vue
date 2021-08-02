<template>
  <div class="table-actions">
    <slot name="prefix"></slot>
    <!--    <FaIconButton-->
    <!--        v-if="showButton(ACTION_ADD)"-->
    <!--        :icon="['fa', 'plus']"-->
    <!--        class="action-btn"-->
    <!--        size="mini"-->
    <!--        tooltip="Add"-->
    <!--        type="success"-->
    <!--        :disabled="selection.length === 0"-->
    <!--        @click="onAdd"-->
    <!--    />-->
    <FaIconButton
        v-if="showButton(ACTION_EDIT)"
        :disabled="selection.length === 0"
        :icon="['fa', 'edit']"
        class="action-btn"
        size="mini"
        tooltip="Edit Selected"
        type="warning"
        @click="onEdit"
    />
    <FaIconButton
        v-if="showButton(ACTION_DELETE)"
        :disabled="selection.length === 0"
        :icon="['fa', 'trash-alt']"
        class="action-btn"
        size="mini"
        tooltip="Delete Selected"
        type="danger"
        @click="onDelete"
    />
    <FaIconButton
        v-if="showButton(TABLE_ACTION_EXPORT)"
        :icon="['fa', 'file-export']"
        class="action-btn"
        size="mini"
        tooltip="Export"
        type="primary"
        @click="onExport"
    />
    <FaIconButton
        v-if="showButton(TABLE_ACTION_CUSTOMIZE_COLUMNS)"
        :icon="['fa', 'arrows-alt']"
        class="action-btn"
        size="mini"
        tooltip="Customize Columns"
        type="primary"
        @click="onCustomizeColumns"
    />
    <slot name="suffix"></slot>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import {ACTION_ADD, ACTION_DELETE, ACTION_EDIT,} from '@/constants/action';
import {TABLE_ACTION_CUSTOMIZE_COLUMNS, TABLE_ACTION_EXPORT,} from '@/constants/table';

export default defineComponent({
  name: 'TableActions',
  components: {
    FaIconButton,
  },
  emits: [
    'edit',
    'delete',
    'export',
    'customize-columns'
  ],
  props: {
    selection: {
      type: Array as PropType<TableData>,
      required: false,
      default: () => {
        return [];
      },
    },
    visibleButtons: {
      type: Array as PropType<BuiltInTableActionButtonName[]>,
      required: false,
      default: () => {
        return [];
      }
    }
  },
  setup(props: TableActionsProps, {emit}) {
    // const onAdd = () => {
    //   emit('click-add');
    // };

    const onEdit = () => {
      emit('edit');
    };

    const onDelete = async () => {
      emit('delete');
    };

    const onExport = () => {
      emit('export');
    };

    const onCustomizeColumns = () => {
      emit('customize-columns');
    };

    const showButton = (name: string): boolean => {
      const {visibleButtons} = props;
      if (visibleButtons.length === 0) return true;
      return visibleButtons.includes(name);
    };

    return {
      ACTION_ADD,
      ACTION_EDIT,
      ACTION_DELETE,
      TABLE_ACTION_EXPORT,
      TABLE_ACTION_CUSTOMIZE_COLUMNS,
      // onAdd,
      onEdit,
      onDelete,
      onExport,
      onCustomizeColumns,
      showButton,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
