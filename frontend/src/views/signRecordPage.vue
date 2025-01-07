<template>
    <div class="container">
      <div class="row d-flex align-items-center justify-content-between pt-3 pb-2">
        <div class="col-md-6 d-flex align-items-center">
          <router-link to="/professor">
            <button class="btn btn-outline-dark h6 me-2 shadow-sm">
              <i class="bi bi-box-arrow-left"></i>
              返回
            </button>
          </router-link>
          <p class="h3 fw-bold">學生出席紀錄/心得單</p>
        </div>
        <div class="col-md-6 d-flex justify-content-end">
          <input v-model="searchKeyword" type="search" placeholder="請輸入關鍵字" />
          <button @click="search" class="btn btn-dark">
            <i class="bi bi-search"></i>
          </button>
        </div>
      </div>
  
      <div class="row table-responsive-md">
        <table class="table table-bordered text-center align-middle">
          <thead>
            <tr class="bg-dark text-white">
              <th scope="col">序號</th>
              <th scope="col">系級</th>
              <th scope="col">學號</th>
              <th scope="col">姓名</th>
              <th scope="col">累計次數</th>
              <th scope="col">更多</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(student, index) in filteredStudents" :key="student.id">
              <th scope="row">{{ index + 1 }}</th>
              <td>{{ student.department }}</td>
              <td>{{ student.studentId }}</td>
              <td>{{ student.name }}</td>
              <td>{{ student.totalAttendance }}</td>
              <td>
                <i
                  class="bi bi-chevron-down"
                  role="button"
                  @click="toggleDetails(index)"
                ></i>
              </td>
            </tr>
            <tr v-if="detailsIndex === index" :key="`details-${index}`">
              <td colspan="8">
                <div>
                  <strong>演講列表:</strong>
                  <table class="table table-sm table-bordered mt-3 text-center align-middle">
                    <thead class="table-light">
                      <tr>
                        <th scope="col">序號</th>
                        <th scope="col">講座名稱</th>
                        <th scope="col">日期</th>
                        <th scope="col">時間</th>
                        <th scope="col">地點</th>
                        <th scope="col">講者</th>
                        <th scope="col">心得單</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(lecture, lectureIndex) in student.lectures" :key="lecture.id">
                        <th scope="row">{{ lectureIndex + 1 }}</th>
                        <td>{{ lecture.title }}</td>
                        <td>{{ lecture.date }}</td>
                        <td>{{ lecture.time }}</td>
                        <td>{{ lecture.location }}</td>
                        <td>{{ lecture.speaker }}</td>
                        <td>
                          <button type="button" class="btn btn-sm">
                            <i class="bi bi-file-text fs-4"></i>
                          </button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        searchKeyword: "",
        students: [
          {
            id: 1,
            department: "碩士一",
            studentId: "g1235xxxx",
            name: "王xx",
            totalAttendance: 1,
            lectures: [
              {
                id: 1,
                title: "AI智變未來 - 擁抱AI科技，共創AI新時代",
                date: "2024.04.10",
                time: "18:00-20:00",
                location: "人文大樓(H104)",
                speaker: "陳國禎",
              },
            ],
          },
        ],
        detailsIndex: null,
      };
    },
    computed: {
      filteredStudents() {
        if (!this.searchKeyword) return this.students;
        return this.students.filter((student) =>
          student.name.includes(this.searchKeyword) ||
          student.studentId.includes(this.searchKeyword)
        );
      },
    },
    methods: {
      search() {
        // 搜尋功能可以依需求擴展
        console.log("搜尋關鍵字:", this.searchKeyword);
      },
      toggleDetails(index) {
        this.detailsIndex = this.detailsIndex === index ? null : index;
      },
    },
  };
  </script>
  
  <style scoped>
  .container {
    background-color: #f0f0f0;
  }
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
  