package model

import (
	"errors"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

type PixKeyRepositoryInterface interface {
	RegisterKey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(kind string, key string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(accountID string) (*Account, error)
}

type PixKey struct {
	Base `valid:"required"`
	Kind string `json:"kind" valid:"notnull"`
	Key string `json:"key" valid:"notnull"`
	AccountID string `json:"account_id" valid:"notnull"`
	Account *Account `valid:"-"`
	Status string `json:"status" valid:"notnull"`
}

func (pixKey *PixKey) isValid() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key (must be email or cpf)")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid status")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewPixKey(kind string, key string, account *Account) (*PixKey, error) {
	pixKey := PixKey{
		Kind: kind,
		Key: key,
		Account: account,
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.isValid()

	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}
