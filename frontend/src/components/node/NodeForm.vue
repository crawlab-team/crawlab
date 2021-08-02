<template>
  <Form
      v-if="form"
      ref="formRef"
      :model="form"
      :selective="isSelectiveForm"
  >
    <!--Row-->
    <FormItem v-if="readonly" :offset="2" :span="2" label="Key" not-editable prop="key">
      <el-input :value="form.key" disabled/>
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="2" label="Name" not-editable prop="name" required>
      <el-input v-model="form.name" :disabled="isFormItemDisabled('name')" placeholder="Name"/>
    </FormItem>
    <FormItem :span="2" label="Tags" prop="tags">
      <TagInput v-model="form.tags" :disabled="isFormItemDisabled('tags')"/>
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="2" label="Type" not-editable prop="type">
      <NodeType :is-master="form.is_master"/>
    </FormItem>
    <FormItem :span="2" label="IP" prop="ip">
      <el-input v-model="form.ip" :disabled="isFormItemDisabled('ip')" placeholder="IP"/>
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="2" label="MAC Address" prop="mac">
      <el-input v-model="form.mac" :disabled="isFormItemDisabled('mac')" placeholder="MAC Address"/>
    </FormItem>
    <FormItem :span="2" label="Hostname" prop="hostname">
      <el-input v-model="form.hostname" :disabled="isFormItemDisabled('hostname')" placeholder="Hostname"/>
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="2" label="Enabled" prop="enabled">
      <Switch v-model="form.enabled" :disabled="isFormItemDisabled('enabled')"/>
    </FormItem>
    <FormItem :span="2" label="Max Runners" prop="max_runners">
      <el-input-number
          v-model="form.max_runners"
          :disabled="isFormItemDisabled('max_runners')"
          :min="0"
          placeholder="Max Runners"
      />
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="4" label="Description" prop="description">
      <el-input
          v-model="form.description"
          :disabled="isFormItemDisabled('description')"
          placeholder="Description"
          type="textarea"
      />
    </FormItem>
  </Form>
  <!--./Row-->
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import useNode from '@/components/node/node';
import TagInput from '@/components/input/TagInput.vue';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';
import NodeType from '@/components/node/NodeType.vue';
import Switch from '@/components/switch/Switch.vue';

export default defineComponent({
  name: 'NodeForm',
  props: {
    readonly: {
      type: Boolean,
    }
  },
  components: {
    Switch,
    NodeType,
    Form,
    FormItem,
    TagInput,
  },
  setup(props, {emit}) {
    // store
    const store = useStore();

    return {
      ...useNode(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
