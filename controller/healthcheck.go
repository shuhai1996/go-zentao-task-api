package controller

type HealthCheckRequest struct {
	Timestamp string `form:"timestamp" binding:"required"`
	Sign      string `form:"sign" binding:"required"`
}