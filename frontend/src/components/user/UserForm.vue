<template>
  <Form
      v-if="form"
      ref="formRef"
      :model="form"
      :rules="formRules"
      :selective="isSelectiveForm"
      class="user-form"
  >
    <!-- Row -->
    <FormItem :span="2" label="Username" prop="username" required>
      <el-input v-model="form.username" :disabled="isFormItemDisabled('username')" placeholder="Username"/>
    </FormItem>
    <FormItem :span="2" label="Password" prop="password" required>
      <el-input
          v-if="isSelectiveForm || !isDetail"
          v-model="form.password"
          :disabled="isFormItemDisabled('password')"
          placeholder="Password"
          type="password"
      />
      <LabelButton
          v-else
          :icon="['fa','lock']"
          label="Change Password"
          type="danger"
          @click="onChangePassword"
      />
    </FormItem>
    <!-- ./Row -->

    <!-- Row -->
    <FormItem :span="2" label="Email" prop="email">
      <el-input v-model="form.email" :disabled="isFormItemDisabled('email')" placeholder="Email" type="email"/>
    </FormItem>
    <FormItem :span="2" label="Role" prop="role" required>
      <el-select v-model="form.role" :disabled="isFormItemDisabled('role')">
        <el-option :value="ROLE_ADMIN" label="Admin"/>
        <el-option :value="ROLE_NORMAL" label="Normal"/>
      </el-select>
    </FormItem>
    <!-- ./Row -->
  </Form>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {useStore} from 'vuex';
import useUser from '@/components/user/user';
import Form from '@/components/form/Form.vue';
import FormItem from '@/components/form/FormItem.vue';
import {ROLE_ADMIN, ROLE_NORMAL} from '@/constants/user';
import LabelButton from '@/components/button/LabelButton.vue';
import {ElMessageBox} from 'element-plus';
import useUserDetail from '@/views/user/detail/userDetail';

export default defineComponent({
  name: 'UserForm',
  components: {
    LabelButton,
    FormItem,
    Form,
  },
  setup() {
    // store
    const ns = 'user';
    const store = useStore();

    const {
      activeId,
    } = useUserDetail();

    const onChangePassword = async () => {
      const {value} = await ElMessageBox.prompt('Please enter the new password', 'Change Password', {
        inputType: 'password',
        inputPlaceholder: 'New password',
        inputValidator: (value: string) => {
          return value?.length < 5 ? 'Invalid password. Length must be no less than 5.' : true;
        }
      });
      return await store.dispatch(`${ns}/changePassword`, {id: activeId.value, password: value});
    };

    const isDetail = computed<boolean>(() => !!activeId.value);

    return {
      ...useUser(store),
      ROLE_ADMIN,
      ROLE_NORMAL,
      onChangePassword,
      isDetail,
    };
  },
});
</script>

<style scoped>
</style>
