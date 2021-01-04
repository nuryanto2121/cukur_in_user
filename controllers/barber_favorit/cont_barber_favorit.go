package contbarberfavorit

import (
	"context"
	"fmt"
	"net/http"
	ibarberfavorit "nuryanto2121/cukur_in_user/interface/barber_favorit"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	"nuryanto2121/cukur_in_user/pkg/logging"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	util "nuryanto2121/cukur_in_user/pkg/utils"

	"github.com/labstack/echo/v4"
)

type contBarberFavorit struct {
	useFavorit ibarberfavorit.Usecase
}

func NewContBarberFavorit(e *echo.Echo, a ibarberfavorit.Usecase) {
	controller := &contBarberFavorit{
		useFavorit: a,
	}
	r := e.Group("/user/favorit")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	// r.GET("/status_order", controller.GetStatusOrder)
	r.GET("", controller.GetList)
	r.POST("", controller.Create)
}

// CreateBarberFavorit :
// @Summary Add Or Delete Barber Favorit
// @Security ApiKeyAuth
// @Tags Barber Favorit
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.AddBarberFavorit true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/favorit [post]
func (u *contBarberFavorit) Create(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		form   models.AddBarberFavorit
	)

	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	err = u.useFavorit.Create(ctx, claims, &form) //usePaket.Create(ctx, claims, &form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", nil)
}

// GetList :
// @Summary GetList Barber Favorit
// @Security ApiKeyAuth
// @Tags Barber Favorit
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "Version Device"
// @Param latitude query number true "Latitude"
// @Param longitude query number true "Longitude"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /user/favorit [get]
func (u *contBarberFavorit) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		//valid      validation.Validation // wajib
		// paramquery   = models.ParamList{} // ini untuk list
		paramquery   = models.ParamListGeo{} // ini untuk list
		responseList = models.ResponseModelList{}
		err          error
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	// logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseErrorList(http.StatusBadRequest, fmt.Sprintf("%v", err), responseList)
	}

	responseList, err = u.useFavorit.GetList(ctx, claims, paramquery)
	if err != nil {
		// return e.JSON(http.StatusBadRequest, err.Error())
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	// return e.JSON(http.StatusOK, ListDataPaket)
	return appE.ResponseList(http.StatusOK, "", responseList)
}
