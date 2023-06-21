package data

import (
	"context"
	"fmt"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"

	user "github.com/Allen-Career-Institute/go-kratos-sample/api/user/v1"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/data/entity"
)

const tableName = "users"

type userRepository struct {
	db  *gorm.DB
	log *log.Helper
}

func NewUserRepository(data *Data, logger log.Logger) biz.UserRepository {
	return &userRepository{data.db, log.NewHelper(logger)}
}

func (u *userRepository) FindById(ctx context.Context, id string) (*entity.UserEntity, error) {
	userEntity := entity.UserEntity{}
	if err := u.db.WithContext(ctx).Table(tableName).First(&userEntity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &userEntity, nil
}

func (u *userRepository) FindByMobileNumber(ctx context.Context, mobileNumber string) (*entity.UserEntity, error) {
	userEntity := entity.UserEntity{}
	if err := u.db.WithContext(ctx).Table(tableName).Where("mobile_number = ?", mobileNumber).First(&userEntity).Error; err != nil {
		return nil, err
	}

	return &userEntity, nil
}

func (u *userRepository) MobileNumberExists(ctx context.Context, mobileNumber string) (exists bool, err error) {
	var count int64
	if err = u.db.WithContext(ctx).Table(tableName).Where("mobile_number = ?", mobileNumber).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *userRepository) Create(ctx context.Context, entity *entity.UserEntity) (err error) {
	exists, err := u.MobileNumberExists(ctx, entity.MobileNumber)
	if err != nil {
		return err
	} else if exists {
		return user.ErrorUserAlreadyExist("mobile_number already exists, mobile_number: %v", entity.MobileNumber)
	}

	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		entity.ID = uuid.New()
		entity.Status = user.Status_STATUS_CREATED.String()
		if err := tx.Table(tableName).Create(&entity).Error; err != nil {
			return err
		}

		return nil
	})

	return
}

func (u *userRepository) Update(ctx context.Context, e *entity.UserEntity) (err error) {
	var existingEntity *entity.UserEntity
	if existingEntity, err = u.FindById(ctx, e.ID.String()); err != nil {
		return user.ErrorUserNotFound("User not found for the userId: %s", e.ID)
	}

	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existingEntity.ID = e.ID
		existingEntity.GivenName = e.GivenName
		existingEntity.FamilyName = e.FamilyName
		existingEntity.Status = e.Status
		if err := tx.Table(tableName).Save(&existingEntity).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (u *userRepository) Delete(ctx context.Context, id string) (err error) {
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(tableName).Delete(&entity.UserEntity{}, uuid.MustParse(id)).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (u *userRepository) List(ctx context.Context, offset uint, limit int) (
	rows []*entity.UserEntity, err error,
) {
	return nil, fmt.Errorf("operation not implemented")
}
