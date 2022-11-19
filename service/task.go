package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT * FROM tasks") // Use DB#Select for multiple entries
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task.html", task)
}

// form to create new task
func NewTaskForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form_new_task.html", gin.H{"Title": "New task registration"})
}

// create new task
func NewTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get form data
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")

	// Insert a task
	result, err := db.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render status
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	if id, err := result.LastInsertId(); err == nil {
		path = fmt.Sprintf("/task/%d", id) // 正常にIDを取得できた場合は /task/<id> へ戻る
	}
	ctx.Redirect(http.StatusFound, path)
}

// form to edit task
func EditTaskForm(ctx *gin.Context) {
	// ID の取得
	// /task/:id/edit
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	// Get target task
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Render edit form
	ctx.HTML(http.StatusOK, "form_edit_task.html",
		gin.H{"Title": fmt.Sprintf("Edit task %d", task.ID), "Task": task})
}

// edit task
func EditTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// /task/:id/edit
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Get form data
	title := ctx.PostForm("title")
	is_done := ctx.PostForm("is_done")
	description := ctx.PostForm("description")

	// 型変換
	is_done_bool, _ := strconv.ParseBool(is_done)

	// Update a task
	result, err := db.Exec("UPDATE tasks SET title=?, is_done=?, description=? WHERE id=?", title, is_done_bool, description, id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render status
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	if rows, err := result.RowsAffected(); err == nil && rows == 1 {
		path = fmt.Sprintf("/task/%d", id) // 正常に1行更新できた場合は /task/<id> へ戻る
	}
	ctx.Redirect(http.StatusFound, path)
}
