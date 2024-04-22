package bukucontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"strconv"
	"strings"
	"time"
)

func BukuGet(c *fiber.Ctx) error {
	var (
		buku    []response.Buku
		rowbuku response.Buku
		cond    string
		count   int
		search  string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	search = c.Query("Search")
	if search != "" {
		cond += ` AND (buku.Judul LIKE '%` + search + `%' OR buku.Penulis LIKE '%` + search + `%' OR buku.Penerbit LIKE '%` + search + `%')`
	}

	bukuQry, err := db.QueryContext(ctx, `
	SELECT buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi, AVG(ulasan_buku.Rating) AS AvgRating
	FROM buku 
	LEFT JOIN ulasan_buku ON buku.Id = ulasan_buku.BukuId
	WHERE buku.StokBuku != '0'
	`+cond+`
	GROUP BY buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi
	`+helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer bukuQry.Close()
	for bukuQry.Next() {
		err := bukuQry.Scan(&rowbuku.Id, &rowbuku.Judul, &rowbuku.Penulis, &rowbuku.Penerbit, &rowbuku.Cover, &rowbuku.BackCover, &rowbuku.JumlahHalaman, &rowbuku.TahunTerbit, &rowbuku.StokBuku, &rowbuku.Deskripsi, &rowbuku.Rating)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		buku = append(buku, rowbuku)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(buku.Id) 
		FROM buku
		WHERE StokBuku != '0'
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}
	

	res := helpers.GetResponse(200, fiber.Map{
		"Buku":  buku,
		"Total": count,
	}, nil)
	return c.JSON(res)
}

func BukuDetail(c *fiber.Ctx) error {
	var (
		buku        []response.Buku
		rowbuku     response.Buku
		kategori    []response.Kategori
		rowkategori response.Kategori
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	bukuQry, err := db.QueryContext(ctx, `
	SELECT buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi, kategori.Kategori, kategori.Id
	FROM kategori_buku
	JOIN buku ON buku.Id = kategori_buku.BukuId
	JOIN kategori ON  kategori.Id = kategori_buku.KategoriId
	WHERE buku.Id = ?
	`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer bukuQry.Close()
	for bukuQry.Next() {
		err := bukuQry.Scan(&rowbuku.Id, &rowbuku.Judul, &rowbuku.Penulis, &rowbuku.Penerbit, &rowbuku.Cover, &rowbuku.BackCover, &rowbuku.JumlahHalaman, &rowbuku.TahunTerbit, &rowbuku.StokBuku, &rowbuku.Deskripsi, &rowkategori.Kategori, &rowkategori.Id)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		buku = append(buku, rowbuku)
		kategori = append(kategori, rowkategori)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Buku":     buku,
		"Kategori": kategori,
	}, nil)
	return c.JSON(res)
}

func BookmarkDetail(c *fiber.Ctx) error {
	var (
		count int
		bukuId *string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*), BukuId 
		FROM koleksi_pribadi
		WHERE UserId = ? AND BukuId = ?
		`, c.Locals("UserId"), c.Params("id")).Scan(&count, &bukuId)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}
	res := helpers.GetResponse(200, fiber.Map{
		"total": count,
		"buku_id": bukuId,
	}, nil)
	return c.JSON(res)
}

func BookmarkDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM koleksi_pribadi WHERE UserId = ? AND BukuId = ?
	`, c.Locals("UserId"), c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}

func BukuPost(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()
	kategori := request.KategoriBuku{}
	buku := request.Buku{}

	if err := c.BodyParser(&buku); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	file, errFile := c.FormFile("cover")

	if errFile != nil {
		return errFile
	}

	fileName := file.Filename
	extension := fileName[strings.LastIndex(fileName, "."):]

	timestamp := time.Now().Format("20060102150405")
	coverName := "cover-" + strings.ReplaceAll(buku.Judul, " ", "") + "-" + timestamp + extension

	c.SaveFile(file, "public/uploads/cover/"+coverName)

	fileNameBackCover := file.Filename
	extensionBackCover := fileNameBackCover[strings.LastIndex(fileNameBackCover, "."):]

	timestampBackCover := time.Now().Format("20060102150405")
	backCoverName := "backcover-" + strings.ReplaceAll(buku.Judul, " ", "") + "-" + timestampBackCover + extensionBackCover

	c.SaveFile(file, "public/uploads/cover/"+backCoverName)

	qry, errBuku := db.ExecContext(ctx, `
	INSERT INTO buku (Judul, Penulis, Penerbit, Cover, BackCover, JumlahHalaman, TahunTerbit, StokBuku, Deskripsi)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, buku.Judul, buku.Penulis, buku.Penerbit, coverName, backCoverName, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi)

	if errBuku != nil {
		res := helpers.GetResponse(500, nil, errBuku)
		return c.Status(res.Status).JSON(res)
	}

	buku.Id, _ = qry.LastInsertId()

	qryKategori, errKategori := db.ExecContext(ctx, `INSERT INTO kategori_buku (KategoriId, BukuId) VALUE (?,?)`, buku.KategoriId, buku.Id)

	if errKategori != nil {
		res := helpers.GetResponse(500, nil, errKategori)
		return c.Status(res.Status).JSON(res)
	}

	kategori.Id, _ = qryKategori.LastInsertId()

	res := helpers.GetResponse(200, fiber.Map{
		"judul":          buku.Judul,
		"penulis":        buku.Penulis,
		"penerbit":       buku.Penerbit,
		"cover":          coverName,
		"back_cover":     backCoverName,
		"jumlah_halaman": buku.JumlahHalaman,
		"tahun_terbit":   buku.TahunTerbit,
		"stok_buku":      buku.StokBuku,
		"deskripsi":      buku.Deskripsi,
		"kategori_id":    buku.KategoriId,
		"created_at":     buku.CreateAt,
		"update_at":      buku.UpdateAt,
	}, nil)
	return c.JSON(res)
}

func BukuPut(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	buku := request.Buku{}
	if err := c.BodyParser(&buku); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	coverName := buku.Cover
	backCoverName := buku.BackCover
	file, errFile := c.FormFile("cover")


	if errFile != nil {
		_, err := db.ExecContext(ctx, `
		UPDATE buku SET Judul = ?, Penulis = ?, Penerbit = ?, JumlahHalaman = ?, TahunTerbit = ?, StokBuku = ?, Deskripsi = ? WHERE Id = ?`, buku.Judul, buku.Penulis, buku.Penerbit, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi, c.Params("id"))
	
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}
	
		buku.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)
	
		res := helpers.GetResponse(200, buku, nil)
		return c.JSON(res)
		}

	fileName := file.Filename
	extension := fileName[strings.LastIndex(fileName, "."):]

	timestamp := time.Now().Format("20060102150405")
	coverName = "cover-" + strings.ReplaceAll(buku.Judul, " ", "") + "-" + timestamp + extension

	c.SaveFile(file, "public/uploads/cover/"+coverName)

	fileNameBackCover := file.Filename
	extensionBackCover := fileNameBackCover[strings.LastIndex(fileNameBackCover, "."):]

	timestampBackCover := time.Now().Format("20060102150405")
	backCoverName = "backcover-" + strings.ReplaceAll(buku.Judul, " ", "") + "-" + timestampBackCover + extensionBackCover

	c.SaveFile(file, "public/uploads/cover/"+backCoverName)

	_, err := db.ExecContext(ctx, `
	UPDATE buku SET Judul = ?, Penulis = ?, Penerbit = ?, Cover = ?, BackCover = ?, JumlahHalaman = ?, TahunTerbit = ?, StokBuku = ?, Deskripsi = ? WHERE Id = ?`, buku.Judul, buku.Penulis, buku.Penerbit, coverName, backCoverName, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi, c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	buku.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, buku, nil)
	return c.JSON(res)
}



func BukuDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM kategori_buku WHERE BukuId = ?`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	_, errbook := db.ExecContext(ctx, `
	DELETE FROM buku WHERE Id = ?`, c.Params("id"))
	if errbook != nil {
		res := helpers.GetResponse(500, nil, errbook)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}
