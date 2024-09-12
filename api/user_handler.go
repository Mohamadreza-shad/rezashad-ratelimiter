package api

import (
	"fmt"
	"net/http"
	"web-metric/service/ratelimiter"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	ratelimiterService *ratelimiter.Service
	validate           *validator.Validate
}

func (h *UserHandler) SetUserManualRateLimit(c *gin.Context) {
	var params ratelimiter.SetUserConfigParams
	err := c.BindJSON(&params)
	if err != nil {
		fmt.Println(err.Error())
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Bad Request: "+err.Error(),
		)
		return
	}
	err = h.validate.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Invalid Request: "+err.Error())
		return
	}
	err = h.ratelimiterService.SetUserConfig(c.Request.Context(), params)
	if err != nil {
		MakeErrorResponseWithoutCode(c.Writer, err)
		return
	}
	MakeSuccessResponse(c.Writer, nil, "user rate limit set successfully")
}

func (h *UserHandler) AllProducts(c *gin.Context) {
	products := h.ratelimiterService.AllProducts(c.Request.Context())
	MakeSuccessResponse(c.Writer, products, "products have been fetched successfull")
}

func NewUserHandler(
	ratelimiterService *ratelimiter.Service,
	validate *validator.Validate,
) *UserHandler {
	return &UserHandler{
		ratelimiterService: ratelimiterService,
		validate:           validate,
	}
}
