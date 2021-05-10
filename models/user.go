package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(db *sqlx.DB) *User {
	user := &User{}
	user.db = db
	user.table = "users"
	user.hasID = true
	return user
}

type UserRow struct {
	ID                int64  `db:"id"`
	Email             string `db:"email"`
	Phone             string `db:"phone"`
	Password          string `db:"password"`
	Roles             string `db:"Roles"`
	Admin             bool
	Dsc               bool
	Investor          bool
	BlogReader        bool
	InvestorRelations bool
}

// type RoleType int

// const (
//     Admin RoleType = iota
//     Dsc
//     Investor
//     BlogReader
// )

// func (r RoleType) String() string {
//     return [...]string{"Admin", "Dsc", "Investor", "BlogReader"}[r]
// }

type User struct {
	Base
}

func (u *User) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*UserRow, error) {
	userId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u.GetById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllUsers(tx *sqlx.Tx) ([]*UserRow, error) {
	users := []*UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v", u.table)
	err := u.db.Select(&users, query)
	return users, err
}

// GetById returns record by id.
func (u *User) GetById(tx *sqlx.Tx, id int64) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", u.table)
	err := u.db.Get(user, query, id)
	u.UpdateRoles(user)
	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT *  FROM %v WHERE email=?", u.table)
	err := u.db.Get(user, query, email)
	u.UpdateRoles(user)
	return user, err
}

// GetByEmail returns record by email but checks password first.
func (u *User) GetUserByEmailAndPassword(tx *sqlx.Tx, email, password string) (*UserRow, error) {
	user, err := u.GetByEmail(tx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	u.UpdateRoles(user)
	return user, err
}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx, email, password, passwordAgain string, phone string) (*UserRow, error) {
	if email == "" {
		return nil, errors.New("Email cannot be blank.")
	}
	if phone == "" {
		return nil, errors.New("Phone is invalid.")
	}
	if password == "" {
		return nil, errors.New("Password cannot be blank.")
	}
	if password != passwordAgain {
		return nil, errors.New("Password is invalid.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["email"] = email
	data["password"] = hashedPassword
	data["phone"] = phone
	data["Roles"] = "BlogReader"

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordById(tx *sqlx.Tx, userId int64, email, password, passwordAgain string, phone string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data["email"] = email
	}
	if phone != "" {
		data["phone"] = phone
	}

	if password != "" && passwordAgain != "" && password == passwordAgain {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return nil, err
		}

		data["password"] = hashedPassword
	}

	if len(data) > 0 {
		_, err := u.UpdateByID(tx, data, userId)
		if err != nil {
			return nil, err
		}
	}

	return u.GetById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllEmails(tx *sqlx.Tx) ([]string, error) {
	emails := []string{}
	query := fmt.Sprintf("SELECT email FROM %v", u.table)
	err := u.db.Select(&emails, query)

	return emails, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *User) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
func (i *User) UpdateRoles(u *UserRow) {
	switch roles := u.Roles; roles {
	case "InvestorRelations,Admin,Dsc,Investor,BlogReader":
		u.InvestorRelations = true
	case "Admin,Dsc,Investor,BlogReader":
		u.Admin = true
	case "Dsc,Investor,BlogReader":
		u.Dsc = true
	case "Investor,BlogReader":
		u.Investor = true
	default:
		u.BlogReader = true
	}
}
