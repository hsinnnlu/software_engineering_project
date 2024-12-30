<template>
    <div class="checking-page">
      <header class="d-flex justify-content-between align-items-center p-3 bg-light border-bottom">
        <a href="#" class="text-primary fw-bold">切換為簽退</a>
        <span>AI智慧未來 - 擁抱AI科技，共創AI新時代</span>
        <span>{{ currentTime }}</span>
        <a href="#" class="text-primary fw-bold">登出</a>
      </header>
  
      <main class="container my-4">
        <div class="row">
          <!-- 左側輸入框與資訊 -->
          <div class="col-md-6">
            <div class="input-group mb-3">
              <input
                type="text"
                class="form-control"
                placeholder="輸入學號"
                v-model="studentId"
              />
              <button class="btn btn-primary" @click="handleCheckIn">簽到</button>
            </div>
            <p class="text-secondary">
              <a href="#" class="text-primary" @click="viewRecords">查看簽到退記錄表</a>
            </p>
            <!-- 學生信息置中區域 -->
            <div class="student-info text-center">
              <h2 class="fw-bold">{{ studentName }}</h2>
              <h3>{{ studentId }}</h3>
              <h4>{{ formattedDate }}</h4>
              <h4>{{ currentTime }}</h4>
            </div>
          </div>
  
          <!-- 右側即時顯示區 -->
          <div class="col-md-6">
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
            <div class="text-end mt-2">
              <button class="btn btn-outline-secondary" @click="syncData">
                🔄 上傳雲端
              </button>
            </div>
          </div>
        </div>
  
        <!-- 底部按鈕 -->
        <div class="mt-4 text-center">
          <button class="btn btn-success mx-2">簽到成功</button>
          <button class="btn btn-danger mx-2">簽到失敗</button>
          <button class="btn btn-secondary mx-2">未簽到</button>
        </div>
      </main>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        lectureName: "",
        studentId: "", // 輸入的學號
        studentName: "學生姓名", // 模擬的學生姓名
        records: [
          // 示例記錄數據
          { id: 1, time: "14:00:30", studentId: "s1135xxxx", studentName: "王xx", synced: true },
        ],
        currentTime: "", // 當前時間
      };
    },
    computed: {
      // 格式化的當前日期
      formattedDate() {
        const date = new Date();
        return `${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, "0")}/${String(
          date.getDate()
        ).padStart(2, "0")}`;
      },
    },
    methods: {
      // 處理簽到功能
      handleCheckIn() {
        if (!this.studentId) {
          alert("請輸入學號");
          return;
        }
        // 添加新的簽到記錄（模擬）
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
      // 查看記錄表（模擬功能）
      viewRecords() {
        alert("查看簽到退記錄表");
      },
      // 模擬同步數據
      syncData() {
        alert("數據已同步到雲端！");
        this.records.forEach((record) => (record.synced = true)); // 標記同步完成
      },
      // 更新當前時間
      updateTime() {
        const now = new Date();
        this.currentTime = now.toLocaleTimeString();
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
    font-size: 28px;
    margin: 0;
  }
  .student-info h3 {
    font-size: 20px;
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
  </style>