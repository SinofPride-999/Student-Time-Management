package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SinofPride-999/student-time-management/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

func (ac *AdminController) AdminDashboard(c *gin.Context) {
    fmt.Println("AdminDashboard handler called")
		
    // Verify admin status again for safety
    session := sessions.Default(c)
    if session.Get("isAdmin") == nil || !session.Get("isAdmin").(bool) {
        c.Redirect(http.StatusFound, "/login")
        return
    }

    // Get statistics with error handling
    var userCount, taskCount int64
    if err := ac.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "error": "Failed to get user count",
        })
        return
    }

    if err := ac.DB.Model(&models.Task{}).Count(&taskCount).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "error": "Failed to get task count",
        })
        return
    }

    // Get all users with their task counts
    var users []struct {
        Username  string
        TaskCount int64
        CreatedAt time.Time
    }
    if err := ac.DB.Model(&models.User{}).
        Select("users.username, count(tasks.id) as task_count, users.created_at").
        Joins("left join tasks on tasks.user_id = users.id").
        Group("users.id").
        Scan(&users).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "error": "Failed to get user statistics",
        })
        return
    }

    // Get recent tasks
    var recentTasks []models.Task
    if err := ac.DB.Preload("User").Order("created_at desc").Limit(10).Find(&recentTasks).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "error": "Failed to get recent tasks",
        })
        return
    }

    c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
        "userCount":    userCount,
        "taskCount":    taskCount,
        "users":        users,
        "recentTasks":  recentTasks,
        "currentRoute": "dashboard",
        "now":         time.Now(),
    })
}
