import 'crawlab-ui/dist/crawlab-ui.css';
import {createApp} from 'crawlab-ui';

(async function () {
  // @ts-ignore
  window.VUE_APP_API_BASE_URL = process.env.VUE_APP_API_BASE_URL;
  await createApp();
})();
