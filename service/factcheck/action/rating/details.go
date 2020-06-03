package rating

import (
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util/render"
	"github.com/factly/dega-server/validation"
	"github.com/go-chi/chi"
)

// details - Get rating by id
// @Summary Show a rating by id
// @Description Get rating by ID
// @Tags Rating
// @ID get-rating-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param space_id path string true "Space ID"
// @Param rating_id path string true "Claimant ID"
// @Success 200 {object} model.Claimant
// @Router /{space_id}/factcheck/ratings/{rating_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	ratingID := chi.URLParam(r, "rating_id")
	id, err := strconv.Atoi(ratingID)

	spaceID := chi.URLParam(r, "space_id")
	sid, err := strconv.Atoi(spaceID)

	if err != nil {
		return
	}

	result := &model.Rating{}

	result.ID = uint(id)

	err = config.DB.Model(&model.Rating{}).Preload("Medium").Where(&model.Rating{
		SpaceID: uint(sid),
	}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, result)
}
