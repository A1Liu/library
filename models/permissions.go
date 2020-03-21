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

type void struct{}

type Permission struct {
	PType     uint64
	Reference uint64
	_v        void
}

func TargetedPermission(permissionType uint64, id uint64) *Permission {
	if permissionType != ValidateSingleAuthor && permissionType != ValidateSingleBook {
		log.Fatal("Gave invalid permission type ", permissionType)
	}

	return &Permission{permissionType, id, void{}}
}

func BroadPermission(permissionType uint64) *Permission {
	if !IsValidPermissionType(permissionType) ||
		permissionType == ValidateSingleAuthor ||
		permissionType == ValidateSingleBook {
		log.Fatal("Gave invalid permission type ", permissionType)
	}

	return &Permission{permissionType, 0, void{}}
}

func (perm *Permission) IsTargeted() bool {
	return IsTargeted(perm.PType)
}

func IsTargeted(permissionType uint64) bool {
	return permissionType == ValidateSingleAuthor || permissionType == ValidateSingleBook
}

func IsValidPermissionType(value uint64) bool {
	return value < 6
}

//func PType(permissionType uint64) string {
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

func PermissionContains(l, r uint64) bool {
	if !IsValidPermissionType(l) || IsValidPermissionType(r) {
		log.Fatal("Got invalid permission: (", l, ", ", r, ")")
	}

	if l == ValidateBooks {
		return r == ValidateBooks || r == ValidateSingleBook
	} else if l == ValidateAuthors {
		return r == ValidateAuthors || r == ValidateSingleAuthor
	}

	return l == r
}
