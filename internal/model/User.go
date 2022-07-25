package model

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// User 用户数据结构
type User struct {
	ID           uint
	Email        string
	Username     string
	Bio          string
	Image        string
	PasswordHash string
}

//UserLogin 登录状态
type UserLogin struct {
	Email    string
	Username string
	Token    string
	Bio      string
	Image    string
}

type UserUpdate struct {
	Email    string
	Username string
	Password string
	Bio      string
	Image    string
}

// HashPassword 加盐加密
func hashPassword(pwd string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(b) // 加盐之后返回
}

// 验证输入和加盐密码
func verifyPassword(hashed, input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input)); err != nil {
		return false
	}
	return true
}

// UserRepo 用户操作
type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error                      // 创建用户
	GetUserByEmail(ctx context.Context, email string) (*User, error)       //通过email获取用户信息
	GetUserByUsername(ctx context.Context, username string) (*User, error) //通过用户名获取用户信息
	GetUserByID(ctx context.Context, id uint) (*User, error)               //通过id获取用户信息
	UpdateUser(ctx context.Context, user *User) (*User, error)             //更新用户
}

// ProfileRepo 简介的更新
type ProfileRepo interface {
	GetProfile(ctx context.Context, username string) (*Profile, error)                                            // 查询简介操作
	FollowUser(ctx context.Context, currentUserID uint, followingID uint) error                                   // 关注操作
	UnfollowUser(ctx context.Context, currentUserID uint, followingID uint) error                                 // 取消关注操作
	GetUserFollowingStatus(ctx context.Context, currentUserID uint, userIDs []uint) (following []bool, err error) // 通过关注状态获取用户信息
}

type UserUsecase struct {
	ur UserRepo
	pr ProfileRepo
}

// Profile 简介的结构
type Profile struct {
	ID        uint
	Username  string
	Bio       string
	Image     string
	Following bool
}

func NewUserUsecase(ur UserRepo,
	pr ProfileRepo) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr}
}

func (uc *UserUsecase) Register(ctx context.Context, username, email, password string) (*UserLogin, error) {
	u := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if err := uc.ur.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &UserLogin{
		Email:    email,
		Username: username,
	}, nil
}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (*UserLogin, error) {
	if len(email) == 0 {
		return nil, errors.New("422, email ，cannot be empty")
	}
	u, err := uc.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !verifyPassword(u.PasswordHash, password) {
		return nil, errors.New("user login failed")
	}

	return &UserLogin{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}, nil
}

//func (uc *UserUsecase) GetCurrentUser(ctx context.Context) (*User, error) {
//	cu := auth.FromContext(ctx)
//	u, err := uc.ur.GetUserByID(ctx, cu.UserID)
//	if err != nil {
//		return nil, err
//	}
//	return u, nil
//}
//
//func (uc *UserUsecase) UpdateUser(ctx context.Context, uu *UserUpdate) (*UserLogin, error) {
//	cu := auth.FromContext(ctx)
//	u, err := uc.ur.GetUserByID(ctx, cu.UserID)
//	if err != nil {
//		return nil, err
//	}
//	u.Email = uu.Email
//	u.Image = uu.Image
//	u.PasswordHash = hashPassword(uu.Password)
//	u.Bio = uu.Bio
//	u, err = uc.ur.UpdateUser(ctx, u)
//	if err != nil {
//		return nil, err
//	}
//	return &UserLogin{
//		Email:    u.Email,
//		Username: u.Username,
//		Bio:      u.Bio,
//		Image:    u.Image,
//	}, nil
//}
