<template>
  <NavActionGroup>
    <NavActionFaIcon :icon="['fa', 'tools']"/>
    <NavActionItem>
      <FaIconButton :icon="['fa', 'play']" tooltip="Run" type="success" @click="onRun"/>
    </NavActionItem>
    <NavActionItem>
      <FaIconButton :icon="['fa', 'clone']" tooltip="Clone" type="info"/>
    </NavActionItem>
    <NavActionItem>
      <FaIconButton :icon="['far', 'star']" plain tooltip="Favorite" type="warning"/>
    </NavActionItem>
  </NavActionGroup>
  <!--TODO: implement-->
  <NavActionGroup v-if="false">
    <NavActionFaIcon :icon="['fab', 'git-alt']"/>
    <NavActionItem>
      <FaIconButton :icon="['fa', 'upload']" tooltip="Upload File" type="primary"/>
    </NavActionItem>
    <NavActionItem>
      <FaIconButton :icon="['fa', 'paper-plane']" tooltip="Commit" type="success"/>
    </NavActionItem>
  </NavActionGroup>

  <!-- Dialogs (handled by store) -->
  <RunSpiderDialog v-if="activeDialogKey === 'run'"/>
  <!-- ./Dialogs -->
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionItem from '@/components/nav/NavActionItem.vue';
import NavActionFaIcon from '@/components/nav/NavActionFaIcon.vue';
import {useStore} from 'vuex';
import useSpider from '@/components/spider/spider';
import RunSpiderDialog from '@/components/spider/RunSpiderDialog.vue';

export default defineComponent({
  name: 'SpiderDetailActionsCommon',
  components: {
    NavActionFaIcon,
    FaIconButton,
    NavActionGroup,
    NavActionItem,
    RunSpiderDialog,
  },
  setup() {
    // store
    const ns = 'spider';
    const store = useStore();

    const onRun = () => {
      store.commit(`${ns}/showDialog`, 'run');
    };

    return {
      ...useSpider(store),
      onRun,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
