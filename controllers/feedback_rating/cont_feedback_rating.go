package contfeedbackrating

import (
	"context"
	"fmt"
	"net/http"
	ifeedbackrating "nuryanto2121/cukur_in_user/interface/feedback_rating"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	"nuryanto2121/cukur_in_user/pkg/logging"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ContFeedbackRating struct {
	useFeedbackRating ifeedbackrating.Usecase
}

func NewContFeedbackRating(e *echo.Echo, a ifeedbackrating.Usecase) {
	controller := &ContFeedbackRating{
		useFeedbackRating: a,
	}

	r := e.Group("/user/feedback_rating")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.POST("", controller.Create)
	r.PUT("/:id", controller.Update)
}

// CreateSaFeedbackRating :
// @Summary Add FeedbackRating
// @Security ApiKeyAuth
// @Tags FeedbackRating
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.AddFeedbackRating true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/feedback_rating [post]
func (u *ContFeedbackRating) Create(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		form   models.AddFeedbackRating
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	err = u.useFeedbackRating.Create(ctx, claims, &form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateSaFeedbackRating :
// @Summary Rubah FeedbackRating
// @Security ApiKeyAuth
// @Tags FeedbackRating
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Param req body models.AddFeedbackRating true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/feedback_rating/{id} [put]
func (u *ContFeedbackRating) Update(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		err    error
		// valid  validation.Validation                 // wajib
		id   = e.Param("id") //kalo bukan int => 0
		form = models.AddFeedbackRating{}
	)

	FeedbackRatingID, _ := strconv.Atoi(id)
	// logger.Info(id)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	// form.UpdatedBy = claims.FeedbackRatingName
	err = u.useFeedbackRating.Update(ctx, claims, FeedbackRatingID, form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusCreated, "Ok", nil)
}
