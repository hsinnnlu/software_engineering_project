<template>
  <div id="app">
    <header v-show="showNav">
      <nav>
        <ul>
          <li><router-link to="/">Home</router-link></li>
          <li><router-link to="/login">Login</router-link></li>
        </ul>
      </nav>
    </header>
    <!-- 動態渲染路由對應的內容 -->
    <router-view />
  </div>
</template>

<script>
export default {
  computed: {
    // 根據路由決定是否顯示 <nav>
    showNav() {
      // 檢查當前路由名稱或路徑
      const hiddenRoutes = ["/login"]; // 定義需要隱藏 <nav> 的路由
      return !hiddenRoutes.includes(this.$route.path);
    },
  },
};
</script>

<style scoped>

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
header {
  background: #42b983;
  padding: 10px 0;
}
nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  justify-content: center;
}
nav ul li {
  margin: 0 15px;
}
nav ul li a {
  color: white;
  text-decoration: none;
  font-weight: bold;
}
nav ul li a:hover {
  text-decoration: underline;
}
</style>
