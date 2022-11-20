package service

import (
	"crypto/sha256"
	"net/http"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

func NewUserForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new_user_form.html", gin.H{"Title": "Register user"})
}

func RegisterUser(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Username or Password is not provided", "Username": username, "Password": password})
		return
	}

	if len(password) < 8 {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password must be at least 8 characters", "Username": username, "Password": password})
		return
	}

	// password double check
	password2 := ctx.PostForm("password2")
	if password != password2 {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password and Password confirmation are not same"})
		return
	}

	// DB 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// 重複チェック
	var duplicate int
	err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if duplicate > 0 { // count の結果が返却されるので
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Username is already taken", "Username": username, "Password": password})
		return
	}

	// DB への保存
	result, err := db.Exec("INSERT INTO users(name, password) VALUES (?, ?)", username, hash(password))
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// 保存状態の確認
	id, _ := result.LastInsertId()
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE id = ?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// change password
func ChangePasswordForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "change_password_form.html", gin.H{"Title": "Change password"})
}

func ChangePassword(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	oldPassword := ctx.PostForm("old_password")
	newPassword := ctx.PostForm("new_password")

	if newPassword == "" {
		ctx.HTML(http.StatusBadRequest, "change_password_form.html", gin.H{"Title": "Change password", "Error": "New password is not provided"})
		return
	}

	if len(newPassword) < 8 {
		ctx.HTML(http.StatusBadRequest, "change_password_form.html", gin.H{"Title": "Change password", "Error": "New password must be at least 8 characters"})
		return
	}

	// db 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// old password is correct?
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE name = ?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if user.Password == nil || string(user.Password) != string(hash(oldPassword)) {
		ctx.HTML(http.StatusBadRequest, "change_password_form.html", gin.H{"Title": "Change password", "Error": "Old password is incorrect"})
		return
	}

	// update password
	_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", hash(newPassword), user.ID)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// private
func hash(password string) []byte {
	const salt = "todolist.go#"
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(password))
	return h.Sum(nil)
}
