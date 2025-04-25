package controllers

import (
	"backend/db"
	issuePoolService "backend/services/issue"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CountIssuePools(c *gin.Context) {
	issuePoolService := issuePoolService.NewIssuePoolService(db.DB)
	pools, err := issuePoolService.CountIssuePools(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get pools", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Pools retrieved successfully", pools))
}
