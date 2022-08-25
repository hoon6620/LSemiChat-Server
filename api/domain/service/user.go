package service

import (
	"app/api/domain/entity"
	"app/api/domain/repository"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type userService struct {
	userRepository repository.UserRepository
	fileRepository repository.FileRepository
}

type UserService interface {
	New(userID, name, mail, image, profile, password string, isAdmin int) (*entity.User, error)
	UpdateProfile(user *entity.User, name, mail, image, profile string) (*entity.User, error)
	UpdateUserID(user *entity.User, userID string) (*entity.User, error)
	UpdatePassword(user *entity.User, password string) (*entity.User, error)
	GetByID(id string) (*entity.User, error)
	GetByUserID(userID string) (*entity.User, error)
	GetByMail(mail string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Delete(id string) error
	GetFollows(id string) ([]*entity.User, error)
	AddFollow(userID, followedUserID string) error
	DeleteFollow(userID, followedUserID string) error
	GetFollowers(id string) ([]*entity.User, error)
}

func NewUserService(ur repository.UserRepository, fr repository.FileRepository) UserService {
	return &userService{
		userRepository: ur,
		fileRepository: fr,
	}
}

func (us *userService) New(userID, name, mail, image, profile, password string, isAdmin int) (*entity.User, error) {
	id, err := GenerateUUID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := &entity.User{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Mail:      mail,
		Image:     image,
		Profile:   profile,
		IsAdmin:   isAdmin,
		LoginAt:   &now,
		CreatedAt: &now,
		UpdatedAt: &now,
		Password:  password,
	}
	err = us.userRepository.Create(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert db")
	}
	err = us.fileRepository.CreateUserDir(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dir")
	}
	return user, nil
}

func (us *userService) UpdateProfile(user *entity.User, name, mail, image, profile string) (*entity.User, error) {
	now := time.Now()
	user.Name = name
	user.Mail = mail
	user.Image = image
	user.Profile = profile
	user.UpdatedAt = &now

	err := us.userRepository.UpdateProfile(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update db")
	}
	return user, nil
}

func (us *userService) UpdateUserID(user *entity.User, userID string) (*entity.User, error) {
	now := time.Now()
	user.UserID = userID
	user.UpdatedAt = &now

	err := us.userRepository.UpdateUserID(user.ID, user.UserID, &now)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update db")
	}
	return user, nil
}

func (us *userService) UpdatePassword(user *entity.User, password string) (*entity.User, error) {
	now := time.Now()
	user.Password = password
	user.UpdatedAt = &now

	err := us.userRepository.UpdateUserID(user.ID, user.Password, &now)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update db")
	}
	return user, nil
}

func (us *userService) GetByID(id string) (*entity.User, error) {
	user, err := us.userRepository.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return user, nil
}

func (us *userService) GetByUserID(userID string) (*entity.User, error) {
	user, err := us.userRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return user, nil
}

func (us *userService) GetByMail(mail string) (*entity.User, error) {
	user, err := us.userRepository.FindByMail(mail)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return user, nil
}

func (us *userService) GetAll() ([]*entity.User, error) {
	users, err := us.userRepository.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}
	return users, nil
}

func (us *userService) Delete(id string) error {
	err := us.userRepository.DeleteByID(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	return nil
}

func (us *userService) GetFollows(id string) ([]*entity.User, error) {
	users, err := us.userRepository.FindFollows(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get follows")
	}
	return users, nil
}

func (us *userService) AddFollow(userID, followedUserID string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "failed to generate uuid")
	}
	err = us.userRepository.AddFollow(id.String(), userID, followedUserID)
	if err != nil {
		return errors.Wrap(err, "failed to add follow")
	}
	return nil
}

func (us *userService) DeleteFollow(userID, followedUserID string) error {
	err := us.userRepository.DeleteFollow(userID, followedUserID)
	if err != nil {
		return errors.Wrap(err, "failed to delete follow")
	}
	return nil
}

func (us *userService) GetFollowers(id string) ([]*entity.User, error) {
	users, err := us.userRepository.FindFollowers(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get followers")
	}
	return users, nil
}
