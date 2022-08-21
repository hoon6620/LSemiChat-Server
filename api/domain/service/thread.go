package service

import (
	"app/api/domain/entity"
	"app/api/domain/repository"
	"time"

	"github.com/pkg/errors"
)

type ThreadService interface {
	New(name, description string, limitUsers, isPublic int, author *entity.User) (*entity.Thread, error)
	GetAll() ([]*entity.Thread, error)
	GetByID(id string) (*entity.Thread, error)
	GetByUserID(userID string) ([]*entity.Thread, error)
	GetOnlyPublic() ([]*entity.Thread, error)
	GetMembersByThreadID(id string) ([]*entity.User, error)
	Update(thread *entity.Thread, name, description string, limitUsers, isPublic int) (*entity.Thread, error)
	Delete(id string) error
	AddMember(threadID, userID string, isAdmin int) error
	RemoveMember(threadID, userID string) error
}

type threadService struct {
	threadRepository repository.ThreadRepository
	fileRepository   repository.FileRepository
}

func NewThreadService(tr repository.ThreadRepository, fr repository.FileRepository) ThreadService {
	return &threadService{
		threadRepository: tr,
		fileRepository:   fr,
	}
}

func (ts *threadService) New(name, description string, limitUsers, isPublic int, author *entity.User) (*entity.Thread, error) {
	id, err := GenerateUUID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate uuid")
	}
	now := time.Now()
	thread := &entity.Thread{
		ID:          id,
		Name:        name,
		Description: description,
		LimitUsers:  limitUsers,
		IsPublic:    isPublic,
		Author:      author,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	if err = ts.threadRepository.Create(thread); err != nil {
		return nil, errors.Wrap(err, "failed to create thread")
	}
	ts.fileRepository.CreateThreadDir(id)
	return thread, nil
}

func (ts *threadService) GetAll() ([]*entity.Thread, error) {
	threads, err := ts.threadRepository.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get threads")
	}
	return threads, nil
}

func (ts *threadService) GetByID(id string) (*entity.Thread, error) {
	thread, err := ts.threadRepository.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get thread")
	}
	return thread, nil
}

func (ts *threadService) GetByUserID(userID string) ([]*entity.Thread, error) {
	threads, err := ts.threadRepository.FindByUserID(userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get threads")
	}
	return threads, nil
}

func (ts *threadService) GetOnlyPublic() ([]*entity.Thread, error) {
	threads, err := ts.threadRepository.FindOnlyPublic()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get threads")
	}
	return threads, nil
}

func (ts *threadService) GetMembersByThreadID(id string) ([]*entity.User, error) {
	members, err := ts.threadRepository.FindMembersByThreadID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get members")
	}
	return members, nil
}

func (ts *threadService) Update(thread *entity.Thread, name, description string, limitUsers, isPublic int) (*entity.Thread, error) {
	now := time.Now()
	thread.UpdatedAt = &now
	thread.Name = name
	thread.Description = description
	thread.LimitUsers = limitUsers
	thread.IsPublic = isPublic
	err := ts.threadRepository.Update(thread)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update thread")
	}
	return thread, nil
}

func (ts *threadService) Delete(id string) error {
	if err := ts.threadRepository.Delete(id); err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	return nil
}

func (ts *threadService) AddMember(threadID, userID string, isAdmin int) error {
	id, err := GenerateUUID()
	if err != nil {
		return errors.Wrap(err, "failed to generate uuid")
	}
	if err = ts.threadRepository.AddMember(id, threadID, userID, isAdmin); err != nil {
		return errors.Wrap(err, "failed to add member")
	}
	return nil
}

func (ts *threadService) RemoveMember(threadID, userID string) error {
	if err := ts.threadRepository.RemoveMember(threadID, userID); err != nil {
		return errors.Wrap(err, "failed to remove member")
	}
	return nil
}
