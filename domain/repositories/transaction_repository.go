package repositories

import (
	"stori-challenge-v1/domain/entities"
	"sync"
)

type TransactionRepository interface {
	Save(transaction entities.Process, wg *sync.WaitGroup) error
}
