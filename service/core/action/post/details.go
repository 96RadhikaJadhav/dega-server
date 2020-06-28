package post

import (
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get post by id
// @Summary Show a post by id
// @Description Get post by ID
// @Tags Post
// @ID get-post-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param post_id path string true "Post ID"
// @Success 200 {object} postData
// @Router /core/posts/{post_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		return
	}

	postID := chi.URLParam(r, "post_id")
	id, err := strconv.Atoi(postID)

	if err != nil {
		return
	}

	result := &postData{}
	result.Categories = make([]model.Category, 0)
	result.Tags = make([]model.Tag, 0)

	categories := []model.PostCategory{}
	tags := []model.PostTag{}

	result.ID = uint(id)

	err = config.DB.Model(&model.Post{}).Preload("Medium").Preload("Format").Where(&model.Post{
		SpaceID: uint(sID),
	}).First(&result.Post).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	// fetch all categories
	config.DB.Model(&model.PostCategory{}).Where(&model.PostCategory{
		PostID: uint(id),
	}).Preload("Category").Preload("Category.Medium").Find(&categories)

	// fetch all tags
	config.DB.Model(&model.PostTag{}).Where(&model.PostTag{
		PostID: uint(id),
	}).Preload("Tag").Find(&tags)

	for _, c := range categories {
		if c.Category.ID != 0 {
			result.Categories = append(result.Categories, c.Category)
		}
	}

	for _, t := range tags {
		if t.Tag.ID != 0 {
			result.Tags = append(result.Tags, t.Tag)
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}