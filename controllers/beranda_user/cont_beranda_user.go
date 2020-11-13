package contberandauser

import (
	"context"
	"fmt"
	"net/http"
	iberandauser "nuryanto2121/cukur_in_user/interface/beranda_user"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	"nuryanto2121/cukur_in_user/pkg/logging"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	util "nuryanto2121/cukur_in_user/pkg/utils"

	"github.com/labstack/echo/v4"
)

type contBerandaUser struct {
	useBerandaUser iberandauser.Usecase
}

func NewContBerandaUser(e *echo.Echo, a iberandauser.Usecase) {
	controller := &contBerandaUser{
		useBerandaUser: a,
	}

	r := e.Group("/user/beranda")
	r.Use(midd.JWT)
	// r.GET("/status_order", controller.GetStatusOrder)
	r.GET("", controller.GetBeranda)
}

// GetBeranda :
// @Summary Get Data Beranda
// @Security ApiKeyAuth
// @Tags User Beranda
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "Version Device"
// @Param latitude number true "Latitude"
// @Param longitude number true "Longitude"
// @Success 200 {object} models.ResponseModelList
// @Router /user/beranda [get]
func (c *contBerandaUser) GetBeranda(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger     = logging.Logger{}
		appE       = tool.Res{R: e} // wajib
		paramquery = models.GeoBarber{}
		paramList  = models.ParamListGeo{} // ini untuk list
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	paramList.GeoBarber = paramquery

	ClosestBarber, err := c.useBerandaUser.GetClosestBarber(ctx, claims, paramList)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), ClosestBarber)
	}
	CapsterRecommend, err := c.useBerandaUser.GetRecomentCapster(ctx, claims, paramList)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), CapsterRecommend)
	}
	BarberRecommend, err := c.useBerandaUser.GetRecomentBarber(ctx, claims, paramList)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), BarberRecommend)
	}
	result := map[string]interface{}{
		"closest_barber":         ClosestBarber,
		"capster_recommendation": CapsterRecommend,
		"barber_recommendation":  BarberRecommend,
	}

	// return e.JSON(http.StatusOK, ListBarbersPost)
	return appE.Response(http.StatusOK, "", result)
}
