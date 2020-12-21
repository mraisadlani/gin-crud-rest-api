package services

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/vanilla/gin-crud/api/dto"
	"github.com/vanilla/gin-crud/api/entity"
	"github.com/vanilla/gin-crud/api/repository"
)

type UserService interface {
	GetProfile(userID string) entity.User
	UpdateProfile(user dto.UserDTO) entity.User
	AllProfile() []entity.User
	DeleteProfile(ids uint64) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (s *userService) UpdateProfile(user dto.UserDTO) entity.User {
	userToUpdate := entity.User{}

	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))

	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	updatedUser := s.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (s *userService) GetProfile(userID string) entity.User {
	return s.userRepository.ProfileUser(userID)
}

func (s *userService) AllProfile() []entity.User {
	return s.userRepository.AllProfileUser()
}

func (s *userService) DeleteProfile(ids uint64) entity.User {
	return s.userRepository.DeleteUser(ids)
}
