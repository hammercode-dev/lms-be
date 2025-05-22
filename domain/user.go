package domain

import (
	"context"
	"net/http"
	"time"
)

type (
	UserRepository interface {
		GetUsers(ctx context.Context) (users []User, err error)
		CreateUser(ctx context.Context, userReq User) (user User, err error)
		FindById(ctx context.Context, id int8) (user User, err error)
		FindByEmail(ctx context.Context, email string) (user User, err error)
		UpdateProfileUser(ctx context.Context, userReq UserUpdateProfile, id int) error
		DeleteUser(ctx context.Context, id int8) error
		LogoutUser(ctx context.Context, token string, expiredAt time.Time) error
		CleanupLogoutToken(ctx context.Context) error
		GetUsersGenericConditions(ctx context.Context, filter GetUserBy) (users []User, err error)
		ResetPassword(ctx context.Context, email, password, token string) error
		ForgotPassword(ctx context.Context, token string, expiredAt time.Time, user User) error
		StoreToken(ctx context.Context, token string, expiredAt time.Time, uid int) error
		UnactivateTokenByUser(ctx context.Context, uid int) error
		CheckActiveToken(ctx context.Context, token string) (logoutToken LogoutToken, err error)
	}
	UserUsecase interface {
		GetUsers(ctx context.Context) (users []User, err error)
		GetUserById(ctx context.Context, id int8) (users User, err error)
		Register(ctx context.Context, userReq User) (user User, err error)
		Login(ctx context.Context, userReq Login) (user User, token string, err error)
		UpdateProfileUser(ctx context.Context, userReq UserUpdateProfile, id int) error
		DeleteUser(ctx context.Context, id int8) error
		Logout(ctx context.Context, token string) error
		ForgotPassword(ctx context.Context, emailForgot ForgotPassword) error
		ResetPassword(ctx context.Context, requestUser ForgotPassword) error
	}

	UserHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
		UpdateProfileUser(w http.ResponseWriter, r *http.Request)
		DeleteUser(w http.ResponseWriter, r *http.Request)
		GetUsers(w http.ResponseWriter, r *http.Request)
		Logout(w http.ResponseWriter, r *http.Request)
		Register(w http.ResponseWriter, r *http.Request)
		GetUserById(w http.ResponseWriter, r *http.Request)
		GetUserProfile(w http.ResponseWriter, r *http.Request)
		ForgotPassword(w http.ResponseWriter, r *http.Request)
		ResetPassword(w http.ResponseWriter, r *http.Request)
	}
	User struct {
		ID          int       `gorm:"primaryKey" json:"id"`
		Username    string    `json:"username" gorm:"type:varchar(255);not null;unique"`
		Email       string    `json:"email" gorm:"type:varchar(255);not null;unique"`
		Password    string    `json:"password" gorm:"type:varchar(255);not null"`
		Role        string    `json:"role"`
		Fullname    string    `json:"fullname"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Gender      string    `json:"gender"`
		PhoneNumber string    `json:"phone_number"`
		Address     string    `json:"address"`
		Github      string    `json:"github"`
		Linkedin    string    `json:"linkedin"`
		PersonalWeb string    `json:"personal_web"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	UserUpdateProfile struct {
		Username    string    `json:"username"`
		Fullname    string    `json:"fullname"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Gender      string    `json:"gender"`
		PhoneNumber string    `json:"phone_number"`
		Address     string    `json:"address"`
		Github      string    `json:"github"`
		Linkedin    string    `json:"linkedin"`
		PersonalWeb string    `json:"personal_web"`
	}

	Register struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	Login struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	GetUserBy struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email" `
		Role     string `json:"role"`
	}

	ForgotPassword struct {
		Email           string `json:"email,omitempty"`
		Token           string `json:"token,omitempty"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
)

func (User) TableName() string {
	return "users"
}

func RegistToUser(r Register) User {
	return User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

func UserUpdateProfileToUser(u UserUpdateProfile) User {
	return User{
		Username:    u.Username,
		Fullname:    u.Fullname,
		DateOfBirth: u.DateOfBirth,
		Gender:      u.Gender,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
		Github:      u.Github,
		Linkedin:    u.Linkedin,
		PersonalWeb: u.PersonalWeb,
	}
}
