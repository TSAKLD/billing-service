package service

import (
	"avitoTZ/entity"
	"avitoTZ/repository"
	"errors"
	"time"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(r *repository.Repository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (us *UserService) AddUser(u entity.User) (entity.User, error) {
	if u.Balance != 0 {
		return us.repo.AddUser(u)
	}

	u.Balance = 0
	return us.repo.AddUser(u)
}

func (us *UserService) UserBalance(id int) (entity.User, error) {
	return us.repo.UserBalance(id)
}

func (us *UserService) BalanceInCurrency(id int, c string) (float64, error) {
	return us.repo.BalanceInCurrency(id, c)
}

func (us *UserService) ChangeUserBalance(tr entity.Transaction) error {
	tr.Date = time.Now()

	if *tr.Amount < 0 {
		tr.SystemCommentary = "Списание средств с баланса пользователя"

		user, err := us.repo.UserBalance(*tr.Retriever)
		if err != nil {
			return err
		}

		balance := user.Balance + *tr.Amount

		if balance < 0 {
			err := errors.New("not enough balance")
			return err
		}

		return us.repo.ChangeUserBalance(tr)
	}

	tr.SystemCommentary = "Пополнение баланса пользователя"
	return us.repo.ChangeUserBalance(tr)
}

func (us *UserService) TransferMoney(tr entity.Transaction) error {
	tr.Date = time.Now()
	tr.SystemCommentary = "Перевод между пользователями"

	sender, err := us.repo.UserBalance(*tr.Sender)
	if err != nil {
		return err
	}

	if sender.Balance < *tr.Amount {
		err := errors.New("not enough balance")
		return err
	}

	_, err = us.repo.UserBalance(*tr.Retriever)
	if err != nil {
		return err
	}

	err = us.repo.TransferMoney(tr)
	if err != nil {
		return err
	}

	return nil
}
