package kategoricontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"strconv"
)

func KategoriGet(c *fiber.Ctx) error {
	var (
		kategori    []response.Kategori
		rowkategori response.Kategori
		cond        string
		count       int
		search      string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	search = c.Query("Search")
	if search != "" {
		cond += `WHERE (Kategori LIKE '%` + search + `%')`
	}

	kategoriQry, err := db.QueryContext(ctx, `
	SELECT Id, Kategori 
	FROM kategori 
	`+cond+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer kategoriQry.Close()
	for kategoriQry.Next() {
		err := kategoriQry.Scan(&rowkategori.Id, &rowkategori.Kategori)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		kategori = append(kategori, rowkategori)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(kategori.Id) 
		FROM kategori
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Kategori": kategori,
		"Total":    count,
	}, nil)
	return c.JSON(res)
}

func KategoriDetail(c *fiber.Ctx) error {
	var (
		kategori    []response.Kategori
		rowkategori response.Kategori
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	kategoriQry, err := db.QueryContext(ctx, `
	SELECT Id, Kategori 
	FROM kategori WHERE Id = ?
	`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer kategoriQry.Close()
	for kategoriQry.Next() {
		err := kategoriQry.Scan(&rowkategori.Id, &rowkategori.Kategori)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		kategori = append(kategori, rowkategori)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Kategori": kategori,
	}, nil)
	return c.JSON(res)
}

func KategoriPost(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	kategori := request.Kategori{}
	if err := c.BodyParser(&kategori); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	qry, err := db.ExecContext(ctx, `
	INSERT INTO kategori (Kategori)
	VALUES (?)
`, kategori.Kategori)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	kategori.Id, _ = qry.LastInsertId()

	res := helpers.GetResponse(200, kategori, nil)
	return c.JSON(res)
}

func KategoriPut(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	kategori := request.Kategori{}
	if err := c.BodyParser(&kategori); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := db.ExecContext(ctx, `
	UPDATE kategori SET Kategori = ? WHERE Id = ?`, kategori.Kategori, c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	kategori.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, kategori, nil)
	return c.JSON(res)
}

func KategoriDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM kategori WHERE Id = ?`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}
