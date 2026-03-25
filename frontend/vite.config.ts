import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    open: false,
    host: true,
    // Docker内ではAPI_TARGET環境変数でbackendサービス名を指定する
    proxy: {
      '/api': {
        target: process.env.API_TARGET,
        changeOrigin: true,
      },
    },
    // Dockerボリュームマウント環境でのホットリロードにポーリングを使用
    watch: {
      usePolling: true,
    },
  },
});
