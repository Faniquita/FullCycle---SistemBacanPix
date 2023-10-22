package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPeding    string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

// Function(PARAM) (RETURN)
type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

// Lista de transações
type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json:"amount" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-" `
	Status            string   `json:"status" valid:"notnull"`
	Description       string   `json:"description" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" valid:"notnull"`
}

// Methode: Validação
func (transaction *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("the amount must be greater than 0 | valor não pode ser menos que 0")
	}

	if transaction.Status != TransactionPeding && transaction.Status != TransactionCompleted && transaction.Status != TransactionError {
		return errors.New("invalid status for the transaction | status da transação invalido")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same | a conta de destinatario não podem ser iguais")
	}

	if err != nil {
		return err
	}
	return nil
}

// Função | *Bank = Ponteiro
func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPeding,
		Description: description,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// Completar a transação
func (transaction *Transaction) Complete() error {
	transaction.Status = TransactionCompleted
	transaction.UpdateAt = time.Now()
	err := transaction.isValid()
	return err
}

// Cancelando a transação
func (transaction *Transaction) Cancel(description string) error {
	transaction.Status = TransactionError
	transaction.UpdateAt = time.Now()
	transaction.Description = description
	err := transaction.isValid()
	return err
}

// Confirmando a transação
func (transaction *Transaction) Confirm() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdateAt = time.Now()
	err := transaction.isValid()
	return err
}
