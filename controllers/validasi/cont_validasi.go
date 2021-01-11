package contvalidasi

import (
	"context"
	"fmt"
	"net/http"
	midd "nuryanto2121/cukur_in_user/middleware"
	app "nuryanto2121/cukur_in_user/pkg"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	repofunction "nuryanto2121/cukur_in_user/repository/function"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ContValidasi struct {
}

func NewContValidasi(e *echo.Echo) {
	controller := &ContValidasi{}

	r := e.Group("/user/validasi")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.GET("/:id", controller.GetDataBy)
}

// GetDataByID :
// @Summary Validasi utk jam oprasional sudah lewat atau belum buka, Param
// @Security ApiKeyAuth
// @Tags Validasi
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID Barber"
// @Success 200 {object} tool.ResponseModel
// @Router /user/validasi/{id} [get]
func (u *ContValidasi) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE = tool.Res{R: e} // wajib
		id   = e.Param("id")  //kalo bukan int => 0
		now  = util.GetTimeNow()
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	fn := &repofunction.FN{
		Claims: claims,
	}

	dataBarber, err := fn.GetBarberData(ID)

	if !fn.InTimeActiveBarber(dataBarber, now) {
		return appE.Response(http.StatusBadRequest, "Mohon maaf , waktu di luar jam oprasional", nil)
	}

	return appE.Response(http.StatusOK, "Ok", nil)
}
