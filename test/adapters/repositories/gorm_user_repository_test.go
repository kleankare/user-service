package test

import (
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kleankare/user-service/internal/adapters/repositories"
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type gormUserRepositoryTestSuite struct {
	suite.Suite
	mock     sqlmock.Sqlmock
	DB       *gorm.DB
	userRepo ports.UserRepository
}

func (suite *gormUserRepositoryTestSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	db, suite.mock, err = sqlmock.New()
	require.NoError(suite.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	suite.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(suite.T(), err)

	suite.userRepo = repositories.NewGormUserRepository(suite.DB)
}

func (suite *gormUserRepositoryTestSuite) TearDownSuite() {
	sqlDB, err := suite.DB.DB()
	if err != nil {
		suite.T().Fatalf("Failed to get database connection: %v", err)
	}

	suite.mock.ExpectClose()
	err = sqlDB.Close()
	if err != nil {
		suite.T().Fatalf("Failed to close database connection: %v", err)
	}
}

func (suite *gormUserRepositoryTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestGormUserRepository(t *testing.T) {
	suite.Run(t, new(gormUserRepositoryTestSuite))
}

func (suite *gormUserRepositoryTestSuite) TestCreateUser_Success() {
	username := "<USERNAME>"
	password := "<PASSWORD>"
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(1, username, password))
	suite.mock.ExpectCommit()

	role, err := suite.userRepo.CreateUser(username, password)

	suite.NoError(err)
	suite.Equal(username, role.Username)
	suite.Equal(password, role.Password)
}

func (suite *gormUserRepositoryTestSuite) TestCreateUser_Error() {
	username := "<USERNAME>"
	password := "<PASSWORD>"
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "users"`).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	_, err := suite.userRepo.CreateUser(username, password)

	suite.Error(err)
}

func (suite *gormUserRepositoryTestSuite) TestReadUsers_Success() {
	firstUserName := "<FIRST_USERNAME>"
	secondUserName := "<SECOND_USERNAME>"
	password := "<PASSWORD>"
	suite.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(1, firstUserName, password).
			AddRow(2, secondUserName, password))

	users, err := suite.userRepo.ReadUsers()

	suite.NoError(err)
	suite.Len(users, 2)
	suite.Equal(firstUserName, users[0].Username)
	suite.Equal(secondUserName, users[1].Username)
	suite.Equal(password, users[0].Password)
}

func (suite *gormUserRepositoryTestSuite) TestReadUsers_Error() {
	suite.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnError(errors.New("some error"))

	_, err := suite.userRepo.ReadUsers()

	suite.Error(err)
}

func (suite *gormUserRepositoryTestSuite) TestReadUsers_Empty() {
	suite.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}))

	users, err := suite.userRepo.ReadUsers()

	suite.NoError(err)
	suite.Len(users, 0)
}

func (suite *gormUserRepositoryTestSuite) TestReadUser_Success() {
	userID := 1
	username := "<USERNAME>"
	password := "<PASSWORD>"
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(userID, username, password))

	role, err := suite.userRepo.ReadUser(strconv.Itoa(userID))

	suite.NoError(err)
	suite.Equal(username, role.Username)
}

func (suite *gormUserRepositoryTestSuite) TestReadUser_Error() {
	userID := 1
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id"`).
		WillReturnError(errors.New("some error"))

	_, err := suite.userRepo.ReadUser(strconv.Itoa(userID))

	suite.Error(err)
}

func (suite *gormUserRepositoryTestSuite) TestReadUser_NotFound() {
	userID := 1
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}))

	role, err := suite.userRepo.ReadUser(strconv.Itoa(userID))

	suite.Error(err)
	suite.Nil(role)
}

func (suite *gormUserRepositoryTestSuite) TestUpdateUser_Success() {
	userID := 1
	username := "<USERNAME>"
	newUserName := "<NEW_USERNAME>"
	password := "<PASSWORD>"
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(userID, username, password))
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "users"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	newUser := domain.User{Username: newUserName}
	updatedUser, err := suite.userRepo.UpdateUser(strconv.Itoa(userID), newUser)

	suite.NoError(err)
	suite.Equal(newUserName, updatedUser.Username)
}

func (suite *gormUserRepositoryTestSuite) TestUpdateUser_Error() {
	userID := 1
	username := "<USERNAME>"
	newUserName := "superadmin"
	password := "<PASSWORD>"
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(userID, username, password))
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "users"`).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	newUser := domain.User{Username: newUserName}
	_, err := suite.userRepo.UpdateUser(strconv.Itoa(userID), newUser)

	suite.Error(err)
}

func (suite *gormUserRepositoryTestSuite) TestUpdateUser_NotFound() {
	userID := 2
	newUserName := "superadmin"
	suite.mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}))

	newUser := domain.User{Username: newUserName}
	_, err := suite.userRepo.UpdateUser(strconv.Itoa(userID), newUser)

	suite.Error(err)
}

func (suite *gormUserRepositoryTestSuite) TestDeleteUser_Success() {
	userID := 1
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "users" SET "deleted_at"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.userRepo.DeleteUser(strconv.Itoa(userID))

	suite.NoError(err)
}

func (suite *gormUserRepositoryTestSuite) TestDeleteUser_Error() {
	userID := 1
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "users" SET "deleted_at`).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	err := suite.userRepo.DeleteUser(strconv.Itoa(userID))

	suite.Error(err)
}
