package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vanilla/gin-crud/api/dto"
	"github.com/vanilla/gin-crud/api/entity"
	"github.com/vanilla/gin-crud/api/payload"
	"github.com/vanilla/gin-crud/api/services"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	err := ctx.ShouldBind(&loginDTO)

	if err != nil {
		response := payload.ErrorResponse("Failed to process request", err.Error(), payload.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)

	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := payload.MessageResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := payload.ErrorResponse("Please check again your credential", "Invalid Credential", payload.EmptyObject{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	err := ctx.ShouldBind(&registerDTO)

	if err != nil {
		respose := payload.ErrorResponse("Failed to process request", err.Error(), payload.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respose)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := payload.ErrorResponse("Failed to process request", "Duplicate email", payload.EmptyObject{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createUser.ID, 10))
		createUser.Token = token
		response := payload.MessageResponse(true, "OK!", createUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
