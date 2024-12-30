import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "axios";


// 引入 Bootstrap 的 CSS 和 JS
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min.js";

// 創建應用實例
const app = createApp(App);

// 配置 Axios 的基地址
axios.defaults.baseURL = "http://localhost:8888"; // Golang 後端伺服器地址

// 全局注入 Axios（可選，視需求）
app.config.globalProperties.$axios = axios;

// 使用路由
app.use(router);

// 掛載應用
app.mount("#app");