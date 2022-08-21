package repository

type FileRepository interface {
	CreateFile(threadID string, fileName string, file []byte) error
	GetFile(threadID string, fileName string) ([]byte, error)
	SetUserIcon(threadID string, file []byte) error
	GetUserIcon(threadID string) ([]byte, error)
	SetThreadIcon(threadID string, file []byte) error
	GetThreadIcon(threadID string) ([]byte, error)
	CreateUserDir(userID string) error
	CreateThreadDir(threadID string) error
}
