<template>
  <div class="checking-page">
    <!-- 頂部區域 -->
    <header class="d-flex justify-content-between align-items-center p-2 bg-light border-bottom">
      <a href="#" class="text-primary fw-bold">切換為簽退</a>
      <span>AI智慧未來 - 擁抱AI科技，共創AI新時代</span>
      <span>{{ currentTime }}</span>
      <a href="#" class="text-primary fw-bold">登出</a>
    </header>

    <!-- 主要內容區 -->
    <main class="container my-4">
      <div class="row h-100 mt-4">
        <div class="col-md-2">
          <div class="input-group mb-3">
            <input
              type="text"
              class="form-control"
              placeholder="輸入學號"
              v-model="studentId"
            />
            <button class="btn btn-primary" @click="handleCheckIn">簽到</button>
          </div>
        </div>
        <!-- 左側輸入框與資訊 -->
        <div class="col-md-8">
          <!-- 學生資訊置中區域 -->
          <div class="student-info text-center">
            <h2 class="fw-bold">{{ studentName }}</h2>
            <h3>{{ studentId }}</h3>
            <h4>{{ formattedDate }}</h4>
            <h4>{{ currentTime }}</h4>
          </div>
        </div>

        <!-- 即時顯示彈跳視窗按鈕 -->
        <div class="col-md-2 d-flex align-items-start justify-content-end">
          <button class="btn btn-outline-primary" @click="showModal = true">
            查看簽到退記錄表
          </button>
        </div>
      </div>

      <!-- 底部按鈕 -->
      <div class="mt-4 text-center">
        <button class="btn btn-success mx-2">簽到成功</button>
        <button class="btn btn-danger mx-2">簽到失敗</button>
      </div>
    </main>

    <!-- 彈跳視窗（Modal） -->
    <div
      class="modal fade"
      tabindex="-1"
      role="dialog"
      :class="{ show: showModal }"
      style="display: block;"
      v-if="showModal"
    >
      <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">簽到退記錄表</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <!-- 原先的即時顯示區域 -->
            <table class="table table-bordered text-center">
              <thead>
                <tr>
                  <th>時間</th>
                  <th>學號</th>
                  <th>姓名</th>
                  <th>同步</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in records" :key="record.id">
                  <td>{{ record.time }}</td>
                  <td>{{ record.studentId }}</td>
                  <td>{{ record.studentName }}</td>
                  <td>
                    <span class="text-success" v-if="record.synced">✓</span>
                    <span class="text-danger" v-else>✗</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="modal-footer">
            <button class="btn btn-outline-secondary me-auto" @click="syncData">
              🔄 上傳雲端
            </button>
            <button type="button" class="btn btn-secondary" @click="closeModal">
              關閉
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- 彈跳視窗 END -->
  </div>
</template>

<script>
export default {
  data() {
    return {
      // 畫面顯示的資料
      lectureName: "",
      studentId: "", // 輸入的學號
      studentName: "學生", // 模擬的學生姓名
      
      // 簽到 資料
      records: [
        // 範例資料
        {
          id: 1, // 內部顯示 id
          time: "14:00:30", // 簽到時間
          studentId: "s1135xxxx", // 學號
          studentName: "王xx", // 名稱
          synced: true, // 是否同步
        },
      ],
      currentTime: "", // 當前時間
      showModal: false, // 控制彈跳視窗顯示/隱藏
    };
  },
  // 會一直重複執行的東東
  computed: {
    // 格式化的當前日期
    formattedDate() {
      const date = new Date();
      return `${date.getFullYear()}/${String(date.getMonth() + 1).padStart(
        2,
        "0"
      )}/${String(date.getDate()).padStart(2, "0")}`;
    },
  },
  methods: {
    // 處理手動簽到功能
    handleCheckIn() {
      if (!this.studentId) {
        alert("請輸入學號");
        return;
      }
      // 添加新的簽到記錄 ！！
      const newRecord = {
        id: this.records.length + 1,
        time: new Date().toLocaleTimeString(),
        studentId: this.studentId,
        studentName: this.studentName,
        synced: false,
      };
      this.records.push(newRecord);
      this.studentId = ""; // 清空輸入框
    },
    // 模擬同步數據 ！！
    syncData() {
      alert("數據已同步到雲端！");
      this.records.forEach((record) => (record.synced = true)); // 標記同步完成
    },
    // 更新當前時間
    updateTime() {
      const now = new Date();
      this.currentTime = now.toLocaleTimeString();
    },
    // 開關彈跳視窗
    closeModal() {
      this.showModal = false;
    },
  },
  mounted() {
    // 每秒更新時間
    this.updateTime();
    setInterval(this.updateTime, 1000);
  },
};
</script>

<style scoped>
.checking-page {
  font-family: Arial, sans-serif;
  background-color: #f8f9fa;
  min-height: 100vh;
}
header {
  font-size: 14px;
}
.student-info {
  text-align: center; /* 確保文本置中 */
  margin-top: 30px; /* 添加一些上邊距 */
}
.student-info h2 {
  font-size: 5em;
  margin: 0;
  margin-top: 30%;
}
.student-info h3 {
  font-size: 3em;
  margin: 0;
  color: gray;
}
.student-info h4 {
  font-size: 18px;
  margin: 5px 0;
}
.table {
  background-color: #fff;
}

/* Bootstrap Modal 手動顯示的處理 (若不透過Bootstrap原生的JavaScript) */
.modal.show {
  display: block;
  background: rgba(0, 0, 0, 0.3); /* 背景遮罩 */
}
.modal-dialog {
  margin-top: 10%; /* 調整彈窗垂直位置 */
}
.btn-close {
  border: none;
  background: none;
  font-size: 1.5rem;
}
</style>