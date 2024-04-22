package peminjamancontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	// "main.go/models/request"
	"main.go/models/response"
	"strconv"
	"errors"
)

type Peminjaman struct {
	Id            int64     `json:"id"`
	BukuId            int64     `json:"buku_id"`
	UserId            int64     `json:"user_id"`
	NamaLengkap            string     `json:"nama_lengkap"`
	Email            string     `json:"email"`
	ProfilePicture            string     `json:"profile_picture"`
	Judul         string    `json:"judul"`
	Penulis       string    `json:"penulis"`
	Penerbit      string    `json:"penerbit"`
	Cover         string    `json:"cover"`
	BackCover     string    `json:"back_cover"`
	StokBuku      string    `json:"stock_buku"`
	Deskripsi     *string   `json:"deskripsi"`
	Rating     *string   `json:"rating"`
	Kategori     *string   `json:"kategori"`
	JumlahHalaman int       `json:"jumlah_halaman"`
	TahunTerbit   int       `json:"tahun_terbit"`
	TanggalPeminjaman       string  `json:"tanggal_peminjaman"`
	TanggalPengembalian string `json:"tanggal_pengembalian"`
	StatusPeminjaman string `json:"status_peminjaman"`
}

type ReqPeminjaman struct {
	Id            int64     `json:"id"`
	// UserId            int64     `json:"user_id"`
	BukuId            int64     `json:"buku_id"`
	TanggalPeminjaman       string  `json:"tanggal_peminjaman"`
	TanggalPengembalian string `json:"tanggal_pengembalian"`
	// StatusPeminjaman string `json:"status_peminjaman"`
}

type ReqStatusPeminjaman struct {
	Id            int64     `json:"id"`
	StatusPeminjaman string `json:"status_peminjaman"`
}

func PeminjamanGet(c *fiber.Ctx) error {
	var (
		// buku    []response.Buku
		// rowbuku response.Buku
		peminjaman         = []Peminjaman{}
        rowpeminjaman      = Peminjaman{}
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

	bukuQry, err := db.QueryContext(ctx, `
	SELECT 
    peminjaman.Id, 
    buku.Id, 
    buku.Judul, 
    buku.Penulis, 
    buku.Penerbit, 
    buku.Cover, 
    buku.BackCover, 
    buku.JumlahHalaman, 
    buku.TahunTerbit, 
    buku.StokBuku, 
    buku.Deskripsi, 
    FormatDate(peminjaman.TanggalPeminjaman) AS TanggalPeminjaman, 
    FormatDate(peminjaman.TanggalPengembalian) AS TanggalPengembalian, 
    peminjaman.StatusPeminjaman, 
    peminjaman.UserId, 
    users.NamaLengkap, 
    users.Email, 
    users.ProfilePicture
FROM 
    peminjaman 
JOIN 
    buku ON buku.Id = peminjaman.BukuId 
JOIN 
    users ON users.Id = peminjaman.UserId`+cond+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer bukuQry.Close()
	for bukuQry.Next() {
		err := bukuQry.Scan(&rowpeminjaman.Id, &rowpeminjaman.BukuId, &rowpeminjaman.Judul, &rowpeminjaman.Penulis, &rowpeminjaman.Penerbit, &rowpeminjaman.Cover, &rowpeminjaman.BackCover, &rowpeminjaman.JumlahHalaman, &rowpeminjaman.TahunTerbit, &rowpeminjaman.StokBuku, &rowpeminjaman.Deskripsi, &rowpeminjaman.TanggalPeminjaman, &rowpeminjaman.TanggalPengembalian, &rowpeminjaman.StatusPeminjaman, &rowpeminjaman.UserId, &rowpeminjaman.NamaLengkap, &rowpeminjaman.Email, &rowpeminjaman.ProfilePicture)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		peminjaman = append(peminjaman, rowpeminjaman)
		// buku = append(peminjaman, rowpeminjaman)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(peminjaman.Id)
		FROM peminjaman JOIN buku ON buku.Id = peminjaman.BukuId
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Buku": peminjaman,
		"Total": count,
	}, nil)
	return c.JSON(res)

	//res := helpers.GetResponse(200, fiber.Map{
	//	"UserId":   c.Locals("UserId").(float64),
	//	"Username": c.Locals("Username").(string),
	//}, nil)
	//return c.JSON(res)
}

func PeminjamanStatusGet(c *fiber.Ctx) error {
	var (
		// buku    []response.Buku
		// rowbuku response.Buku
		peminjaman         = []Peminjaman{}
        rowpeminjaman      = Peminjaman{}
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

	bukuQry, err := db.QueryContext(ctx, `
	SELECT peminjaman.Id, buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi, peminjaman.TanggalPeminjaman, peminjaman.TanggalPengembalian, peminjaman.StatusPeminjaman, peminjaman.UserId, users.NamaLengkap, users.Email, users.ProfilePicture
	FROM peminjaman JOIN buku ON buku.Id = peminjaman.BukuId JOIN users ON users.Id = peminjaman.UserId WHERE peminjaman.StatusPeminjaman =?
	`, c.Params("status"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer bukuQry.Close()
	for bukuQry.Next() {
		err := bukuQry.Scan(&rowpeminjaman.Id, &rowpeminjaman.BukuId, &rowpeminjaman.Judul, &rowpeminjaman.Penulis, &rowpeminjaman.Penerbit, &rowpeminjaman.Cover, &rowpeminjaman.BackCover, &rowpeminjaman.JumlahHalaman, &rowpeminjaman.TahunTerbit, &rowpeminjaman.StokBuku, &rowpeminjaman.Deskripsi, &rowpeminjaman.TanggalPeminjaman, &rowpeminjaman.TanggalPengembalian, &rowpeminjaman.StatusPeminjaman, &rowpeminjaman.UserId, &rowpeminjaman.NamaLengkap, &rowpeminjaman.Email, &rowpeminjaman.ProfilePicture)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		peminjaman = append(peminjaman, rowpeminjaman)
		// buku = append(peminjaman, rowpeminjaman)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(peminjaman.Id)
		FROM peminjaman JOIN buku ON buku.Id = peminjaman.BukuId
		WHERE peminjaman.StatusPeminjaman =?
	`, c.Params("status")).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Buku": peminjaman,
		"Total": count,
	}, nil)
	return c.JSON(res)

	//res := helpers.GetResponse(200, fiber.Map{
	//	"UserId":   c.Locals("UserId").(float64),
	//	"Username": c.Locals("Username").(string),
	//}, nil)
	//return c.JSON(res)
}

func PeminjamanStatusGetDetail(c *fiber.Ctx) error {
	var (
		// buku    []response.Buku
		// rowbuku response.Buku
		peminjaman         = []Peminjaman{}
        rowpeminjaman      = Peminjaman{}
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

	bukuQry, err := db.QueryContext(ctx, `
	SELECT
    peminjaman.Id,
    buku1.Id AS bukuId,
    buku1.Judul,
    buku1.Penulis,
    buku1.Penerbit,
    buku1.Cover,
    buku1.BackCover,
    buku1.JumlahHalaman,
    buku1.TahunTerbit,
    buku1.StokBuku,
    buku1.Deskripsi,
    peminjaman.TanggalPeminjaman,
    peminjaman.TanggalPengembalian,
    peminjaman.StatusPeminjaman,
    peminjaman.UserId,
    users.NamaLengkap,
    users.Email,
    users.ProfilePicture,
    kategori.Kategori
FROM
    peminjaman
    JOIN buku AS buku1 ON buku1.Id = peminjaman.BukuId
    JOIN users ON users.Id = peminjaman.UserId
    JOIN kategori_buku ON buku1.Id = kategori_buku.BukuId
    JOIN kategori ON kategori.Id = kategori_buku.KategoriId
WHERE
    (peminjaman.StatusPeminjaman != 'reject' AND peminjaman.StatusPeminjaman != 'done')
    AND peminjaman.UserId = ?
	`, c.Locals("UserId"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer bukuQry.Close()
	for bukuQry.Next() {
		err := bukuQry.Scan(&rowpeminjaman.Id, &rowpeminjaman.BukuId, &rowpeminjaman.Judul, &rowpeminjaman.Penulis, &rowpeminjaman.Penerbit, &rowpeminjaman.Cover, &rowpeminjaman.BackCover, &rowpeminjaman.JumlahHalaman, &rowpeminjaman.TahunTerbit, &rowpeminjaman.StokBuku, &rowpeminjaman.Deskripsi, &rowpeminjaman.TanggalPeminjaman, &rowpeminjaman.TanggalPengembalian, &rowpeminjaman.StatusPeminjaman, &rowpeminjaman.UserId, &rowpeminjaman.NamaLengkap, &rowpeminjaman.Email, &rowpeminjaman.ProfilePicture, &rowpeminjaman.Kategori)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		peminjaman = append(peminjaman, rowpeminjaman)
		// buku = append(peminjaman, rowpeminjaman)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(peminjaman.Id)
		FROM peminjaman JOIN buku ON buku.Id = peminjaman.BukuId
		WHERE (StatusPeminjaman != 'reject' AND StatusPeminjaman != 'done') AND peminjaman.UserId =?
	`, c.Locals("UserId")).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Buku": peminjaman,
		"Total": count,
	}, nil)
	return c.JSON(res)

	//res := helpers.GetResponse(200, fiber.Map{
	//	"UserId":   c.Locals("UserId").(float64),
	//	"Username": c.Locals("Username").(string),
	//}, nil)
	//return c.JSON(res)
}

func PeminjamanDetail(c *fiber.Ctx) error {
	var (
		buku    []response.Buku
		rowbuku response.Buku
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	bukuQry, err := db.QueryContext(ctx, `
	SELECT buku.Id, buku.Judul, buku.Penulis, buku.Penerbit, buku.Cover, buku.BackCover, buku.JumlahHalaman, buku.TahunTerbit, buku.StokBuku, buku.Deskripsi
	FROM peminjaman JOIN buku ON buku.Id = peminjaman.BukuId WHERE peminjaman.UserId = ?
	`, c.Locals("UserId"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer bukuQry.Close()
	for bukuQry.Next() {
		err := bukuQry.Scan(&rowbuku.Id, &rowbuku.Judul, &rowbuku.Penulis, &rowbuku.Penerbit, &rowbuku.Cover, &rowbuku.BackCover, &rowbuku.JumlahHalaman, &rowbuku.TahunTerbit, &rowbuku.StokBuku, &rowbuku.Deskripsi)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		buku = append(buku, rowbuku)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Buku": buku,
	}, nil)
	return c.JSON(res)
}

func RequestBuku(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	peminjaman := ReqPeminjaman{}
	if err := c.BodyParser(&peminjaman); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	var userChecker int
	db.QueryRowContext(ctx, `
	SELECT COUNT(*) FROM peminjaman
	WHERE (StatusPeminjaman != 'reject' AND StatusPeminjaman != 'done') AND UserId = ? 
    `, c.Locals("UserId")).Scan(&userChecker)

	if userChecker != 0 {
		res := helpers.GetResponse(409, nil, errors.New("selesaikan peminjaman buku sebelumnya terlebih dahulu"))
		return c.Status(res.Status).JSON(res)
	}
	
	qry, err := db.ExecContext(ctx, `
	INSERT INTO peminjaman (BukuId, TanggalPeminjaman, TanggalPengembalian, StatusPeminjaman, UserId)
	VALUES (?,?,?,?,?)
`, peminjaman.BukuId, peminjaman.TanggalPeminjaman, peminjaman.TanggalPengembalian, "request", c.Locals("UserId").(float64))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	peminjaman.Id, _ = qry.LastInsertId()

	res := helpers.GetResponse(200, fiber.Map{
		"buku_id": peminjaman.BukuId,
		"tanggal_peminjaman": peminjaman.TanggalPeminjaman,
		"tanggal_pengembalian": peminjaman.TanggalPengembalian,
		"status": "request",
	}, nil)
	return c.JSON(res)
}

func ChangeStatusPeminjaman(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	peminjaman := ReqStatusPeminjaman{}
	if err := c.BodyParser(&peminjaman); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

_, err := db.ExecContext(ctx, `
	UPDATE peminjaman SET StatusPeminjaman = ? WHERE Id = ?`, peminjaman.StatusPeminjaman, c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	peminjaman.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, fiber.Map{
		"id":      peminjaman.Id,
		"status_peminjaman": peminjaman.StatusPeminjaman,
	}, nil)
	return c.JSON(res)
}

func PeminjamanDelete(c *fiber.Ctx) error {
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
