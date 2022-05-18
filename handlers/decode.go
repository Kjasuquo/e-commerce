package handlers

import (
	servererrors "github.com/decadevs/shoparena/handlers/serverErrors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// decode decodes the body of c into v
func (h *Handler) Decode(c *gin.Context, v interface{}) []string {
	if err := c.ShouldBindJSON(v); err != nil {
		var errs []string
		verr, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fieldErr := range verr {
				errs = append(errs, servererrors.NewFieldError(fieldErr).String())
			}
		} else {
			errs = append(errs, "internal server error")
		}
		return errs
	}
	return nil
}