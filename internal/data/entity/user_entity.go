package entity

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	user "github.com/Allen-Career-Institute/go-kratos-sample/api/user/v1"
)

type UserEntity struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key"`
	MobileNumber string    `gorm:"type:varchar(20);not null;index:idx_mobile_number,unique;"`
	GivenName    string    `gorm:"type:varchar(255);"`
	FamilyName   string    `gorm:"type:varchar(255);"`
	Status       string    `gorm:"type:varchar(45);"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// This function is invoked before creating any UserEntity
func (userEntity *UserEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if userEntity.ID == uuid.Nil {
		userEntity.ID = uuid.New()
	}
	return
}

func (userEntity *UserEntity) FromPB(userInfo *user.UserInfo) {
	if userInfo.Name != "" || len(userInfo.Name) > 0 {
		userEntity.ID = uuid.MustParse(userInfo.Name)
	}

	userEntity.MobileNumber = userInfo.MobileNumber
	userEntity.GivenName = userInfo.GivenName
	userEntity.FamilyName = userInfo.FamilyName

	if userInfo.Status == user.Status_STATUS_UNSPECIFIED {
		userEntity.Status = user.Status_STATUS_CREATED.String()
	} else {
		userEntity.Status = userInfo.Status.String()
	}
}

func (userEntity *UserEntity) ToPB(userInfo *user.UserInfo) {
	userInfo.Name = userEntity.ID.String()
	userInfo.MobileNumber = userEntity.MobileNumber
	userInfo.GivenName = userEntity.GivenName
	userInfo.FamilyName = userEntity.FamilyName
	userInfo.Status = user.Status(user.Status_value[userEntity.Status])
	userInfo.CreatedAt = timestamppb.New(userEntity.CreatedAt)
	userInfo.UpdatedAt = timestamppb.New(userEntity.UpdatedAt)
}
