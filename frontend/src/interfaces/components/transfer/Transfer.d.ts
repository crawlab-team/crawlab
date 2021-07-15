import {DataItem, Key} from 'element-plus/lib/el-transfer/src/transfer';

declare global {
  interface TransferProps {
    value: Key[];
    data: DataItem[];
    titles: string[];
    buttonTexts: string[];
    buttonTooltips: string[];
  }

  interface DataMap {
    [key: string]: DataItem;
  }
}
