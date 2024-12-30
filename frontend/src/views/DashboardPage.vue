<template>
  <div class="dashboard">
    <div class="row">
      <!-- 左側個人信息 -->
      <div class="col-md-3 sidebar bg-light p-3">
        <div class="text-center">
          <img src="https://via.placeholder.com/80" alt="User Avatar" class="rounded-circle mb-3" />
          <h4>{{ username }}</h4>
          <p>{{ roleLabel }}</p>
        </div>
        <ul class="nav flex-column">
          <li class="nav-item">
            <router-link to="/" class="nav-link text-primary">首頁</router-link>
          </li>
          <li v-if="role === 'student'" class="nav-item">
            <router-link to="/status" class="nav-link text-primary">聽課狀況</router-link>
          </li>
          <li v-if="role === 'student'" class="nav-item">
            <router-link to="/lectures" class="nav-link text-primary">講座資訊</router-link>
          </li>
          <li v-if="role === 'admin'" class="nav-item">
            <router-link to="/admin" class="nav-link text-primary">管理後台</router-link>
          </li>
          <li v-if="role === 'teacher'" class="nav-item">
            <router-link to="/teacher" class="nav-link text-primary">教師頁面</router-link>
          </li>
        </ul>
      </div>

      <!-- 右側內容 -->
      <div class="col-md-9 content p-3">
        <div v-if="role === 'student'" class="student-content">
          <h5>學生專屬內容</h5>
          <p>這裡顯示學生專屬資料或功能。</p>
        </div>
        <div v-if="role === 'admin'" class="admin-content">
          <h5>管理員專屬內容</h5>
          <p>這裡顯示管理員的功能面板。</p>
        </div>
        <div v-if="role === 'teacher'" class="teacher-content">
          <h5>教師專屬內容</h5>
          <p>這裡顯示教師專屬功能。</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    username: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      role: "", // 用於存儲用戶角色
    };
  },
  computed: {
    roleLabel() {
      switch (this.role) {
        case "student":
          return "學生";
        case "admin":
          return "管理員";
        case "teacher":
          return "教師";
        default:
          return "未知角色";
      }
    },
  },
  created() {
    // 從 localStorage 或 Vuex 獲取角色
    this.role = localStorage.getItem("role") || "guest";
  },
};
</script>

<style scoped>
.sidebar {
  height: 100vh;
  border-right: 1px solid #ddd;
}
</style>