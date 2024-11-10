package filter

import (
	"github.com/eris-apple/eaapi"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Filter struct {
	Limit  int
	Offset int
	Order  string
}

func GetDefaultsFilter(filter *Filter, prefix ...string) *Filter {
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	if filter.Order == "" {
		if len(prefix) > 0 {
			filter.Order = prefix[0] + ".id desc"
		} else {
			filter.Order = "id desc"
		}
	}

	return &Filter{
		Limit:  filter.Limit,
		Offset: filter.Offset,
		Order:  filter.Order,
	}
}

func GetDefaultsFilterFromQuery(ctx *gin.Context, prefix ...string) *Filter {
	var limitInt int
	var offsetInt int

	limitString := ctx.Query("limit")
	if limitString != "" {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			eaapi.NewHandlerError(ctx, eaapi.NewBadRequestError("bad request params"))
			return nil
		}

		limitInt = limit
	}

	offsetString := ctx.Query("offset")
	if offsetString != "" {
		offset, err := strconv.Atoi(offsetString)
		if err != nil {
			eaapi.NewHandlerError(ctx, eaapi.NewBadRequestError("bad request params"))
			return nil
		}
		offsetInt = offset
	}

	order := ctx.Query("order")
	if order == "" {
		order = "id desc"
	}

	if len(prefix) > 0 {
		order = prefix[0] + "." + order
	}

	f := &Filter{
		Limit:  limitInt,
		Offset: offsetInt,
		Order:  order,
	}

	return GetDefaultsFilter(f)
}
