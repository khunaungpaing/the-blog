package dto

import "github.com/khunaungpaing/the-blog-api/models"

type RequestPost struct {
	ID         uint              `json:"id"`
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	Slug       string            `json:"slug"`
	Status     string            `json:"status"`
	Categories []RequestCategory `json:"categories"`
	Tags       []RequestTag      `json:"tags"`
	Media      RequestMedia      `json:"media"`
}

func (rp RequestPost) ToModel(post models.Post) models.Post {
	post.ID = rp.ID
	post.Title = rp.Title
	post.Content = rp.Content
	post.Slug = rp.Slug
	post.Status = rp.Status
	post.Categories = RequestCategoryToModelList(rp.Categories)
	post.Tags = RequestTagToModelList(rp.Tags)
	if post.Media == nil {
		post.Media = new(models.Media)
	}
	post.Media = mediaToModel(post.Media, rp.Media)
	return post
}

type RequestCategory struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (rc RequestCategory) ToModel(category models.Category) models.Category {
	category.ID = rc.ID
	category.Name = rc.Name
	category.Description = rc.Description
	return category
}

type RequestTag struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (rt RequestTag) ToModel(tag models.Tag) models.Tag {
	tag.ID = rt.ID
	tag.Name = rt.Name
	tag.Description = rt.Description
	return tag
}

type RequestMedia struct {
	ID       uint   `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	MimeType string `json:"mime_type"`
}

func mediaToModel(media *models.Media, rm RequestMedia) *models.Media {
	media.Filename = rm.Filename
	media.Path = rm.Path
	media.MimeType = rm.MimeType
	return media
}

type RequestComment struct {
	Content string `json:"content"`
}

func RequestCategoryToModelList(list []RequestCategory) []models.Category {
	var categories []models.Category
	for _, rc := range list {
		categories = append(categories, rc.ToModel(models.Category{}))
	}
	return categories
}

func RequestTagToModelList(list []RequestTag) []models.Tag {
	var tags []models.Tag
	for _, rt := range list {
		tags = append(tags, rt.ToModel(models.Tag{}))
	}
	return tags
}
