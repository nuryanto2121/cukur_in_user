package contadvertise

import (
	"context"
	"fmt"
	"net/http"
	iadvertise "nuryanto2121/cukur_in_user/interface/advertise"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	"nuryanto2121/cukur_in_user/pkg/logging"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	"strconv"

	"github.com/labstack/echo/v4"
)

type contAdvertise struct {
	useAdvertise iadvertise.Usecase
}

func NewContAdvertise(e *echo.Echo, a iadvertise.Usecase) {
	controller := &contAdvertise{
		useAdvertise: a,
	}

	r := e.Group("/user/advertise")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)

}

// GetDataByID :
// @Summary GetById Advertise
// @Security ApiKeyAuth
// @Tags Advertise
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Success 200 {object} models.ResponseModelList
// @Router /user/advertise/{id} [get]
func (u *contAdvertise) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = tool.Res{R: e} // wajib
		id     = e.Param("id")  //kalo bukan int => 0
	)
	ID, err := strconv.Atoi(id)
	logger.Info(ID)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	data, err := u.useAdvertise.GetDataBy(ctx, claims, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Advertise
// @Security ApiKeyAuth
// @Tags Advertise
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /user/advertise [get]
func (u *contAdvertise) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE         = tool.Res{R: e}     // wajib
		paramquery   = models.ParamList{} // ini untuk list
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

	responseList, err = u.useAdvertise.GetList(ctx, claims, paramquery)
	if err != nil {
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	return appE.ResponseList(http.StatusOK, "", responseList)
}
