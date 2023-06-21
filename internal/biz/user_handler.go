package biz

import (
	"context"
	user "github.com/Allen-Career-Institute/go-kratos-sample/api/user/v1"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/data/entity"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(user.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// UserHandler Handler code for user service exposed features to encapsulate business logic
// across multiple repos into the layer.
type UserHandler struct {
	repo UserRepository
	log  *log.Helper
}

type UserRepository interface {
	BaseRepository[entity.UserEntity, string, uint]
	FindByMobileNumber(ctx context.Context, mobileNumber string) (*entity.UserEntity, error)
	MobileNumberExists(ctx context.Context, mobileNumber string) (exists bool, err error)
}

// NewUserHandler creates a new instance of the handler for user entity related operations.
func NewUserHandler(repo UserRepository, logger log.Logger) *UserHandler {
	return &UserHandler{repo: repo, log: log.NewHelper(logger)}
}

// GetUser gets a User by id, and returns the found UserInfo.
func (uc *UserHandler) GetUser(ctx context.Context, id string) (*user.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("GetUser: %v", id)
	ue, err := uc.repo.FindById(ctx, id)
	if err != nil {
		log.Errorf("user not found: %v", id)
		return nil, ErrUserNotFound
	}
	u := &user.UserInfo{}
	ue.ToPB(u)
	return u, nil
}

// CreateUser creates a User, and returns the new UserInfo.
func (uc *UserHandler) CreateUser(ctx context.Context, u *user.UserInfo) (*user.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u)
	ue := &entity.UserEntity{}
	ue.FromPB(u)
	if err := uc.repo.Create(ctx, ue); err != nil {
		log.Errorf("user not created: %v", u)
		return nil, err
	}
	ue.ToPB(u)
	return u, nil
}

// UpdateUser updates a User, and returns the updated UserInfo.
func (uc *UserHandler) UpdateUser(ctx context.Context, u *user.UserInfo) (*user.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %v", u)
	ue := &entity.UserEntity{}
	ue.FromPB(u)
	if err := uc.repo.Update(ctx, ue); err != nil {
		log.Errorf("user not updated: %v", u)
		return nil, err
	}
	ue.ToPB(u)
	return u, nil
}

// DeleteUser deletes a User, and returns the error if any.
func (uc *UserHandler) DeleteUser(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteUser by id: %v", id)
	if err := uc.repo.Delete(ctx, id); err != nil {
		log.Errorf("user not deleted: %v", id)
		return err
	}

	return nil
}
