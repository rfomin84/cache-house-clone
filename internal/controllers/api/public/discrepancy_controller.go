package public

import (
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DiscrepancyController struct {
	DiscrepancyState *managers.DiscrepancyState
}

func (c *DiscrepancyController) Index(ctx echo.Context) error {
	validate := validator.New()
	validate.RegisterValidation("date", ValidateDate)
	validate.RegisterValidation("start_date", ValidateStartDate)

	queryParams, _ := url.ParseQuery(ctx.Request().URL.String())

	billingTypes := queryParams["billing_types"]

	type requestData struct {
		StartDate string `json:"start_date" validate:"required,date,start_date"`
		EndDate   string `json:"end_date" validate:"required,date"`
	}

	r := &requestData{
		StartDate: ctx.QueryParam("start_date"),
		EndDate:   ctx.QueryParam("end_date"),
	}
	err := validate.Struct(r)
	if err != nil {
		type errMsg struct {
			Mgs string `json:"mgs"`
		}
		errSlice := make([]errMsg, 0)
		for _, err := range err.(validator.ValidationErrors) {

			msgText := err.Field() + " should be " + err.Tag()
			if err.Tag() == "start_date" {
				deadline := carbon.Now().SubMonths(6)
				msgText = err.Field() + " date must be later than " + deadline.ToDateTimeString()
			}

			errSlice = append(errSlice, errMsg{
				Mgs: msgText,
			})
		}

		errMsgJson := struct {
			Errors []errMsg `json:"errors"`
		}{
			Errors: errSlice,
		}

		return ctx.JSON(http.StatusBadRequest, errMsgJson)
	}

	var feedType managers.FeedType
	if len(queryParams["is_dsp"]) > 0 {
		value, err := strconv.Atoi(queryParams["is_dsp"][0])
		if err != nil {
			c.DiscrepancyState.Logger.Error(err.Error())
			return ctx.String(http.StatusInternalServerError, "Error")
		}
		feedType = managers.FeedType(value)
	} else {
		feedType = managers.All
	}

	startDate := carbon.Parse(ctx.QueryParam("start_date"))
	endDate := carbon.Parse(ctx.QueryParam("end_date"))

	result := c.DiscrepancyState.GetDiscrepancies(startDate.Carbon2Time(), endDate.Carbon2Time(), billingTypes, feedType)

	ctx.Response().Header().Set(echo.HeaderContentType, "application/json")
	return ctx.JSON(http.StatusOK, result)
}

// ValidateDate implements validator.Func
func ValidateDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

func ValidateStartDate(fl validator.FieldLevel) bool {
	startDate := carbon.Now().SubMonths(6)
	date := carbon.Parse(fl.Field().String())

	return startDate.Lte(date)
}
