package models

type Product struct {
	ID           string          `json:"id" binding:"required"`
	Name         string          `json:"name" binding:"required"`
	Description  string          `json:"description" binding:"required"`
	ThumbnailURL string          `json:"thumbnailUrl"`
	Reviews      []ProductReview `json:"reviews" binding:"required"`
}

type ProductReview struct {
	ReviewerName      string `json:"name" binding:"required"`
	ReviewDescription string `json:"review"`
	ReviewRating      int    `json:"rating" binding:"required,gte=0,lte=5"`
}
