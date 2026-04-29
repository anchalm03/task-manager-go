package comments

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

type CommentView struct {
	Id        uuid.UUID  `json:"id"`
	Text      string     `json:"text"`
	TaskID    uuid.UUID  `json:"task_id"`
	UserID    *uuid.UUID `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func toView(cm *models.Comment) CommentView {
	return CommentView{
		Id:        cm.Id,
		Text:      cm.Text,
		TaskID:    cm.TaskID,
		UserID:    cm.UserID,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	}
}

type CommentResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         *CommentView         `json:"data,omitempty"`
}

type CommentsResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         []CommentView        `json:"data"`
}

type EmptyResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
}

// ---- request shapes ----

type CreateRequest struct {
	Text string `json:"text" binding:"required"`
}

// ---- helpers ----

func msgFor(code errorcodes.ErrorCode, override string) string {
	if override != "" {
		return override
	}
	return code.Message()
}

func writeComment(c *gin.Context, code errorcodes.ErrorCode, override string, data *CommentView) {
	c.JSON(code.HttpStatusCode(), CommentResponse{
		Success:      code == errorcodes.NoError,
		ErrorCode:    code,
		ErrorMessage: msgFor(code, override),
		Data:         data,
	})
}

func writeComments(c *gin.Context, code errorcodes.ErrorCode, override string, data []CommentView) {
	c.JSON(code.HttpStatusCode(), CommentsResponse{
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

// canAccessTask returns true if the user is the project owner or the task assignee.
func canAccessTask(task *models.Task, userID uuid.UUID) (bool, error) {
	project, err := models.GetProjectById(task.ProjectID)
	if err != nil {
		return false, err
	}
	isOwner := project.OwnerID != nil && *project.OwnerID == userID
	isAssignee := task.AssignedTo != nil && *task.AssignedTo == userID
	return isOwner || isAssignee, nil
}

// ---- handlers ----

// POST /api/tasks/:id/comments
func CreateComment(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeComment(c, errorcodes.Unauthorized, "", nil)
		return
	}

	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeComment(c, errorcodes.BadRequest, "invalid task id", nil)
		return
	}

	task, err := models.GetTaskById(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeComment(c, errorcodes.NotFound, "task not found", nil)
			return
		}
		writeComment(c, errorcodes.InternalServerError, "", nil)
		return
	}

	allowed, err := canAccessTask(task, userID)
	if err != nil {
		writeComment(c, errorcodes.InternalServerError, "", nil)
		return
	}
	if !allowed {
		writeComment(c, errorcodes.Forbidden, "", nil)
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeComment(c, errorcodes.BadRequest, "", nil)
		return
	}

	comment := models.Comment{
		Text:   req.Text,
		TaskID: taskID,
		UserID: &userID,
	}
	if err := models.CreateComment(&comment); err != nil {
		writeComment(c, errorcodes.InternalServerError, "", nil)
		return
	}
	view := toView(&comment)
	writeComment(c, errorcodes.NoError, "", &view)
}

// GET /api/tasks/:id/comments
func ListCommentsByTask(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeComments(c, errorcodes.Unauthorized, "", nil)
		return
	}

	taskID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeComments(c, errorcodes.BadRequest, "invalid task id", nil)
		return
	}

	task, err := models.GetTaskById(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeComments(c, errorcodes.NotFound, "task not found", nil)
			return
		}
		writeComments(c, errorcodes.InternalServerError, "", nil)
		return
	}

	allowed, err := canAccessTask(task, userID)
	if err != nil {
		writeComments(c, errorcodes.InternalServerError, "", nil)
		return
	}
	if !allowed {
		writeComments(c, errorcodes.Forbidden, "", nil)
		return
	}

	list, err := models.GetCommentsByTaskId(taskID)
	if err != nil {
		writeComments(c, errorcodes.InternalServerError, "", nil)
		return
	}
	views := make([]CommentView, 0, len(list))
	for i := range list {
		views = append(views, toView(&list[i]))
	}
	writeComments(c, errorcodes.NoError, "", views)
}

// DELETE /api/comments/:id — comment author or project owner can delete
func DeleteComment(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeEmpty(c, errorcodes.Unauthorized, "")
		return
	}

	commentID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeEmpty(c, errorcodes.BadRequest, "invalid comment id")
		return
	}

	comment, err := models.GetCommentById(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeEmpty(c, errorcodes.NotFound, "")
			return
		}
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}

	isAuthor := comment.UserID != nil && *comment.UserID == userID
	isOwner := false
	if !isAuthor {
		task, err := models.GetTaskById(comment.TaskID)
		if err == nil {
			project, err := models.GetProjectById(task.ProjectID)
			if err == nil && project.OwnerID != nil && *project.OwnerID == userID {
				isOwner = true
			}
		}
	}
	if !isAuthor && !isOwner {
		writeEmpty(c, errorcodes.Forbidden, "")
		return
	}

	if err := models.DeleteComment(commentID); err != nil {
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}
	writeEmpty(c, errorcodes.NoError, "")
}
