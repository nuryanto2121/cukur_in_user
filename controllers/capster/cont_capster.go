package contcapster

import (
	"context"
	"fmt"
	"net/http"
	icapsters "nuryanto2121/cukur_in_user/interface/capster"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	"nuryanto2121/cukur_in_user/pkg/logging"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ContCapster struct {
	useCapster icapsters.Usecase
}

func NewContCapster(e *echo.Echo, a icapsters.Usecase) {
	controller := &ContCapster{
		useCapster: a,
	}

	r := e.Group("/user/capster")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags Capster
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Param latitude query number true "Latitude"
// @Param longitude query number true "Longitude"
// @Success 200 {object} tool.ResponseModel
// @Router /user-service/user/capster/{id} [get]
func (u *ContCapster) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger    = logging.Logger{}
		appE      = tool.Res{R: e} // wajib
		id        = e.Param("id")  //kalo bukan int => 0
		GeoBarber = models.GeoBarber{}
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	httpCode, errMsg := app.BindAndValid(e, &GeoBarber)
	logger.Info(util.Stringify(GeoBarber))
	if httpCode != 200 {
		return appE.Response(http.StatusBadRequest, errMsg, nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	data, err := u.useCapster.GetDataBy(ctx, claims, ID, GeoBarber)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Capster
// @Security ApiKeyAuth
// @Tags Capster
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param latitude query number true "Latitude"
// @Param longitude query number true "Longitude"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /user-service/user/capster [get]
func (u *ContCapster) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger       = logging.Logger{}
		appE         = tool.Res{R: e}        // wajib
		paramquery   = models.ParamListGeo{} // ini untuk list
		responseList = models.ResponseModelList{}
		err          error
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseErrorList(http.StatusBadRequest, fmt.Sprintf("%v", err), responseList)
	}

	responseList, err = u.useCapster.GetList(ctx, claims, paramquery)
	if err != nil {
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	return appE.ResponseList(http.StatusOK, "", responseList)
}
