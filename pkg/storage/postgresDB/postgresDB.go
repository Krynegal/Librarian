package postgresDB

import (
	"database/sql"
	"fmt"
	"github.com/Krynegal/Librarian.git/pkg/storage"
	_ "github.com/lib/pq"
	"strings"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) GetBooksByTitle(bookName string) ([]storage.BookInfo, error) {
	rows, err := d.db.Query(`SELECT lower(Book.Name), Bookcase.Description, Section.Number, Shelf.Number FROM Book 
    	JOIN Shelf ON Shelf.ID = Book.ID_Shelf 
		JOIN Section ON Section.ID = Shelf.ID_Section 
		JOIN Bookcase ON Bookcase.ID = Section.ID_Bookcase 
		WHERE Name = $1`, strings.ToLower(bookName))
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		cerr := rows.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	return sliceOfBooks(rows)
}

func (d *Database) GetBooksByAuthor(authorLastname string) ([]storage.BookInfo, error) {
	rows, err := d.db.Query(`SELECT lower(Name), Bookcase.Description, Section.Number, Shelf.Number FROM Author
    	JOIN ID_Book_ID_Author ON ID_Book_ID_Author.ID_Author = Author.ID
    	JOIN Book ON Book.ID = ID_Book_ID_Author.ID_Book
    	JOIN Shelf ON Shelf.ID = Book.ID_Shelf
    	JOIN Section ON Section.ID = Shelf.ID_Section
    	JOIN Bookcase ON Bookcase.ID = Section.ID_Bookcase
    	WHERE Author.Lastname = $1`, strings.ToLower(authorLastname))
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := rows.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	return sliceOfBooks(rows)
}

func sliceOfBooks(rows *sql.Rows) ([]storage.BookInfo, error) {
	var books []storage.BookInfo
	for rows.Next() {
		bi := storage.BookInfo{}
		err := rows.Scan(&bi.Name, &bi.Bookcase, &bi.SectionNumber, &bi.ShelfNumber)
		if err != nil {
			return nil, err
		}
		books = append(books, bi)
	}
	return books, nil
}
