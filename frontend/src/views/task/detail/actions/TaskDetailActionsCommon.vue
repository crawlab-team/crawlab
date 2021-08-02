<template>
  <NavActionGroup class="task-detail-actions-common">
    <NavActionFaIcon :icon="['fa', 'tools']"/>
    <NavActionItem>
      <TaskStatus :status="form.status" size="normal"/>
    </NavActionItem>
    <NavActionItem>
      <FaIconButton :icon="['fa', 'redo']" tooltip="Restart" type="warning" @click="onRestart"/>
    </NavActionItem>
    <NavActionItem>
      <FaIconButton v-if="cancellable" :icon="['fa', 'pause']" tooltip="Cancel" type="info" @click="onCancel"/>
      <FaIconButton v-else :icon="['fa', 'trash-alt']" tooltip="Delete" type="danger" @click="onDelete"/>
    </NavActionItem>
  </NavActionGroup>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionItem from '@/components/nav/NavActionItem.vue';
import NavActionFaIcon from '@/components/nav/NavActionFaIcon.vue';
import {useStore} from 'vuex';
import {isCancellable} from '@/utils/task';
import useTask from '@/components/task/task';
import {ElMessage, ElMessageBox} from 'element-plus';
import useRequest from '@/services/request';
import useTaskDetail from '@/views/task/detail/taskDetail';
import {useRouter} from 'vue-router';
import TaskStatus from '@/components/task/TaskStatus.vue';

const {
  post,
} = useRequest();

export default defineComponent({
  name: 'TaskDetailActionsCommon',
  components: {
    TaskStatus,
    NavActionFaIcon,
    FaIconButton,
    NavActionGroup,
    NavActionItem,
  },
  setup() {
    // router
    const router = useRouter();

    // store
    const ns = 'task';
    const store = useStore();

    // use task
    const {
      form,
    } = useTask(store);

    // use task detail
    const {
      activeId,
    } = useTaskDetail();

    // restart
    const onRestart = async () => {
      await ElMessageBox.confirm('Are you sure to restart?', 'Restart', {type: 'warning'});
      await post(`/tasks/${activeId.value}/restart`);
      await ElMessage.success('Restarted successfully');
      await store.dispatch(`${ns}/getById`, activeId.value);
    };

    // cancel
    const onCancel = async () => {
      await ElMessageBox.confirm('Are you sure to cancel?', 'Cancel', {type: 'warning'});
      await ElMessage.info('Attempt to cancel');
      await post(`/tasks/${activeId.value}/cancel`);
      await store.dispatch(`${ns}/getById`, activeId.value);
    };

    // delete
    const onDelete = async () => {
      await ElMessageBox.confirm('Are you sure to delete?', 'Delete', {
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      });
      await store.dispatch(`${ns}/deleteById`, activeId.value);
      await router.push('/tasks');
    };

    // cancellable
    const cancellable = computed<boolean>(() => isCancellable(form.value?.status));

    return {
      ...useTask(store),
      onRestart,
      onCancel,
      onDelete,
      cancellable,
    };
  },
});
</script>

<style scoped>
.task-detail-actions-common >>> .task-status {
  margin-right: 10px;
}
</style>
