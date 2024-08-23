package models

import (
	"time"

	"github.com/elhaqeeem/paket/internal/utils"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type CommonFields struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

type Post struct {
	CommonFields
	UserID string         `json:"userId" binding:"required"`
	Title  string         `json:"title" binding:"required"`
	Slug   string         `json:"slug" gorm:"unique"`
	Body   string         `json:"body" binding:"required"`
	Tags   pq.StringArray `json:"tags" gorm:"type:varchar(64)[]"`
}

// SetSlugAndTags sets the Slug and Tags fields for Post
func (p *Post) SetSlugAndTags() {
	p.Slug = slug.Make(p.Title)
	p.Tags = utils.ToTagSlice(p.Tags)
}

// BeforeSave is a GORM hook that runs before saving a Post
func (p *Post) BeforeSave(tx *gorm.DB) (err error) {
	p.SetSlugAndTags()
	return
}

type PostComment struct {
	CommonFields
	UserID string `json:"userId" binding:"required"`
	PostID uint   `json:"postId" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

type PostVote struct {
	CommonFields
	UserID string `json:"userId" binding:"required" gorm:"primaryKey"`
	PostID uint   `json:"postId" binding:"required" gorm:"primaryKey"`
	Value  int    `json:"value" binding:"required"`
}

type PostSave struct {
	CommonFields
	UserID string `json:"userId" binding:"required" gorm:"primaryKey"`
	PostID uint   `json:"postId" binding:"required" gorm:"primaryKey"`
}
