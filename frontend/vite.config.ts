import {resolve} from 'path';
import {defineConfig, splitVendorChunkPlugin} from 'vite';
import vue from '@vitejs/plugin-vue';
import dynamicImport from 'vite-plugin-dynamic-import';
import {visualizer} from 'rollup-plugin-visualizer';
import externalGlobals from 'rollup-plugin-external-globals';
import {externalizeDeps} from 'vite-plugin-externalize-deps';

export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (id.includes('node_modules')) {
            if (id.includes('@fortawesome')) return '@fortawesome';
            if (id.includes('element-plus')) return 'element-plus';
            if (id.includes('zrender')) return 'zrender';
            if (id.includes('echarts')) return 'echarts';
            if (id.includes('codemirror')) return 'codemirror';
            if (id.includes('atom-material-icons')) return 'atom-material-icons';
            if (id.includes('crawlab-ui')) return 'crawlab-ui';;
            return 'vendor.[hash]';
          }
        }
      },
      external: [
        // 'codemirror',
        // 'echarts',
      ],
      plugins: [
        // @ts-ignore
        // externalGlobals({
        //   // codemirror: 'CodeMirror',
        //   echarts: 'echarts',
        // })
      ],
    }
  },
  resolve: {
    dedupe: ['vue', 'element-plus', 'codemirror'],
    alias: [
      {find: '@', replacement: resolve(__dirname, 'src')},
    ],
    extensions: [
      '.js',
      '.ts',
      '.jsx',
      '.tsx',
      '.json',
      '.vue',
      '.scss',
    ]
  },
  plugins: [
    vue(),
    dynamicImport(),
    // splitVendorChunkPlugin(),
    // externalizeDeps(),
    // @ts-ignore
    visualizer({
      open: true,
      // open: false,
    }),
  ],
  server: {
    cors: true,
  },
});
