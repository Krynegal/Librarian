package storage

type Storage interface {
	GetBooksByTitle(bookName string) (string, error)
	GetBooksByAuthor(authorLastname string) (string, error)
}