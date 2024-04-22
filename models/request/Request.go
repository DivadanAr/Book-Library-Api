package request

import "time"

type Claims struct {
	Username string `json:"username"`
	SiteId   string `json:"siteid"`
	Exp      string `json:"exp"`
}

type Roles struct {
	Id       int64  `json:"id"`
	NameRole string `json:"name_role"`
}

type RoleUser struct {
	Id     int64  `json:"id"`
	RoleId string `json:"role_id"`
	UserId string `json:"user_id"`
}

type Buku struct {
	Id            int64     `json:"id"`
	Judul         string    `json:"judul"`
	Penulis       string    `json:"penulis"`
	Penerbit      string    `json:"penerbit"`
	Cover         string    `json:"cover"`
	BackCover     string    `json:"back_cover"`
	JumlahHalaman int       `json:"jumlah_halaman"`
	StokBuku      string    `json:"stok_buku"`
	Deskripsi      string    `json:"deskripsi"`
	TahunTerbit   int       `json:"tahun_terbit"`
	KategoriId    int       `json:"kategori_id"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}

type Kategori struct {
	Id       int64  `json:"id"`
	Kategori string `json:"kategori"`
}

type KategoriBuku struct {
	Id         int64 `json:"id"`
	KategoriId int   `json:"kategori_id"`
	BukuId     int   `json:"buku_id"`
}

type KoleksiPribadi struct {
	Id     int64 `json:"id"`
	BukuId int   `json:"buku_id"`
}

type UlasanBuku struct {
	Id     int64  `json:"id"`
	UserId int    `json:"user_id"`
	BukuId int    `json:"buku_id"`
	Ulasan string `json:"ulasan"`
	Rating int    `json:"rating"`
}
