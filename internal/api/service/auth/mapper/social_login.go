package mapper

import (
	"mandarine/internal/api/persistence/model"
	"mandarine/pkg/oauth"
)

func MapUserInfoToUserEntity(userInfo oauth.UserInfo) *model.UserEntity {
	return &model.UserEntity{
		Username:        userInfo.Username,
		Email:           userInfo.Email,
		IsEnabled:       true,
		IsEmailVerified: userInfo.IsEmailVerified,
		IsPasswordTemp:  true,
	}
}
