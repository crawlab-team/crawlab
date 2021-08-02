<template>
  <Form
      v-if="form"
      ref="formRef"
      :model="form"
      :rules="formRules"
      :selective="isSelectiveForm"
  >
    <FormItem :span="2" label="Name" not-editable prop="name" required>
      <el-input v-model="form.name" :disabled="isFormItemDisabled('name')" placeholder="Name"/>
    </FormItem>
    <FormItem :span="2" label="Tags" prop="tags">
      <TagInput v-model="form.tags" :disabled="isFormItemDisabled('tags')"/>
    </FormItem>
    <FormItem :span="4" label="Description" prop="description">
      <el-input
          v-model="form.description"
          :disabled="isFormItemDisabled('description')"
          placeholder="Description"
          type="textarea"
      />
    </FormItem>
  </Form>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import useProject from '@/components/project/project';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';
import TagInput from '@/components/input/TagInput.vue';

export default defineComponent({
  name: 'ProjectForm',
  components: {TagInput, FormItem, Form},
  setup() {
    // store
    const store = useStore();

    return {
      ...useProject(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
