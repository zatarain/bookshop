package models

import "gorm.io/gorm"

type DataAccessInterface interface {
	Create(interface{}) *gorm.DB
}
