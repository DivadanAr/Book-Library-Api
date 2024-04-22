package routes

import (
	"github.com/gofiber/fiber/v2"
	"main.go/controllers/authcontroller"
	"main.go/controllers/bukucontroller"
	"main.go/controllers/kategoribukucontroller"
	"main.go/controllers/kategoricontroller"
	"main.go/controllers/koleksipribadicontroller"
	"main.go/controllers/ulasanbukucontroller"
	"main.go/controllers/peminjamancontroller"
	"main.go/controllers/datacontroller"
	"main.go/helpers"
)

func Init(r *fiber.App) {
	api := r.Group("/Api")
	auth := api.Group("/Auth")
	buku := api.Group("/buku")
	data := api.Group("/data")
	kategori := api.Group("/kategori")
	kategoriBuku := api.Group("/kategori-buku")
	koleksiPribadi := api.Group("/koleksi-pribadi")
	ulasanBuku := api.Group("/ulasan-buku")
	peminjaman := api.Group("/peminjaman")

	auth.Get("Users", authcontroller.UserGet)
	auth.Delete("Users/:id", authcontroller.UserDelete)
	auth.Static("/Get-Profile/", "./public/uploads/profile")
	auth.Post("Register", authcontroller.Register)
	auth.Post("Login", authcontroller.Login)
	auth.Put("Setup-Address", helpers.VerifyToken, authcontroller.SetupAddress)
	auth.Put("Setup-Profile", helpers.VerifyToken, authcontroller.SetupProfile)
	auth.Put("Setup-Profile-picture", helpers.VerifyToken, authcontroller.SetupPhotoProfile)

	auth.Get("/Roles/Get", authcontroller.RolesGet)
	auth.Get("/Roles/Detail/:id", authcontroller.RolesDetail)
	auth.Post("/Roles/Store", authcontroller.RolesPost)
	auth.Put("/Roles/Put/:id", authcontroller.RolesPut)
	auth.Delete("/Roles/Delete/:id", authcontroller.RolesDelete)

	buku.Get("/Get", bukucontroller.BukuGet)
	buku.Get("/Get-Koleksi/:id", helpers.VerifyToken, bukucontroller.BookmarkDetail)
	buku.Delete("/Delete-Koleksi/:id", helpers.VerifyToken, bukucontroller.BookmarkDelete)
	buku.Static("/Get-Cover/", "./public/uploads/cover")
	buku.Get("/Detail/:id", bukucontroller.BukuDetail)
	buku.Post("/Store", bukucontroller.BukuPost)
	buku.Put("/Put/:id", bukucontroller.BukuPut)
	buku.Delete("/Delete/:id", bukucontroller.BukuDelete)

	kategori.Get("/Get", kategoricontroller.KategoriGet)
	kategori.Get("/Detail/:id", kategoricontroller.KategoriDetail)
	kategori.Post("/Store", kategoricontroller.KategoriPost)
	kategori.Put("/Put/:id", kategoricontroller.KategoriPut)
	kategori.Delete("/Delete/:id", kategoricontroller.KategoriDelete)

	kategoriBuku.Get("/Get", kategoribukucontroller.KategoriBukuGet)
	kategoriBuku.Get("/Detail/:id", kategoribukucontroller.KategoriBukuDetail)
	kategoriBuku.Get("/Detail/:id", kategoribukucontroller.KategoriBukuDetail)
	kategoriBuku.Post("/Store", kategoribukucontroller.KategoriBukuPost)
	kategoriBuku.Put("/Put/:id", kategoribukucontroller.KategoriBukuPut)
	kategoriBuku.Delete("/Delete/:id", kategoribukucontroller.KategoriBukuDelete)

	koleksiPribadi.Get("/Get", koleksipribadicontroller.KoleksiPribadiGet)
	koleksiPribadi.Get("/Detail", helpers.VerifyToken, koleksipribadicontroller.KoleksiPribadiDetail)
	koleksiPribadi.Post("/Store", helpers.VerifyToken, koleksipribadicontroller.KoleksiPribadiPost)
	koleksiPribadi.Put("/Put/:id", helpers.VerifyToken, koleksipribadicontroller.KoleksiPribadiPut)
	koleksiPribadi.Delete("/Delete/:id", koleksipribadicontroller.KoleksiPribadiDelete)

	ulasanBuku.Get("/Get", ulasanbukucontroller.UlasanBukuGet)
	ulasanBuku.Get("/Detail/:id", ulasanbukucontroller.UlasanBukuDetail)
	ulasanBuku.Post("/Store", helpers.VerifyToken, ulasanbukucontroller.UlasanBukuPost)
	ulasanBuku.Put("/Put/:id", helpers.VerifyToken, ulasanbukucontroller.UlasanBukuPut)
	ulasanBuku.Delete("/Delete/:id", ulasanbukucontroller.UlasanBukuDelete)

	peminjaman.Get("/Get", peminjamancontroller.PeminjamanGet)
	peminjaman.Get("/Get/Detail", helpers.VerifyToken, peminjamancontroller.PeminjamanStatusGetDetail)
	peminjaman.Get("/Get/:status", peminjamancontroller.PeminjamanStatusGet)
	peminjaman.Get("/Detail/:id", peminjamancontroller.PeminjamanDetail)
	peminjaman.Post("/req-buku", helpers.VerifyToken, peminjamancontroller.RequestBuku)
	peminjaman.Put("/change-status/:id", peminjamancontroller.ChangeStatusPeminjaman)
	peminjaman.Delete("/Delete/:id", peminjamancontroller.PeminjamanDelete)

	data.Get("/Get", datacontroller.CountDataGet)
}
