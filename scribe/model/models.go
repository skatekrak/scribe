package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" swaggertype:"string"`
}

type Lang struct {
	IsoCode   string         `gorm:"primaryKey" json:"isoCode"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" swaggertype:"string"`
	ImageURL  string         `json:"imageUrl"`

	Sources []Source `json:"-"`
} // @name Lang

type Source struct {
	Model

	RefreshedAt *time.Time `json:"refreshedAt"`
	Order       int        `gorm:"index" json:"order"`
	SourceType  string     `json:"sourceType"`
	LangIsoCode string     `json:"-"`
	Lang        Lang       `json:"lang"`
	Title       string     `json:"title"`
	ShortTitle  string     `json:"shortTitle"`
	IconURL     string     `json:"iconUrl"`
	CoverURL    string     `json:"coverUrl"`
	Description string     `json:"description"`
	SkateSource bool       `gorm:"default:true" json:"skateSource"`
	WebsiteURL  string     `json:"websiteUrl"`
	PublishedAt *time.Time `json:"publishedAt"`
	SourceID    string     `gorm:"unique,index" json:"sourceId"` // Vimeo, Youtube or Feedly ID, depending on the type

	Contents []Content `json:"-"`
} // @name Source

type Content struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" swaggertype:"string"`

	SourceID uint   `json:"-"`
	Source   Source `json:"source"`

	ContentID    string    `gorm:"uniqueIndex" json:"contentId"` // Youtube or Vimeo ID or Feedly ID
	PublishedAt  time.Time `json:"publishedAt"`
	Title        string    `json:"title"`
	ContentURL   string    `json:"contentUrl"` // Youtube or Vimeo video url or article URL
	ThumbnailURL string    `json:"thumbnailUrl"`
	RawSummary   string    `json:"rawSummary"`
	Summary      string    `json:"summary"`
	RawContent   string    `json:"rawContent"`
	Content      string    `json:"content"`
	Author       *string   `json:"author"` // For feedly article
	Type         string    `json:"type"`
} // @name Content

func (c *Content) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}

type Config struct {
	Key       string         `gorm:"primaryKey" json:"key"`
	Value     sql.NullString `json:"value"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updatedAt"`
}
