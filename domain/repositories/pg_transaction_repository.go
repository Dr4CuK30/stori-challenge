package repositories

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"stori-challenge-v1/domain/entities"
	"sync"
)

type PgTransactionRepository struct {
	db *sql.DB
}

func (r *PgTransactionRepository) Save(process entities.Process, wg *sync.WaitGroup) error {
	defer r.db.Close()
	db := r.getConnection()
	var id int
	err := db.QueryRow("INSERT INTO process (origin, origin_name) VALUES ($1, $2) RETURNING id",
		process.Origin, process.OriginName).Scan(&id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, transaction := range process.Transactions {
		fmt.Println("Inserting transaction")
		_, err := db.Exec("INSERT INTO transactions (id, amount, day, month, process_id) VALUES ($1, $2, $3, $4, $5)",
			transaction.Id, transaction.Amount, transaction.Day, transaction.Month, id)
		if err != nil {
			fmt.Println("Error inserting transaction", err)
			return err
		}
	}
	wg.Done()
	return nil
}

func (r *PgTransactionRepository) getConnection() *sql.DB {
	if r.db == nil {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName)
		fmt.Println("Connecting to the database")
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			fmt.Println("Error connecting to the database")
			panic(err)
		}
		fmt.Println("Connected to the database")
		r.db = db
	}
	return r.db
}
