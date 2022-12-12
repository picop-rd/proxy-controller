package controller

import (
	"errors"
	"net/http"

	"github.com/hiroyaonoe/bcop-proxy-controller/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/usecase"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Route struct {
	uc *usecase.Route
}

func NewRoute(uc *usecase.Route) *Route {
	return &Route{uc: uc}
}

func (r *Route) Put(c echo.Context) error {
	routes := []entity.Route{}

	if err := c.Bind(&routes); err != nil {
		log.Debug().Err(err).Msg("illegal body")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := r.uc.Register(c.Request().Context(), routes)
	if err != nil {
		arrRoutes := zerolog.Arr()
		for _, r := range routes {
			arrRoutes = arrRoutes.Object(r)
		}
		if errors.Is(err, entity.ErrInvalid) {
			log.Debug().Err(err).Array("routes", arrRoutes).Msg("illegal routes")
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		log.Error().Err(err).Array("routes", arrRoutes).Msg("unexpected error PUT /routes")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)

}
