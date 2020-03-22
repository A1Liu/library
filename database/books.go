package database

import (
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
)

func SelectBooks(pageIndex uint64) ([]models.Book, error) {
	books := make([]models.Book, 50)[:0]

	rows, err := psql.Select("*").
		From("books").
		Where(sq.Lt{"id": 50 * (pageIndex + 1)}).
		Where(sq.GtOrEq{"id": 50 * pageIndex}).
		Limit(50).
		RunWith(globalDb).
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

func InsertBook(user *models.User, title string, description string) (uint64, error) {
	row := psql.Insert("books").
		Columns("suggested_by", "validated_at", "title", "description").
		Values(user.NilId(), nil, title, description).
		Suffix("RETURNING \"id\"").
		RunWith(globalDb).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func InsertValidateBook(user *models.User,
	title string, description string) (uint64, error) {
	row := psql.Insert("books").
		Columns("suggested_by", "validated_by", "title", "description").
		Values(user.NilId(), user.NilId(), title, description).
		Suffix("RETURNING \"id\"").
		RunWith(globalDb).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func ValidateBook(user *models.User, bookId uint64) error {
	_, err := psql.Update("books").
		Set("validated_by", user.NilId()).
		Where(sq.Eq{"id": bookId}).
		Where(sq.Eq{"validated_by": nil}).
		RunWith(globalDb).
		Query()

	return err
}

func GetBook(bookId uint64) (*models.Book, error) {
	row := psql.Select("*").
		From("books").
		Where(sq.Eq{"id": bookId}).
		RunWith(globalDb).
		QueryRow()

	var book models.Book
	err := row.Scan(&book.Id, &book.SuggestedAt, &book.SuggestedBy,
		&book.ValidatedAt, &book.ValidatedBy, &book.Title, &book.Description, &book.ImageId)
	return &book, err
}

func MergeBookInto(from uint64, into uint64) error {
	_, err := psql.Update("written_by").
		Set("book_id", into).
		Where(sq.Eq{"book_id": from}).
		RunWith(globalDb).
		Exec()

	if err != nil {
		return err
	}

	_, err = psql.Delete("books").Where(sq.Eq{"id": from}).RunWith(globalDb).Exec()
	return err
}
