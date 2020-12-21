package services

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/mashingan/smapping"
	"github.com/vanilla/gin-crud/api/dto"
	"github.com/vanilla/gin-crud/api/entity"
	"github.com/vanilla/gin-crud/api/repository"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}

func (s *authService) VerifyCredential(email string, password string) interface{} {
	res := s.userRepository.VerifyCredential(email, password)

	if v, ok := res.(entity.User); ok {
		comparedPassword := comparedPassword(v.Password, []byte(password))

		if v.Email == email && comparedPassword {
			return res
		}

		return false
	}

	return false
}

func (s *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}

	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))

	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	res := s.userRepository.InsertUser(userToCreate)

	return res
}

func (s *authService) FindByEmail(email string) entity.User {
	return s.userRepository.FindByEmail(email)
}

func (s *authService) IsDuplicateEmail(email string) bool {
	res := s.userRepository.IsDuplicateEmail(email)

	return !(res.Error == nil)
}

func comparedPassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
