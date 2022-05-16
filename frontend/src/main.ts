import 'crawlab-ui/dist/crawlab-ui.css';
import 'vue';
import {createApp} from 'crawlab-ui';

(async function () {
  await createApp({
    // @ts-ignore
    initBaiduTongji: window['VUE_APP_INIT_BAIDU_TONGJI'] !== 'false',
    // @ts-ignore
    initUmeng: window['VUE_APP_INIT_UMENG'] !== 'false',
  });
})();
