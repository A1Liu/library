package models

import (
	"errors"
	"fmt"
	"log"
)

const (
	ValidateBooks        uint64 = 0
	ValidateAuthors      uint64 = 1
	ElevateUsers         uint64 = 2
	DemoteUsers          uint64 = 3
	ValidateSingleAuthor uint64 = 4
	ValidateSingleBook   uint64 = 5
)

type Permission struct {
	Type uint64
	Ref  uint64
}

func TargetedPermission(permissionType uint64, id uint64) *Permission {
	if !IsTargeted(permissionType) {
		log.Fatal("Gave invalid permission type ", permissionType)
	}

	return &Permission{permissionType, id}
}

func BroadPermission(permissionType uint64) *Permission {
	if !IsValidPermissionType(permissionType) || IsTargeted(permissionType) {
		log.Fatal("Gave invalid permission type ", permissionType)
	}

	return &Permission{permissionType, 0}
}

func IsTargeted(permissionType uint64) bool {
	return permissionType > 3 && permissionType < 6
}

func IsValidPermissionType(value uint64) bool {
	return value < 6
}

//func Type(permissionType uint64) string {
//	switch permissionType {
//	case ValidateBooks:
//		return "ValidateBooks"
//	case ValidateAuthors:
//		return "ValidateAuthors"
//	case ElevateUsers:
//		return "ElevateUsers"
//	case DemoteUsers:
//		return "DemoteUsers"
//	case ValidateSingleAuthor:
//		return "ValidateSingleAuthor"
//	case ValidateSingleBook:
//		return "ValidateSingleBook"
//	}
//
//	log.Fatal("This should never happen")
//	return ""
//}

func BuildPermission(permissionType string, reference uint64) (*Permission, error) {
	switch permissionType {
	case "ValidateBooks":
		return BroadPermission(ValidateBooks), nil
	case "ValidateAuthors":
		return BroadPermission(ValidateAuthors), nil
	case "ElevateUsers":
		return BroadPermission(ElevateUsers), nil
	case "DemoteUsers":
		return BroadPermission(DemoteUsers), nil
	case "ValidateSingleAuthor":
		return TargetedPermission(ValidateSingleAuthor, reference), nil
	case "ValidateSingleBook":
		return TargetedPermission(ValidateSingleBook, reference), nil
	}

	return nil, errors.New(fmt.Sprintf("permissionType %s is invalid", permissionType))
}

func PermissionContains(l, r Permission) bool {
	if l.Type == ValidateBooks {
		return r.Type == ValidateBooks || r.Type == ValidateSingleBook
	} else if l.Type == ValidateAuthors {
		return r.Type == ValidateAuthors || r.Type == ValidateSingleAuthor
	}

	return l == r
}
