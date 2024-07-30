package tt

import (
	"errors"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name     string `gorm:"unique;index"`
	Active   bool
	Projects []Project
}

type ClientList []Client

func clientInit() {
	if err := db.AutoMigrate(&Client{}); err == nil && db.Migrator().HasTable(&Client{}) {
		if err := db.First(&Client{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c := clientNew()
			c.Name = "Default"
			c.Active = true
			result := db.Create(&c)

			if result.Error != nil {
				log.Fatalf(ErrorString, CharError, result.Error)
			}
		}
	}
}

func clientNew() Client {
	return Client{}
}

func clientListNew() ClientList {
	return ClientList{}
}

func (clientList *ClientList) listAll() {
	result := db.Find(&clientList)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (clientList *ClientList) listByName(pattern string) int64 {
	like := "%" + pattern + "%"
	result := db.Where("name like ?", like).Find(&clientList)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
	return (result.RowsAffected)
}

func (client *Client) create(name string) {
	client.Name = name
	client.Active = true
	result := db.Create(&client)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (client *Client) getById(id uint) {
	result := db.First(&client, id)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (client *Client) getByIdOrName(pattern string) {
	if digitCheck.MatchString(pattern) {
		clientId, _ := strconv.ParseUint(pattern, 10, 32)
		client.getById(uint(clientId))
	} else {
		cl := clientListNew()
		rows := cl.listByName(pattern)
		if rows != 1 {
			log.Fatalf(ErrorString, CharError, errors.New(errorUnambiguously))
		}
		client.getById(cl[0].ID)
	}
}

func (client *Client) modify(name string) {
	client.Name = name
	result := db.Save(&client)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (client *Client) activate() {
	client.Active = true
	result := db.Save(&client)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (client *Client) deactivate() {
	client.Active = false
	result := db.Save(&client)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}
