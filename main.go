package main

import (
	"github.com/labstack/echo/v4"
	"todo_list/connect_db"
	verify "todo_list/middleware"
	"todo_list/user"
)

func main() {
	e := echo.New()
	db := connect_db.DBSql{}
	db.New()
	e.POST("/user/register", user.Register)
	e.POST("/user/login", user.Login)
	e.GET("/user/me", user.GetUserIntToken, verify.VerifyJWT)
	e.PUT("/user/me", user.UpdateUser, verify.VerifyJWT)
	e.DELETE("/user/me", user.DeleteUser, verify.VerifyJWT)
	e.POST("/user/me/avatar", user.UploadImage, verify.VerifyJWT)
	e.POST("/task", user.AddTask, verify.VerifyJWT)
	e.GET("/task", user.GetAllTask, verify.VerifyJWT)
	e.GET("task/:id", user.GetTaskByID, verify.VerifyJWT)
	e.GET("/task/getTask", user.GetTasksByCompleted, verify.VerifyJWT)
	e.PUT("/task/:id", user.UpdateById, verify.VerifyJWT)
	e.DELETE("/task/:id", user.DeleteById, verify.VerifyJWT)
	e.GET("/task/byPagination", user.GetTasksPagination, verify.VerifyJWT)
	e.Logger.Fatal(e.Start(":1323"))
}
