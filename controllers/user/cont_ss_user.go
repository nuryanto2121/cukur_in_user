package contuser

import (
	"context"
	"fmt"
	"net/http"
	iusers "nuryanto2121/cukur_in_user/interface/user"
	midd "nuryanto2121/cukur_in_user/middleware"
	"nuryanto2121/cukur_in_user/models"
	app "nuryanto2121/cukur_in_user/pkg"
	tool "nuryanto2121/cukur_in_user/pkg/tools"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"

	_ "nuryanto2121/cukur_in_user/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type ContUser struct {
	useUser iusers.Usecase
}

func NewContUser(e *echo.Echo, a iusers.Usecase) {
	controller := &ContUser{
		useUser: a,
	}
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health_check", controller.HealthCheck)

	r := e.Group("/user/user")
	r.Use(midd.JWT)
	r.Use(midd.Versioning)
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)
	r.POST("", controller.Create)
	r.POST("/change_password", controller.ChangePassword)
	r.PUT("/:id", controller.Update)
	r.DELETE("/:id", controller.Delete)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags User
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Success 200 {object} tool.ResponseModel
// @Router /user/user/{id} [get]
func (u *ContUser) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		id   = e.Param("id")  //kalo bukan int => 0
		// valid  validation.Validation                 // wajib
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	data, err := u.useUser.GetDataBy(ctx, claims, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList User
// @Security ApiKeyAuth
// @Tags User
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /user/user [get]
func (u *ContUser) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		//valid      validation.Validation // wajib
		paramquery   = models.ParamList{} // ini untuk list
		responseList = models.ResponseModelList{}
		err          error
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	// logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
	}
	// claims, err := app.GetClaims(e)
	// if err != nil {
	// 	return appE.ResponseErrorList(http.StatusBadRequest, fmt.Sprintf("%v", err), responseList)
	// }
	// if !claims.IsAdmin {
	// 	paramquery.InitSearch = " id_created = " + strconv.Itoa(claims.UserID)
	// }
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseErrorList(http.StatusBadRequest, fmt.Sprintf("%v", err), responseList)
	}
	responseList, err = u.useUser.GetList(ctx, claims, paramquery)
	if err != nil {
		// return e.JSON(http.StatusBadRequest, err.Error())
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	// return e.JSON(http.StatusOK, ListDataUser)
	return appE.ResponseList(http.StatusOK, "", responseList)
}

// CreateSaUser :
// @Summary Add User
// @Security ApiKeyAuth
// @Tags User
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.AddUser true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/user [post]
func (u *ContUser) Create(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger     = logging.Logger{} // wajib
		appE    = tool.Res{R: e} // wajib
		sysUser models.SsUser
		form    models.AddUser
	)

	// user := e.Get("user").(*jwt.Token)
	// claims := user.Claims.(*util.Claims)
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	// logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	err := mapstructure.Decode(form, &sysUser)
	if err != nil {
		return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)

	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	// sysUser.UserInput = claims.UserID
	err = u.useUser.Create(ctx, claims, &sysUser)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", sysUser)
}

// ChangePassword :
// @Summary Rubah Password
// @Security ApiKeyAuth
// @Tags User
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.ChangePassword true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/user/change_password [post]
func (u *ContUser) ChangePassword(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{} // wajib
		appE = tool.Res{R: e} // wajib
		err  error
		// valid  validation.Validation                 // wajib
		form = models.ChangePassword{}
	)
	// user := e.Get("user").(*jwt.Token)
	// claims := user.Claims.(*util.Claims)

	// logger.Info(id)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	// logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	// form.UpdatedBy = claims.UserName
	err = u.useUser.ChangePassword(ctx, claims, form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateSaUser :
// @Summary Rubah Profile
// @Security ApiKeyAuth
// @Tags User
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Param req body models.UpdateUser true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /user/user/{id} [put]
func (u *ContUser) Update(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{} // wajib
		appE = tool.Res{R: e} // wajib
		err  error
		// valid  validation.Validation                 // wajib
		id   = e.Param("id") //kalo bukan int => 0
		form = models.UpdateUser{}
	)
	// user := e.Get("user").(*jwt.Token)
	// claims := user.Claims.(*util.Claims)

	MenuID, _ := strconv.Atoi(id)
	// logger.Info(id)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	// logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	// form.UpdatedBy = claims.UserName
	err = u.useUser.Update(ctx, claims, MenuID, form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusCreated, "Ok", nil)
}

// DeleteSaUser :
// @Summary Delete User
// @Security ApiKeyAuth
// @Tags User
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Success 200 {object} tool.ResponseModel
// @Router /user/user/{id} [delete]
func (u *ContUser) Delete(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		id   = e.Param("id")  //kalo bukan int => 0
		// valid  validation.Validation                 // wajib
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	err = u.useUser.Delete(ctx, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", nil)
}

func (u *ContUser) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "success")
}
