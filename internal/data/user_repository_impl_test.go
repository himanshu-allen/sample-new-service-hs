package data

import (
	"context"
	"errors"
	user "github.com/Allen-Career-Institute/go-kratos-sample/api/user/v1"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/data/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

func TestUserRepository_FindById(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()
	db.Table("users").Create(&testUser)

	userRepository := &userRepository{db: db}
	actualUser, err := userRepository.FindById(context.Background(), testUser.ID.String())
	actualUser.UpdatedAt = actualUser.UpdatedAt.Local()
	actualUser.CreatedAt = actualUser.CreatedAt.Local()

	require.NoError(t, err)
	require.Equal(t, testUser, actualUser)
}

func TestUserRepository_FindById_IdDoesNotExist(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	userRepository := &userRepository{db: db}
	_, err = userRepository.FindById(context.Background(), testUser.ID.String())

	require.Equal(t, errors.New("record not found"), err)
}

func TestUserRepository_FindByMobileNumber(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()
	db.Table("users").Create(&testUser)
	expectedUser := testUser

	userRepository := &userRepository{db: db}
	actualUser, err := userRepository.FindByMobileNumber(context.Background(), testUser.MobileNumber)
	actualUser.UpdatedAt = actualUser.UpdatedAt.Local()
	actualUser.CreatedAt = actualUser.CreatedAt.Local()

	require.NoError(t, err)
	require.Equal(t, expectedUser, actualUser)
}

func TestUserRepository_FindByMobileNumber_MobileNumberDoesNotExist(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	userRepository := &userRepository{db: db}
	_, err = userRepository.FindByMobileNumber(context.Background(), testUser.MobileNumber)

	require.Equal(t, errors.New("record not found"), err)
}

func TestUserRepository_Create(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	userRepository := &userRepository{db: db}
	err = userRepository.Create(context.Background(), testUser)

	savedUserEntity := entity.UserEntity{}
	db.Table("users").Where("mobile_number = ?", testUser.MobileNumber).First(&savedUserEntity)

	require.NoError(t, err)
	require.Equal(t, savedUserEntity.Status, "STATUS_CREATED")
}

func TestUserRepository_Create_FailsWhenExists(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()
	db.Table("users").Create(&testUser)

	userRepository := &userRepository{db: db}
	err = userRepository.Create(context.Background(), testUser)

	require.Equal(t, user.ErrorUserAlreadyExist("mobile_number already exists, mobile_number: 9876543210"), err)
}

func TestUserRepository_Update_Status(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()
	db.Table("users").Create(&testUser)

	testUser.Status = user.Status_STATUS_ACTIVE.String()
	userRepository := &userRepository{db: db}
	err = userRepository.Update(context.Background(), testUser)

	updatedUserEntity := entity.UserEntity{}
	db.Table("users").Where("mobile_number = ?", testUser.MobileNumber).First(&updatedUserEntity)

	require.NoError(t, err)
	require.Equal(t, "STATUS_ACTIVE", updatedUserEntity.Status)
}

func TestUserRepository_Update_Name(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()
	db.Table("users").Create(&testUser)

	testUser.GivenName = "Test"
	testUser.FamilyName = "Name"
	userRepository := &userRepository{db: db}
	err = userRepository.Update(context.Background(), testUser)

	updatedUserEntity := entity.UserEntity{}
	db.Table("users").Where("mobile_number = ?", testUser.MobileNumber).First(&updatedUserEntity)

	require.NoError(t, err)
	require.Equal(t, "STATUS_CREATED", updatedUserEntity.Status)
	require.Equal(t, "Test", updatedUserEntity.GivenName)
	require.Equal(t, "Name", updatedUserEntity.FamilyName)
}

func TestUserRepository_Update_FailsWhenDoesNotExist(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	userRepository := &userRepository{db: db}
	err = userRepository.Update(context.Background(), testUser)

	require.Equal(t, user.ErrorUserNotFound("User not found for the userId: %s", testUser.ID), err)
}

func TestUserRepository_Delete(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()
	db.Table("users").Create(&testUser)

	userRepository := &userRepository{db: db}
	err = userRepository.Delete(context.Background(), testUser.ID.String())

	deletedUserEntity := entity.UserEntity{}
	err2 := db.Table("users").Where("mobile_number = ?", testUser.MobileNumber).First(&deletedUserEntity).Error

	require.NoError(t, err)
	require.Error(t, errors.New("record not found"), err2)
}

func TestUserRepository_Delete_FailsWhenDoesNotExist(t *testing.T) {
	db, testUser, err := setup()
	if err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}
	defer teardown(db)

	testUser.ID = uuid.New()
	testUser.Status = user.Status_STATUS_CREATED.String()

	userRepository := &userRepository{db: db}
	err = userRepository.Delete(context.Background(), testUser.ID.String())

	require.Error(t, errors.New("record not found"), err)
}

func setup() (db *gorm.DB, testUser *entity.UserEntity, err error) {
	db, err = gorm.Open(sqlite.Open("file:../../../test.db?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to open sqlite for unit testing: %v", err)
		return nil, nil, err
	}

	err = db.Table("users").Migrator().CreateTable(&entity.UserEntity{})
	if err != nil {
		return nil, nil, err
	}
	db.Migrator().HasTable("users")

	testUser = &entity.UserEntity{
		ID:           [16]byte{},
		MobileNumber: "9876543210",
		GivenName:    "John",
		FamilyName:   "Doe",
		Status:       "",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	return db, testUser, nil
}

func teardown(db *gorm.DB) {
	err := db.Migrator().DropTable("users")
	if err != nil {
		return
	}
}
