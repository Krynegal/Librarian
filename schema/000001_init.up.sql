CREATE TABLE Bookcase
(
    ID          INT PRIMARY KEY,
    Description VARCHAR(30)
);

CREATE TABLE Section
(
    ID          INT PRIMARY KEY,
    Number      int,
    ID_Bookcase INT REFERENCES Bookcase (ID)
);

CREATE TABLE Shelf
(
    ID         INT PRIMARY KEY,
    Number     int,
    ID_Section INT REFERENCES Section (ID)
);

CREATE TABLE Genre
(
    ID   INT PRIMARY KEY,
    Name VARCHAR(30)
);

CREATE TABLE Book
(
    ID       INT PRIMARY KEY,
    ID_Shelf INT REFERENCES Shelf (ID),
    ID_Genre INT REFERENCES Genre (ID),
    Name     VARCHAR(30)
);

CREATE TABLE Author
(
    ID       INT PRIMARY KEY,
    Lastname VARCHAR(30),
    Initials VARCHAR(30)
);

CREATE TABLE ID_Book_ID_Author
(
    ID_Book   INT REFERENCES Book (ID),
    ID_Author INT REFERENCES Author (ID)
);