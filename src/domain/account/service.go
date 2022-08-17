package account

import (
	"pangya/src/internal/database"
	"pangya/src/models"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	CreateAccount(req CreateAccountRequest) (uint, error)
	FindAccountByUsernameAndPassword(username string, password string) (models.Account, bool)
	SetNicknamebyId(nickname string, id uint) error
	ExistsNickname(nickname string) (bool, error)
}

var svc Service

type accountService struct {
	db *gorm.DB
}

// ExistsNickname implements Service
func (as *accountService) ExistsNickname(nickname string) (bool, error) {
	var count int64
	g := as.db.Model(&models.Account{}).Where("nickname = ?", nickname).Count(&count)
	if g.Error != nil {
		return false, g.Error
	}
	return count != 0, nil
}

// SetNicknamebyId implements Service
func (as *accountService) SetNicknamebyId(nickname string, id uint) error {
	return as.db.Model(&models.Account{}).Where("id = ?", id).Update("nickname", nickname).Error
}

// CreateAccount implements Service
func (as *accountService) CreateAccount(req CreateAccountRequest) (uint, error) {
	acc := models.Account{
		Username: req.Username,
		Password: strings.ToLower(req.Password),
	}
	result := as.db.Create(&acc)
	if result.Error != nil {
		return 0, result.Error
	}
	return acc.ID, nil
}

// FindAccountByUsernameAndPassword implements Service
func (as *accountService) FindAccountByUsernameAndPassword(username string, password string) (models.Account, bool) {
	var acc models.Account
	result := as.db.First(&acc, "username = ? AND password = ?", username, strings.ToLower(password))
	if result.RowsAffected == 0 {
		return models.Account{}, false
	}
	return acc, true
}

func Svc() Service {
	if svc == nil {
		svc = &accountService{
			db: database.Get(),
		}
	}
	return svc
}
