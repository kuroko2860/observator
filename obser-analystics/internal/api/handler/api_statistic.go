package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"kuroko.com/analystics/internal/model"
)

// @Summary		Get Api Statistic
// @Description	Get Api Statistic
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	true	"Service Name"
// @Param			endpoint		query		string	true	"Endpoint"
// @Param			method			query		string	true	"Method"
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Param			unit			query		string	true	"Unit"
// @Success		200				{object}	model.ApiStatistic
// @Failure		500				{object}	model.Error
// @Router			/api-statistics [get]
func (h *Handler) GetApiStatisticHandler(c echo.Context) error {
	serviceName := c.QueryParam("service_name")
	endpoint := c.QueryParam("endpoint")
	method := c.QueryParam("method")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	unit := c.QueryParam("unit")

	res, err := h.service.GetApiStatisticService(c.Request().Context(), serviceName, endpoint, method, from, to, unit)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Long Api
// @Description	Get Long Api
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Param			threshold		query		string	true	"Threshold"
// @Success		200				{object}	[]any
// @Failure		500				{object}	model.Error
// @Router			/long-api [get]
func (h *Handler) GetLongApiHandler(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	threshold := c.QueryParam("threshold")

	res, err := h.service.GetLongApiService(c.Request().Context(), from, to, threshold)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Called Api
// @Description	Get Called Api
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Param			username		query		string	true	"Username"
// @Success		200				{object}	[]any
// @Failure		500				{object}	model.Error
// @Router			/called-api [get]
func (h *Handler) GetCalledApiHandler(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	username := c.QueryParam("username")

	res, err := h.service.GetCalledApiService(c.Request().Context(), from, to, username)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Top Called Api
// @Description	Get Top Called Api
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Param			limit			query		string	true	"Limit"
// @Success		200				{object}	[]any
// @Failure		500				{object}	model.Error
// @Router			/top-called-api [get]
func (h *Handler) GetTopCalledApiHandler(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	limit := c.QueryParam("limit")

	res, err := h.service.GetTopCalledApi(c.Request().Context(), from, to, limit)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Http api by service Api
// @Description	Get Http api by service Api
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Param			service_name			query		string	true	"Service Name"
// @Success		200				{object}	[]any
// @Failure		500				{object}	model.Error
// @Router			/http-api [get]
func (h *Handler) GetHttpApiByServiceHandler(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	serviceName := c.QueryParam("service_name")

	res, err := h.service.GetHttpApiByService(c.Request().Context(), from, to, serviceName)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Alert
// @Description	Get Alert
// @Tags			alert
// @Accept			json
// @Produce		json
// @Success		200				{object}	[]model.AlertGetObject
// @Failure		500				{object}	model.Error
// @Router			/get-alert [get]
func (h *Handler) GetAlertHandler(c echo.Context) error {
	res, err := h.service.FindAllAlertGet(c.Request().Context())
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Uri List
// @Description	Get Uri List
// @Tags			alert
// @Accept			json
// @Produce		json
// @Success		200				{object}	[]model.URIObject
// @Failure		500				{object}	model.Error
// @Router			/uri-list [get]
func (h *Handler) GetUriListHandler(c echo.Context) error {
	res, err := h.service.FindURI(c.Request().Context())
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Ignore Alert
// @Description	Ignore Alert
// @Tags			alert
// @Accept			json
// @Produce		json
// @Param			id				path		string	true	"Id"
// @Success		200				{object}	model.AlertGetObject
// @Failure		500				{object}	model.Error
// @Router			/ignore-alert/{id} [patch]
func (h *Handler) IgnoreAlertHandler(c echo.Context) error {
	id := c.Param("id")
	err := h.service.IgnoreAlertGet(c.Request().Context(), id)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, model.AlertGetObject{})
}

// @Summary		Online Time
// @Description	Online Time
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			start_time			query		string	true	"Start Time"
// @Param			end_time			query		string	true	"End Time"
// @Param			user_id			query		string	true	"User Id"
// @Success		200				{object}	[]model.OnlineTimeOutput
// @Failure		500				{object}	model.Error
// @Router			/online-time [get]
func (h *Handler) OnlineTimeHandler(c echo.Context) error {
	startTimeStr := c.QueryParam("start_time")
	startTime, _ := strconv.Atoi(startTimeStr)
	endTimeStr := c.QueryParam("end_time")
	endTime, _ := strconv.Atoi(endTimeStr)
	userId := c.QueryParam("user_id")

	res, err := h.service.CheckOnlineTime(c.Request().Context(), model.TimeInput{
		StartTime: int64(startTime),
		EndTime:   int64(endTime),
	}, userId)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Online User
// @Description	Online User
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			start_time			query		string	true	"Start Time"
// @Param			end_time			query		string	true	"End Time"
// @Success		200				{object}	[]string
// @Failure		500				{object}	model.Error
// @Router			/online-user [get]
func (h *Handler) OnlineUserHandler(c echo.Context) error {
	startTimeStr := c.QueryParam("start_time")
	startTime, _ := strconv.Atoi(startTimeStr)
	endTimeStr := c.QueryParam("end_time")
	endTime, _ := strconv.Atoi(endTimeStr)

	res, err := h.service.CheckOnlineUser(c.Request().Context(), model.TimeInput{
		StartTime: int64(startTime),
		EndTime:   int64(endTime),
	})
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Service Statistic
// @Description	Service Statistic
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			date			query		string	true	"Date"
// @Param			service			query		string	true	"Service"
// @Success		200				{object}	[]model.ServiceStatisticObject
// @Failure		500				{object}	model.Error
// @Router			/service-statistic [get]
func (h *Handler) ServiceStatisticHandler(c echo.Context) error {
	date := c.QueryParam("date")
	svc := c.QueryParam("service")
	var rs []model.ServiceStatisticObject
	if svc == "" {
		rs, _ = h.service.FindServiceStatisticByDate(c.Request().Context(), date)
	} else {
		rs, _ = h.service.FindServiceStatisticByDateAndName(c.Request().Context(), date, svc)
	}
	return c.JSON(200, rs)
}

// @Summary		Uri Statistic
// @Description	Uri Statistic
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			date			query		string	true	"Date"
// @Param			uri			query		string	true	"Uri"
// @Success		200				{object}	[]model.URIStatisticObject
// @Failure		500				{object}	model.Error
// @Router			/uri-statistic [get]
func (h *Handler) UriStatisticHandler(c echo.Context) error {
	date := c.QueryParam("date")
	uri := c.QueryParam("uri")
	var rs []model.URIStatisticObject
	var err error
	if uri == "" {
		rs, err = h.service.FindURIStatisticByDate(c.Request().Context(), date)
	} else {
		rs, err = h.service.FindURIStatisticByDateAndUri(c.Request().Context(), date, uri)
	}

	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}
	return c.JSON(200, rs)
}

// @Summary		Usage
// @Description	Usage
// @Tags			api
// @Accept			json
// @Produce		json
// @Param			start_time			query		string	true	"Start Time"
// @Param			end_time			query		string	true	"End Time"
// @Param			service			query		string	true	"Service"
// @Param			path			query		string	true	"Path"
// @Success		200				{object}	[]string
// @Failure		500				{object}	model.Error
// @Router			/usage [get]
func (h *Handler) UsageHandler(c echo.Context) error {
	startTimeStr := c.QueryParam("start_time")
	startTime, _ := strconv.Atoi(startTimeStr)
	endTimeStr := c.QueryParam("end_time")
	endTime, _ := strconv.Atoi(endTimeStr)

	svc := c.QueryParam("service")
	path := c.QueryParam("path")

	var rs []string
	var err error
	if path == "" {
		rs, err = h.service.CheckUserFromTo(c.Request().Context(), model.TimeInput{
			StartTime: int64(startTime),
			EndTime:   int64(endTime),
		}, svc)

	} else {
		rs, err = h.service.CheckUserFromToWithPath(c.Request().Context(), model.TimeInput{
			StartTime: int64(startTime),
			EndTime:   int64(endTime),
		}, svc, path)
	}

	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, rs)
}
