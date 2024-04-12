package response

import "time"

type Response struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Token   string      `json:"Token,omitempty"`
	Data    interface{} `json:"Data"`
}

type Users struct {
	Id          int32  `json:"id"`
	Username    string `json:"username"`
	NamaLengkap string `json:"nama_lengkap"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

type Roles struct {
	Id       int64  `json:"id"`
	NameRole string `json:"name_role"`
}

type Buku struct {
	Id            int64     `json:"id"`
	Judul         string    `json:"judul"`
	Penulis       string    `json:"penulis"`
	Penerbit      string    `json:"penerbit"`
	Cover         string    `json:"cover"`
	BackCover     string    `json:"back_cover"`
	JumlahHalaman int       `json:"jumlah_halaman"`
	TahunTerbit   int       `json:"tahun_terbit"`
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
	UserId int   `json:"user_id"`
	BukuId int   `json:"buku_id"`
}

type UlasanBuku struct {
	Id     int64  `json:"id"`
	UserId int    `json:"user_id"`
	BukuId int    `json:"buku_id"`
	Ulasan string `json:"ulasan"`
	Rating int    `json:"rating"`
}
