package router

import (
	"api/internal/entity"
	"api/internal/usecase"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type searchUsecase interface {
	Execute(ctx context.Context, filter usecase.SearchInput) ([]entity.User, error)
}

type searchRouter struct {
	uc searchUsecase
}

func NewSearchRouter(uc searchUsecase) *searchRouter {
	return &searchRouter{
		uc: uc,
	}
}

func (r *searchRouter) SearchRouter(api *gin.RouterGroup) {
	api.GET("/users", func(c *gin.Context) {
		var (
			firstName = c.Query("first_name")
			lastName  = c.Query("last_name")
			email     = c.Query("email_address")

			fields = c.Query("fields")
		)

		if err := isValidFields(fields); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		input := usecase.SearchInput{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,

			Fields: fields,
		}

		users, err := r.uc.Execute(c.Request.Context(), input)
		if err != nil {
			c.Error(err)
			return
		}

		response := make([]SearchResponse, 0, len(users))
		for _, u := range users {
			r := setFields(&u, strings.Split(fields, ",")...)
			response = append(response, r)
		}

		c.JSON(http.StatusOK, response)
	})
}

type SearchResponse struct {
	ID           *int64  `json:"id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Email        *string `json:"email_address,omitempty"`
	ParentUserID *int64  `json:"parent_user_id,omitempty"`
}

var validFields = map[string]struct{}{
	"id":             {},
	"first_name":     {},
	"last_name":      {},
	"email_address":  {},
	"created_at":     {},
	"deleted_at":     {},
	"merged_at":      {},
	"parent_user_id": {},
}

func isValidFields(str string) error {
	fields := strings.Split(str, ",")
	for _, f := range fields {
		if _, ok := validFields[f]; !ok {
			return fmt.Errorf("invalid field option %s", f)
		}
	}
	return nil
}

func setFields(user *entity.User, fields ...string) SearchResponse {
	s := SearchResponse{}
	for _, field := range fields {
		switch strings.TrimSpace(field) {
		case "id":
			s.ID = &user.ID
		case "first_name":
			s.FirstName = &user.FirstName
		case "last_name":
			s.LastName = &user.LastName
		case "email_address":
			s.Email = &user.Email
			// case "created_at":
			// 	s.CreatedAt = &user.CreatedAt
		}
	}
	return s
}
