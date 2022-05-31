package repository

import (
	"avitoTZ/bootstrap"
	"avitoTZ/entity"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func New(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func DBConnect(c *bootstrap.Config) (*sql.DB, error) {
	info := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)

	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}

func (r *Repository) AddUser(u entity.User) (entity.User, error) {
	q := "insert into users(name, balance) values($1, $2) returning id, name, balance"

	err := r.db.QueryRow(q, u.Name, u.Balance).Scan(&u.ID, &u.Name, &u.Balance)
	if err != nil {
		return entity.User{}, err
	}

	return u, err
}

func (r *Repository) UserBalance(id int) (entity.User, error) {
	var u entity.User

	q := "select balance from users where id = $1"

	err := r.db.QueryRow(q, id).Scan(&u.Balance)
	if err != nil {
		return entity.User{}, err
	}

	return u, err
}

func (r *Repository) ChangeUserBalance(tr entity.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := "update users set balance=balance+$1 where id=$2"

	_, err = tx.Exec(q, tr.Amount, tr.Retriever)
	if err != nil {
		return err
	}

	err = r.UserBalanceTransaction(tr)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TransferMoney(tr entity.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := "update users set balance = balance - $1 where id = $2"
	_, err = tx.Exec(q, tr.Amount, tr.Sender)
	if err != nil {
		return err
	}

	q = "update users set balance = balance + $1 where id = $2"
	_, err = tx.Exec(q, tr.Amount, tr.Retriever)
	if err != nil {
		return err
	}

	err = r.TransferTransaction(tr)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (r *Repository) UserBalanceTransaction(tr entity.Transaction) error {
	q := "insert into transactions(user_id, system_commentary, user_commentary, amount, date)" +
		" values($1, $2, $3, $4, $5)"

	_, err := r.db.Exec(q, tr.Retriever, tr.SystemCommentary, tr.UserCommentary, tr.Amount, tr.Date)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TransferTransaction(tr entity.Transaction) error {
	q := "insert into transactions(user_id, sender_id, system_commentary, user_commentary, amount, date)" +
		" values($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(q, tr.Retriever, tr.Sender, tr.SystemCommentary, tr.UserCommentary, tr.Amount, tr.Date)
	if err != nil {
		return err
	}

	q = "insert into transactions(user_id, retriever_id, system_commentary, user_commentary, amount, date)" +
		" values($1, $2, $3, $4, $5, $6)"

	_, err = r.db.Exec(q, tr.Sender, tr.Retriever, tr.SystemCommentary, tr.UserCommentary, *tr.Amount*-1, tr.Date)
	if err != nil {
		return err
	}

	return nil
}
