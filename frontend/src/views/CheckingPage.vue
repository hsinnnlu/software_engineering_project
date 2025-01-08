<template>
  <div class="checking-page" :style="{ backgroundColor: backgroundColor }">
    <!-- 頂部區域 -->
    <header class="d-flex justify-content-between align-items-center p-2 border-bottom">
      <a v-if="this.isSignIn" href="#" class="text-primary fw-bold" @:click="switchStatus">切換為簽退</a>
      <a v-else href="#" class="text-primary fw-bold" @:click="switchStatus">切換為簽到</a>
      <span>{{ lectureName }}</span>
      <span>{{ currentTime }}</span>
      <a href="#" class="text-primary fw-bold" @click="endChecking">退出</a>
    </header>

    <!-- 主要內容區 -->
    <main class="container my-4">
      <div class="row h-100 mt-4">
        <!-- 左上：輸入學號 + 按鈕 -->
        <div class="col-md-2">
          <div class="input-group mb-3">
            <input
              type="text"
              class="form-control"
              placeholder="輸入學號"
              v-model="studentId"
            />
            <button class="btn btn-primary" @click="handleCheckIn">
              簽到
            </button>
          </div>
        </div>

        <!-- 中間：學生資訊置中區域 -->
        <div class="col-md-8">
          <div class="student-info text-center">
            <h2 class="fw-bold">{{ studentName }}</h2>
            <h3>{{ studentId }}</h3>
            <h4>{{ formattedDate }}</h4>
            <h4>{{ currentTime }}</h4>
          </div>
        </div>

        <!-- 右上：查看簽到退記錄表(Modal)按鈕 -->
        <div class="col-md-2 d-flex align-items-start justify-content-end">
          <button class="btn btn-outline-primary" @click="showModal = true">
            查看簽到退記錄表
          </button>
        </div>
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
            <!-- 即時顯示區域 -->
            <table class="table table-bordered text-center">
              <thead>
                <tr>
                  <th>狀態</th>
                  <th>時間</th>
                  <th>學號</th>
                  <th>姓名</th>
                  <th>同步</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in records" :key="record.id">
                  <td>{{ record.status }}</td>
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
import axios from 'axios';

export default {
  props: {
    lecture_name: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      // 畫面顯示的資料
      lectureName: this.lecture_name,
      studentId: "", // 輸入的學號
      studentName: "", // 學生名稱：待改進（後端必須回傳名字才能紀錄）
      lecture_id: "",
      isSignIn: true, 
      backgroundColor: "#f8f9fa",
      
      // 簽到資料
      records: [
        // 範例資料
        {
          id: 1,            // 內部顯示 id
          status: "簽到",
          time: "14:00:30", // 簽到時間
          studentId: "s1135xxxx", // 學號
          studentName: "王xx",    // 名稱
          synced: true,     // 是否同步
        },
      ],
      currentTime: "",   // 當前時間
      showModal: false,  // 控制彈跳視窗顯示/隱藏
    };
  },
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
    async handleCheckIn() {
      if (!this.studentId) {
        alert("請輸入學號");
        return;
      }
      const status = (this.isSignIn) ? "in" : "out";


      // 先把資料加入前端顯示的 records
      const newRecord = {
        id: this.records.length + 1,
        status: (status=="in") ? "簽到" : "簽退",
        time: new Date().toLocaleTimeString(),
        studentId: this.studentId,
        // studentName: this.studentName,
        synced: false,
      };

      this.records.push(newRecord);
      // 再呼叫 API 將資料送到後端
      try {
        const res = await this.sendCheckInRequest(status);
        
        // 如果後端同步，則更新同步狀況
        this.records[this.records.length - 1].synced = true;
        this.studentName = res?.data?.username;
        this.records[this.records.length - 1].studentName = this.studentName;

        this.backgroundColor = "#A5E696";        
        setTimeout(()=>{this.backgroundColor="#f8f9fa"}, 2000)
      } catch (error) {
        this.backgroundColor = "#FFA1A1"

        this.studentName="簽到失敗";
        setTimeout(()=>{this.studentName="", this.backgroundColor="#f8f9fa", this.studentId=""}, 2000); 
        console.error(error);
      }

      // 延遲顯示

    },
    // 實際送出簽到的請求 (sign_in, sign_out)
    async sendCheckInRequest(Status) {
      // 準備要送出的資料結構 { sign_in_time, sign_out_time, status }
      const now = new Date().toISOString();
      const token = localStorage.getItem("token");
      var requestBody;
      const route = (Status==="in") ? "sign-in" : "sign-out"
      if(Status == "in"){
          requestBody = {
            sign_in_time: now,
            sign_out_time: null,  // 簽到時先給空字串或 null
            status: Status,  // 用來區分簽到或簽退
          };
      } else{
        requestBody = {
            sign_in_time: null,
            sign_out_time: now,  // 簽到時先給空字串或 null
            status: Status,  // 用來區分簽到或簽退
          };
      }

      try{
        const res = await axios.post(
          `/${this.lecture_id}/`+`${this.studentId}/${route}`,
          requestBody,
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
        return res
      } catch(error){
        throw new Error("Network response was not ok");
      }
    },
    // 退出簽到，回到主畫面
    endChecking(){
      const username = localStorage.getItem('username');
      this.$router.push({ 
        path: `${username}`, 
      });
    },

    // 獲取伺服器上的資料
    async fetchCheckingList() {
      const token = localStorage.getItem("token");

      // 清空現有記錄，避免數據重複
      this.records = [];

      try {
        // 發送 GET 請求
        const res = await axios.get(`/list-checking/${this.lecture_id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (res.data && Array.isArray(res.data.records)) {
          res.data.records.forEach((rec) => {
            // 簽到列
            if (rec.sign_in_time) {
              this.records.push({
                id: this.records.length + 1, // 唯一 ID
                status: "簽到",
                time: rec.sign_in_time,
                studentId: rec.user_id || "未知學號",
                studentName: rec.user_name || "未知名稱",
                synced: true, // 假設數據已同步
              });
            }

            // 簽退列
            if (rec.sign_out_time) {
              this.records.push({
                id: this.records.length + 1, // 唯一 ID
                status: "簽退",
                time: rec.sign_out_time,
                studentId: rec.user_id || "未知學號",
                studentName: rec.user_name || "未知名稱",
                synced: true, // 假設數據已同步
              });
            }
          });
        } else {
          console.log("沒東西")
        }
      } catch (error) {
        // 處理請求錯誤
        console.error("Error fetching checking list:", error);

        // 顯示錯誤訊息
        alert(
          error.response?.data?.error ||
            "無法獲取簽到列表，請稍後再試！"
        );
      }
    },
    

    // 模擬「查看簽到退記錄表」的資料同步
    syncData() {
      alert("數據已同步到雲端！");
      this.records.forEach((record) => (record.synced = true)); // 標記同步完成
    },

    // 開關彈跳視窗
    closeModal() {
      this.showModal = false;
    },

    // 切換模式
    switchStatus(){
      this.isSignIn = !this.isSignIn;
    },

    // 更新當前時間
    updateTime() {
      const now = new Date();
      this.currentTime = now.toLocaleTimeString();
    },
  }, // methods end
  mounted() {
    // 每秒更新時間
    this.updateTime();
    setInterval(this.updateTime, 1000);
  },
  created(){
    this.lecture_id = this.$route.query.lecture_id;
    this.lectureName = this.$route.query.lecture_name;
    this.fetchCheckingList();
  }
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