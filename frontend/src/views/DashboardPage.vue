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
          <li v-if="role === 'manager'" class="nav-item">
            <router-link to="/addpage" class="nav-link text-primary">加入講座</router-link>
          </li>
          <li v-if="role === 'assistant'" class="nav-item">
              <router-link to="/record" class="nav-link text-primary">管理簽到退系統</router-link>
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
        <div v-if="role === 'manager'" class="admin-content">
          <h5>管理員專屬內容</h5>
          <p>這裡顯示管理員的功能面板。</p>
        </div>
        <div v-if="role === 'professor'" class="teacher-content">
          <h5>教授專屬內容</h5>
          <p>這裡顯示教授專屬功能。</p>
        </div>
        <div v-if="role === 'assistant'" class="teacher-content">
          <h5>講座管理人員內容</h5>
          <p>查看講座列表的欄位</p>

          <!-- 檢查是否有講座資料 -->
          <div v-if="lectures.length > 0">
            <table class="table table-bordered">
              <thead>
                <tr>
                  <th>講座名稱</th>
                  <th>時間</th>
                  <th>講師</th>
                  <th>地點</th>
                  <th>進入簽到系統</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="lecture in lectures" :key="lecture.id">
                  <td>{{ lecture.name }}</td>
                  <td>{{ lecture.timestamp }}</td>
                  <td>{{ lecture.speaker }}</td>
                  <td>{{ lecture.location }}</td>
                  <td>
                    <button class="btn btn-sm btn-primary" @click="startChecking(lecture.id, lecture.name)">進入</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 如果沒有講座資料 -->
          <div v-else>
            <p>目前沒有任何講座資料。</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';


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

      lectures: []
    };
  },
  methods:{
    async fetchLectures(){
      const token = localStorage.getItem('token');
      try{
        const response = await axios.post(
          "/lecture",
          {},
          { headers: { Authorization: `Bearer ${token}` } }
        )
        this.lectures = response.data.lecture;
      } catch(error){
        console.error("Failed to fetch lectures:", error);
        this.lectures = []; // 重置資料
        this.errorMessage =
          error.response?.data?.error || "無法獲取講座資料，請稍後再試。";
      }
    },
    // 切換到 CheckingPage
    startChecking(lectureId, lecture_name) {
      this.$router.push({ 
        path: "/checking", 
        query: { 
          lecture_id: lectureId,
          lecture_name: lecture_name
        } 
      });
    },
  },
  computed: {
    roleLabel() {
      switch (this.role) {
        case "student":
          return "學生";
        case "manager":
          return "管理員";
        case "professor":
          return "教授";
        case "assistant":
          return "講座管理人員"
        default:
          return "未知角色";
      }
    },

  },
  created() {
    // 從 localStorage 或 Vuex 獲取角色
    this.role = localStorage.getItem("role") || "guest";
    this.fetchLectures();
    console.log(this.lectures);
  },
};
</script>

<style scoped>
.sidebar {
  height: 100vh;
  border-right: 1px solid #ddd;
}
</style>