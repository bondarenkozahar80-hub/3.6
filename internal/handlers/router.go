package handlers

import (
	"context"
	"time"

	"github.com/bondarenkozahar80-hub/3.6/internal/middleware"
	"github.com/bondarenkozahar80-hub/3.6/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type transactionCRUDer interface {
	Create(ctx context.Context, t *model.Transaction) (int, error)
	GetByID(ctx context.Context, id int) (*model.Transaction, error)
	GetAll(ctx context.Context) ([]model.Transaction, error)
	Update(ctx context.Context, t *model.Transaction) error
	Delete(ctx context.Context, id int) error
	GetByPeriod(ctx context.Context, from, to time.Time) ([]model.Transaction, error)
	GetAllSorted(ctx context.Context, sortField, order string) ([]model.Transaction, error)
}

type transactionAnalyticsGetter interface {
	GetSum(ctx context.Context, from, to time.Time) (int64, error)
	GetAvg(ctx context.Context, from, to time.Time) (float64, error)
	GetCount(ctx context.Context, from, to time.Time) (int64, error)
	GetMedian(ctx context.Context, from, to time.Time) (float64, error)
	GetPercentile90(ctx context.Context, from, to time.Time) (float64, error)
	GroupByDay(ctx context.Context, from, to time.Time) (map[string]int64, error)
	GroupByWeek(ctx context.Context, from, to time.Time) (map[string]int64, error)
	GroupByMonth(ctx context.Context, from, to time.Time) (map[string]int64, error)
	GroupByCategory(ctx context.Context, from, to time.Time) (map[string]int64, error)
}

type Router struct {
	Router           *ginext.Engine
	tCRUDer          transactionCRUDer
	tAnalyticsGetter transactionAnalyticsGetter
}

func New(router *ginext.Engine, cruder transactionCRUDer, analyticGetter transactionAnalyticsGetter) *Router {
	return &Router{
		Router:           router,
		tCRUDer:          cruder,
		tAnalyticsGetter: analyticGetter,
	}
}

func (r *Router) Routes() {

	r.Router.Use(middleware.CORSMiddleware())
	r.Router.Use(middleware.LoggerMiddleware())

	r.Router.POST("/items", r.createTransactionHandler)
	r.Router.GET("/items", r.getAllTransactionsHandler)
	r.Router.GET("/items/:id", r.getTransactionByIDHandler)
	r.Router.PUT("/items/:id", r.updateTransactionHandler)
	r.Router.DELETE("/items/:id", r.deleteTransactionHandler)
	r.Router.GET("/items/period", r.getTransactionsByPeriodHandler)
	r.Router.GET("/items/sorted", r.getAllSortedHandler)
	r.Router.GET("/analytics", r.getAnalyticsHandler)
	r.Router.GET("/analytics/day", r.groupByDayHandler)
	r.Router.GET("/analytics/week", r.groupByWeekHandler)
	r.Router.GET("/analytics/month", r.groupByMonthHandler)
	r.Router.GET("/analytics/category", r.groupByCategoryHandler)
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")

}
