<template>
  <el-select v-if="type === 'select'" :value="value" @change="onChange">
    <el-option v-for="op in data.options" :key="op" :label="op" :value="op"/>
  </el-select>
  <el-input-number
      v-else-if="type === 'input-number'"
      :min="data.min !== undefined ? data.min : 0"
      :step="data.step !== undefined ? data.step : 1"
      :value="value"
      @change="onChange"
  />
  <el-switch v-else-if="type === 'switch'" :value="value" @change="onChange"/>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {getOptionDefinition} from '@/utils/codemirror';

export default defineComponent({
  name: 'FileEditorSettingsFormItem',
  props: {
    value: {
      type: Object,
      required: false,
    },
    name: {
      type: String,
      required: true,
    },
  },
  setup(props, {emit}) {
    const def = computed<FileEditorOptionDefinition | undefined>(() => {
      const {name} = props;
      return getOptionDefinition(name);
    });
    const type = computed<FileEditorOptionDefinitionType | undefined>(() => def.value?.type);
    const data = computed<any>(() => def.value?.data);

    const onChange = (value: any) => {
      emit('input', value);
    };

    return {
      type,
      data,
      onChange,
    };
  },
});
</script>
