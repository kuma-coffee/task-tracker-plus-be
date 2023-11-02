package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) *sessionsRepo {
	return &sessionsRepo{db}
}

func (u *sessionsRepo) AddSessions(session model.Session) error {
	u.db.Create(&session)

	return nil
}

func (u *sessionsRepo) DeleteSession(token string) error {
	err := u.db.Where("token = ?", token).Delete(&model.Session{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *sessionsRepo) UpdateSessions(session model.Session) error {
	err := u.db.Table("sessions").Where("email = ?", session.Email).Updates(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *sessionsRepo) SessionAvailEmail(email string) (model.Session, error) {
	temp := []model.Session{}
	err := u.db.Table("sessions").Select("*").Scan(&temp).Error
	if err != nil {
		return model.Session{}, err
	}

	for _, v := range temp {
		if v.Email == email {
			return v, nil
		}
	}

	return model.Session{}, errors.New("session not found")
}

func (u *sessionsRepo) SessionAvailToken(token string) (model.Session, error) {
	temp := []model.Session{}
	err := u.db.Table("sessions").Select("*").Scan(&temp).Error
	if err != nil {
		return model.Session{}, err
	}

	for _, v := range temp {
		if v.Token == token {
			return v, nil
		}
	}

	return model.Session{}, errors.New("session not found")
}

func (u *sessionsRepo) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		err := u.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, err
	}

	return session, nil
}

func (u *sessionsRepo) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
