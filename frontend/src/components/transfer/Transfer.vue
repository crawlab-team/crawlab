<template>
  <div class="transfer">
    <TransferPanel
        :checked="leftChecked"
        :data="leftData"
        :title="titles[0]"
        class="transfer-panel-left"
        @check="onLeftCheck"
        @drag="onLeftDrag"
    />
    <div class="actions">
      <Button :disabled="leftDisabled" :tooltip="leftTooltip || 'Move to Left'" size="large" @click="onLeftMove">
        <div class="btn-content">
          <font-awesome-icon :icon="['fa', 'angle-left']" style="margin-right: 5px"/>
          {{ buttonTexts[0] }}
        </div>
      </Button>
      <Button :disabled="rightDisabled" :tooltip="rightTooltip || 'Move to Right'" size="large" @click="onRightMove">
        <div class="btn-content">
          {{ buttonTexts[1] }}
          <font-awesome-icon :icon="['fa', 'angle-right']" style="margin-left: 5px"/>
        </div>
      </Button>
    </div>
    <TransferPanel
        :checked="rightChecked"
        :data="rightData"
        :title="titles[1]"
        class="transfer-panel-right"
        @check="onRightCheck"
        @drag="onRightDrag"
    />
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, ref} from 'vue';
import {DataItem, Key} from 'element-plus/lib/el-transfer/src/transfer';
import TransferPanel from '@/components/transfer/TransferPanel.vue';
import Button from '@/components/button/Button.vue';

export default defineComponent({
  name: 'Transfer',
  components: {Button, TransferPanel},
  props: {
    value: {
      type: Array,
      required: false,
      default: () => {
        return [];
      }
    },
    data: {
      type: Array,
      required: false,
      default: () => {
        return [];
      }
    },
    titles: {
      type: Array,
      required: false,
      default: () => {
        return [];
      }
    },
    buttonTexts: {
      type: Array,
      required: false,
      default: () => {
        return [];
      }
    },
    buttonTooltips: {
      type: Array,
      required: false,
      default: () => {
        return [];
      }
    },
  },
  emits: [
    'change',
  ],
  setup(props, {emit}) {
    const dataMap = computed<DataMap>(() => {
      const {data} = props as TransferProps;
      const map = {} as DataMap;
      data.forEach(d => {
        map[d.key] = d;
      });
      return map;
    });

    const leftChecked = ref<Key[]>([]);
    const leftData = computed<DataItem[]>(() => {
      const {value, data} = props as TransferProps;
      return data.filter(d => !value.includes(d.key));
    });
    const leftTooltip = computed<string>(() => {
      const {buttonTooltips} = props as TransferProps;
      return buttonTooltips[0];
    });
    const onLeftCheck = (value: Key[]) => {
      leftChecked.value = value;
    };

    const rightChecked = ref<Key[]>([]);
    const rightData = computed<DataItem[]>(() => {
      const {value} = props as TransferProps;
      return value.map(key => dataMap.value[key]);
    });
    const rightTooltip = computed<string>(() => {
      const {buttonTooltips} = props as TransferProps;
      return buttonTooltips[1];
    });
    const onRightCheck = (value: Key[]) => {
      rightChecked.value = value;
    };

    const leftDisabled = computed<boolean>(() => rightChecked.value.length === 0);
    const rightDisabled = computed<boolean>(() => leftChecked.value.length === 0);

    const change = (value: Key[]) => {
      emit('change', value);
    };

    const onLeftMove = () => {
      const {value} = props as TransferProps;
      const newValue = value.filter(d => !rightChecked.value.includes(d));
      change(newValue);
      rightChecked.value = [];
    };
    const onLeftDrag = (items: DataItem[]) => {
      const {value} = props as TransferProps;
      const itemKeys = items.map(d => d.key);
      const newValue = value.filter(d => !itemKeys.includes(d));
      change(newValue);
    };

    const onRightMove = () => {
      const {value} = props as TransferProps;
      const newValue = value.concat(leftChecked.value);
      change(newValue);
      leftChecked.value = [];
    };
    const onRightDrag = (items: DataItem[]) => {
      const newValue = items.map(d => d.key);
      change(newValue);
    };

    return {
      leftChecked,
      leftData,
      leftDisabled,
      leftTooltip,
      onLeftCheck,
      onLeftMove,
      onLeftDrag,
      rightChecked,
      rightData,
      rightDisabled,
      rightTooltip,
      onRightCheck,
      onRightMove,
      onRightDrag,
    };
  },
});
</script>

<style lang="scss" scoped>
.transfer {
  width: 100%;
  min-height: 480px;
  display: flex;
  align-items: center;

  .actions {
    display: flex;
    align-items: center;

    .button {
      .btn-content {
        display: flex;
        align-items: center;
      }
    }
  }
}
</style>
<style scoped>
</style>
