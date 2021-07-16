<template>
  <Form
      v-if="form"
      ref="formRef"
      :model="form"
      :rules="formRules"
      :selective="isSelectiveForm"
      class="schedule-form"
  >
    <!-- Row -->
    <FormItem :span="2" label="Name" prop="name" required>
      <el-input v-model="form.name" :disabled="isFormItemDisabled('name')" placeholder="Name"/>
    </FormItem>
    <FormItem :span="2" label="Spider" prop="spider_id" required>
      <el-select
          v-model="form.spider_id"
          :disabled="isFormItemDisabled('spider_id')"
      >
        <el-option
            v-for="op in allSpiderSelectOptions"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem :span="2" label="Cron Expression" prop="cron" required>
      <el-input v-model="form.cron" :disabled="isFormItemDisabled('cron')" placeholder="Cron Expression"/>
    </FormItem>
    <FormItem :not-editable="isSelectiveForm" :span="2" label="Cron Info">
      <div class="nav-btn">
        <ScheduleCron :cron="form.cron" icon-only size="small"/>
      </div>
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem :span="2" label="Command" prop="cmd">
      <InputWithButton
          v-model="form.cmd"
          :button-icon="['fa', 'edit']"
          :disabled="isFormItemDisabled('cmd')"
          button-label="Edit"
          placeholder="Command"
      />
    </FormItem>
    <FormItem :span="2" label="Param" prop="param">
      <InputWithButton
          v-model="form.param"
          :button-icon="['fa', 'edit']"
          :disabled="isFormItemDisabled('param')"
          button-label="Edit"
          placeholder="Params"
      />
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem :span="2" label="Default Mode" prop="mode">
      <el-select
          v-model="form.mode"
          :disabled="isFormItemDisabled('mode')"
      >
        <el-option
            v-for="op in modeOptions"
            :key="op.value"
            :label="op.label"
            :value="op.value"
        />
      </el-select>
    </FormItem>
    <FormItem :span="2" label="Enabled" prop="enabled" required>
      <Switch v-model="form.enabled" @change="onEnabledChange"/>
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
          :disabled="isFormItemDisabled('node_tags')"
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
          :disabled="form.mode === TASK_MODE_SELECTED_NODE_TAGS && isFormItemDisabled('node_ids')"
          :options="allNodeSelectOptions"
      />
    </FormItem>

    <!-- Row -->
    <FormItem :span="4" label="Description" prop="description">
      <el-input
          v-model="form.description"
          :disabled="isFormItemDisabled('description')"
          placeholder="Description"
          type="textarea"
      />
    </FormItem>
    <!-- ./Row -->
  </Form>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import useSchedule from '@/components/schedule/schedule';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';
import useSpider from '@/components/spider/spider';
import {TASK_MODE_SELECTED_NODE_TAGS, TASK_MODE_SELECTED_NODES} from '@/constants/task';
import useNode from '@/components/node/node';
import CheckTagGroup from '@/components/tag/CheckTagGroup.vue';
import InputWithButton from '@/components/input/InputWithButton.vue';
import Switch from '@/components/switch/Switch.vue';
import {ElMessage} from 'element-plus';
import ScheduleCron from '@/components/schedule/ScheduleCron.vue';

export default defineComponent({
  name: 'ScheduleForm',
  components: {
    ScheduleCron,
    Switch,
    FormItem,
    Form,
    CheckTagGroup,
    InputWithButton,
  },
  setup() {
    // store
    const ns = 'schedule';
    const store = useStore();

    // use node
    const {
      allListSelectOptions: allNodeSelectOptions,
      allTags: allNodeTags,
    } = useNode(store);

    // use spider
    const {
      allListSelectOptions: allSpiderSelectOptions,
    } = useSpider(store);

    // use schedule
    const {
      form,
    } = useSchedule(store);

    // on enabled change
    const onEnabledChange = async (value: boolean) => {
      if (value) {
        await store.dispatch(`${ns}/enable`, form.value._id);
        ElMessage.success('Enabled successfully');
      } else {
        await store.dispatch(`${ns}/disable`, form.value._id);
        ElMessage.success('Disabled successfully');
      }
      await store.dispatch(`${ns}/getList`);
    };

    return {
      ...useSchedule(store),

      allSpiderSelectOptions,
      allNodeSelectOptions,
      allNodeTags,
      TASK_MODE_SELECTED_NODES,
      TASK_MODE_SELECTED_NODE_TAGS,
      onEnabledChange,
    };
  },
});
</script>

<style scoped>
</style>
