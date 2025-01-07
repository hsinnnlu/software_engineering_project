<template>
    <div class="container mt-3">
      <!-- 新增講座 -->
      <button class="btn btn-primary mb-3" @click="showModal = true">
        新增講座
      </button>
  
      <!-- 新增講座 Modal -->
      <div
        v-if="showModal"
        class="modal fade show"
        style="display: block; background-color: rgba(0, 0, 0, 0.6)"
      >
        <div class="modal-dialog">
          <div class="modal-content custom-modal">
            <div class="modal-header">
              <h5 class="modal-title">新增講座</h5>
              <button type="button" class="btn-close" @click="showModal = false"></button>
            </div>
            <div class="modal-body">
              <label>講座名稱</label>
              <input v-model="newLecture.title" class="form-control" />
  
              <label>日期</label>
              <input v-model="newLecture.date" class="form-control" />
  
              <label>時間</label>
              <input v-model="newLecture.time" class="form-control" />
  
              <label>地點</label>
              <input v-model="newLecture.place" class="form-control" />
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="showModal = false">取消</button>
              <button class="btn btn-primary" @click="addLecture">確定新增</button>
            </div>
          </div>
        </div>
      </div>
      <!-- end 新增講座 Modal -->
  
      <!-- 編輯講座 Modal -->
      <div
        v-if="showEditModal"
        class="modal fade show"
        style="display: block; background-color: rgba(0, 0, 0, 0.6)"
      >
        <div class="modal-dialog">
          <div class="modal-content custom-modal">
            <div class="modal-header">
              <h5 class="modal-title">編輯講座</h5>
              <button type="button" class="btn-close" @click="cancelEditLecture"></button>
            </div>
            <div class="modal-body">
              <label>講座名稱</label>
              <input v-model="editLectureForm.title" class="form-control" />
  
              <label>日期</label>
              <input v-model="editLectureForm.date" class="form-control" />
  
              <label>時間</label>
              <input v-model="editLectureForm.time" class="form-control" />
  
              <label>地點</label>
              <input v-model="editLectureForm.place" class="form-control" />
  
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="cancelEditLecture">取消</button>
              <button class="btn btn-primary" @click="confirmEditLecture">確定修改</button>
            </div>
          </div>
        </div>
      </div>
      <!-- end 編輯講座 Modal -->
  
      <!-- 編輯參與人員 Modal -->
      <div
        v-if="showParticipantEditModal"
        class="modal fade show"
        style="display: block; background-color: rgba(0, 0, 0, 0.6)"
      >
        <div class="modal-dialog">
          <div class="modal-content custom-modal">
            <div class="modal-header">
              <h5 class="modal-title">編輯參與人員</h5>
              <button type="button" class="btn-close" @click="cancelEditParticipant"></button>
            </div>
            <div class="modal-body">
              <label>學生</label>
              <input v-model="editParticipantForm.studentName" class="form-control" />
              <label>指導教授</label>
              <input v-model="editParticipantForm.professor" class="form-control" />
              <label>學號</label>
              <input v-model="editParticipantForm.studentID" class="form-control" />
              <label>簽到</label>
              <input v-model="editParticipantForm.signIn" class="form-control" />
              <label>簽退</label>
              <input v-model="editParticipantForm.signOut" class="form-control" />
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="cancelEditParticipant">取消</button>
              <button class="btn btn-primary" @click="confirmEditParticipant">確定修改</button>
            </div>
          </div>
        </div>
      </div>
      <!-- end 編輯參與人員 Modal -->
  
      <!-- 講座列表 -->
      <table class="table table-bordered text-center thcolor">
        <thead>
          <tr>
            <th style="width: 50px">序號</th>
            <th>講座名稱</th>
            <th>日期</th>
            <th>時間</th>
            <th>地點</th>
            <th style="width: 80px">編輯</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="(lecture, index) in lectures" :key="lecture.id">
            <tr style="background-color: #fff">
              <td @click="toggleExpanded(lecture.id)" style="cursor: pointer">{{ index + 1 }}</td>
              <td>{{ lecture.title }}</td>
              <td>{{ lecture.date }}</td>
              <td>{{ lecture.time }}</td>
              <td>{{ lecture.place }}</td>
              <td @click.stop="editLecture(lecture.id)">編輯</td>
            </tr>
  
            <tr v-if="expandedLectureId === lecture.id" style="background-color: #f8f9fa">
              <td colspan="7">
                <table class="table table-bordered text-center">
                  <thead style="background-color: #e9ecef">
                    <tr>
                      <th>學生</th>
                      <th>指導教授</th>
                      <th>學號</th>
                      <th>簽到</th>
                      <th>簽退</th>
                      <th style="width: 80px">編輯</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(p, pIndex) in lecture.participants" :key="pIndex">
                      <td>{{ p.studentName }}</td>
                      <td>{{ p.professor }}</td>
                      <td>{{ p.studentID }}</td>
                      <td>{{ p.signIn }}</td>
                      <td>{{ p.signOut }}</td>
                      <td @click.stop="editParticipant(p, lecture.id, pIndex)">編輯</td>
                    </tr>
                  </tbody>
                </table>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </template>
  
  <script>
import axios from 'axios';

  export default {
    data() {
      return {
        showModal: false,
        newLecture: { title: "", date: "", time: "", place: ""},
        showEditModal: false,
        editLectureForm: { id: null, title: "", date: "", time: "", place: ""},
        showParticipantEditModal: false,
        editParticipantForm: {
          studentName: "",
          professor: "",
          studentID: "",
          signIn: "",
          signOut: "",
          lectureId: null,
          pIndex: null,
        },
        lectures: [],
        expandedLectureId: null,
      };
    },
    methods: {
      // 請求講座列表
      async fetchLectures(){
        
        // 先請求講座列表
        const token = localStorage.getItem("token");
        try{
          const response = await axios.post(
            "/lecture",
            null,
            { headers: { Authorization: `Bearer ${token}` }} ,
          )
          // 獲取講座數據
          const fetchedLectures = response.data.lecture || [];

          // 清空當前列表並更新為新數據
          this.lectures = []; // 清空舊數據
          fetchedLectures.forEach((lec) => {

            let date = "未知日期";
            let time = "未知時間";

            // 確保 timestamp 存在並可解析
            if (lec.timestamp && typeof lec.timestamp === "string") {
              const [parsedDate, parsedTime] = lec.timestamp.split("T");
              date = parsedDate || "未知日期";
              time = parsedTime || "未知時間";
            }
            // 確保格式一致
            this.lectures.push({
              id: lec.id,
              title: lec.name,
              date: date,
              time: time,
              place: lec.location,

              participants: lec.participants || [], // 確保參與人員為空數組而不是 null
            });
          });
        } catch(error) {
          console.error("Failed to fetch lectures:", error);
            alert("無法獲取講座數據，請稍後再試！");
        }
        
      },
      

      toggleExpanded(id) {
        this.expandedLectureId = this.expandedLectureId === id ? null : id;
      },
      addLecture() {
        if (!this.newLecture.title) return;
        const newId = this.lectures.length
          ? Math.max(...this.lectures.map((l) => l.id)) + 1
          : 1;
        this.lectures.push({ ...this.newLecture, id: newId, participants: [] });
        this.newLecture = { title: "", date: "", time: "", place: ""};
        this.showModal = false;
      },
      // 編輯講座
      editLecture(id) {
        const lecture = this.lectures.find((l) => l.id === id);
        if (!lecture) return;
        this.editLectureForm = JSON.parse(JSON.stringify(lecture));
        this.showEditModal = true;
      },
      // 確認編輯
      async confirmEditLecture() {
        // 找到需要編輯的講座索引
        const idx = this.lectures.findIndex((l) => l.id === this.editLectureForm.id);
        if (idx === -1) {
          alert("未找到對應的講座！");
          return;
        }

        const token = localStorage.getItem("token");
        console.log(this.editLecture);
        try {
          // 發送請求到後端
          const response = await axios.post(
            "/edit-lecture",
            this.editLectureForm, // 發送表單數據
            {
              headers: { Authorization: `Bearer ${token}` },
            }
          );

          // 從後端獲取更新後的數據，並更新到本地列表
          const updatedLecture = response.data.lecture;
          this.lectures.splice(idx, 1, updatedLecture);

          // 取消編輯模式
          this.cancelEditLecture();
          alert("講座修改成功！");
        } catch (error) {
          console.error("Failed to update lecture:", error);
          alert("講座修改失敗，請稍後再試！");
        }
      },
      cancelEditLecture() {
        this.showEditModal = false;
        this.editLectureForm = { id: null, title: "", date: "", time: "", place: ""};
      }, 
      // 編輯參與人員
      editParticipant(p, lectureId, pIndex) {
        this.editParticipantForm = { ...JSON.parse(JSON.stringify(p)), lectureId, pIndex };
        this.showParticipantEditModal = true;
      },
      // 確認編輯部分
      confirmEditParticipant() {
        const lecture = this.lectures.find((l) => l.id === this.editParticipantForm.lectureId);
        if (!lecture) return;
        const idx = this.editParticipantForm.pIndex;
        if (idx < 0 || idx >= lecture.participants.length) return;
        lecture.participants.splice(idx, 1, {
          studentName: this.editParticipantForm.studentName,
          professor: this.editParticipantForm.professor,
          studentID: this.editParticipantForm.studentID,
          signIn: this.editParticipantForm.signIn,
          signOut: this.editParticipantForm.signOut,
        });
        this.cancelEditParticipant();
      },
      cancelEditParticipant() {
        this.showParticipantEditModal = false;
        this.editParticipantForm = {
          studentName: "",
          professor: "",
          studentID: "",
          signIn: "",
          signOut: "",
          lectureId: null,
          pIndex: null,
        };
      },
    },
    created(){
      this.fetchLectures();
    }
  };
  </script>
  
<style scoped>
.custom-modal {
  background-color: #333;
  color: #fff;
}
</style>