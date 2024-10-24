import { resolve } from 'path';
import { defineConfig, UserConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import dynamicImport from 'vite-plugin-dynamic-import';
import { visualizer } from 'rollup-plugin-visualizer';

export default defineConfig(({ mode }) => {
  const config: UserConfig = {
    build: {
      rollupOptions: {
        output: {
          manualChunks(id, meta) {
            if (id.includes('node_modules')) {
              const arr = id.toString().split('node_modules/');
              const modulePath = arr[arr.length - 1];
              return modulePath?.split('/')?.[0];
            }
            if (id.includes('three.min.js')) {
              return 'three';
            }
          },
        },
      },
    },
    optimizeDeps: {
      include: ['element-plus', 'axios'],
    },
    resolve: {
      dedupe: ['vue', 'vue-router', 'vuex', 'axios', 'element-plus'],
      alias: {
        '@': resolve(__dirname, 'src'),
      },
      extensions: ['.js', '.ts', '.jsx', '.tsx', '.json', '.vue', '.scss'],
    },
    plugins: [vue(), dynamicImport()],
    server: {
      cors: true,
    },
  };

  if (mode === 'analyze') {
    config.plugins.push(visualizer({ open: true, gzipSize: true }));
  }

  return config;
});
