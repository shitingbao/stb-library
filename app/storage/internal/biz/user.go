package biz

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pborman/uuid"

	"golang.org/x/crypto/scrypt"
)

type UserBase struct {
	Name       string
	Avatar     string
	Email      string
	Phone      string
	Salt       string
	CreateTime string
}

//User 用户对象
type User struct {
	UserBase
	Pwd []byte //对应mysql的varbinary,末尾不会填充，不能使用binary，因为不足会使用ox0填充导致取出的时候多18位的0
}

func (User) TableName() string {
	return "user"
}

type ArgUser struct {
	UserBase
	Pwd string
}

type UserRepo interface {
	GetUser(username string) (*User, error)
	IsExistUser(username string) (bool, error)
	RegisterUser(token, username string)
	DelUser(ctx context.Context, token string)
	// Login(context.Context, *ArgUser) error
	// Logout(context.Context, *User) error
}

type Usercase struct {
	repo UserRepo
	log  *log.Helper
}

func (u *Usercase) LoginOut(ctx context.Context, token string) error {
	u.repo.DelUser(ctx, token)
	return nil
}

func (u *Usercase) Login(ctx context.Context, pa *ArgUser) error {
	if pa.Name == "" || pa.Pwd == "" {
		return errors.New("name or password cant not nil")
	}
	usr, err := u.repo.GetUser(pa.Name)
	if err != nil {
		return err
	}
	ie, err := u.Equal(pa.Pwd, pa.Salt, usr.Pwd)
	if err != nil {
		return err
	}
	if !ie {
		return errors.New("password have error")
	}
	token := uuid.NewUUID().String()
	u.repo.RegisterUser(token, pa.Name)
	return nil
}

//前端的hex字符串
func (u *Usercase) huexEncode(md5Pwd string) ([]byte, error) {
	decoded, err := hex.DecodeString(md5Pwd)
	if err != nil {
		// logrus.WithFields(logrus.Fields{"decode": err}).Error("hex")
		return nil, err
	}
	return decoded, nil
}

//BuildIserSalt 随机获取用户中一段+uuid生成随机盐，防止代码泄密密码生成过程被破解
func (u *Usercase) BuildIserSalt(user string) string {
	rand.Seed(time.Now().UnixNano())
	sl := rand.Intn(len(user))
	return user[sl:] + base64.RawURLEncoding.EncodeToString(uuid.NewUUID())
}

//buildUserPassword 根据密码文本和盐生成密文
func (u *Usercase) buildUserPassword(pwdMd5, salt []byte) ([]byte, error) {
	return scrypt.Key(pwdMd5, salt, 16384, 8, 1, 32)
}

//Equal 密文验证
func (u *Usercase) Equal(pwd, salt string, uPwd []byte) (bool, error) {
	bPwd, err := u.BuildPas(pwd, salt)
	if err != nil {
		return false, err
	}
	return bytes.Equal(bPwd, uPwd), nil
}

//BuildPas 解析前端的hex密码文本，并调用密文生成函数
func (u *Usercase) BuildPas(pwd, salt string) ([]byte, error) {
	h, err := u.huexEncode(pwd)
	if err != nil {
		return nil, err
	}
	bPwd, err := u.buildUserPassword(h, []byte(salt))
	if err != nil {
		// logrus.WithFields(logrus.Fields{"pwd": err}).Error("validPwdMd5")
		return nil, err
	}
	return bPwd, nil
}
