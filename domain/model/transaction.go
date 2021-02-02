package model

import (
	"errors"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	TransactionPending string = "pending"
	TransactionCompleted string = "completed"
	TransactionConfirmed string = "confirmed"
	TransactionError string = "error"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transactions []*Transaction
}

type Transaction struct {
	Base `valid:"required"`
	AccountFrom *Account `valid:"-"`
	Amount float64 `json:"amount" valid:"notnull"`
	PixKeyTo *PixKey `valid:"-"`
	Status string `json:"status" valid:"notnull"`
	Description string `json:"description" valid:"notnull"`
	CancelDescription string `json:"cancel_description" valid:"-"`
}

func (transaction *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if transaction.Status != TransactionPending && transaction.Status != TransactionCompleted && transaction.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the source and destination account can not be the same")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, status string, description string) (*Transaction, error) {
	transaction := Transaction {
		AccountFrom: accountFrom,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		Status: TransactionPending,
		Description: description,
		CancelDescription: "",
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (transaction *Transaction) Complete() error {
	transaction.Status = TransactionCompleted
	transaction.UpdatedAt = time.Now()

	return transaction.isValid()
}

func (transaction *Transaction) Confirm() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdatedAt = time.Now()

	return transaction.isValid()
}

func (transaction *Transaction) Cancel(description string) error {
	transaction.Status = TransactionError
	transaction.CancelDescription = description
	transaction.UpdatedAt = time.Now()

	return transaction.isValid()
}