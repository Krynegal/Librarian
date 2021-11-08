package postgresDB

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"tg/pkg/storage"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}



func (d *Database) GetBooksByTitle(bookName string) (string, error) {
	rows, err := d.db.Query("SELECT Book.Name, Bookcase.Description, Section.Number, Shelf.Number FROM Book JOIN Shelf ON Shelf.ID = Book.ID_Shelf JOIN Section ON Section.ID = Shelf.ID_Section JOIN Bookcase ON Bookcase.ID = Section.ID_Bookcase WHERE Name = $1", bookName)
	if err != nil {
		fmt.Println(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	return getJSONBooks(rows)
}

func (d *Database) GetBooksByAuthor(authorLastname string) (string, error) {
	rows, err := d.db.Query("SELECT Name, Bookcase.Description, Section.Number, Shelf.Number FROM Author JOIN ID_Book_ID_Author ON ID_Book_ID_Author.ID_Author = Author.ID JOIN Book ON Book.ID = ID_Book_ID_Author.ID_Book JOIN Shelf ON Shelf.ID = Book.ID_Shelf JOIN Section ON Section.ID = Shelf.ID_Section JOIN Bookcase ON Bookcase.ID = Section.ID_Bookcase WHERE Author.Lastname = $1", authorLastname)
	if err != nil {
		fmt.Println(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	return getJSONBooks(rows)
}

func getJSONBooks(rows *sql.Rows) (string, error) {
	var books []storage.BookInfo
	for rows.Next() {
		bi := storage.BookInfo{}
		err := rows.Scan(&bi.Name, &bi.Bookcase, &bi.SectionNumber, &bi.ShelfNumber)
		if err != nil {
			fmt.Println(err)
		}
		books = append(books, bi)
	}
	booksJSON, err := json.Marshal(books)
	if err != nil {
		fmt.Println(err)
	}

	return string(booksJSON), err
}
