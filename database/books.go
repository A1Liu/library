package database

import (
	"database/sql"
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
)

func SelectBooks(db *sql.DB, pageIndex uint64) ([]models.Book, error) {
	books := make([]models.Book, 50)[:0]

	rows, err := psql.Select("*").
		From("books").
		Where(sq.Lt{"id": 50 * (pageIndex + 1)}).
		Where(sq.GtOrEq{"id": 50 * pageIndex}).
		Limit(50).
		RunWith(db).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var book models.Book
	for rows.Next() {
		err := rows.Scan(&book.Id, &book.SuggestedAt, &book.SuggestedBy,
			&book.ValidatedAt, &book.ValidatedBy, &book.Title, &book.Description)
		if err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func InsertBook(db *sql.DB, userId *uint64, title string, description string) (uint64, error) {
	row := psql.Insert("books").
		Columns("suggested_by", "validated_at", "title", "description").
		Values(userId, nil, title, description).
		RunWith(db).
		Suffix("RETURNING \"id\"").
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func InsertValidateBook(db *sql.DB, userId uint64, title string, description string) (uint64, error) {
	row := psql.Insert("books").
		Columns("suggested_by", "validated_by", "title", "description").
		Values(userId, userId, title, description).
		RunWith(db).
		Suffix("RETURNING \"id\"").
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func ValidateBook(db *sql.DB, userId uint64, bookId uint64) error {
	_, err := psql.Update("books").
		Set("validated_by", userId).
		Where(sq.Eq{"id": bookId}).
		Where(sq.Eq{"validated_by": nil}).
		RunWith(db).
		Query()

	return err
}

func GetBook(db *sql.DB, bookId uint64) (*models.Book, error) {
	row := psql.Select("*").
		From("books").
		Where(sq.Eq{"id": bookId}).
		RunWith(db).
		QueryRow()

	var book models.Book
	err := row.Scan(&book.Id, &book.SuggestedAt, &book.SuggestedBy,
		&book.ValidatedAt, &book.ValidatedBy, &book.Title, &book.Description)
	return &book, err
}
