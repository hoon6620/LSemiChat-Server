package repository

import (
	"app/api/domain/entity"
	"app/api/domain/repository"
	"app/api/infrastructure/database"

	"github.com/pkg/errors"
)

type threadRepository struct {
	sqlHandler database.SQLHandler
}

func NewThreadRepository(sh database.SQLHandler) repository.ThreadRepository {
	return &threadRepository{
		sqlHandler: sh,
	}
}

func (tr *threadRepository) Create(thread *entity.Thread) error {
	_, err := tr.sqlHandler.Exec(`
		INSERT INTO threads(id, name, description, limit_users, user_id, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		thread.ID,
		thread.Name,
		thread.Description,
		thread.LimitUsers,
		thread.Author.ID,
		thread.IsPublic,
		thread.CreatedAt,
		thread.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert db")
	}
	return nil
}

func (tr *threadRepository) FindAll() ([]*entity.Thread, error) {
	rows, err := tr.sqlHandler.Query(`
		SELECT id, name, description, limit_users, user_id, is_public, created_at, updated_at
		FROM threads
	`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select")
	}
	var threads []*entity.Thread
	for rows.Next() {
		var thread entity.Thread
		var author entity.User
		if err = rows.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.LimitUsers, &author.ID, &thread.IsPublic, &thread.CreatedAt, &thread.UpdatedAt); err != nil {
			if rows.CheckNoRows(err) {
				return nil, nil
			}
			return nil, errors.Wrap(err, "failed to scan")
		}
		thread.Author = &author
		threads = append(threads, &thread)
	}
	return threads, nil
}

func (tr *threadRepository) FindByID(id string) (*entity.Thread, error) {
	row := tr.sqlHandler.QueryRow(`
		SELECT id, name, description, limit_users, user_id, is_public, created_at, updated_at
		FROM threads
		WHERE id=?
	`, id)
	var thread entity.Thread
	var author entity.User
	if err := row.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.LimitUsers, &author.ID, &thread.IsPublic, &thread.CreatedAt, &thread.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "failed to scan")
	}
	thread.Author = &author
	return &thread, nil
}

func (tr *threadRepository) FindByUserID(userID string) ([]*entity.Thread, error) {
	rows, err := tr.sqlHandler.Query(`
		SELECT t.id, t.name, t.description, t.limit_users, t.user_id, t.is_public, t.created_at, t.updated_at
		FROM threads AS t
		JOIN users_threads AS ut
		ON t.id = ut.thread_id
		WHERE ut.user_id=?
	`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select")
	}

	var threads []*entity.Thread
	for rows.Next() {
		var thread entity.Thread
		var author entity.User
		if err = rows.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.LimitUsers, &author.ID, &thread.IsPublic, &thread.CreatedAt, &thread.UpdatedAt); err != nil {
			if rows.CheckNoRows(err) {
				return nil, nil
			}
			return nil, errors.Wrap(err, "failed to scan")
		}
		thread.Author = &author
		threads = append(threads, &thread)
	}
	return threads, nil
}

func (tr *threadRepository) FindOnlyPublic() ([]*entity.Thread, error) {
	rows, err := tr.sqlHandler.Query(`
		SELECT id, name, description, limit_users, user_id, is_public, created_at, updated_at
		FROM threads
		WHERE is_public=1
	`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select")
	}
	var threads []*entity.Thread
	for rows.Next() {
		var thread entity.Thread
		var author entity.User
		if err = rows.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.LimitUsers, &author.ID, &thread.IsPublic, &thread.CreatedAt, &thread.UpdatedAt); err != nil {
			if rows.CheckNoRows(err) {
				return nil, nil
			}
			return nil, errors.Wrap(err, "failed to scan")
		}
		thread.Author = &author
		threads = append(threads, &thread)
	}
	return threads, nil
}

func (tr *threadRepository) FindMembersByThreadID(id string) ([]*entity.User, error) {
	rows, err := tr.sqlHandler.Query(`
		SELECT u.id, u.user_id, u.name, u.image, u.profile, u.mail, u.login_at, u.created_at, u.updated_at
		FROM users_threads AS r
		INNER JOIN users AS u
		ON u.id=r.user_id
		WHERE r.thread_id=? 
	`, id)
	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(&user.ID, &user.UserID, &user.Name, &user.Image, &user.Profile, &user.Mail, &user.LoginAt, &user.CreatedAt, &user.UpdatedAt); err != nil {
			if rows.CheckNoRows(err) {
				return nil, nil
			}
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (tr *threadRepository) Update(thread *entity.Thread) error {
	_, err := tr.sqlHandler.Exec(`
		UPDATE threads
		SET name=?, description=?, limit_users=?, is_public=?, updated_at=?
		WHERE id=?
	`,
		thread.Name,
		thread.Description,
		thread.LimitUsers,
		thread.IsPublic,
		thread.UpdatedAt,
		thread.ID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update db")
	}
	return nil
}

func (tr *threadRepository) AddMember(id, threadID, userID string, isAdmin int) error {
	_, err := tr.sqlHandler.Exec(`
		INSERT INTO users_threads(id, user_id, thread_id, is_admin)
		VALUES (?, ?, ?, ?)
	`,
		id,
		userID,
		threadID,
		isAdmin,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert relation")
	}
	return nil
}

func (tr *threadRepository) RemoveMember(threadID, userID string) error {
	_, err := tr.sqlHandler.Exec(`
		DELETE FROM users_threads
		WHERE user_id=? and thread_id=?
	`, userID, threadID)
	if err != nil {
		return errors.Wrap(err, "failed to delete relation")
	}
	return nil
}

func (tr *threadRepository) Delete(id string) error {
	// NOTE: users_threadsのrelationの全切りしてる。nullとかのがいくね...?
	_, err := tr.sqlHandler.Exec(`
		DELETE FROM users_threads
		WHERE thread_id=?
	`, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete relations")
	}
	_, err = tr.sqlHandler.Exec(`
		DELETE FROM threads
		WHERE id=?
	`, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	return nil
}
