package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/vanilla/gin-crud/api/dto"

	"github.com/gin-gonic/gin"
	"github.com/vanilla/gin-crud/api/entity"
	"github.com/vanilla/gin-crud/api/payload"
	"github.com/vanilla/gin-crud/api/services"
)

type BookController interface {
	GetAll(context *gin.Context)
	InsertBook(context *gin.Context)
	FindByID(context *gin.Context)
	UpdateBook(context *gin.Context)
	DeleteBook(context *gin.Context)
}

type bookController struct {
	bookService services.BookService
	jwtService  services.JWTService
}

func NewBookController(bookService services.BookService, jwtService services.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) GetAll(context *gin.Context) {
	var books []entity.Book = c.bookService.All()

	res := payload.MessageResponse(true, "OK", books)
	context.JSON(http.StatusOK, res)
}

func (c *bookController) InsertBook(context *gin.Context) {
	var bookDTO dto.BookCreateDTO
	err := context.ShouldBind(&bookDTO)

	if err != nil {
		res := payload.ErrorResponse("Failed to process request", err.Error(), payload.EmptyObject{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)

		if err == nil {
			bookDTO.UserID = convertedUserID
		}

		result := c.bookService.Insert(bookDTO)
		response := payload.MessageResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *bookController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := payload.ErrorResponse("No param id was found", err.Error(), payload.EmptyObject{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = c.bookService.FindByID(id)
	if book.ID == 0 {
		res := payload.ErrorResponse("Data not found", "No data with given id", payload.EmptyObject{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		res := payload.MessageResponse(true, "OK", book)
		context.JSON(http.StatusOK, res)
	}
}

func (c *bookController) UpdateBook(context *gin.Context) {
	var bookDTO dto.BookUpdateDTO
	err := context.ShouldBind(&bookDTO)

	if err != nil {
		res := payload.ErrorResponse("Failed to process request", err.Error(), payload.EmptyObject{})
		context.JSON(http.StatusBadRequest, res)
	}

	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, bookDTO.ID) {
		id, err := strconv.ParseUint(userID, 10, 64)

		if err == nil {
			bookDTO.UserID = id
		}

		result := c.bookService.Update(bookDTO)
		res := payload.MessageResponse(true, "OK", result)
		context.JSON(http.StatusOK, res)
	} else {
		res := payload.ErrorResponse("You dont have permission", "You are not the owner", payload.EmptyObject{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) DeleteBook(context *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)

	if err != nil {
		response := payload.ErrorResponse("Failed tou get id", "No param id were found", payload.EmptyObject{})
		context.JSON(http.StatusBadRequest, response)
	}
	book.ID = id
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		res := payload.MessageResponse(true, "Deleted", payload.EmptyObject{})
		context.JSON(http.StatusOK, res)
	} else {
		res := payload.ErrorResponse("You dont have permission", "You are not the owner", payload.EmptyObject{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)

	if err != nil {
		panic(err.Error())
	}

	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
