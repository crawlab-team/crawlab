import {SPIDER_TYPE_CONFIGURABLE, SPIDER_TYPE_CUSTOMIZED,} from '@/constants/spider';

declare global {
  interface SpiderTypeProps {
    type: SpiderType;
  }

  type SpiderType = SPIDER_TYPE_CUSTOMIZED | SPIDER_TYPE_CONFIGURABLE;
}
