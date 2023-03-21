package user

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo_list/connect_db"
	"todo_list/models"
)

func md5Hash(string string) string {
	hash := md5.Sum([]byte(string))
	return hex.EncodeToString(hash[:])
}

func Register(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		log.Print(err.Error())
		return err
	}
	db := connect_db.GetDB()
	passwordHex := md5Hash(u.Password)
	insert, err := db.Prepare("insert into users (Id,Name,Email,Password,Age) values (?,?,?,?,?)")
	if err != nil {
		log.Print(err.Error())
	}
	_, err = insert.Exec(u.Id, u.Name, u.Email, passwordHex, u.Age)
	if err != nil {
		log.Print(err.Error())
	}
	return c.JSON(http.StatusCreated, u)
}
func Login(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		log.Print(err.Error())
		return err
	}
	db := connect_db.GetDB()

	passwordHex := md5Hash(u.Password)

	var pwd string
	var id int
	err := db.QueryRow("select id, password from users where email =?", u.Email).Scan(&id, &pwd)
	if err != nil {
		log.Print(err.Error())
	}
	if pwd == passwordHex {
		claims := &models.JwtCustomClaims{
			Id:       id,
			Email:    u.Email,
			Password: true,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("mysecretkey"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
	}
	return c.JSON(http.StatusUnauthorized, "Sorry, unable to verify your information")
}
func GetUserIntToken(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	userEmail := claims.Email

	db := connect_db.GetDB()

	var u = &models.User{}

	err := db.QueryRow("select *from users where email=?", userEmail).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Age, &u.Avatar)
	if err != nil {
		log.Print(err.Error())
	}

	return c.JSON(http.StatusOK, u)
}
func UpdateUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	userEmail := claims.Email

	db := connect_db.GetDB()

	u := new(models.User)
	err := c.Bind(u)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	passwordHex := md5Hash(u.Password)
	_, err = db.Exec("update users set id=?, name=?,email=?,password=?,age=? where email=?", u.Id, u.Name, u.Email, passwordHex, u.Age, userEmail)
	if err != nil {
		log.Print(err.Error())
	}
	return c.JSON(http.StatusOK, u)
}
func DeleteUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	userEmail := claims.Email

	db := connect_db.GetDB()

	_, err := db.Exec("DELETE FROM users WHERE email = ?", userEmail)
	if err != nil {
		log.Print(err.Error())
	}
	return c.JSON(200, "delete succesful")
}
func UploadImage(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)

	u := new(models.User)

	err := c.Bind(&u)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	db := connect_db.GetDB()
	// Lưu dữ liệu vào cơ sở dữ liệu SQL
	_, err = db.Exec("update  user set avatar=? where email=?", u.Avatar, claims.Email)
	if err != nil {
		log.Println("======error update===", err.Error())
		respErr := models.Error{
			Message: err.Error(),
			Status:  500,
			Type:    "Update",
			Code:    50001,
		}
		return c.JSON(respErr.Status, respErr)
	}

	return c.JSON(http.StatusOK, u.Avatar)
}
func AddTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	userId := claims.Id
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return err
	}
	db := connect_db.GetDB()
	result, err := db.Exec("INSERT INTO tasks (user_id, description , completed , created_at) VALUES ( ?, ?,?,?)", userId, task.Description, task.Completed, time.Now())
	if err != nil {
		log.Println("======error update===", err.Error())
		respErr := models.Error{
			Message: err.Error(),
			Status:  500,
			Type:    "???",
		}
		return c.JSON(respErr.Status, respErr)
	}
	idInsert, err := result.LastInsertId()
	if err != nil {
		log.Println("======error update===", err.Error())
		respErr := models.Error{
			Message: err.Error(),
			Status:  500,
			Type:    "???",
		}
		return c.JSON(respErr.Status, respErr)
	}
	task.Id = int(idInsert)
	task.UserId = userId
	return c.JSON(http.StatusCreated, task)
}
func GetAllTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	userId := claims.Id

	tasks := []models.Task{}
	db := connect_db.GetDB()
	rows, err := db.Query("SELECT * FROM tasks WHERE user_id = ?", userId)
	if err != nil {
		log.Print(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.Id, &task.UserId, &task.Description, &task.Completed, &task.CreateAt)
		if err != nil {
			log.Print(err.Error())
		}
		tasks = append(tasks, task)
	}

	return c.JSON(http.StatusOK, tasks)
}
func GetTaskByID(c echo.Context) error {
	// Lấy giá trị Id từ parameter
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid task Ic",
		})
	}
	db := connect_db.GetDB()

	task := models.Task{}
	err = db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskId).Scan(&task.Id, &task.UserId, &task.Description, &task.Completed, &task.CreateAt)
	if err != nil {
		log.Print(err.Error())
	}

	return c.JSON(http.StatusOK, task)
}
func GetTasksByCompleted(c echo.Context) error {
	//lấy tham số truy vấn
	completed, err := strconv.ParseBool(c.QueryParam("completed"))

	tasks := make([]models.Task, 0)

	db := connect_db.GetDB()
	rows, err := db.Query("SELECT * FROM tasks WHERE completed = ?", completed)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.UserId, &task.Description, &task.Completed, &task.CreateAt)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		tasks = append(tasks, task)
	}

	return c.JSON(http.StatusOK, tasks)
}
func UpdateById(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid task Ic",
		})
	}
	db := connect_db.GetDB()

	task := new(models.Task)
	err = c.Bind(task)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	existingTask := &models.Task{}
	err = db.QueryRow("select * from tasks where id = ?", taskId).Scan(&existingTask.Id, &existingTask.UserId, &existingTask.Description, &existingTask.Completed, &existingTask.CreateAt)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Task not found")
	}
	existingTask.Completed = task.Completed
	_, err = db.Exec("update tasks set completed=? where id=?", existingTask.Completed, taskId)
	if err != nil {
		log.Print(err.Error())
	}

	return c.JSON(http.StatusOK, existingTask)
}
func DeleteById(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid task ",
		})
	}
	db := connect_db.GetDB()

	_, err = db.Exec("delete from tasks where id=?", taskId)
	if err != nil {
		log.Print(err.Error())
	}
	return c.JSON(http.StatusOK, "delete successfully")
}
func GetTasksPagination(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	userId := claims.Id
	page_size, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		log.Print(err.Error())
	}
	// Kết nối cơ sở dữ liệu
	db := connect_db.GetDB()
	var total int
	err = db.QueryRow("select count(*) as total from tasks where id=?", userId).Scan(&total)
	if err != nil {
		log.Print(err.Error())
	}
	var page int
	if total%page_size == 1 {
		page = (total)/page_size + 1
	} else if total%page_size == 0 {
		page = (total) / page_size
	}
	offset := (page - 1) * page_size
	rows, err := db.Query("select * from tasks order by id limit ? offset ?", page_size, offset)
	if err != nil {
		log.Print(err.Error())
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.UserId, &task.Description, &task.Completed, &task.CreateAt); err != nil {
			return err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		log.Print(err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}
