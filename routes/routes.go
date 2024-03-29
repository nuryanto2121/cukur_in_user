package routes

import (
	postgresgorm "nuryanto2121/cukur_in_user/pkg/postgregorm"
	// sqlxposgresdb "nuryanto2121/cukur_in_user/pkg/postgresqlxdb"
	"nuryanto2121/cukur_in_user/pkg/setting"

	_saauthcont "nuryanto2121/cukur_in_user/controllers/auth"
	_authuse "nuryanto2121/cukur_in_user/usecase/auth"

	_saFilecont "nuryanto2121/cukur_in_user/controllers/fileupload"
	_repoFile "nuryanto2121/cukur_in_user/repository/ss_fileupload"
	_useFile "nuryanto2121/cukur_in_user/usecase/ss_fileupload"

	_saBarberFavoritcont "nuryanto2121/cukur_in_user/controllers/barber_favorit"
	_repoBarberFavorit "nuryanto2121/cukur_in_user/repository/barber_favorit"
	_useBarberFavorit "nuryanto2121/cukur_in_user/usecase/barber_favorit"

	_BerandaUsercont "nuryanto2121/cukur_in_user/controllers/beranda_user"
	_useBerandaUser "nuryanto2121/cukur_in_user/usecase/beranda_user"

	_saBarbercont "nuryanto2121/cukur_in_user/controllers/barber"
	_repoBarber "nuryanto2121/cukur_in_user/repository/barber"
	_repoBarberCapster "nuryanto2121/cukur_in_user/repository/barber_capster"
	_useBarber "nuryanto2121/cukur_in_user/usecase/barber"

	_contUser "nuryanto2121/cukur_in_user/controllers/user"
	_repoUser "nuryanto2121/cukur_in_user/repository/ss_user"
	_useUser "nuryanto2121/cukur_in_user/usecase/ss_user"

	_contCapster "nuryanto2121/cukur_in_user/controllers/capster"
	_repoCapster "nuryanto2121/cukur_in_user/repository/capster"
	_useCapster "nuryanto2121/cukur_in_user/usecase/capster"

	_contOrder "nuryanto2121/cukur_in_user/controllers/c_order"
	_repoOrderd "nuryanto2121/cukur_in_user/repository/c_order_d"
	_repoOrder "nuryanto2121/cukur_in_user/repository/c_order_h"
	_useOrder "nuryanto2121/cukur_in_user/usecase/c_order"

	_contBarberPaket "nuryanto2121/cukur_in_user/controllers/barber_paket"
	_repoBarberPaket "nuryanto2121/cukur_in_user/repository/barber_paket"
	_useBarberPaket "nuryanto2121/cukur_in_user/usecase/barber_paket"

	_contBookingCapster "nuryanto2121/cukur_in_user/controllers/booking_capster"
	_repoBookingCapster "nuryanto2121/cukur_in_user/repository/booking_capster"
	_useBookingCapster "nuryanto2121/cukur_in_user/usecase/booking_capster"

	_contNotification "nuryanto2121/cukur_in_user/controllers/notification"
	_repoNotification "nuryanto2121/cukur_in_user/repository/notification"
	_useNotification "nuryanto2121/cukur_in_user/usecase/notification"

	_contAdvertise "nuryanto2121/cukur_in_user/controllers/advertise"
	_repoAdvertise "nuryanto2121/cukur_in_user/repository/advertise"
	_useAdvertise "nuryanto2121/cukur_in_user/usecase/advertise"

	_contFeedbackRating "nuryanto2121/cukur_in_user/controllers/feedback_rating"
	_repoFeedbackRating "nuryanto2121/cukur_in_user/repository/feedback_rating"
	_useFeedbackRating "nuryanto2121/cukur_in_user/usecase/feedback_rating"

	_contValidasi "nuryanto2121/cukur_in_user/controllers/validasi"

	"time"

	"github.com/labstack/echo/v4"
)

//Echo :
type EchoRoutes struct {
	E *echo.Echo
}

func (e *EchoRoutes) InitialRouter() {
	timeoutContext := time.Duration(setting.FileConfigSetting.Server.ReadTimeout) * time.Second

	repoFile := _repoFile.NewRepoFileUpload(postgresgorm.Conn)
	useFile := _useFile.NewSaFileUpload(repoFile, timeoutContext)
	_saFilecont.NewContFileUpload(e.E, useFile)

	repoFeedbackRating := _repoFeedbackRating.NewRepoFeedbackRating(postgresgorm.Conn)
	useFeedbackRating := _useFeedbackRating.NewFeedbackRating(repoFeedbackRating, timeoutContext)
	_contFeedbackRating.NewContFeedbackRating(e.E, useFeedbackRating)

	repoAdvertise := _repoAdvertise.NewRepoAdvertise(postgresgorm.Conn)
	useAdvertise := _useAdvertise.NewUseAdvertise(repoAdvertise, repoFile, timeoutContext)
	_contAdvertise.NewContAdvertise(e.E, useAdvertise)

	repoUser := _repoUser.NewRepoSysUser(postgresgorm.Conn)
	useUser := _useUser.NewUserSysUser(repoUser, repoFile, timeoutContext)
	_contUser.NewContUser(e.E, useUser)

	repoBarberFavorit := _repoBarberFavorit.NewRepoBarberFavorit(postgresgorm.Conn)
	useBarberFavorit := _useBarberFavorit.NewBarberFavorit(repoBarberFavorit, timeoutContext)
	_saBarberFavoritcont.NewContBarberFavorit(e.E, useBarberFavorit)

	repoBookingCapster := _repoBookingCapster.NewRepoBookingCapster(postgresgorm.Conn)
	useBookingCapster := _useBookingCapster.NewBookingCapster(repoBookingCapster, timeoutContext)
	_contBookingCapster.NewContBookingCapster(e.E, useBookingCapster)

	repoBarberPaket := _repoBarberPaket.NewRepoBarberPaket(postgresgorm.Conn)
	useBarberPaket := _useBarberPaket.NewBarberPaket(repoBarberPaket, timeoutContext)
	_contBarberPaket.NewContBarberPaket(e.E, useBarberPaket)

	repoBarberCapster := _repoBarberCapster.NewRepoBarberCapster(postgresgorm.Conn)
	repoBarber := _repoBarber.NewRepoBarber(postgresgorm.Conn)
	useBarber := _useBarber.NewUserMBarber(repoBarber, repoBarberPaket, repoBarberCapster, repoFile, timeoutContext)
	_saBarbercont.NewContBarber(e.E, useBarber)

	useBerandaUser := _useBerandaUser.NewUseBerandaUser(useBarber, repoFile, repoBarberCapster, repoBarber, useAdvertise, timeoutContext)
	_BerandaUsercont.NewContBerandaUser(e.E, useBerandaUser)

	repoCapster := _repoCapster.NewRepoCapsterCollection(postgresgorm.Conn)
	useCapster := _useCapster.NewUserMCapster(repoCapster, repoUser, repoBarberCapster, repoFile, timeoutContext)
	_contCapster.NewContCapster(e.E, useCapster)

	repoNotif := _repoNotification.NewRepoNotification(postgresgorm.Conn)
	useNotif := _useNotification.NewUseNotification(repoNotif, timeoutContext)
	_contNotification.NewContNotification(e.E, useNotif)

	repoOrderD := _repoOrderd.NewRepoOrderD(postgresgorm.Conn)
	repoOrder := _repoOrder.NewRepoOrderH(postgresgorm.Conn)
	useOrder := _useOrder.NewUserMOrder(repoOrder, repoOrderD, repoBarber,
		repoBookingCapster, useNotif,
		repoFeedbackRating, repoNotif,
		timeoutContext)
	_contOrder.NewContOrder(e.E, useOrder)

	_contValidasi.NewContValidasi(e.E)
	//_saauthcont
	// repoAuth := _repoAuth.NewRepoOptionDB(postgresgorm.Conn)
	useAuth := _authuse.NewUserAuth(repoUser, repoFile, timeoutContext)
	_saauthcont.NewContAuth(e.E, useAuth)

}
