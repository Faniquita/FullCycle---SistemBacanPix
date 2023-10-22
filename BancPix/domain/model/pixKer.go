package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// Function(PARAM) (RETURN)
type PixKeyRepositoryInterface interface {
	RegisterKey(pixkey *PixKey) (*PixKey, error)
	FinKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}

type PixKey struct {
	Base      `valid:"required"`
	Kind      string   `json:"kind" valid:"notnull"` // Tipo da chave
	Key       string   `json:"key" valid:"notnull"`
	AccountID string   `json:"account_id" valid:"notnull"`
	Account   *Account `valid:"notnull"`
	Status    string   `json:"status" valid:"notnull"`
}

// Methode: Validação
func (pixKey *PixKey) isValid() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key | tipo de chave invalida")
	}

	if pixKey.Status != "active" && pixKey.Kind != "inactive" {
		return errors.New("invalid status | status invalido")
	}

	if err != nil {
		return err
	}
	return nil
}

// Função | *Bank = Ponteiro
func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey{
		Kind:    kind,
		Key:     key,
		Account: account,
		Status:  "active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.isValid()
	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}
