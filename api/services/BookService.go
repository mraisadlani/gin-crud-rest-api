package services

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/vanilla/gin-crud/api/dto"
	"github.com/vanilla/gin-crud/api/entity"
	"github.com/vanilla/gin-crud/api/repository"
)

type BookService interface {
	Insert(b dto.BookCreateDTO) entity.Book
	Update(b dto.BookUpdateDTO) entity.Book
	Delete(b entity.Book)
	All() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (s *bookService) Insert(b dto.BookCreateDTO) entity.Book {
	book := entity.Book{}

	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := s.bookRepository.InsertBook(book)
	return res
}

func (s *bookService) Update(b dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}

	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := s.bookRepository.UpdateBook(book)
	return res
}

func (s *bookService) Delete(b entity.Book) {
	s.bookRepository.DeleteBook(b)
}

func (s *bookService) All() []entity.Book {
	return s.bookRepository.AllBook()
}

func (s *bookService) FindByID(bookID uint64) entity.Book {
	return s.bookRepository.FindBookByID(bookID)
}

func (s *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := s.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
