package kategoribukucontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"strconv"
)

func KategoriBukuGet(c *fiber.Ctx) error {
	var (
		kategoriBuku    []response.KategoriBuku
		rowkategoriBuku response.KategoriBuku
		cond            string
		count           int
		search          string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	search = c.Query("Search")
	if search != "" {
		cond += `WHERE (BukuId LIKE '%` + search + `% OR Kategori LIKE %` + search + `% ')`
	}

	kategoriBukuQry, err := db.QueryContext(ctx, `
	SELECT Id, BukuId, KategoriId 
	FROM kategori_buku 
	`+cond+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer kategoriBukuQry.Close()
	for kategoriBukuQry.Next() {
		err := kategoriBukuQry.Scan(&rowkategoriBuku.Id, &rowkategoriBuku.BukuId, &rowkategoriBuku.KategoriId)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		kategoriBuku = append(kategoriBuku, rowkategoriBuku)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(kategori_buku.Id) 
		FROM kategori_buku
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"KategoriBuku": kategoriBuku,
		"Total":        count,
	}, nil)
	return c.JSON(res)
}

func KategoriBukuDetail(c *fiber.Ctx) error {
	var (
		kategoriBuku    []response.KategoriBuku
		rowkategoriBuku response.KategoriBuku
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	kategoriBukuQry, err := db.QueryContext(ctx, `
	SELECT Id, BukuId, KategoriId 
	FROM kategori_buku WHERE Id = ?
	`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer kategoriBukuQry.Close()
	for kategoriBukuQry.Next() {
		err := kategoriBukuQry.Scan(&rowkategoriBuku.Id, &rowkategoriBuku.BukuId, &rowkategoriBuku.KategoriId)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		kategoriBuku = append(kategoriBuku, rowkategoriBuku)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"KategoriBuku": kategoriBuku,
	}, nil)
	return c.JSON(res)
}

func KategoriBukuPost(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	kategoriBuku := request.KategoriBuku{}
	if err := c.BodyParser(&kategoriBuku); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	qry, err := db.ExecContext(ctx, `
	INSERT INTO kategori_buku (BukuId,KategoriId)
	VALUES (?,?)
`, kategoriBuku.BukuId, kategoriBuku.KategoriId)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	kategoriBuku.Id, _ = qry.LastInsertId()

	res := helpers.GetResponse(200, kategoriBuku, nil)
	return c.JSON(res)
}

func KategoriBukuPut(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	kategoriBuku := request.KategoriBuku{}
	if err := c.BodyParser(&kategoriBuku); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := db.ExecContext(ctx, `
	UPDATE kategori_buku SET BukuId = ?, KategoriId = ? WHERE Id = ?`, kategoriBuku.BukuId, kategoriBuku.KategoriId, c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	kategoriBuku.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, kategoriBuku, nil)
	return c.JSON(res)
}

func KategoriBukuDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM kategori_buku WHERE Id = ?`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}
