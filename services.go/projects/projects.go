package projects

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

type ProjectView struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OwnerID     *uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toView(p *models.Project) ProjectView {
	return ProjectView{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type ProjectResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         *ProjectView         `json:"data,omitempty"`
}

type ProjectsResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         []ProjectView        `json:"data"`
}

type EmptyResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
}

// ---- request shapes ----

type CreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// ---- helpers ----

func msgFor(code errorcodes.ErrorCode, override string) string {
	if override != "" {
		return override
	}
	return code.Message()
}

func writeProject(c *gin.Context, code errorcodes.ErrorCode, override string, data *ProjectView) {
	c.JSON(code.HttpStatusCode(), ProjectResponse{
		Success:      code == errorcodes.NoError,
		ErrorCode:    code,
		ErrorMessage: msgFor(code, override),
		Data:         data,
	})
}

func writeProjects(c *gin.Context, code errorcodes.ErrorCode, override string, data []ProjectView) {
	c.JSON(code.HttpStatusCode(), ProjectsResponse{
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

// POST /api/projects
func CreateProject(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeProject(c, errorcodes.Unauthorized, "", nil)
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeProject(c, errorcodes.BadRequest, "", nil)
		return
	}

	project := models.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     &userID,
	}
	if err := models.CreateProject(&project); err != nil {
		writeProject(c, errorcodes.InternalServerError, "", nil)
		return
	}

	view := toView(&project)
	writeProject(c, errorcodes.NoError, "", &view)
}

// GET /api/projects — projects owned by caller
func ListProjects(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeProjects(c, errorcodes.Unauthorized, "", nil)
		return
	}

	projects, err := models.GetAllProjectsByUserId(userID)
	if err != nil {
		writeProjects(c, errorcodes.InternalServerError, "", nil)
		return
	}

	views := make([]ProjectView, 0, len(projects))
	for i := range projects {
		views = append(views, toView(&projects[i]))
	}
	writeProjects(c, errorcodes.NoError, "", views)
}

// GET /api/projects/:id
func GetProject(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeProject(c, errorcodes.Unauthorized, "", nil)
		return
	}

	projectID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeProject(c, errorcodes.BadRequest, "invalid project id", nil)
		return
	}

	project, err := models.GetProjectById(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeProject(c, errorcodes.NotFound, "", nil)
			return
		}
		writeProject(c, errorcodes.InternalServerError, "", nil)
		return
	}

	if project.OwnerID == nil || *project.OwnerID != userID {
		writeProject(c, errorcodes.Forbidden, "", nil)
		return
	}

	view := toView(project)
	writeProject(c, errorcodes.NoError, "", &view)
}

// PATCH /api/projects/:id
func UpdateProject(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeProject(c, errorcodes.Unauthorized, "", nil)
		return
	}

	projectID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeProject(c, errorcodes.BadRequest, "invalid project id", nil)
		return
	}

	project, err := models.GetProjectById(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeProject(c, errorcodes.NotFound, "", nil)
			return
		}
		writeProject(c, errorcodes.InternalServerError, "", nil)
		return
	}
	if project.OwnerID == nil || *project.OwnerID != userID {
		writeProject(c, errorcodes.Forbidden, "", nil)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeProject(c, errorcodes.BadRequest, "", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if len(updates) == 0 {
		writeProject(c, errorcodes.BadRequest, "no fields to update", nil)
		return
	}

	if err := models.UpdateProject(projectID, updates); err != nil {
		writeProject(c, errorcodes.InternalServerError, "", nil)
		return
	}

	updated, err := models.GetProjectById(projectID)
	if err != nil {
		writeProject(c, errorcodes.InternalServerError, "", nil)
		return
	}
	view := toView(updated)
	writeProject(c, errorcodes.NoError, "", &view)
}

// DELETE /api/projects/:id
func DeleteProject(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeEmpty(c, errorcodes.Unauthorized, "")
		return
	}

	projectID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		writeEmpty(c, errorcodes.BadRequest, "invalid project id")
		return
	}

	project, err := models.GetProjectById(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeEmpty(c, errorcodes.NotFound, "")
			return
		}
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}
	if project.OwnerID == nil || *project.OwnerID != userID {
		writeEmpty(c, errorcodes.Forbidden, "")
		return
	}

	if err := models.DeleteProject(projectID); err != nil {
		writeEmpty(c, errorcodes.InternalServerError, "")
		return
	}
	writeEmpty(c, errorcodes.NoError, "")
}
