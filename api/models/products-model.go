package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductRequestModel struct {
	ID           string `json:"id"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type ProductReviewRequestModel struct {
	ReviewerName      string `json:"name" binding:"required"`
	ReviewDescription string `json:"review"`
	ReviewRating      int    `json:"rating" binding:"required,gte=0,lte=5"`
}

type Product struct {
	ID           primitive.ObjectID `bson:"id" json:"id" binding:"required"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Description  string             `bson:"description" json:"description" binding:"required"`
	ThumbnailURL string             `bson:"thumbnail_url" json:"thumbnailUrl"`
	Reviews      []ProductReview    `bson:"reviews" json:"reviews" binding:"required"`
}

type ProductReview struct {
	ReviewerName      string `bson:"name" json:"name" binding:"required"`
	ReviewDescription string `bson:"review" json:"review"`
	ReviewRating      int    `bson:"rating" json:"rating" binding:"required,gte=0,lte=5"`
}
