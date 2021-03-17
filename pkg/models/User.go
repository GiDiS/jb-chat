package models

import "strconv"

type Uid int32

const NoUser Uid = 0
const ToAllUsers Uid = -1
const ToConnection Uid = -2

type Email string

type User struct {
	UserId    Uid    `json:"uid" db:"uid"`
	Nickname  string `json:"nickname" db:"nickname"`
	Title     string `json:"title" db:"title"`
	Email     string `json:"email" db:"email"`
	AvatarUrl string `json:"avatarUrl" db:"avatar_url"`
}

func (uid Uid) String() string {
	return strconv.Itoa(int(uid))
}

func (u *User) Update(newUser User) (bool, error) {
	updated := false
	if u.Title != newUser.Title {
		u.Title = newUser.Title
		updated = true
	}
	if u.Nickname != newUser.Nickname {
		u.Nickname = newUser.Nickname
		updated = true
	}

	if u.AvatarUrl != newUser.AvatarUrl {
		u.AvatarUrl = newUser.AvatarUrl
		updated = true
	}

	return updated, nil
}
