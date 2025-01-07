<template>
  <div id="app">
    <header v-show="showNav">
      <nav>
        <ul>
          <li><router-link to="/">Home</router-link></li>
          <li><router-link to="/login">開發者登入按鈕</router-link></li>
        </ul>
      </nav>
    </header>
    <!-- 動態渲染路由對應的內容 -->
    <router-view/>
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
  background: #212529;
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
/* 簽到退頁面中的button */
.custom-btn {
    background-color: white;      
    color: #0d6efd;              
    border-color: white;          
}
.custom-btn:hover {
    background-color: #0d6efd;     
    color: white;  
    border-color: #0d6efd;              
}

/* 是否雲端同步之標示 */
.greencircle {
    width: 20px;            
    height: 20px;             
    background-color: rgb(148, 237, 137);; 
    border-radius: 50%;  
    display: inline-block;
}
.redcircle {
    width: 20px;            
    height: 20px;             
    background-color: rgb(232, 118, 118);; 
    border-radius: 50%;  
    display: inline-block;
}

/* 表格設計 */
.table-bordered th, .table-bordered td{
    border: 1.5px solid #545454;
}

/* 首頁設計 */
.custom-modal {
    background-color: #333;
    color: #fff;
    border-radius: 8px;
}
.custom-modal .modal-header,
.custom-modal .modal-footer {
    border: none;
}
</style>
