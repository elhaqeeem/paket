package handlers

import (
	"net/http"

	"github.com/elhaqeeem/paket/internal/errors"
	"github.com/elhaqeeem/paket/internal/postms/models"
	"github.com/elhaqeeem/paket/internal/postms/services"
	"github.com/gin-gonic/gin"
)

func handleServiceError(err error, c *gin.Context) {
	switch err.(type) {
	case *errors.NotFound:
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "message": err.Error()})
	case *errors.DeleteIsMissingID:
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
	case *errors.CursorDecodingError:
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
	default:
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "message": "An internal server error occurred"})
	}
}

func getPostServiceFromContext(c *gin.Context) services.PostService {
	value, _ := c.Get("postService")
	return value.(services.PostService)
}

func NotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
}

func CreatePost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	post := &models.Post{}
	err := c.ShouldBindJSON(post)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	err = postService.CreatePost(post)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	post := &models.Post{}
	err := c.ShouldBindJSON(post)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	var existingPost models.Post
	existingPost, err = postService.GetPost(postID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	post.ID = uint(postID)
	post.CreatedAt = existingPost.CreatedAt
	err = postService.UpdatePost(post)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, post)
}

func CreatePostComment(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postComment := &models.PostComment{}
	err := c.ShouldBindJSON(postComment)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if _, err := postService.GetPost(uint64(postComment.PostID)); err != nil {
		// TODO: Specifically check err is a NotFound error
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": "Can not create comment for non-existent post"})
		return
	}
	err = postService.CreatePostComment(postComment)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusCreated, postComment)
}

func UpdatePostComment(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postCommentID := uint64(c.GetInt64("ID"))
	postComment := &models.PostComment{}
	err := c.ShouldBindJSON(postComment)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	var existingPostComment models.PostComment
	existingPostComment, err = postService.GetPostComment(postCommentID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	existingPostComment.Body = postComment.Body
	err = postService.UpdatePostComment(&existingPostComment)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, existingPostComment)
}

func GetPostCommentsForPost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	postComments, err := postService.GetPostCommentsForPost(postID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, postComments)
}

func GetPost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	post, err := postService.GetPost(postID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, post)
}

func GetPosts(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	cursor := c.Query("cursor")
	userID := c.Query("userId")
	tag := c.Query("tag")
	posts, nextCursor, err := postService.GetPosts(cursor, userID, tag)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"nextCursor": nextCursor, "results": posts})
}

func DeletePost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	post, err := postService.GetPost(postID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	if err := postService.DeletePost(&post); err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetPostComment(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postCommentID := uint64(c.GetInt64("ID"))
	postComment, err := postService.GetPostComment(postCommentID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, postComment)
}

func DeletePostComment(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postCommentID := uint64(c.GetInt64("ID"))
	postComment, err := postService.GetPostComment(postCommentID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	if err := postService.DeletePostComment(&postComment); err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func CreatePostVote(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postVote := &models.PostVote{}
	err := c.ShouldBindJSON(postVote)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if postVote.Value >= 0 {
		postVote.Value = 1
	} else if postVote.Value < 0 {
		postVote.Value = -1
	}
	// Check the post actually exists
	postExists, err := postService.PostExists(uint64(postVote.PostID))
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if !postExists {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": "Post matching id does not exist"})
		return
	}
	// If a vote already exists, return it
	existingPostVote, _ := postService.GetPostVote(uint64(postVote.PostID), postVote.UserID)
	if existingPostVote.UserID != "" {
		c.JSON(http.StatusOK, existingPostVote)
		return
	}
	err = postService.CreatePostVote(postVote)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusCreated, postVote)
}

func GetPostVoteTotalForPost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	total := postService.GetPostVoteTotalForPost(postID)
	c.JSON(http.StatusOK, gin.H{"total": total})
}

func GetPostVoteUsersForPost(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	userIDs := postService.GetPostVoteUsersForPost(postID)
	c.JSON(http.StatusOK, &userIDs)
}

func CreatePostSave(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postSave := models.PostSave{}
	err := c.ShouldBindJSON(&postSave)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	// Check the post actually exists
	postExists, err := postService.PostExists(uint64(postSave.PostID))
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if !postExists {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": "Post matching id does not exist"})
		return
	}
	var isNew bool
	postSave, isNew, err = postService.CreatePostSave(&postSave)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	var status int
	if isNew {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}
	c.JSON(status, postSave)
}

func GetPostSaves(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postID := uint64(c.GetInt64("ID"))
	postSaves, err := postService.GetPostSaves(postID, "")
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusCreated, postSaves)
}

func DeletePostSave(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	postSaveID := uint64(c.GetInt64("ID"))
	postSave, err := postService.GetPostSave(postSaveID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	if err := postService.DeletePostSave(&postSave); err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetTags(c *gin.Context) {
	postService := getPostServiceFromContext(c)
	tags, err := postService.GetTags()
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, tags)
}
