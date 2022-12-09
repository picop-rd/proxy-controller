package controller

import (
	"errors"
	"net/http"

	"github.com/hiroyaonoe/bcop-proxy-controller/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/usecase"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Proxy struct {
	uc *usecase.Proxy
}

func NewProxy(uc *usecase.Proxy) *Proxy {
	return &Proxy{uc: uc}
}

func (p *Proxy) Register(c echo.Context) error {
	proxyID := c.Param("proxy-id")
	if len(proxyID) == 0 {
		log.Debug().Msg("illegal param")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var proxy entity.Proxy
	if err := c.Bind(&proxy); err != nil {
		log.Debug().Err(err).Msg("illegal body")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	proxy.ProxyID = proxyID

	err := p.uc.Register(c.Request().Context(), proxy)
	if err != nil {
		if errors.Is(err, entity.ErrInvalid) {
			log.Debug().Err(err).Object("proxy", proxy).Msg("illegal proxy")
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		log.Error().Err(err).Object("proxy", proxy).Msg("unexpected error PUT /proxy/:proxy-id/register")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (p *Proxy) Activate(c echo.Context) error {
	proxyID := c.Param("proxy-id")
	if len(proxyID) == 0 {
		log.Debug().Msg("illegal param")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	routes, err := p.uc.Activate(c.Request().Context(), proxyID)
	if err != nil {
		log.Error().Err(err).Str("proxyID", proxyID).Msg("unexpected error PUT /proxy/:proxy-id/activate")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, routes)
}

func (p *Proxy) Delete(c echo.Context) error {
	proxyID := c.Param("proxy-id")
	if len(proxyID) == 0 {
		log.Debug().Msg("illegal param")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := p.uc.Delete(c.Request().Context(), proxyID)
	if err != nil {
		log.Error().Err(err).Str("proxyID", proxyID).Msg("unexpected error DELETE /proxy/:proxy-id")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
