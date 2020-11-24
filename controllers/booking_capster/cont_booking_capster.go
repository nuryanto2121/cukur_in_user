package contbookingbarber

import (
	"context"
	"fmt"
	"net/http"
	ibookingcapster "nuryanto2121/cukur_in_user/interface/booking_capster"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	"nuryanto2121/cukur_in_user/pkg/logging"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	util "nuryanto2121/cukur_in_user/pkg/utils"

	"github.com/labstack/echo/v4"
)

type ContBookingCapster struct {
	useBookingCapster ibookingcapster.Usecase
}

func NewContBookingCapster(e *echo.Echo, a ibookingcapster.Usecase) {
	controller := &ContBookingCapster{
		useBookingCapster: a,
	}

	r := e.Group("/user/booking_capster")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.POST("", controller.Create)
	r.GET("", controller.GetDataBy)
}

// CreateBookingCapster :
// @Summary Add Booking Capster
// @Security ApiKeyAuth
// @Tags BookingCapster
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.AddBookingCapster true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/booking_capster [post]
func (u *ContBookingCapster) Create(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		form   models.AddBookingCapster
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

	err = u.useBookingCapster.Create(ctx, claims, &form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", nil)
}

// GetDataByBookingCapster :
// @Summary GetData Booking Capster
// @Security ApiKeyAuth
// @Tags BookingCapster
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "Version Device"
// @Param req query models.AddBookingCapster true "ss"
// @Success 200 {object} models.ResponseModelList
// @Router /user/booking_capster [get]
func (u *ContBookingCapster) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger     = logging.Logger{}
		appE       = tool.Res{R: e} // wajib
		paramquery = models.AddBookingCapster{}
		err        error
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", errMsg), nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	data, err := u.useBookingCapster.GetDataBy(ctx, claims, paramquery)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}
