package database

import (
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
	"time"
)

func SelectAuthors(pageIndex uint64) ([]models.Author, error) {
	authors := make([]models.Author, 50)[:0]

	rows, err := psql.Select("*").
		From("authors").
		Where(sq.Lt{"id": 50 * (pageIndex + 1)}).
		Where(sq.GtOrEq{"id": 50 * pageIndex}).
		Limit(50).
		RunWith(globalDb).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var author models.Author
	for rows.Next() {
		err := rows.Scan(&author.Id, &author.SuggestedAt, &author.SuggestedBy,
			&author.ValidatedAt, &author.ValidatedBy, &author.FirstName, &author.LastName)
		if err != nil {
			return authors, err
		}

		authors = append(authors, author)
	}

	return authors, rows.Err()
}

func InsertAuthor(user *models.User,
	firstName string, lastName string) (uint64, error) {

	firstName, lastName = firstName[:16], lastName[:16]

	row := psql.Insert("authors").
		Columns("suggested_by", "validated_at", "first_name", "last_name").
		Values(user.NilId(), nil, firstName, lastName).
		Suffix("RETURNING \"id\"").
		RunWith(globalDb).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func ValidateAuthor(user *models.User, authorId uint64) error {
	_, err := psql.Update("authors").
		Set("validated_at", time.Now()).
		Set("validated_by", user.Id).
		Where(sq.Eq{"id": authorId}).
		RunWith(globalDb).
		Exec()

	return err
}

func InsertWrittenBy(user *models.User,
	authorId, bookId uint64) (uint64, error) {

	row := psql.Insert("written_by").
		Columns("suggested_by", "validated_at", "author_id", "book_id").
		Values(user.NilId(), nil, authorId, bookId).
		Suffix("RETURNING \"id\"").
		RunWith(globalDb).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func ValidateWrittenBy(user *models.User, authorId, bookId uint64) error {
	_, err := psql.Update("written_by").
		Set("validated_at", time.Now()).
		Set("validated_by", user.Id).
		Where(sq.Eq{"author_id": authorId, "book_id": bookId}).
		RunWith(globalDb).
		Exec()

	return err
}

func GetAuthor(authorId uint64) (*models.Author, error) {
	row := psql.Select("*").
		From("authors").
		Where(sq.Eq{"id": authorId}).
		RunWith(globalDb).
		QueryRow()

	var author models.Author
	err := row.Scan(&author.Id, &author.SuggestedAt, &author.SuggestedBy,
		&author.ValidatedAt, &author.ValidatedBy, &author.FirstName, &author.LastName)
	return &author, err
}

func MergeAuthorInto(from uint64, into uint64) error {
	_, err := psql.Update("written_by").
		Set("author_id", into).
		Where(sq.Eq{"author_id": from}).
		RunWith(globalDb).
		Exec()

	if err != nil {
		return err
	}

	_, err = psql.Delete("authors").Where(sq.Eq{"id": from}).RunWith(globalDb).Exec()
	return err
}
