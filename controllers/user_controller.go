package controllers

import (
	"fmt"
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
	// Create a struct that matches your capitalized form fields
	type LoginForm struct {
		Username string `form:"Username" binding:"required"`
		Password string `form:"Password" binding:"required"`
		IsAdmin  string `form:"is_admin"`
	}

	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid form data"})
		return
	}

	// Debug print to verify values are received
	fmt.Printf("Login attempt - Username: %s, IsAdmin: %s\n", form.Username, form.IsAdmin)

	// Admin login flow
	if form.IsAdmin == "true" {
		if form.Username == "SinofPride-999" && form.Password == "darkplace" {
			session := sessions.Default(c)
			session.Set("user", 0) // Special admin ID
			session.Set("username", "SinofPride-999")
			session.Set("isAdmin", true)
			if err := session.Save(); err != nil {
				c.HTML(http.StatusInternalServerError, "login.html", gin.H{
					"error": "Failed to save session",
				})
				return
			}
			c.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid admin credentials",
		})
		return
	}

	// Regular user login flow
	var dbUser models.User
	if err := uc.DB.Where("username = ?", form.Username).First(&dbUser).Error; err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := dbUser.CheckPassword(form.Password); err != nil {
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

func (uc *UserController) AdminLogin(c *gin.Context) {
	fmt.Println("AdminLogin handler called")

	type AdminLoginForm struct {
		Username string `form:"Username" binding:"required"`
		Password string `form:"Password" binding:"required"`
	}

	var form AdminLoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid form data"})
		return
	}

	// Hardcoded admin credentials check
	if form.Username != "SinofPride-999" || form.Password != "darkplace" {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid admin credentials"})
		return
	}

	// Set session with special admin flag
	session := sessions.Default(c)
	session.Set("user", 0) // Using 0 as special admin ID
	session.Set("username", "SinofPride-999")
	session.Set("isAdmin", true)
	session.Save()

	c.Redirect(http.StatusFound, "/admin/dashboard")
}
