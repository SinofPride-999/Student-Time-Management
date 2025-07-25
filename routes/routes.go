package routes

import (
	"html/template"
	"net/http"

	"github.com/SinofPride-999/student-time-management/controllers"
	"github.com/SinofPride-999/student-time-management/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tmpl *template.Template

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize templates properly
	tmpl = template.New("")
	tmpl, err := tmpl.ParseGlob("templates/pages/*.html")
	if err != nil {
		panic("Failed to parse templates: " + err.Error())
	}
	r.SetHTMLTemplate(tmpl)

	// Core middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Session middleware
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Static files
	r.Use(static.Serve("/static", static.LocalFile("./static", true)))

	// User routes
	userController := controllers.NewUserController(db)
	r.GET("/", userController.ShowLogin)
	r.GET("/signup", userController.ShowSignup)
	r.POST("/signup", userController.Signup)
	r.GET("/login", userController.ShowLogin)
	r.POST("/login", userController.Login)
	r.GET("/logout", userController.Logout)

	// Auth-protected routes
	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthRequired())
	authGroup.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user").(uint)
		c.Set("userID", userID)
		c.Next()
	})

	taskController := controllers.NewTaskController(db)
	authGroup.GET("/dashboard", taskController.GetDashboard)
	authGroup.POST("/tasks", taskController.CreateTask)
	authGroup.POST("/tasks/:id", taskController.UpdateTask)
	authGroup.POST("/tasks/:id/delete", taskController.DeleteTask)
	authGroup.POST("/tasks/:id/toggle", taskController.ToggleTaskComplete)
	authGroup.GET("/tasks/all", taskController.GetAllTasks)

	// Admin routes
	adminGroup := r.Group("/admin")
	adminGroup.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			isAdmin := session.Get("isAdmin")
			
			if isAdmin == nil || !isAdmin.(bool) {
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
			}
			c.Next()
	})

	adminController := controllers.NewAdminController(db)
	adminGroup.GET("/dashboard", adminController.AdminDashboard)

	// Only one admin login route needed
	r.POST("/admin/login", userController.AdminLogin)

	return r
}
