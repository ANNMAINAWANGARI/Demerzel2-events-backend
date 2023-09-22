package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
	"log"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CreateNewComment(newComment *models.NewComment, userId string) (*models.Comment, error) {
	comment := models.Comment{
		Body:      newComment.Body,
		Images:    newComment.Images,
		EventId:   newComment.EventId,
		CreatorId: userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func UpdateCommentById(updateReq *models.UpdateComment, userId string) (*models.Comment, error) {
	var comment *models.Comment
	result := db.DB.Where("id = ?", updateReq.Id).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return comment, nil // Return nil when the user is not found
		}
		return comment, result.Error // Return the actual error for other errors
	}

	if comment.CreatorId != userId {
		return comment, errors.New("you are not authorized to update this comment")
	}

	comment.Body = updateReq.Body
	if err := db.DB.Save(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}

func GetComments(eventId string) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := db.DB.Where("event_id = ?", eventId).Preload("Creator").Find(&comments).Error
	if err != nil {
		log.Println("Error fetching comments from db")
		return comments, err
	}

	return comments, nil
}

func DeleteCommentById(commentId string, userId string) error {
	var comment models.Comment
	result := db.DB.Where("id = ?", commentId).First(&comment)
	fmt.Println("HEYYY", commentId, result)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // Return nil when the user is not found
		}
		return result.Error // Return the actual error for other errors
	}

	if comment.CreatorId != userId {
		return errors.New("you are not authorized to delete this comment")
	}

	if err := db.DB.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}
