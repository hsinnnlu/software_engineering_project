<template>
  <div class="d-flex align-items-center justify-content-center vh-100">
    <div class="row align-items-center">
      <!-- 左側 Logo 與標題 -->
      <div class="col-md-5 text-center">
        <img src="../assets/THU.png" alt="Tunghai University Logo" class="thu-logo mb-3" />
        <h1 class="h4 fw-bold">東海大學資訊工程學系</h1>
        <h2 class="h6">TUNGHAI UNIVERSITY COMPUTER SCIENCE DEPARTMENT</h2>
      </div>

      <!-- 右側登入表單 -->
      <div class="col-md-6">
        <form @submit.prevent="handleLogin" class="login-form p-4 shadow-sm rounded">
          <h3 class="text-center">登入</h3>
          <div class="form-group mb-3">
            <input
              type="text"
              class="form-control"
              placeholder="使用者名稱"
              v-model="form.username"
              required
            />
          </div>
          <div class="form-group mb-3">
            <input
              type="password"
              class="form-control"
              placeholder="密碼"
              v-model="form.password"
              required
            />
          </div>
          <button type="submit" class="btn btn-dark w-100">登入</button>
          <div class="text-start mt-2">
            <a href="#" class="text-primary h6" data-bs-toggle="modal" data-bs-target="#forgetpasswd">
              忘記密碼/首次登入?
            </a>
          </div>
        </form>
      </div>
    </div>

    <!-- 忘記密碼 Modal -->
    <div
      class="modal fade"
      id="forgetpasswd"
      tabindex="-1"
      aria-labelledby="forgetpasswdLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="forgetpasswdLabel">請輸入您的 tmail 帳號</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <input
              type="email"
              id="guestEmail"
              class="form-control"
              placeholder="輸入您的 tmail 帳號"
              v-model="guestEmail"
            />
          </div>
          <div class="modal-footer">
            <button class="btn btn-dark" @click="sendCode">送出</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      form: {
        username: "",
        password: "",
      },
      guestEmail: "", // 忘記密碼輸入的 email
      errorMessage: "", // 來自後端的訊息
    };
  },
  methods: {
    async handleLogin() {
      try {
        const response = await axios.post("/login", this.form);

        // 如果登入成功
        // alert("登入成功！");
        localStorage.setItem("token", response.data.token); // 儲存 JWT Token
        let role;
        switch(response.data.permission){
          case "1":
            role = "student";
            break;
          case "2":
            role = "manager";
            break;
          case "3":
            role = "professor";
            break;
          case "4":
            role = "assistant";
            break;
          default:
            role = "guest";
            break;
        }
        localStorage.setItem("role", role) // 存身份
        

        // this.$router.push("/")

        // 包含 user_id
        const user_id = response.data.user_id;
        // 動態跳轉到包含 user_id 的路由
        this.$router.push(`/${user_id}/`);
      } catch (error) {
        // 錯誤處理
        if (error.response && error.response.data) {
          // 從後端返回的錯誤訊息
          this.errorMessage = error.response.data.message || "登入失敗，請檢查帳號或密碼。";
        } else {
          // 其他錯誤
          this.errorMessage = "伺服器連線錯誤，請稍後再試。";
        }
      }
    },
    async sendCode() {
      try {
        if (!this.guestEmail) {
          alert("請輸入有效的 tmail 帳號！");
          return;
        }
        console.log({email: this.guestEmail})
        const response = await axios.post("/sendemail", { email: this.guestEmail });

        if (response.data.success) { 
          alert("重置密碼的鏈接已發送到您的信箱！");
          
        } else {
          alert("發送失敗：" + response.data.message);
        }
      } catch (error) {
        console.error("Error:", error);
        alert("伺服器錯誤，請稍後再試！");
      }
    },
  },
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 1rem;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.thu-logo {
  width: 150px; /* 設置寬度 */
  height: auto; /* 高度根據寬度自適應，保持圖片比例 */
  border-radius: 10%; /* 可選：讓圖片的邊框有輕微的圓角 */
  object-fit: cover; /* 如果圖片比例過長或過寬，進行裁剪以適應容器 */
  margin: 0 auto; /* 保持圖片水平居中 */
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); /* 添加輕微陰影效果 */
}
form div {
  margin-bottom: 1rem;
}
button {
  padding: 0.5rem 1rem;
}
</style>