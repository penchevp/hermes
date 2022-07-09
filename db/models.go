package db

import (
	"github.com/google/uuid"
	"time"
)

type NotificationChannel struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string
}

type Customer struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string
}

type CustomerNotificationChannels struct {
	ID                           int32 `gorm:"primarykey"`
	CustomerID                   string
	Customer                     Customer
	NotificationChannelTypeID    string
	NotificationChannelType      NotificationChannel
	NotificationChannelLookupKey string
}

type Notification struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	From      string
	Text      string
}
