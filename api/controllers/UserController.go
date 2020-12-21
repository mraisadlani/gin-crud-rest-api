package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vanilla/gin-crud/api/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vanilla/gin-crud/api/payload"
	"github.com/vanilla/gin-crud/api/services"
)

type UserController interface {
	GetProfile(context *gin.Context)
	UpdateProfile(context *gin.Context)
	GetAllProfile(context *gin.Context)
	DeleteProfile(context *gin.Context)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) GetProfile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.userService.GetProfile(id)

	res := payload.MessageResponse(true, "OK", user)
	context.JSON(http.StatusOK, res)
}

func (c *userController) UpdateProfile(context *gin.Context) {
	var userDTO dto.UserDTO
	err := context.ShouldBind(&userDTO)

	if err != nil {
		respose := payload.ErrorResponse("Failed to process request", err.Error(), payload.EmptyObject{})
		context.AbortWithStatusJSON(http.StatusBadRequest, respose)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}

	userDTO.ID = id
	u := c.userService.UpdateProfile(userDTO)
	res := payload.MessageResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}

func (c *userController) GetAllProfile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])

	fmt.Println(id)

	u := c.userService.AllProfile()
	res := payload.MessageResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}

func (c *userController) DeleteProfile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])

	fmt.Println(id)

	ids, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		respose := payload.ErrorResponse("Failed to process request", err.Error(), payload.EmptyObject{})
		context.AbortWithStatusJSON(http.StatusBadRequest, respose)
		return
	}

	u := c.userService.DeleteProfile(ids)
	res := payload.MessageResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}
