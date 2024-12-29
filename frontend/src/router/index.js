import { createRouter, createWebHistory } from "vue-router";
import LoginPage from "@/views/LoginPage.vue";
import DashboardPage from "@/views/DashboardPage.vue";
import CheckingPage from "@/views/CheckingPage.vue";


const routes = [
//   { path: "/", name: "HomePage", component: HomePage },
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
  }
];


const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
  });
  
export default router;