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
    <FormItem :span="2" label="Color" prop="color" required>
      <ColorPicker
          v-model="form.color"
          :predefine="predefinedColors"
          class="color-picker"
          show-alpha
      />
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
import {defineComponent, readonly} from 'vue';
import {useStore} from 'vuex';
import useTag from '@/components/tag/tag';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';
import {getPredefinedColors} from '@/utils/color';
import ColorPicker from '@/components/color/ColorPicker.vue';

export default defineComponent({
  name: 'TagForm',
  components: {ColorPicker, FormItem, Form},
  setup() {
    // store
    const store = useStore();

    // predefined colors
    const predefinedColors = readonly<string[]>(getPredefinedColors());

    return {
      ...useTag(store),
      predefinedColors,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
