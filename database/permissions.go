package database

import (
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
)

func AddPermission(giver *models.User, userId uint64,
	permission *models.Permission) (uint64, error) {
	row := psql.Insert("permissions").
		Columns("given_to", "authorized_by", "permission_to", "metadata").
		Values(userId, giver.Id, permission.Type, permission.Ref).
		Suffix("RETURNING \"id\"").
		RunWith(globalDb).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
}

func HasPermissions(user *models.User, permissions []models.Permission) (bool, error) {
	if user.UserGroup == models.AdminUser { // @Responsibility should this be somewhere else?
		return true, nil
	}

	rows, err := psql.Select("permission_to", "metadata").
		From("permissions").
		Where(sq.Eq{"given_to": user.Id}).
		RunWith(globalDb).
		Query()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	permissionMap := make(map[models.Permission]bool, len(permissions))
	for _, p := range permissions {
		permissionMap[p] = true
	}

	var permission models.Permission
	for rows.Next() {
		err := rows.Scan(&permission.Type, &permission.Ref)
		if err != nil {
			return false, err
		}

		for _, p := range permissions {
			if _, ok := permissionMap[p]; ok && models.PermissionContains(permission, p) {
				delete(permissionMap, p)
			}
		}
	}

	return len(permissionMap) == 0, rows.Err()
}

func GetPermissions(user *models.User) ([]models.Permission, error) {
	rows, err := psql.Select("permission_to", "metadata").
		From("permissions").
		Where(sq.Eq{"given_to": user.Id}).
		RunWith(globalDb).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make([]models.Permission, 2)
	var p models.Permission
	for rows.Next() {
		err := rows.Scan(&p.Type, &p.Ref)
		if err != nil {
			return perms, err
		}
	}

	return perms, rows.Err()
}

func RemovePermissions(userId uint64, permission models.Permission) error {
	rows, err := psql.Select("id", "permission_to", "metadata").
		From("permissions").
		Where(sq.Eq{"given_to": userId}).
		RunWith(globalDb).
		Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	ids := make([]uint64, 2)
	var id uint64
	var p models.Permission
	for rows.Next() {
		err := rows.Scan(&id, &p.Type, &p.Ref)
		if err != nil {
			return err
		}

		if models.PermissionContains(permission, p) {
			ids = append(ids, id)
		}
	}

	_, err = psql.Delete("permissions").Where(sq.Eq{"id": ids}).Exec()
	return err
}
