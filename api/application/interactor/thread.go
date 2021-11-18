package interactor

import (
	"app/api/domain/entity"
	"app/api/domain/service"

	"github.com/pkg/errors"
)

type ThreadInteractor interface {
	Create(name, description string, limitUsers, isPublic int, authorID string) (*entity.Thread, error)
	GetAll() ([]*entity.Thread, error)
	GetByID(id string) (*entity.Thread, error)
	GetOnlyPublic() ([]*entity.Thread, error)
	GetMembersByThreadID(id string) ([]*entity.User, error)
	Update(id, name, description string, limitUsers, isPublic int) (*entity.Thread, error)
	Delete(id string) error
	AddMember(threadID, userID string) error
	RemoveMember(threadID, userID string) error
	ForceToLeave(requestUserID, threadID, leavedUserID string) error
	IsParticipated(string, string) bool
}

type threadInteractor struct {
	threadService   service.ThreadService
	userService     service.UserService
	tagService      service.TagService
	categoryService service.CategoryService
}

func NewThreadInteractor(ts service.ThreadService, us service.UserService, tas service.TagService, cs service.CategoryService) ThreadInteractor {
	return &threadInteractor{
		threadService:   ts,
		userService:     us,
		tagService:      tas,
		categoryService: cs,
	}
}

func (ti *threadInteractor) Create(name, description string, limitUsers, isPublic int, authorID string) (*entity.Thread, error) {
	author, err := ti.userService.GetByUserID(authorID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get author")
	}
	thread, err := ti.threadService.New(name, description, limitUsers, isPublic, author)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new thread")
	}
	thread.Author = author
	err = ti.threadService.AddMember(thread.ID, author.ID, 1) // TODO: めっちゃハードコーディングやん
	if err != nil {
		return nil, errors.Wrap(err, "failed to add member")
	}
	return thread, nil
}

func (ti *threadInteractor) GetAll() ([]*entity.Thread, error) {
	threads, err := ti.threadService.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get threads")
	}
	result := make([]*entity.Thread, 0, len(threads))
	for _, thread := range threads {
		author, _ := ti.userService.GetByID(thread.Author.ID)
		thread.Author = author
		userTags, _ := ti.tagService.GetByUserUUID(author.ID)
		thread.Author.Tags = AddCategoryToTag(userTags, ti.categoryService)
		threadTags, _ := ti.tagService.GetByThreadID(thread.ID)
		thread.Tags = AddCategoryToTag(threadTags, ti.categoryService)
		result = append(result, thread)
	}
	return result, nil
}

func (ti *threadInteractor) GetByID(id string) (*entity.Thread, error) {
	thread, err := ti.threadService.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get thread")
	}
	author, err := ti.userService.GetByID(thread.Author.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get author")
	}
	thread.Author = author
	userTags, err := ti.tagService.GetByUserUUID(author.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tags of user")
	}
	threadTags, err := ti.tagService.GetByThreadID(thread.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tags of thread")
	}
	thread.Author.Tags = AddCategoryToTag(userTags, ti.categoryService)
	thread.Tags = AddCategoryToTag(threadTags, ti.categoryService)
	return thread, nil

}

func (ti *threadInteractor) GetOnlyPublic() ([]*entity.Thread, error) {
	threads, err := ti.threadService.GetOnlyPublic()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get threads")
	}
	result := make([]*entity.Thread, 0, len(threads))
	for _, thread := range threads {
		author, _ := ti.userService.GetByID(thread.Author.ID)
		thread.Author = author
		userTags, _ := ti.tagService.GetByUserUUID(author.ID)
		thread.Author.Tags = AddCategoryToTag(userTags, ti.categoryService)
		threadTags, _ := ti.tagService.GetByThreadID(thread.ID)
		thread.Tags = AddCategoryToTag(threadTags, ti.categoryService)
		result = append(result, thread)
		result = append(result, thread)
	}
	return result, nil
}

func (ti *threadInteractor) GetMembersByThreadID(id string) ([]*entity.User, error) {
	members, err := ti.threadService.GetMembersByThreadID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get members")
	}
	result := make([]*entity.User, 0, len(members))
	for _, member := range members {
		tags, _ := ti.tagService.GetByUserUUID(member.ID)
		member.Tags = AddCategoryToTag(tags, ti.categoryService)
		result = append(result, member)
	}
	return result, nil
}

func (ti *threadInteractor) Update(id, name, description string, limitUsers, isPublic int) (*entity.Thread, error) {
	oldThread, err := ti.threadService.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get thread")
	}
	author, err := ti.userService.GetByID(oldThread.Author.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get author")
	}
	thread, err := ti.threadService.Update(oldThread, name, description, limitUsers, isPublic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update thread")
	}
	thread.Author = author
	userTags, err := ti.tagService.GetByUserUUID(author.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tags of user")
	}
	threadTags, err := ti.tagService.GetByThreadID(thread.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tags of thread")
	}
	thread.Author.Tags = AddCategoryToTag(userTags, ti.categoryService)
	thread.Tags = AddCategoryToTag(threadTags, ti.categoryService)
	return thread, nil
}

func (ti *threadInteractor) Delete(id string) error {
	if err := ti.threadService.Delete(id); err != nil {
		return errors.Wrap(err, "failed to delete thread")
	}
	return nil
}

func (ti *threadInteractor) AddMember(threadID, userID string) error {
	_, err := ti.threadService.GetByID(threadID)
	if err != nil {
		return errors.Wrap(err, "failed to get thread")
	}
	user, err := ti.userService.GetByUserID(userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}
	if err = ti.threadService.AddMember(threadID, user.ID, 0); err != nil { // TODO: めっちゃハードコーディングやん
		return errors.Wrap(err, "failed to add member")
	}
	return nil
}

func (ti *threadInteractor) RemoveMember(threadID, userID string) error {
	user, err := ti.userService.GetByUserID(userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}
	if err := ti.threadService.RemoveMember(threadID, user.ID); err != nil {
		return errors.Wrap(err, "failed to remove member")
	}
	return nil
}

func (ti *threadInteractor) ForceToLeave(requestUserID, threadID, leavedUserID string) error {
	// TODO: implement
	return nil
}

func (ti *threadInteractor) IsParticipated(ThreadID string, UserID string) bool {
	members, err := ti.threadService.GetMembersByThreadID(ThreadID)
	if err != nil {
		return false
	}
	for _, member := range members {
		if member.ID == UserID {
			return true
		}
	}
	return false
}
