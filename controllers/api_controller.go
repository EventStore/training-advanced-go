package controllers

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	DateLayout = "2006-01-02"
)

type SlotsController struct {
	availableSlotsRepository readmodel.AvailableSlotsRepository
	dispatcher               *infrastructure.Dispatcher
	eventStore               infrastructure.EventStore
}

func NewSlotsController(d *infrastructure.Dispatcher, a readmodel.AvailableSlotsRepository, e infrastructure.EventStore) *SlotsController {
	return &SlotsController{
		dispatcher:               d,
		availableSlotsRepository: a,
		eventStore:               e,
	}
}

func (c *SlotsController) Register(prefix string, e *echo.Echo) {
	e.GET(path.Join(prefix, "/slots/today/available"), c.AvailableTodayHandler)
	e.GET(path.Join(prefix, "/slots/:date/available"), c.AvailableHandler)

	e.POST(path.Join(prefix, "/doctor/schedule"), c.ScheduleDayHandler)
	e.POST(path.Join(prefix, "/slots/:dayId/cancel-booking"), c.CancelBookingHandler)
	e.POST(path.Join(prefix, "/slots/:dayId/book"), c.BookSlotHandler)
	e.POST(path.Join(prefix, "/calendar/:date/day-started"), c.CalendarDayStartedHandler)

}

func (c *SlotsController) AvailableTodayHandler(ctx echo.Context) error {
	today, err := time.Parse(DateLayout, "2020-08-01")
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return c.availableHandler(ctx, today)
}

func (c *SlotsController) AvailableHandler(ctx echo.Context) error {
	date, err := time.Parse(DateLayout, ctx.Param("date"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return c.availableHandler(ctx, date)
}

func (c *SlotsController) availableHandler(ctx echo.Context, date time.Time) error {
	availableSlots, err := c.availableSlotsRepository.GetSlotsAvailableOn(date)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	response := make([]AvailableSlotResponse, 0)
	for _, a := range availableSlots {
		response = append(response, AvailableSlotResponseFrom(a))
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *SlotsController) ScheduleDayHandler(ctx echo.Context) error {
	req := ScheduleDayRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	scheduleDay, err := req.ToCommand()
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	metadata, err := c.GetCommandMetadata(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = c.dispatcher.Dispatch(scheduleDay, metadata)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	url := fmt.Sprintf("/api/slots/%s/available", scheduleDay.Date.Format(DateLayout))
	return ctx.Redirect(http.StatusFound, url)
}

func (c *SlotsController) CancelBookingHandler(ctx echo.Context) error {
	req := CancelSlotBookingRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	metadata, err := c.GetCommandMetadata(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = c.dispatcher.Dispatch(req.ToCommand(ctx.Param("dayId")), metadata)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.String(http.StatusOK, "successfully cancelled booking")
}

func (c *SlotsController) BookSlotHandler(ctx echo.Context) error {
	req := BookSlotRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	metadata, err := c.GetCommandMetadata(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = c.dispatcher.Dispatch(req.ToCommand(ctx.Param("dayId")), metadata)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.String(http.StatusOK, "slot successfully booked")
}

func (c *SlotsController) CalendarDayStartedHandler(ctx echo.Context) error {
	date, err := time.Parse(DateLayout, ctx.Param("date"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	metadata, err := c.GetCommandMetadata(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = c.eventStore.AppendEventsToAny("doctorday-time-events", metadata, events.NewCalendarDayStarted(date))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.String(http.StatusOK, "calendar day successfully started")
}

func (c *SlotsController) GetCommandMetadata(ctx echo.Context) (infrastructure.CommandMetadata, error) {
	correlationId, err := getHeaderUUIDValue(ctx, "X-CorrelationId")
	if err != nil {
		return infrastructure.CommandMetadata{}, nil
	}

	causationId, err := getHeaderUUIDValue(ctx, "X-CausationId")
	if err != nil {
		return infrastructure.CommandMetadata{}, nil
	}

	return infrastructure.NewCommandMetadata(correlationId, causationId), nil
}

func getHeaderUUIDValue(ctx echo.Context, name string) (uuid.UUID, error) {
	v := ctx.Request().Header.Get("X-CorrelationId")
	if v == "" {
		return uuid.UUID{}, fmt.Errorf("please provide an %s header", name)
	}

	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("please provide a valid %s header", name)
	}

	return id, nil
}
