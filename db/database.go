package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB // 全局變量, 小寫db是外地注入變量

func InitDB(dataSourceName string) error{
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Successfully connected to database!")
	return err
}

// func GetUserById(db *sql.DB, user_id string) (*models.User, error) {
// 	user := models.User{}
// 	query := "SELECT user_id, password_hash, user_permission, user_name FROM users WHERE user_id = ?"
// 	err := DB.QueryRow(query, user_id).Scan(&user.Id, &user.Password_hash, &user.Permission, &user.Name)
// 	if err != nil {
// 		// 如果找不到使用者
// 		if err == sql.ErrNoRows {
// 			fmt.Println("user not found")
// 			return nil, fmt.Errorf("user not found")
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func GetUserByEmail(email string) (*models.User, error) {
// 	user := models.User{}
// 	query := "SELECT user_id FROM users WHERE user_email = ?"
// 	err := DB.QueryRow(query, email).Scan(&user.Id)
// 	if err != nil {
// 		// 如果找不到使用者
// 		if err == sql.ErrNoRows {
// 			fmt.Println("user not found")
// 			return nil, fmt.Errorf("user not found")
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func GetpasswordByEmail(email string) (*models.User, error) {
// 	user := models.User{}
// 	query := "SELECT password_hash FROM users WHERE user_email = ?"
// 	err := DB.QueryRow(query, email).Scan(&user.Password_hash)
// 	if err != nil {
// 		// 如果找不到密碼
// 		if err == sql.ErrNoRows {
// 			fmt.Println("user not found")
// 			return nil, fmt.Errorf("user not found")
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // 更新用戶密碼
// func UpdatePasswordByEmail(db *sql.DB, email, hashedPassword string) error {
// 	query := "UPDATE users SET password_hash = ? WHERE user_email = ?"
// 	_, err := db.Exec(query, hashedPassword, email)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func UpdatePasswordByUserid(db *sql.DB, user_id, hashedPassword string) error {
// 	query := "UPDATE users SET password_hash = ? WHERE user_id = ?"
// 	_, err := db.Exec(query, hashedPassword, user_id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // 新增講座 ， timestamp格式還沒確定
// func InsertLecture(c *gin.Context, lecture models.Lecture) error {

// 	// 准备插入数据
// 	query := "INSERT INTO Lectures(lecture_id, lecture_name, lecture_speaker, lecture_timestamp, lecture_manager, lecture_location) values(?, ?, ?, ?, ?, ?)"
// 	stmt, err := DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	fmt.Printf("lecture ID is: %d\n", lecture.Id)

// 	// 插入数据
// 	_, err = stmt.Exec(lecture.Id, lecture.Name, lecture.Speaker, lecture.Timestamp, lecture.Manager, lecture.Location)
// 	if err != nil {
// 		fmt.Printf("Insert lecture failed, err: %v\n", err)
// 		return err
// 	}

// 	fmt.Print("Insert lecture success\n")

// 	c.JSON(200, gin.H{
// 		"message": "success for insert lecture",
// 	})

// 	return nil
// }

// func GetLectureById(db *sql.DB, lecture_id int) (*models.Lecture, error) {
// 	lecture := models.Lecture{}
// 	query := "SELECT lecture_id, lecture_name, lecture_speaker, lecture_timestamp, lecture_manager, lecture_location, status FROM Lectures WHERE lecture_id = ?"
// 	err := DB.QueryRow(query, lecture_id).Scan(&lecture.Id, &lecture.Name, &lecture.Speaker, &lecture.Timestamp, &lecture.Manager, &lecture.Location)
// 	if err != nil {
// 		// 如果找不到講座
// 		if err == sql.ErrNoRows {
// 			fmt.Println("lecture not found")
// 			return nil, fmt.Errorf("lecture not found")
// 		}
// 		return nil, err
// 	}
// 	return &lecture, nil
// }
// func GetLecturesByStatus(db *sql.DB, lecture_status ...string) ([]models.Lecture, error) {
// 	lectures := []models.Lecture{}
// 	query := "SELECT lecture_id, lecture_name, lecture_speaker, lecture_timestamp, lecture_manager, lecture_location, status FROM Lectures WHERE status = ?"

// 	for _, status := range lecture_status {
// 		rows, err := DB.Query(query, status)
// 		if err != nil {
// 			fmt.Printf("skip: %s\n", status)
// 			continue
// 		}
// 		defer rows.Close()

// 		for rows.Next() {
// 			l := models.Lecture{}
// 			err := rows.Scan(&l.Id, &l.Name, &l.Speaker, &l.Timestamp, &l.Manager, &l.Location, &l.Status)
// 			if err != nil {
// 				fmt.Printf("err: %s\n", err)
// 				return nil, err
// 			}
// 			lectures = append(lectures, l)
// 		}
// 	}
// 	return lectures, nil

// }

// func GetAnnouncementList(db *sql.DB) ([]models.Announce, error) {
// 	announces := []models.Announce{}
// 	query := "SELECT announce_id, announce_title, announce_content, announce_date FROM Announcements"
// 	rows, err := DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		announce := models.Announce{}
// 		err := rows.Scan(&announce.Id, &announce.Title, &announce.Content, &announce.Date)
// 		if err != nil {
// 			return nil, err
// 		}
// 		announces = append(announces, announce)
// 	}
// 	return announces, nil
// }

// func GetAttendanceRecord(db *sql.DB, userID string) ([]models.User_attendance, error) {
// 	user_attendances := []models.User_attendance{}

// 	// 修改查詢語句，依據 user 進行過濾
// 	query := `
// 		SELECT
// 			Lecture.lecture_timestamp,
// 			Lecture.lecture_name,
// 			User_Lecture.sign_in_time,
// 			User_Lecture.sign_out_time,
// 			Users.qty_lecture
// 		FROM
// 			User_Lecture
// 		JOIN
// 			Lecture ON User_Lecture.lecture_id = Lecture.lecture_id
// 		JOIN
// 			Users ON User_Lecture.user_id = Users.user_id
// 		WHERE
// 			User_Lecture.user_id = ?;
// 	`

// 	// 執行查詢時，將 userID 傳入
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// 逐行處理查詢結果
// 	for rows.Next() {
// 		user_attendance := models.User_attendance{}
// 		err := rows.Scan(&user_attendance.Name, &user_attendance.Timestamp, &user_attendance.Sign_in_time, &user_attendance.Sign_out_time, &user_attendance.Qty_lecture)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// 將每筆紀錄加到切片中
// 		user_attendances = append(user_attendances, user_attendance)
// 	}

// 	return user_attendances, nil
// }

// func GetLectureInformation(db *sql.DB, userID string) ([]models.Lecture, error) {
// 	Lectures := []models.Lecture{}

// 	// 沒有備註的部分
// 	query := "SELECT lecture_name, lecture_timestamp, lecture_location, lecture_speaker FROM Lecture WHERE user = ?"

// 	// 執行查詢時，將 userID 傳入
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// 逐行處理查詢結果
// 	for rows.Next() {
// 		Lecture := models.Lecture{}
// 		err := rows.Scan(&Lecture.Name, &Lecture.Timestamp, &Lecture.Location, &Lecture.Speaker)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// 將每筆紀錄加到切片中
// 		Lectures = append(Lectures, Lecture)
// 	}

// 	return Lectures, nil
// }
