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

// Helper to debug loaded templates
func debugTemplates() gin.HandlerFunc {
	return func(c *gin.Context) {
		var names []string
		for _, t := range tmpl.Templates() {
			names = append(names, t.Name())
		}
		c.JSON(http.StatusOK, gin.H{"templates": names})
	}
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 1. FIRST register debug routes
	r.GET("/debug/templates", func(c *gin.Context) {
		var names []string
		for _, t := range tmpl.Templates() {
			names = append(names, t.Name())
		}
		c.JSON(200, gin.H{"templates": names})
	})

	// Parse templates globally
	tmpl = template.Must(tmpl.ParseGlob("templates/pages/*.html"))
	r.SetHTMLTemplate(tmpl)

	// Test route
	r.GET("/test", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{"error": nil})
	})

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
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/dashboard") })
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

	return r
}
