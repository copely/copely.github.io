import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { fileURLToPath, URL } from 'node:url';

function resolveBasePath(raw?: string): string {
  const trimmed = raw?.trim();
  if (!trimmed) {
    return '/';
  }
  const withLeading = trimmed.startsWith('/') ? trimmed : `/${trimmed}`;
  return withLeading.endsWith('/') ? withLeading : `${withLeading}/`;
}

export default defineConfig({
  base: resolveBasePath(process.env.VITE_BASE_PATH),
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
