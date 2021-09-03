package posts

import (
	"context"
)

type Service interface {
	getPost(ctx context.Context, req getPostsRequest) (res getPostsResponse, err error)
	createPost(ctx context.Context, req createPostRequest) (res createPostResponse, err error)

	getPostComment(ctx context.Context, req getPostCommentRequest) (res getPostCommentResponse, err error)
	createPostComment(ctx context.Context, req createPostCommentRequest) (res createPostCommentResponse, err error)

	togglePostLike(ctx context.Context, req togglePostLikeRequest) (res togglePostLikeResponse, err error)
	toggleCommentLike(ctx context.Context, req toggleCommentLikeRequest) (res toggleCommentLikeResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return service{repo: r}
}

func (s service) getPost(ctx context.Context, req getPostsRequest) (res getPostsResponse, err error) {
	p, err := s.repo.GetPost(ctx, req.PostID)
	if err != nil {
		res.Err = err
		return
	}

	res.Post = p
	return
}

func (s service) createPost(ctx context.Context, req createPostRequest) (res createPostResponse, err error) {
	images := []string{}
	videos := []string{}

	for i := range req.Images {
		id, err := s.repo.PostFile(ctx, req.Images[i])
		if err != nil {
			return res, err
		}
		images = append(images, id)
		req.Images[i].Close()
	}

	for i := range req.Videos {
		id, err := s.repo.PostFile(ctx, req.Videos[i])
		if err != nil {
			return res, err
		}
		videos = append(videos, id)
		req.Videos[i].Close()
	}

	id, t, err := s.repo.CreatePost(ctx, req.UserID, req.Text, string(req.Visibility), images, videos)

	res.Post = &Post{
		ID:       id,
		UserID:   req.UserID,
		Text:     req.Text,
		Images:   images,
		Videos:   videos,
		Likes:    []string{req.UserID},
		Comments: []string{},
		Created:  t.String(),
	}

	return
}

func (s service) getPostComment(ctx context.Context, req getPostCommentRequest) (res getPostCommentResponse, err error) {
	c, err := s.repo.GetPostComment(ctx, req.CommentID)
	if err != nil {
		res.Err = err
		return
	}

	res.Comment = c
	return
}

func (s service) createPostComment(ctx context.Context, req createPostCommentRequest) (res createPostCommentResponse, err error) {
	id, t, err := s.repo.CreatePostComment(ctx, req.PostID, req.UserID, req.Text)
	if err != nil {
		res.Err = err
		return
	}

	res.Comment = &Comment{
		ID:      id,
		PostID:  req.PostID,
		UserID:  req.UserID,
		Text:    req.Text,
		Likes:   []string{req.UserID},
		Created: t.String(),
	}
	return res, nil
}

func (s service) togglePostLike(ctx context.Context, req togglePostLikeRequest) (res togglePostLikeResponse, err error) {
	err = s.repo.TogglePostLike(ctx, req.PostID, req.UserID)
	if err != nil {
		res.Err = err
	}

	return
}

func (s service) toggleCommentLike(ctx context.Context, req toggleCommentLikeRequest) (res toggleCommentLikeResponse, err error) {
	err = s.repo.ToggleCommentLike(ctx, req.CommentID, req.UserID)
	if err != nil {
		res.Err = err
	}

	return
}
