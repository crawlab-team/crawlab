<template>
  <Form v-if="form" ref="formRef" :model="form" class="task-form">
    <!-- Row -->
    <FormItem :offset="2" :span="2" label="Spider" prop="spider_id">
      <el-select
          v-model="form.spider_id"
          :disabled="isFormItemDisabled('spider_id') || readonly"
      >
        <el-option
            v-for="op in allSpiderSelectOptions"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
      <FaIconButton
          v-if="readonly"
          :icon="['fa', 'external-link-alt']"
          class="nav-btn"
          tooltip="Go to Spider"
          @click="onGoToSpider"
      />
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem v-if="readonly" :offset="2" :span="2" label="Node" prop="node_id">
      <el-input v-if="noNodeId" disabled placeholder="Unassigned"/>
      <el-select
          v-else
          v-model="form.node_id"
          disabled
      >
        <el-option
            v-for="op in allNodeSelectOptions"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
      <FaIconButton
          v-if="readonly"
          :icon="['fa', 'external-link-alt']"
          class="nav-btn"
          tooltip="Go to Spider"
          :disabled="noNodeId"
          @click="onGoToNode"
      />
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem v-if="readonly" :span="4" label="Status" prop="status">
      <TaskStatus :status="form.status" size="small"/>
      <Tag
          v-if="form.status === 'error'"
          :icon="['fa', 'exclamation']"
          :label="form.error"
          class="error-message"
          size="small"
          tooltip="Task error message"
          type="danger"
      />
      <Tag
          v-else-if="cancellable"
          :icon="['fa', 'pause']"
          class="cancel-btn"
          clickable
          label="Cancel"
          size="small"
          tooltip="Cancel task"
          type="info"
          @click="onCancel"
      />
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem :span="2" label="Command" prop="cmd" required>
      <InputWithButton
          v-model="form.cmd"
          :button-icon="['fa', 'edit']"
          :disabled="isFormItemDisabled('cmd') || readonly"
          button-label="Edit"
          placeholder="Command"
      />
    </FormItem>
    <FormItem :span="2" label="Param" prop="param">
      <InputWithButton
          v-model="form.param"
          :button-icon="['fa', 'edit']"
          :disabled="isFormItemDisabled('param') || readonly"
          button-label="Edit"
          placeholder="Params"
      />
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem :span="2" label="Mode" prop="mode" required>
      <el-select
          v-model="form.mode"
          :disabled="isFormItemDisabled('mode') || readonly"
      >
        <el-option
            v-for="op in modeOptions"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
    </FormItem>
    <FormItem :span="2" label="Priority" prop="priority" required>
      <el-select
          v-model="form.priority"
          :disabled="isFormItemDisabled('priority') || readonly"
      >
        <el-option
            v-for="op in priorityOptions"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
    </FormItem>
    <!-- ./Row -->

    <FormItem
        v-if="form.mode === TASK_MODE_SELECTED_NODE_TAGS"
        :span="4"
        label="Selected Tags"
        prop="node_tags"
        required
    >
      <CheckTagGroup
          v-model="form.node_tags"
          :disabled="isFormItemDisabled('node_tags') || readonly"
          :options="allNodeTags"
      />
    </FormItem>

    <FormItem
        v-if="[TASK_MODE_SELECTED_NODES, TASK_MODE_SELECTED_NODE_TAGS].includes(form.mode)"
        :span="4"
        label="Selected Nodes"
        required
    >
      <CheckTagGroup
          v-model="form.node_ids"
          :disabled="(form.mode === TASK_MODE_SELECTED_NODE_TAGS && isFormItemDisabled('node_ids')) || readonly"
          :options="allNodeSelectOptions"
      />
    </FormItem>
  </Form>
</template>

<script lang="ts">
import {computed, defineComponent, watch} from 'vue';
import {useStore} from 'vuex';
import useSpider from '@/components/spider/spider';
import useNode from '@/components/node/node';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';
import InputWithButton from '@/components/input/InputWithButton.vue';
import CheckTagGroup from '@/components/tag/CheckTagGroup.vue';
import {TASK_MODE_SELECTED_NODE_TAGS, TASK_MODE_SELECTED_NODES} from '@/constants/task';
import useRequest from '@/services/request';
import useTask from '@/components/task/task';
import TaskStatus from '@/components/task/TaskStatus.vue';
import Tag from '@/components/tag/Tag.vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import {useRouter} from 'vue-router';
import {isCancellable} from '@/utils/task';
import {ElMessage, ElMessageBox} from 'element-plus';
import {isZeroObjectId} from '@/utils/mongo';
import useTaskDetail from '@/views/task/detail/taskDetail';

const {
  post,
} = useRequest();

export default defineComponent({
  name: 'TaskForm',
  components: {
    FaIconButton,
    Tag,
    TaskStatus,
    Form,
    FormItem,
    InputWithButton,
    CheckTagGroup,
  },
  props: {
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    // router
    const router = useRouter();

    // store
    const ns = 'task';
    const store = useStore();

    // use node
    const {
      allListSelectOptionsWithEmpty: allNodeSelectOptions,
      allTags: allNodeTags,
    } = useNode(store);

    // use spider
    const {
      allListSelectOptions: allSpiderSelectOptions,
    } = useSpider(store);

    // use task
    const {
      form,
      allSpiderDict,
      modeOptionsDict,
    } = useTask(store);

    // use task detail
    const {
      activeId,
    } = useTaskDetail();

    // use request
    const {
      get,
    } = useRequest();

    // watch spider id
    watch(() => {
      const task = form.value as Task;
      return task.spider_id;
    }, async () => {
      const task = form.value as Task;
      if (!task.spider_id) return;
      const res = await get<any, Spider>(`/spiders/${task.spider_id}`);
      task.cmd = res.data.cmd;
      task.param = res.data.param;
    });

    const getSpiderName = (id: string) => {
      const spider = allSpiderDict.value.get(id) as Spider;
      return spider?.name;
    };

    const getModeName = (id: string) => {
      const op = modeOptionsDict.value.get(id) as SelectOption;
      return op?.label;
    };

    const onGoToSpider = () => {
      router.push(`/spiders/${form.value.spider_id}`);
    };

    const onGoToNode = () => {
      router.push(`/nodes/${form.value.node_id}`);
    };

    const cancellable = computed<boolean>(() => isCancellable(form.value.status));

    const onCancel = async () => {
      await ElMessageBox.confirm('Are you sure to cancel?', 'Cancel', {type: 'warning'});
      await ElMessage.info('Attempt to cancel');
      try {
        await post(`/tasks/${activeId.value}/cancel`);
      } finally {
        await store.dispatch(`${ns}/getById`, activeId.value);
      }
    };

    const noNodeId = computed<boolean>(() => isZeroObjectId(form.value.node_id));

    return {
      ...useTask(store),

      // custom
      TASK_MODE_SELECTED_NODES,
      TASK_MODE_SELECTED_NODE_TAGS,
      allNodeSelectOptions,
      allNodeTags,
      allSpiderSelectOptions,
      getSpiderName,
      getModeName,
      onGoToSpider,
      onGoToNode,
      cancellable,
      onCancel,
      noNodeId,
    };
  },
});
</script>

<style scoped>
.task-form >>> .nav-btn {
  position: absolute;
  padding-left: 10px;
}

.task-form >>> .error-message,
.task-form >>> .cancel-btn {
  margin-left: 10px;
}

.task-form >>> .cancel-btn:hover {
  opacity: 0.8;
}
</style>
