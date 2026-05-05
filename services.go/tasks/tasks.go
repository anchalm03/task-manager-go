package tasks

import (
	"errors"
	"task_manager/errorcodes"
	"task_manager/middlewares"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// ---- view & response shapes ----

type TaskView struct {
	Id          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
	ProjectID   uuid.UUID  `json:"project_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toView(t *models.Task) TaskView {
	return TaskView{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		DueDate:     t.DueDate,
		AssignedTo:  t.AssignedTo,
		ProjectID:   t.ProjectID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

type TaskResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         *TaskView            `json:"data,omitempty"`
}

type TasksResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         []TaskView           `json:"data"`
}

type EmptyResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
}

// ---- request shapes ----

type CreateRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *string    `json:"assigned_to"` // user UUID as string
}

type UpdateRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

// AssignRequest is the body for POST /api/tasks/:id/assign.
// A null user_id unassigns the task.
type AssignRequest struct {
	UserID *string `json:"user_id"`
}

// ---- helpers ----

func msgFor(code errorcodes.ErrorCode, override string) string {
	if override != "" {
		return override
	}
	return code.Message()
}

func writeTask(c *gin.Context, code errorcodes.ErrorCode, override string, data *TaskView) {
	c.JSON(code.HttpStatusCode(), TaskResponse{
		Success:      code == errorcodes.NoError,
		ErrorCode:    code,
		ErrorMessage: msgFor(code, override),
		Data:         data,
	})
}

func writeTasks(c *gin.Context, code errorcodes.ErrorCode, override string, data []TaskView) {
	c.JSON(code.HttpStatusCode(), TasksResponse{
		Success:      code == errorcodes.NoError,
		ErrorCode:    code,
		ErrorMessage: msgFor(code, override),
		Data:         data,
	})
}

func writeEmpty(c *gin.Context, code errorcodes.ErrorCode, override string) {
	c.JSON(code.HttpStatusCode(), EmptyResponse{
		Success:      code == errorcodes.NoError,
		ErrorCode:    code,
		ErrorMessage: msgFor(code, override),
	})
}

// ---- handlers ----

// POST /api/projects/:id/tasks — only project owner can create a task
func CreateTask(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeTask(c, errorcodes.Unauthorized, "", nil)
		return
	}

	projectID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeTask(c, errorcodes.BadRequest, "invalid project id", nil)
		return
	}

	project, err := models.GetProjectById(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeTask(c, errorcodes.NotFound, "project not found", nil)
			return
		}
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}
	if project.OwnerID == nil || *project.OwnerID != userID {
		writeTask(c, errorcodes.Forbidden, "", nil)
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeTask(c, errorcodes.BadRequest, "", nil)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		ProjectID:   projectID,
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		assigneeID, err := uuid.FromString(*req.AssignedTo)
		if err != nil {
			writeTask(c, errorcodes.BadRequest, "invalid assigned_to", nil)
			return
		}
		if _, err := models.GetUserById(assigneeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeTask(c, errorcodes.BadRequest, "assignee user does not exist", nil)
				return
			}
			writeTask(c, errorcodes.InternalServerError, "", nil)
			return
		}
		task.AssignedTo = &assigneeID
	}

	if err := models.CreateTask(&task); err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}
	view := toView(&task)
	writeTask(c, errorcodes.NoError, "", &view)
}

// GET /api/projects/:id/tasks — project owner only (for now)
func ListTasksByProject(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeTasks(c, errorcodes.Unauthorized, "", nil)
		return
	}

	projectID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeTasks(c, errorcodes.BadRequest, "invalid project id", nil)
		return
	}

	project, err := models.GetProjectById(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeTasks(c, errorcodes.NotFound, "project not found", nil)
			return
		}
		writeTasks(c, errorcodes.InternalServerError, "", nil)
		return
	}
	if project.OwnerID == nil || *project.OwnerID != userID {
		writeTasks(c, errorcodes.Forbidden, "", nil)
		return
	}

	list, err := models.GetTasksByProjectId(projectID)
	if err != nil {
		writeTasks(c, errorcodes.InternalServerError, "", nil)
		return
	}
	views := make([]TaskView, 0, len(list))
	for i := range list {
		views = append(views, toView(&list[i]))
	}
	writeTasks(c, errorcodes.NoError, "", views)
}

// GET /api/me/tasks — tasks assigned to the caller
func ListMyTasks(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeTasks(c, errorcodes.Unauthorized, "", nil)
		return
	}

	list, err := models.GetTasksByUserId(userID)
	if err != nil {
		writeTasks(c, errorcodes.InternalServerError, "", nil)
		return
	}
	views := make([]TaskView, 0, len(list))
	for i := range list {
		views = append(views, toView(&list[i]))
	}
	writeTasks(c, errorcodes.NoError, "", views)
}

// PATCH /api/tasks/:id — project owner OR assignee can update status/due/title/description.
// Assignment changes go through POST /api/tasks/:id/assign.
func UpdateTask(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeTask(c, errorcodes.Unauthorized, "", nil)
		return
	}

	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeTask(c, errorcodes.BadRequest, "invalid task id", nil)
		return
	}

	task, err := models.GetTaskById(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeTask(c, errorcodes.NotFound, "", nil)
			return
		}
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}

	project, err := models.GetProjectById(task.ProjectID)
	if err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}
	isOwner := project.OwnerID != nil && *project.OwnerID == userID
	isAssignee := task.AssignedTo != nil && *task.AssignedTo == userID
	if !isOwner && !isAssignee {
		writeTask(c, errorcodes.Forbidden, "", nil)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeTask(c, errorcodes.BadRequest, "", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
	}
	if len(updates) == 0 {
		writeTask(c, errorcodes.BadRequest, "no fields to update", nil)
		return
	}

	if err := models.UpdateTask(taskID, updates); err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}

	updated, err := models.GetTaskById(taskID)
	if err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}
	view := toView(updated)
	writeTask(c, errorcodes.NoError, "", &view)
}

// POST /api/tasks/:id/assign — only the project owner can assign/unassign.
// Body: { "user_id": "<uuid>" }  → assign to that user
// Body: { "user_id": null }      → unassign
func AssignTask(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeTask(c, errorcodes.Unauthorized, "", nil)
		return
	}

	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeTask(c, errorcodes.BadRequest, "invalid task id", nil)
		return
	}

	task, err := models.GetTaskById(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeTask(c, errorcodes.NotFound, "", nil)
			return
		}
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}

	project, err := models.GetProjectById(task.ProjectID)
	if err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}
	if project.OwnerID == nil || *project.OwnerID != userID {
		writeTask(c, errorcodes.Forbidden, "only project owner can assign tasks", nil)
		return
	}

	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeTask(c, errorcodes.BadRequest, "", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.UserID == nil {
		updates["assigned_to"] = nil
	} else {
		assigneeID, err := uuid.FromString(*req.UserID)
		if err != nil {
			writeTask(c, errorcodes.BadRequest, "invalid user_id", nil)
			return
		}
		if _, err := models.GetUserById(assigneeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeTask(c, errorcodes.BadRequest, "assignee user does not exist", nil)
				return
			}
			writeTask(c, errorcodes.InternalServerError, "", nil)
			return
		}
		updates["assigned_to"] = assigneeID
	}

	if err := models.UpdateTask(taskID, updates); err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}

	updated, err := models.GetTaskById(taskID)
	if err != nil {
		writeTask(c, errorcodes.InternalServerError, "", nil)
		return
	}
	view := toView(updated)
	writeTask(c, errorcodes.NoError, "", &view)
}

// DELETE /api/tasks/:id — project owner only
func DeleteTask(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeEmpty(c, errorcodes.Unauthorized, "")
		return
	}

	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeEmpty(c, errorcodes.BadRequest, "invalid task id")
		return
	}

	task, err := models.GetTaskById(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeEmpty(c, errorcodes.NotFound, "")
			return
		}
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}

	project, err := models.GetProjectById(task.ProjectID)
	if err != nil {
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}
	if project.OwnerID == nil || *project.OwnerID != userID {
		writeEmpty(c, errorcodes.Forbidden, "")
		return
	}

	if err := models.DeleteTask(taskID); err != nil {
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}
	writeEmpty(c, errorcodes.NoError, "")
}
