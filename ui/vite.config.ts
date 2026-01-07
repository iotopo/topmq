import path from 'node:path'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'

import Vue from '@vitejs/plugin-vue'
import Unocss from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'

import { defineConfig } from 'vite'

const srcPath = path.resolve(__dirname, 'src')

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    host: '0.0.0.0',
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },

  build: {
    target: 'esnext',
    // target: 'chrome89',
    minify: true,
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        minifyInternalExports: false,
      },
    },
  },
  resolve: {
    alias: {
      'vue-i18n': 'vue-i18n/dist/vue-i18n.esm-bundler.js',
      '@/': `${srcPath}/`,
    },
  },

  css: {
    preprocessorOptions: {
      scss: {
        // additionalData: `@use "@/styles/element/index.scss" as *;`,
        // api: "modern-compiler",
      },
    },
  },

  plugins: [
    Vue(),
    VueI18nPlugin({}),
    AutoImport({
      // Auto import functions from Vue, e.g. ref, reactive, toRef...
      // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
      imports: ['vue', 'vue-router', '@vueuse/core'],

      // Auto import functions from Element Plus, e.g. ElMessage, ElMessageBox... (with style)
      // 自动导入 Element Plus 相关函数，如：ElMessage, ElMessageBox... (带样式)
      resolvers: [ElementPlusResolver({})],
      vueTemplate: true,
      dts: path.resolve(srcPath, 'auto-imports.d.ts'),
    }),
    Components({
      // allow auto load markdown components under `./src/components/`
      extensions: ['vue', 'md'],
      // allow auto import and register components used in markdown
      include: [/\.vue$/, /\.vue\?vue/, /\.md$/],
      resolvers: [
        // Auto register icon components
        IconsResolver({
          enabledCollections: ['ep', 'mdi'],
        }),
        ElementPlusResolver({
          importStyle: 'sass',
        }),
      ],
      dts: path.resolve(srcPath, 'components.d.ts'),
    }),

    Icons({
      autoInstall: true,
    }),

    // https://github.com/antfu/unocss
    // see uno.config.ts for config
    Unocss(),
  ],
})
