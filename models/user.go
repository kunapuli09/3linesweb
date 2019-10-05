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
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
	Admin    bool   `db:"admin"`
	Dsc      bool
	FundOne  bool
	FundTwo  bool
}

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
	user.Dsc = isDsc(user.Email)
	user.FundOne = isFundI(user.Email)
	user.FundTwo = isFundII(user.Email)
	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT *  FROM %v WHERE email=?", u.table)
	err := u.db.Get(user, query, email)
	user.Dsc = isDsc(user.Email)
	user.FundOne = isFundI(user.Email)
	user.FundTwo = isFundII(user.Email)
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
	user.Dsc = isDsc(user.Email)
	user.FundOne = isFundI(user.Email)
	user.FundTwo = isFundII(user.Email)
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
	data["admin"] = 0

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
func isDsc(a string) bool {
	//hardcode roles temporarily
	dsc_team := []string{
		"rmaddhi@gmail.com",
		"roy.rajiv@gmail.com",
		"arun.taman@gmail.com",
		"sgosala99@gmail.com",
		"dsc@3lines.vc",
	}
	for _, b := range dsc_team {
		if b == a {
			return true
		}
	}
	return false
}
func isFundI(a string) bool {
	//hardcode roles temporarily
	fundone_insvestors := []string{
		"naga_mulukutla@yahoo.com",
		"vamseekc@yahoo.com",
		"igsvenkat@gmail.com",
		"karun15@gmail.com",
		"bmallikarjun@hotmail.com",
		"maddali.srinivas@gmail.com",
		"mudigondag@yahoo.com",
		"skondam@gmail.com",
		"kiran_misc@yahoo.com",
		"kaladhara@gmail.com",
		"arun.taman@gmail.com",
		"dileep.kasam@gmail.com",
		"rajesh.gundu30@gmail.com",
		"prasadds@hotmail.com",
		"ashwinakurian@gmail.com",
		"lganeshbabu@gmail.com",
		"vamseea@yahoo.com",
		"venkatesh.pallipadi@gmail.com",
		"rnalla@dsgsys.com",
		"rveeranki@dsgsys.com",
		"sri@createchsys.com",
		"rmaddhi@gmail.com",
		"sumanth.asap@gmail.com",
		"sara95@mac.com",
		"jdodda@gmail.com",
		"domakuntla.srinivas@gmail.com",
		"baskrack@gmail.com",
		"vlkrishna@gmail.com",
		"ukmohan@me.com",
		"kevin.morningstar@gmail.com",
		"satish_vegesna@yahoo.com",
		"ad_rao@yahoo.com",
		"mohanmuthu@yahoo.com",
		"phanikola@gmail.com",
		"kalagara_rama@yahoo.com",
		"gvrao98@yahoo.com",
		"ksreddy007in@yahoo.com",
		"saishashank@gmail.com",
		"padma.nimmala@gmail.com",
		"nimmalavenkat@gmail.com",
	}
	for _, b := range fundone_insvestors {
		if b == a {
			return true
		}
	}
	return false
}
func isFundII(a string) bool {
	//hardcode roles temporarily
	fundtwo_insvestors := []string{
		"sgosala99@gmail.com",
		"bens@hotmail.com",
		"smallina@yahoo.com",
		"rpeddamallu@gmail.com",
		"venkata.konkala@gmail.com",
		"ramu433@yahoo.com",
		"immanni@gmail.com",
		"arun.taman@gmail.com",
		"sadhu.behera@gmail.com",
		"niraj_desai@yahoo.com",
		"vkachhia@yahoo.com",
		"hemmathur@yahoo.com",
		"bobkusal@outlook.com",
		"hemmathur@yahoo.com",
		"ssatrasa@gmail.com",
		"tamilselvant@yahoo.com",
		"hkapasi@gmail.com",
		"kalyanmuddasani@gmail.com",
		"avinashreddy1@gmail.com",
		"ckgovula@gmail.com",
		"roy.rajiv@gmail.com",
		"vkrishna28@gmail.com",
		"dparekh@adt.com",
		"rajan.modi@oracle.com",
		"drmoditejas@gmail.com",
		"drtejasmodi@yahoo.com",
		"baskrack@gmail.com",
		"igsvenkat@gmail.com",
		"vamseekc@yahoo.com",
		"ssheik007@hotmail.com",
		"satish_vegesna@yahoo.com",
		"kiran_misc@yahoo.com",
		"mudigondag@yahoo.com",
		"maddali.srinivas@gmail.com",
		"rmaddhi@gmail.com",
		"skondam@gmail.com",
		"rajesh.gundu30@gmail.com",
		"dileep.kasam@gmail.com",
		"padma.nimmala@gmail.com",
		"nimmalavenkat@gmail.com",
		"venkatesh.pallipadi@gmail.com",
	}
	for _, b := range fundtwo_insvestors {
		if b == a {
			return true
		}
	}
	return false
}
