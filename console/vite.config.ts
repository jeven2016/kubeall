import {defineConfig} from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "src"),
            "@bindings": path.resolve(__dirname, "bindings", "tools-gui3", "pkg", "service")
        }
    },
    server: {
        proxy: {
            // 简单字符串形式：将 '/api' 开头的请求代理到目标地址
            '/api': 'http://localhost:8080',

            // 更详细的配置方式
            // '/api': {
            //   target: 'http://localhost:5000', // 后端服务地址
            //   changeOrigin: true,              // 修改请求头中的 Host 为目标地址
            //   rewrite: (path) => path.replace(/^\/api/, '') // 移除请求路径中的 '/api' 前缀
            // }
        }
    },
});
