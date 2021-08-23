import Vue from 'vue';

declare global {
  interface Window {
    initCanvas?: Function;
    resetCanvas?: Function;
    _hmt?: Array;
    'vue3-sfc-loader': { loadModule };
    Vue: Vue;
  }
}
