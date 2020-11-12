package iberandauser

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Usecase interface {
	GetClosestBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList)
	GetRecomentCapster(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList)
	GetRecomentBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList)
}
