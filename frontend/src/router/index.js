import { createRouter, createWebHistory } from "vue-router";
import LoginPage from "@/views/LoginPage.vue";
import DashboardPage from "@/views/DashboardPage.vue";
import CheckingPage from "@/views/CheckingPage.vue";
import signRecordPage from "@/views/signRecordPage.vue";
import AddPage from "@/views/AddPage.vue";

const routes = [
  { path: "/", 
    redirect: () => {
      // 從 localStorage 獲取用戶名
      const username = localStorage.getItem("username") || "Guest";
      return `/${username}`;
    },
  },
  { path: "/login", name: "LoginPage", component: LoginPage },
  {
    path: "/:username",
    name: "Dashboard",
    component: DashboardPage,
    props: true, // 傳遞 username 作為 props
  },
  {
    path: "/checking",
    name: "Checking",
    component: CheckingPage,
    props: true
  },
  {
    path: "/record",
    name: "signRecordPage",
    component: signRecordPage,
    props: true,
  },
  {
    path: "/addpage",
    name: "AddPage",
    component: AddPage,
    props: true
  }
];


const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
  });
  
export default router;