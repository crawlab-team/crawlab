import {DataItem, Key} from 'element-plus/lib/el-transfer/src/transfer';

declare global {
  interface TransferPanelProps {
    checked: Key[];
    data: DataItem[];
    title: string;
  }
}
