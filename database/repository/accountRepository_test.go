package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"social_network_project/entities"
	"strings"
	"testing"
	"time"
)

var u = &entities.Account{
	ID:          uuid.New().String(),
	Username:    "marcelito001",
	Name:        "Marcelo Sabido",
	Description: "Eu Marcelo, Eu Marcelo",
	Email:       "marcelo111@gmail.com",
	Password:    "23042",
	CreatedAt:   time.Now().UTC().Format("2006-01-02"),
	UpdatedAt:   time.Now().UTC().Format("2006-01-02"),
	Deleted:     false,
}

var body = map[string]interface{}{
	"username":    "maciel",
	"name":        "Nicole Miguel Maciel ",
	"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer porta vehicula purus bibendum pretium.",
	"email":       "ralph333@gmail.com",
	"password":    "2222",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestAccountRepository_InsertAccount(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		INSERT INTO account (id, username, name, description, email, password, created_at, updated_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(u.ID, u.Username, u.Name, u.Description, u.Email, u.Password,
		u.CreatedAt, u.UpdatedAt, u.Deleted).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertAccount(u)
	assert.Error(t, err)
}

func TestAccountRepository_FindAccountPasswordByEmail(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT password 
		FROM account
		WHERE email = $1
		AND deleted = false`

	rows := sqlmock.NewRows([]string{"password", "email"}).
		AddRow(u.Password, u.Email)

	mock.ExpectQuery(query).WithArgs(u.Email).WillReturnRows(rows)

	id, err := repo.FindAccountIDbyEmail(u.Email)
	assert.Empty(t, id)
	assert.Error(t, err)
}

func TestAccountRepository_FindAccountIDbyEmail(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM account
		WHERE email = $1`

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(u.ID, u.Email)

	mock.ExpectQuery(query).WithArgs(u.Email).WillReturnRows(rows)

	id, err := repo.FindAccountIDbyEmail(u.Email)
	assert.Empty(t, id)
	assert.Error(t, err)
}

func TestAccountRepository_FindAccountbyID(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id, username, name, description, email, password, created_at, updated_at, deleted
		FROM account
		WHERE id = $1`

	rows := sqlmock.NewRows([]string{"id", "username", "name", "description", "email", "password", "created_at",
		"updated_at", "deleted"}).
		AddRow(u.ID, u.Username, u.Name, u.Description, u.Email, u.Password, u.CreatedAt, u.UpdatedAt, u.Deleted)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	account, err := repo.FindAccountByID(&u.ID)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestAccountRepository_dinamicQueryChangeAccountDataByID(t *testing.T) {

	var values []interface{}
	var where []string

	for key, value := range body {
		values = append(values, value)
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, key, value))
	}

	stringQueryExpected := "UPDATE account SET " + strings.Join(where, ", ") + " WHERE id = $1 AND deleted = false"

	stringQuery := dinamicQueryChangeAccountDataByID(body)

	assert.Equal(t, len(stringQueryExpected), len(stringQuery))
}

func TestAccountRepository_ChangeAccountDataByID(t *testing.T) {

	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := dinamicQueryChangeAccountDataByID(body)

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ChangeAccountDataByID(&u.ID, body)
	assert.Error(t, err)
}

func TestAccountRepository_DeleteAccountByID(t *testing.T) {

	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE account 
		SET deleted = true 
		WHERE id = $1`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteAccountByID(&u.ID)
	assert.Error(t, err)
}

func TestAccountRepository_ExistsAccountByID(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM account
		WHERE id = $1
		AND deleted = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(u.ID)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	exist, err := repo.ExistsAccountByID(&u.ID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestAccountRepository_ExistsAccountByUsername(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM account
		WHERE username = $1
		AND deleted = false`

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(u.ID, u.Username)

	mock.ExpectQuery(query).WithArgs(u.Email).WillReturnRows(rows)

	exist, err := repo.ExistsAccountByUsername(&u.Username)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestAccountRepository_ExistsAccountByEmail(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM account
		WHERE email = $1
		AND deleted = false`

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(u.ID, u.Email)

	mock.ExpectQuery(query).WithArgs(u.Email).WillReturnRows(rows)

	exist, err := repo.ExistsAccountByEmail(&u.Email)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestAccountRepositoryStruct_InsertAccountFollow(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		INSERT INTO account_follow (account_id, account_id_followed)
		VALUES ($1, $2)`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(u.ID, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertAccountFollow(&u.ID, &u.ID)
	assert.Error(t, err)
}

func TestAccountRepositoryStruct_FindAccountFollowingByAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()
	page := "1"
	query := `
		SELECT account.id, account.username, account."name", account.description, account.email,
		account."password", account."password", account.created_at , account.updated_at, account.deleted 
		FROM account_follow
		INNER JOIN account ON account_follow.account_id = account.id
		WHERE account_follow.account_id = $1
		AND account_follow.unfollowed = false`

	rows := sqlmock.NewRows([]string{"account_id_followed", "account_follow"}).
		AddRow(u.ID, u.ID)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	exist, err := repo.FindAccountFollowingByAccountID(&u.ID, &page)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestAccountRepositoryStruct_FindAccountFollowersByAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()
	page := "1"
	query := `
	SELECT account.id, account.username, account.name, account.description, account.email,
	account.password, account.created_at , account.updated_at, account.deleted
	FROM account_follow
	INNER JOIN account ON account_follow.account_id_followed = account.id
	WHERE account_follow.account_id_followed = $1
	AND account_follow.unfollowed = false`

	rows := sqlmock.NewRows([]string{"account_id_followed", "account_follow"}).
		AddRow(u.ID, u.ID)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	exist, err := repo.FindAccountFollowersByAccountID(&u.ID, &page)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestAccountRepositoryStruct_DeleteAccountFollow(t *testing.T) {

	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE account_follow 
		SET unfollowed = true 
		WHERE account_id= $1
		AND account_id_followed = $2`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteAccountFollow(&u.ID, &u.ID)
	assert.Error(t, err)
}

func TestAccountRepositoryStruct_ExistsFollowByAccountIDAndAccountFollowedID(t *testing.T) {
	db, mock := NewMock()
	repo := AccountRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT account_id
		FROM account_follow
		WHERE account_id = $1
		AND account_id_followed = $2`

	rows := sqlmock.NewRows([]string{"account_id", "account_id_followed"}).
		AddRow(u.ID, u.ID)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	exist, err := repo.ExistsFollowByAccountIDAndAccountFollowedID(&u.ID, &u.ID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}
