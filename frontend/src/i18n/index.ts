import {createI18n} from 'vue-i18n';
import en from './lang/en';
import zh from './lang/zh';

const i18n: any = createI18n({
  locale: localStorage.getItem('lang') || 'zh',
  messages: {
    en,
    zh
  },
  silentTranslationWarn: true
});

export default i18n;
