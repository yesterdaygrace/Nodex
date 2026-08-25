package models

import (
	"time"

	"gorm.io/gorm"
)

// Post is a blog post row.
//
// CreatedAt/UpdatedAt give the UI a "last-updated" timestamp and a stable
// newest-first ordering for pagination. Status is "published", "draft" or
// "archived" (defaults to "published", §34 Archive). DeletedAt enables
// soft-delete so that DELETE moves a post to Trash instead of destroying it;
// GORM's default scope then hides trashed rows. Folder/Tags support Nodex
// sidebar organization (§11-12): folder is a single folder name, tags is a
// comma-separated list (e.g. "go,backend").
type Post struct {
	Id        int            `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Status    string         `json:"status" gorm:"default:published"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Folder    string         `json:"folder"`
	Tags      string         `json:"tags"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
