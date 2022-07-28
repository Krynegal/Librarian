package storage

type Storage interface {
	GetBooksByTitle(bookName string) ([]BookInfo, error)
	GetBooksByAuthor(authorLastname string) ([]BookInfo, error)
}
