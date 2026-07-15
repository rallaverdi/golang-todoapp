package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string // can be null that's why *string
}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid `FullName` length %d:%w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLength := len([]rune(*u.PhoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf("invalid `PhoneNumber` length %d:%w", phoneNumberLength, core_errors.ErrInvalidArgument)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `PhoneNumber` format:%w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(fullName Nullable[string], phoneNumber Nullable[string]) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("`FullName` cannot be NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u
	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}

//-------------------------------------REDIS CACHE-------------------------------------------------//

type UsersFilter struct {
	PageSize   int  `json:"page_size"`
	PageNumber int  `json:"page_number"`
	UserID     *int `json:"user_id"`
}

func NewUsersFilter(pageSize, pageNumber int, userID *int) UsersFilter {
	if userID == nil {
		return UsersFilter{
			PageSize:   pageSize,
			PageNumber: pageNumber,
			UserID:     nil,
		}
	}

	return UsersFilter{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		UserID:     userID,
	}
}

func (f UsersFilter) Validate() error {
	if f.PageSize < 0 {
		return fmt.Errorf("invalid `PageSize` %d: %w", f.PageSize, core_errors.ErrInvalidArgument)
	}

	if f.PageNumber < 0 {
		return fmt.Errorf("invalid `PageNumber` %d: %w", f.PageNumber, core_errors.ErrInvalidArgument)
	}

	return nil
}
