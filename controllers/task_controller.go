package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SinofPride-999/student-time-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskForm struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	DueTime     string `form:"due_time" binding:"required"`
	Priority    int    `form:"priority"`
}

type TaskController struct {
	DB *gorm.DB
}

func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{DB: db}
}

func (tc *TaskController) GetDashboard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Get the user's username
	var user models.User
	if err := tc.DB.First(&user, userID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{"error": "User not found"})
		return
	}

	var tasks []models.Task

	// Get today's tasks
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	tc.DB.Where("user_id = ? AND due_time BETWEEN ? AND ?", userID, startOfDay, endOfDay).Order("due_time asc").Find(&tasks)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"tasks":    tasks,
		"now":      now,
		"username": user.Username,
	})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var form TaskForm

	// Debug: Print raw form data
	fmt.Println("Raw form data:", c.Request.PostForm)

	if err := c.ShouldBind(&form); err != nil {
		fmt.Println("Form bind error:", err)
		c.HTML(http.StatusBadRequest, "dashboard.html", gin.H{"error": "Invalid task data: " + err.Error()})
		return
	}

	// Debug: Print parsed form
	fmt.Printf("Parsed form: %+v\n", form)

	// Parse due_time with location
	dueTime, err := time.ParseInLocation("2006-01-02T15:04", form.DueTime, time.Local)
	if err != nil {
		fmt.Println("Time parse error:", err)
		c.HTML(http.StatusBadRequest, "dashboard.html", gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	task := models.Task{
		Title:       form.Title,
		Description: form.Description,
		DueTime:     dueTime,
		Priority:    form.Priority,
		UserID:      userID,
	}

	// Debug: Print task before save
	fmt.Printf("Task to create: %+v\n", task)

	// Enable GORM debug logging
	db := tc.DB.Debug()
	if err := db.Create(&task).Error; err != nil {
		fmt.Println("Database error:", err)
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
			"error": "Could not create task: " + err.Error(),
		})
		return
	}

	// Debug: Print after save
	fmt.Printf("Task created successfully with ID: %d\n", task.ID)

	c.Redirect(http.StatusFound, "/dashboard")
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var form TaskForm

	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "dashboard.html", gin.H{"error": "Invalid task data"})
		return
	}

	dueTime, err := time.Parse("2006-01-02T15:04", form.DueTime)
	if err != nil {
		c.HTML(http.StatusBadRequest, "dashboard.html", gin.H{"error": "Invalid date format"})
		return
	}

	// Fetch the existing task first
	var task models.Task
	if err := tc.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		c.HTML(http.StatusNotFound, "dashboard.html", gin.H{"error": "Task not found"})
		return
	}

	// Update task fields
	task.Title = form.Title
	task.Description = form.Description
	task.DueTime = dueTime
	task.Priority = form.Priority

	if err := tc.DB.Save(&task).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{"error": "Could not update task"})
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	if err := tc.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).Delete(&models.Task{}).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func (tc *TaskController) ToggleTaskComplete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var task models.Task

	if err := tc.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	task.Completed = !task.Completed
	tc.DB.Save(&task)
	c.Redirect(http.StatusFound, "/dashboard")
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Get the user's username
	var user models.User
	if err := tc.DB.First(&user, userID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "tasks_all.html", gin.H{"error": "User not found"})
		return
	}

	var tasks []models.Task

	// Debug: Print user ID
	fmt.Printf("Fetching tasks for user ID: %d\n", userID)

	// Fetch all tasks for this user
	if err := tc.DB.Debug().Where("user_id = ?", userID).Order("due_time asc").Find(&tasks).Error; err != nil {
		fmt.Printf("Error fetching tasks: %v\n", err)
		c.HTML(http.StatusInternalServerError, "tasks_all.html", gin.H{"error": "Could not retrieve tasks"})
		return
	}

	// Debug: Print found tasks
	fmt.Printf("Found %d tasks for user %d\n", len(tasks), userID)
	for i, task := range tasks {
		fmt.Printf("Task %d: %+v\n", i+1, task)
	}

	c.HTML(http.StatusOK, "tasks_all.html", gin.H{
		"allTasks": tasks,
		"username": user.Username,
	})
}
