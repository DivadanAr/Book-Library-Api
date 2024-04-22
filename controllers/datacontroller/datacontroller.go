package datacontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
)

type UserTotals struct {
	TotalLakiLaki  int `json:"total_laki_laki"`
	TotalPerempuan int `json:"total_perempuan"`
	TotalUsers  int `json:"total_users"`
}

type DataTotals struct {
	TotalRequest  int `json:"total_request"`
	TotalApprove  int `json:"total_approve"`
	TotalReject   int `json:"total_reject"`
	TotalBorrowed int `json:"total_borrowed"`
	TotalReturn   int `json:"total_return"`
	TotalDone     int `json:"total_done"`
}

func CountDataGet(c *fiber.Ctx) error {
	var (
		users       UserTotals
		dataTotals  DataTotals
		rowUsers    UserTotals
		rowData     DataTotals
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	userQry, err := db.QueryContext(ctx, `
	CALL GetTotalUser()
	`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer userQry.Close()
	for userQry.Next() {
		err := userQry.Scan(&rowUsers.TotalPerempuan, &rowUsers.TotalLakiLaki, &rowUsers.TotalUsers)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		users = rowUsers
	}

	dataQry, err := db.QueryContext(ctx, `
	CALL GetTotalPeminjamanByStatus()
	`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer dataQry.Close()
	for dataQry.Next() {
		err := dataQry.Scan(&rowData.TotalRequest, &rowData.TotalApprove, &rowData.TotalReject,
			&rowData.TotalBorrowed, &rowData.TotalReturn, &rowData.TotalDone)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		dataTotals = rowData
	}

	res := helpers.GetResponse(200, fiber.Map{
		"users": users,
		"peminjaman":  dataTotals,
	}, nil)
	return c.JSON(res)
}
