package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductRequestModel struct {
	ID           string `json:"id"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type ProductReviewRequestModel struct {
	ReviewerName      string  `json:"name" binding:"required"`
	ReviewDescription string  `json:"review"`
	ReviewRating      float64 `json:"rating" binding:"required,gte=0,lte=5"`
}

type ProductPaginationResponse struct {
	Total    int64                  `json:"total"`
	Page     int64                  `json:"page"`
	Limit    int64                  `json:"limit"`
	LastPage int64                  `json:"lastPage"`
	Products []ProductResponseModel `json:"data"`
}

type ProductResponseModel struct {
	ID            primitive.ObjectID `bson:"_id" json:"id" binding:"required"`
	Name          string             `bson:"name" json:"name" binding:"required"`
	Description   string             `bson:"description" json:"description" binding:"required"`
	ThumbnailURL  string             `bson:"thumbnail_url" json:"thumbnailUrl"`
	AverageRating float64            `bson:"average_rating" json:"averageRating"`
}

type Product struct {
	ID           primitive.ObjectID `bson:"_id" json:"id" binding:"required"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Description  string             `bson:"description" json:"description" binding:"required"`
	ThumbnailURL string             `bson:"thumbnail_url" json:"thumbnailUrl"`
}

type ProductReview struct {
	ID                primitive.ObjectID `bson:"_id" json:"id" binding:"required"`
	ProductID         primitive.ObjectID `bson:"productID" json:"productId" binding:"required"`
	ReviewerName      string             `bson:"name" json:"name" binding:"required"`
	ReviewDescription string             `bson:"review" json:"review"`
	ReviewRating      float64            `bson:"rating" json:"rating" binding:"required,gte=0,lte=5"`
}
