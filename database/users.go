package database

import (
	"errors"

	"webserver/types"
)

func (db *Database) GetUserByID(id float64) (*types.User, error) {
	var user types.User
	trx := db.Instance.Where("ID=?", id)

	err := trx.First(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func (db *Database) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	trx := db.Instance.Where("username=?", username)

	err := trx.First(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func (db *Database) UpdatePassword(user *types.User, newPassword []byte) error {
	trx := db.Instance.Model(&user).Update("password", newPassword)
	if trx.Error != nil {
		return trx.Error
	}

	return nil
}
