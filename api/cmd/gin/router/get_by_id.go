package router

import (
	"api/internal/entity"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getByIDUsecase interface {
	Execute(ctx context.Context, id string) (*entity.User, error)
}

type getByIDRouter struct {
	uc getByIDUsecase
}

func NewGetByIDRouter(uc getByIDUsecase) *getByIDRouter {
	return &getByIDRouter{
		uc: uc,
	}
}

func (r *getByIDRouter) GetByIDRouter(api *gin.RouterGroup) {
	api.GET("/users/:id", func(c *gin.Context) {
		var (
			ctx = c.Request.Context()
			id  = c.Param("id")
		)

		user, err := r.uc.Execute(ctx, id)
		if err != nil {
			c.Error(err)
			return
		}

		if user == nil {
			c.Status(http.StatusNotFound)
			return
		}

		resp := UserByIDResponse{
			ID:           user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			ParentUserID: user.ParentUserID,
		}

		c.JSON(http.StatusOK, resp)
	})
}

type UserByIDResponse struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email_address"`
	ParentUserID *int64 `json:"parent_user_id"`
}
