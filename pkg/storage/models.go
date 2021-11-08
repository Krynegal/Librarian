package storage

type BookInfo struct {
	Name          string `json:"name"`
	Bookcase      string `json:"bookcase"`
	SectionNumber int    `json:"section_number"`
	ShelfNumber   int    `json:"shelf_number"`
}
