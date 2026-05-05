package admin

type DashboardFilter struct {
	UserID    *int64 `form:"user_id"`
	StartDate string `form:"start_date"` // format: 2024-01-01
	EndDate   string `form:"end_date"`
}
