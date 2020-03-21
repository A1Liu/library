package database

import (
	"database/sql"
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
)

func AddPermission(db *sql.DB, giver *models.User, userId uint64,
	permission *models.Permission) (uint64, error) {
	row := psql.Insert("permissions").
		Columns("given_to", "authorized_by", "permission_to", "metadata").
		Values(userId, giver.Id, permission.PType, permission.Reference).
		Suffix("RETURNING \"id\"").
		RunWith(db).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func HasPermissions(db *sql.DB, user *models.User, permissions []models.Permission) (bool, error) {
	if user.UserGroup == models.AdminUser { // @Responsibility should this be somewhere else?
		return true, nil
	}

	rows, err := psql.Select("permission_to").
		From("permissions").
		Where(sq.Eq{"given_to": user.Id}).
		RunWith(db).
		Query()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	permissionMap := make(map[uint64]bool, len(permissions))
	for _, p := range permissions {
		permissionMap[p.PType] = true
	}

	var pOwned uint64
	for rows.Next() {
		err := rows.Scan(&pOwned)
		if err != nil {
			return false, err
		}

		for _, p := range permissions {
			if _, ok := permissionMap[p.PType]; ok && models.PermissionContains(pOwned, p.PType) {
				delete(permissionMap, p.PType)
			}
		}
	}

	return len(permissionMap) == 0, rows.Err()
}

func RemovePermissions(db *sql.DB, userId uint64, permission *models.Permission) error {
	rows, err := psql.Select("id", "permission_to").
		From("permissions").
		Where(sq.Eq{"given_to": userId}).
		RunWith(db).
		Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	ids := make([]uint64, 2)
	var id, permissionTo uint64
	for rows.Next() {
		err := rows.Scan(&id, &permissionTo)
		if err != nil {
			return err
		}

		if models.PermissionContains(permission.PType, permissionTo) {
			ids = append(ids, id)
		}
	}

	_, err = psql.Delete("permissions").Where(sq.Eq{"id": ids}).Exec()
	return err
}
