package auth

import (
	"context"
	"fmt"
	"jb_chat/pkg/auth"
	"jb_chat/pkg/models"
	"jb_chat/pkg/store"
	"strconv"
	"strings"
)

type SignIn interface {
	SignIn(ctx context.Context, req AuthSignInRequest) (AuthSignInResponse, error)
}

func (a *authImpl) SignIn(ctx context.Context, req AuthSignInRequest) (resp AuthSignInResponse, err error) {
	var userRef *models.User
	var token string

	if req.Service == "google" {
		userRef, token, err = a.signInGoogle(ctx, req)
	} else if req.Service == "token" {
		userRef, token, err = a.signInToken(ctx, req)
	} else {
		a.logger.Debugf("Sign-in: %v", req)
		return resp, ErrUnknownAuthService
	}

	if userRef != nil {
		resp.SetMe(&models.UserInfo{User: *userRef, Status: models.UserStatusUnknown})
		resp.SetToken(token)
	}

	return resp, nil
}

func (a *authImpl) signInGoogle(ctx context.Context, req AuthSignInRequest) (*models.User, string, error) {
	a.logger.Debugf("%+v", req)
	profile, err := auth.GetProfileByAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, "", ErrAuthRequired
	}
	reqUser := profile.ToUser()
	user, err := a.upsertUser(ctx, reqUser)
	if err != nil {
		return nil, "", err
	}
	// @todo replace to JWT
	token := fmt.Sprintf("uid:%v", user.UserId)
	return &user, token, nil
}

func (a *authImpl) signInToken(ctx context.Context, req AuthSignInRequest) (*models.User, string, error) {
	if req.AccessToken == "Tyrion" {
		user, err := a.usersStore.GetByEmail(ctx, "tyrion.lannister@lannister.got")
		if err != nil {
			return nil, "", err
		}
		return &user, req.AccessToken, nil
	} else if strings.HasPrefix(req.AccessToken, "uid:") {
		uidStr := strings.TrimPrefix(req.AccessToken, "uid:")
		uid, err := strconv.Atoi(uidStr)
		if err != nil {
			return nil, "", err
		}
		user, err := a.usersStore.GetByUid(ctx, models.Uid(uid))
		if err != nil {
			return nil, "", err
		}
		return &user, req.AccessToken, nil
	}
	return nil, "", ErrAuthRequired
}

func (a *authImpl) upsertUser(ctx context.Context, newUser models.User) (models.User, error) {
	user, err := a.usersStore.GetByEmail(ctx, newUser.Email)

	if err != nil && err != store.ErrUserNotFound {
		return newUser, err
	}
	if err == store.ErrUserNotFound {
		uid, err := a.usersStore.Register(ctx, newUser)
		newUser.UserId = uid
		return newUser, err
	} else {
		_, err := a.usersStore.Save(ctx, user)
		return user, err
	}
}
