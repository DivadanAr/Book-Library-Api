package koleksipribadicontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"strconv"
)

func KoleksiPribadiGet(c *fiber.Ctx) error {
	var (
		koleksiPribadi    []response.Buku
		rowkoleksiPribadi response.Buku
		cond              string
		count             int
		search            string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	search = c.Query("Search")
	if search != "" {
		cond += `WHERE (BukuId LIKE '%` + search + `% OR UserId LIKE %` + search + `% ')`
	}

	koleksiPribadiQry, err := db.QueryContext(ctx, `
	SELECT buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi
	FROM koleksi_pribadi JOIN buku ON buku.Id = koleksi_pribadi.BukuId
	`+cond+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer koleksiPribadiQry.Close()
	for koleksiPribadiQry.Next() {
		err := koleksiPribadiQry.Scan(&rowkoleksiPribadi.Id, &rowkoleksiPribadi.Judul, &rowkoleksiPribadi.Penulis, &rowkoleksiPribadi.Penerbit, &rowkoleksiPribadi.Cover, &rowkoleksiPribadi.BackCover, &rowkoleksiPribadi.JumlahHalaman, &rowkoleksiPribadi.TahunTerbit, &rowkoleksiPribadi.StokBuku, &rowkoleksiPribadi.Deskripsi)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		koleksiPribadi = append(koleksiPribadi, rowkoleksiPribadi)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(koleksi_pribadi.Id)
		FROM koleksi_pribadi JOIN buku ON buku.Id = koleksi_pribadi.BukuId
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"KoleksiPribadi": koleksiPribadi,
		"Total":          count,
	}, nil)
	return c.JSON(res)

	//res := helpers.GetResponse(200, fiber.Map{
	//	"UserId":   c.Locals("UserId").(float64),
	//	"Username": c.Locals("Username").(string),
	//}, nil)
	//return c.JSON(res)
}

func KoleksiPribadiDetail(c *fiber.Ctx) error {
	var (
		koleksiPribadi    []response.Buku
		rowkoleksiPribadi response.Buku
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	koleksiPribadiQry, err := db.QueryContext(ctx, `
	SELECT buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi
	FROM koleksi_pribadi JOIN buku ON buku.Id = koleksi_pribadi.BukuId WHERE koleksi_pribadi.UserId = ?
	`, c.Locals("UserId"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer koleksiPribadiQry.Close()
	for koleksiPribadiQry.Next() {
		err := koleksiPribadiQry.Scan(&rowkoleksiPribadi.Id, &rowkoleksiPribadi.Judul, &rowkoleksiPribadi.Penulis, &rowkoleksiPribadi.Penerbit, &rowkoleksiPribadi.Cover, &rowkoleksiPribadi.BackCover, &rowkoleksiPribadi.JumlahHalaman, &rowkoleksiPribadi.TahunTerbit, &rowkoleksiPribadi.StokBuku, &rowkoleksiPribadi.Deskripsi)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		koleksiPribadi = append(koleksiPribadi, rowkoleksiPribadi)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"KoleksiPribadi": koleksiPribadi,
	}, nil)
	return c.JSON(res)
}

func KoleksiPribadiPost(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	koleksiPribadi := request.KoleksiPribadi{}
	if err := c.BodyParser(&koleksiPribadi); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	qry, err := db.ExecContext(ctx, `
	INSERT INTO koleksi_pribadi (BukuId,UserId)
	VALUES (?,?)
`, koleksiPribadi.BukuId, c.Locals("UserId").(float64))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	koleksiPribadi.Id, _ = qry.LastInsertId()

	res := helpers.GetResponse(200, fiber.Map{
		"id":      koleksiPribadi.Id,
		"buku_id": koleksiPribadi.BukuId,
		"user_id": c.Locals("UserId").(float64),
	}, nil)
	return c.JSON(res)
}

func KoleksiPribadiPut(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	koleksiPribadi := request.KoleksiPribadi{}
	if err := c.BodyParser(&koleksiPribadi); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := db.ExecContext(ctx, `
	UPDATE koleksi_pribadi SET BukuId = ?, UserId = ? WHERE Id = ?`, koleksiPribadi.BukuId, c.Locals("UserId"), c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	koleksiPribadi.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, koleksiPribadi, nil)
	return c.JSON(res)
}

func KoleksiPribadiDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM koleksi_pribadi WHERE Id = ?`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}
