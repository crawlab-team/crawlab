<template>
  <Form
      v-if="form"
      ref="formRef"
      :model="form"
      :selective="isSelectiveForm"
  >
    <!--Row-->
    <FormItem :span="2" :offset="2" label="Name" not-editable prop="name" required>
      <el-input v-model="form.name" disabled placeholder="Name" />
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="2" label="Execute Command" prop="cmd">
      <el-input
          v-model="form.cmd"
          disabled
          placeholder="cmd"
      />
    </FormItem>
    <!--./Row-->

    <!--Row-->
    <FormItem :span="4" label="Description" prop="description">
      <el-input
          v-model="form.description"
          disabled
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
import usePlugin from '@/components/plugin/plugin';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';

export default defineComponent({
  name: 'PluginForm',
  props: {
    readonly: {
      type: Boolean,
    }
  },
  components: {
    Form,
    FormItem,
  },
  setup(props, {emit}) {
    // store
    const store = useStore();

    return {
      ...usePlugin(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
