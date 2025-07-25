package controllers

import (
	"net/http"

	"github.com/SinofPride-999/student-time-management/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) ShowSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func (uc *UserController) Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "Invalid form data"})
		return
	}

	// Hash password
	if err := user.HashPassword(user.Password); err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{"error": "Could not hash password"})
		return
	}

	// Create user
	if err := uc.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "User already exists"})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

func (uc *UserController) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (uc *UserController) Login(c *gin.Context) {
	var inputUser models.User
	var dbUser models.User

	if err := c.ShouldBind(&inputUser); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid form data"})
		return
	}

	// Find user
	if err := uc.DB.Where("username = ?", inputUser.Username).First(&dbUser).Error; err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := dbUser.CheckPassword(inputUser.Password); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	// Set session
	session := sessions.Default(c)
	session.Set("user", dbUser.ID)
	session.Set("username", dbUser.Username)
	session.Save()

	c.Redirect(http.StatusFound, "/dashboard")
}

func (uc *UserController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}
