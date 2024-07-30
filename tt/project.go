package tt

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name     string `gorm:"index:idx_project4customer,unique"`
	ClientID uint   `gorm:"index:idx_project4customer,unique"`
	Client   Client
	Active   bool
}

type ProjectList []Project

func projectInit() {
	if err := db.AutoMigrate(&Project{}); err == nil && db.Migrator().HasTable(&Project{}) {
		if err := db.First(&Project{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c := clientNew()
			c.getById(1)

			p := projectNew()
			p.Name = "Default"
			p.Active = true
			p.Client = c
			result := db.Create(&p)

			if result.Error != nil {
				log.Fatalf(ErrorString, CharError, result.Error)
			}
		}
	}
}

func projectNew() Project {
	return Project{}
}

func projectListNew() ProjectList {
	return ProjectList{}
}

func (projectList *ProjectList) listAll() {
	result := db.Preload("Client").Find(&projectList)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}
func (projectList *ProjectList) listByName(pattern string) int64 {
	like := "%" + pattern + "%"
	result := db.Preload("Client").Where("name like ?", like).Find(&projectList)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
	return (result.RowsAffected)
}

func (projectList *ProjectList) listByClient(clientId uint) {
	result := db.Where("ClientId = ?", clientId).Preload("Client").Find(&projectList)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}

}

func (projectList *ProjectList) listByClientAndName(clientId uint, pattern string) int64 {
	like := "%" + pattern + "%"
	result := db.Preload("Client").Where("ClientId = ?", clientId).Where("name like ?", like).Find(&projectList)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
	return (result.RowsAffected)
}

func (project *Project) create(name string, clientId uint) {
	project.Name = name
	client := clientNew()
	client.getById(clientId)
	project.Client = client
	project.Active = true
	result := db.Create(&project)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (project *Project) getById(id uint) {
	result := db.Preload("Client").First(&project, id)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (project *Project) modify(name string) {
	project.Name = name
	result := db.Save(&project)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (project *Project) activate() {
	project.Active = true
	result := db.Save(&project)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}

func (project *Project) deactivate() {
	project.Active = false
	result := db.Save(&project)

	if result.Error != nil {
		log.Fatalf(ErrorString, CharError, result.Error)
	}
}
