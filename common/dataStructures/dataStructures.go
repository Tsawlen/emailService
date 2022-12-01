package dataStructures

import (
	"time"

	"github.com/google/uuid"
)

type IncomingSubMessage struct {
	Id          uuid.UUID `json:"id"`
	Topic       string    `json:"topic"`
	Message     []byte    `json:"message"`
	Service     string    `json:"service"`
	ReceivedAt  time.Time `json:"receivedAt"`
	DeliveredAt time.Time `json:"deliveredAt"`
}

type EmailMessage struct {
	ToUser  int    `json:"toUser"`
	Type    string `json:"type"`
	Message []byte `json:"message"`
}

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdatedTime"`
	First_name      string    `json:"firstName"`
	Name            string    `json:"name"`
	Gender          string    `json:"gender"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Street          string    `json:"street"`
	HouseNumber     string    `json:"houseNumber"`
	TelephoneNumber string    `json:"telephoneNumber"`
	ProfilPicture   []byte    `json:"profilePicture"`
	Confirmed       bool      `json:"confirmed"`
	Active          bool      `json:"active"`
	Password        string    `json:"password"`
	SearchedSkills  []*Skill  `json:"searchedSkills" gorm:"many2many:user_searchedSkills"`
	AchievedSkills  []*Skill  `json:"achievedSkills" gorm:"many2many:user_achievedSkills"`
	CityIdentifier  int
	City            *City `json:"city" gorm:"foreignKey:CityIdentifier"`
}

type City struct {
	PLZ       uint      `json:"plz" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdatedTime"`
	Place     string    `json:"place"`
}

type Skill struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdatedTime"`
	Name           string    `json:"name"`
	Level          string    `json:"level"`
	UsersSearching []*User   `json:"usersSearching" gorm:"many2many:user_searchedSkills"`
	UsersAchieved  []*User   `json:"usersAchieved" gorm:"many2many:user_achievedSkills"`
}
