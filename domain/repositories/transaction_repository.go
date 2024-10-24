package repositories

import "stori-challenge-v1/domain/entities"

type TransactionRepository interface {
	Save(transaction entities.Process) error
}
