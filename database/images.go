package database

import (
	sq "github.com/Masterminds/squirrel"
)

func InsertImage(image []byte, extension string) (uint64, error) {
	row := psql.Insert("images").
		Columns("extension", "data").
		Values(extension, image).
		Suffix("RETURNING \"id\"").
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetImage(id uint64) ([]byte, string, error) {
	row := psql.Select("data", "extension").
		From("images").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING \"id\"").
		QueryRow()

	var data []byte
	var extension string
	err := row.Scan(&data, &extension)
	if err != nil {
		return nil, "", err
	}
	return data, extension, nil
}
