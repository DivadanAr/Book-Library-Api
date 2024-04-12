package ulasanbukucontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"strconv"
)

func UlasanBukuGet(c *fiber.Ctx) error {
	var (
		ulasanBuku    []response.UlasanBuku
		rowulasanBuku response.UlasanBuku
		cond          string
		count         int
		search        string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	search = c.Query("Search")
	if search != "" {
		cond += `WHERE (BukuId LIKE '%` + search + `% OR UserId LIKE %` + search + `% OR Ulasan LIKE %` + search + `% OR Rating LIKE %` + search + `%')`
	}

	ulasanBukuQry, err := db.QueryContext(ctx, `
	SELECT Id, BukuId, UserId, Ulasan, Rating
	FROM ulasan_buku
	`+cond+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer ulasanBukuQry.Close()
	for ulasanBukuQry.Next() {
		err := ulasanBukuQry.Scan(&rowulasanBuku.Id, &rowulasanBuku.BukuId, &rowulasanBuku.UserId, &rowulasanBuku.Ulasan, &rowulasanBuku.Rating)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		ulasanBuku = append(ulasanBuku, rowulasanBuku)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(ulasan_buku.Id)
		FROM ulasan_buku
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"UlasanBuku": ulasanBuku,
		"Total":      count,
	}, nil)
	return c.JSON(res)

}

func UlasanBukuDetail(c *fiber.Ctx) error {
	var (
		ulasanBuku    []response.UlasanBuku
		rowulasanBuku response.UlasanBuku
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	ulasanBukuQry, err := db.QueryContext(ctx, `
	SELECT Id, BukuId, UserId, Ulasan, Rating 
	FROM ulasan_buku WHERE Id = ?
	`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer ulasanBukuQry.Close()
	for ulasanBukuQry.Next() {
		err := ulasanBukuQry.Scan(&rowulasanBuku.Id, &rowulasanBuku.BukuId, &rowulasanBuku.UserId, &rowulasanBuku.Ulasan, &rowulasanBuku.Rating)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		ulasanBuku = append(ulasanBuku, rowulasanBuku)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"UlasanBuku": ulasanBuku,
	}, nil)
	return c.JSON(res)
}

func UlasanBukuPost(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	ulasanBuku := request.UlasanBuku{}
	if err := c.BodyParser(&ulasanBuku); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	qry, err := db.ExecContext(ctx, `
	INSERT INTO ulasan_buku (UserId, BukuId, Ulasan, Rating)
	VALUES (?,?,?,?)
`, c.Locals("UserId").(float64), ulasanBuku.BukuId, ulasanBuku.Ulasan, ulasanBuku.Rating)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	ulasanBuku.Id, _ = qry.LastInsertId()

	res := helpers.GetResponse(200, fiber.Map{
		"id":      ulasanBuku.Id,
		"buku_id": ulasanBuku.BukuId,
		"ulasan":  ulasanBuku.Ulasan,
		"rating":  ulasanBuku.Rating,
		"user_id": c.Locals("UserId").(float64),
	}, nil)
	return c.JSON(res)
}

func UlasanBukuPut(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	ulasanBuku := request.UlasanBuku{}
	if err := c.BodyParser(&ulasanBuku); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := db.ExecContext(ctx, `
	UPDATE ulasan_buku SET BukuId = ?, UserId = ?, Ulasan = ?, Rating = ? WHERE Id = ?`, ulasanBuku.BukuId, c.Locals("UserId").(float64), ulasanBuku.Ulasan, ulasanBuku.Rating, c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	ulasanBuku.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, fiber.Map{
		"id":      ulasanBuku.Id,
		"buku_id": ulasanBuku.BukuId,
		"ulasan":  ulasanBuku.Ulasan,
		"rating":  ulasanBuku.Rating,
		"user_id": c.Locals("UserId").(float64),
	}, nil)
	return c.JSON(res)
}

func UlasanBukuDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM ulasan_buku WHERE Id = ?`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}
