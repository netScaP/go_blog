package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Tag ...
type Tag struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name     string    `gorm:"size:255;not null;unique" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare ...
func (p *Tag) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate ...
func (p *Tag) Validate() error {

	if p.Name == "" {
		return errors.New("Required Name")
	}

	return nil
}

// SaveTag ...
func (p *Tag) SaveTag(db *gorm.DB) (*Tag, error) {
	var err error
	err = db.Debug().Model(&Tag{}).Create(&p).Error
	if err != nil {
		return &Tag{}, err
	}
	return p, nil
}

// TagPaginateStruct ...
type TagPaginateStruct struct {
	Total uint64 `json:"total"`
	Limit uint64 `json:"limit"`
	Skip  uint64 `json:"skip"`
	Data  []Tag `json:"data"`
}
// TagQueryStruct ...
type TagQueryStruct struct {
	Limit uint64
	Skip  uint64
}

// FindTags ...
func (p *Tag) FindTags(db *gorm.DB, query *TagQueryStruct) (TagPaginateStruct, error) {
	response := TagPaginateStruct{}

	if query.Skip < 0 {
		response.Skip = 0
	} else {
		response.Skip = query.Skip
	}

	if query.Limit < 0 {
		response.Limit = 0
	} else if query.Limit > 50 {
		response.Limit = 50
	} else {
		response.Limit = query.Limit
	}

	dbChain := db.Debug().Model(&Tag{}).Where("id > ?", 0)
	dbChain.Count(&response.Total)
	tags := []Tag{}
	err := dbChain.Limit(query.Limit).Offset(query.Skip).Find(&tags).Error
	if err != nil {
		return response, err
	}
	response.Data = tags

	return response, nil
}

// FindAllTags ...
func (p *Tag) FindAllTags(db *gorm.DB) (*[]Tag, error) {
	var err error
	tags := []Tag{}
	err = db.Debug().Model(&Tag{}).Limit(100).Find(&tags).Error
	if err != nil {
		return &[]Tag{}, err
	}
	return &tags, nil
}

// FindTagByID ...
func (p *Tag) FindTagByID(db *gorm.DB, pid uint64) (*Tag, error) {
	var err error
	err = db.Debug().Model(&Tag{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Tag{}, err
	}
	return p, nil
}

// UpdateTag ...
func (p *Tag) UpdateTag(db *gorm.DB) (*Tag, error) {

	var err error
	err = db.Debug().Model(&Tag{}).Where("id = ?", p.ID).Updates(Tag{Name: p.Name, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Tag{}, err
	}
	return p, nil
}

// DeleteTag ...
func (p *Tag) DeleteTag(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Tag{}).Where("id = ? and author_id = ?", pid, uid).Take(&Tag{}).Delete(&Tag{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Tag not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
