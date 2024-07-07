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

type gormRoleRepositoryTestSuite struct {
	suite.Suite
	mock     sqlmock.Sqlmock
	DB       *gorm.DB
	roleRepo ports.RoleRepository
}

func (suite *gormRoleRepositoryTestSuite) SetupSuite() {
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

	suite.roleRepo = repositories.NewGormRoleRepository(suite.DB)
}

func (suite *gormRoleRepositoryTestSuite) TearDownSuite() {
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

func (suite *gormRoleRepositoryTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestGormRoleRepository(t *testing.T) {
	suite.Run(t, new(gormRoleRepositoryTestSuite))
}

func (suite *gormRoleRepositoryTestSuite) TestCreateRole_Success() {
	roleName := "admin"
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "roles"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, roleName))
	suite.mock.ExpectCommit()

	role, err := suite.roleRepo.CreateRole(roleName)

	suite.NoError(err)
	suite.Equal(roleName, role.Name)
}

func (suite *gormRoleRepositoryTestSuite) TestCreateRole_Error() {
	roleName := "admin"
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "roles"`).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	_, err := suite.roleRepo.CreateRole(roleName)

	suite.Error(err)
}

func (suite *gormRoleRepositoryTestSuite) TestReadRoles_Success() {
	firstRoleName := "admin"
	secondRoleName := "user"
	suite.mock.ExpectQuery(`SELECT \* FROM "roles"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, firstRoleName).
			AddRow(2, secondRoleName))

	roles, err := suite.roleRepo.ReadRoles()

	suite.NoError(err)
	suite.Len(roles, 2)
	suite.Equal(firstRoleName, roles[0].Name)
	suite.Equal(secondRoleName, roles[1].Name)
}

func (suite *gormRoleRepositoryTestSuite) TestReadRoles_Error() {
	suite.mock.ExpectQuery(`SELECT \* FROM "roles"`).
		WillReturnError(errors.New("some error"))

	_, err := suite.roleRepo.ReadRoles()

	suite.Error(err)
}

func (suite *gormRoleRepositoryTestSuite) TestReadRoles_Empty() {
	suite.mock.ExpectQuery(`SELECT \* FROM "roles"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	roles, err := suite.roleRepo.ReadRoles()

	suite.NoError(err)
	suite.Len(roles, 0)
}

func (suite *gormRoleRepositoryTestSuite) TestReadRole_Success() {
	roleID := 1
	roleName := "admin"
	suite.mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(roleID, roleName))

	role, err := suite.roleRepo.ReadRole(strconv.Itoa(roleID))

	suite.NoError(err)
	suite.Equal(roleName, role.Name)
}

func (suite *gormRoleRepositoryTestSuite) TestReadRole_Error() {
	roleID := 1
	suite.mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id"`).
		WillReturnError(errors.New("some error"))

	_, err := suite.roleRepo.ReadRole(strconv.Itoa(roleID))

	suite.Error(err)

}

func (suite *gormRoleRepositoryTestSuite) TestReadRole_NotFound() {
	roleID := 1
	suite.mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	role, err := suite.roleRepo.ReadRole(strconv.Itoa(roleID))

	suite.Error(err)
	suite.Nil(role)
}

func (suite *gormRoleRepositoryTestSuite) TestUpdateRole_Success() {
	roleID := 1
	roleName := "admin"
	newRoleName := "superadmin"
	suite.mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(roleID, roleName))
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "roles"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	newRole := domain.Role{Name: newRoleName}
	updatedRole, err := suite.roleRepo.UpdateRole(strconv.Itoa(roleID), newRole)

	suite.NoError(err)
	suite.Equal(newRoleName, updatedRole.Name)
}

func (suite *gormRoleRepositoryTestSuite) TestUpdateRole_Error() {
	roleID := 1
	roleName := "admin"
	newRoleName := "superadmin"
	suite.mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(roleID, roleName))
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "roles"`).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	newRole := domain.Role{Name: newRoleName}
	_, err := suite.roleRepo.UpdateRole(strconv.Itoa(roleID), newRole)

	suite.Error(err)
}

func (suite *gormRoleRepositoryTestSuite) TestUpdateRole_NotFound() {
	roleID := 2
	newRoleName := "superadmin"
	suite.mock.ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"."id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	newRole := domain.Role{Name: newRoleName}
	_, err := suite.roleRepo.UpdateRole(strconv.Itoa(roleID), newRole)

	suite.Error(err)
}

func (suite *gormRoleRepositoryTestSuite) TestDeleteRole_Success() {
	roleID := 1
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "roles" SET "deleted_at"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.roleRepo.DeleteRole(strconv.Itoa(roleID))

	suite.NoError(err)
}

func (suite *gormRoleRepositoryTestSuite) TestDeleteRole_Error() {
	roleID := 1
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "roles" SET "deleted_at`).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	err := suite.roleRepo.DeleteRole(strconv.Itoa(roleID))

	suite.Error(err)
}
