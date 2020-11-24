package contbarberpaket

import (
	"context"
	"fmt"
	"net/http"
	ibarberpaket "nuryanto2121/cukur_in_user/interface/barber_paket"
	midd "nuryanto2121/cukur_in_user/middleware"
	app "nuryanto2121/cukur_in_user/pkg"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ContBarberPaket struct {
	userBarberPaket ibarberpaket.Usecase
}

func NewContBarberPaket(e *echo.Echo, a ibarberpaket.Usecase) {
	controller := &ContBarberPaket{
		userBarberPaket: a,
	}

	r := e.Group("/user/barber_paket")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.GET("/:id", controller.GetDataBy)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags Barber Paket
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID Barber"
// @Success 200 {object} tool.ResponseModel
// @Router /user/barber_paket/{id} [get]
func (u *ContBarberPaket) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE = tool.Res{R: e} // wajib
		id   = e.Param("id")  //kalo bukan int => 0
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	data, err := u.userBarberPaket.GetDataBy(ctx, claims, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}
