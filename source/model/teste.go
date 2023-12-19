package model

import (
    "gorm.io/gorm"
    "time"
)

type Teste struct {
    ID        uint   `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Teste) TableName() string {
    return "teste"
 }
