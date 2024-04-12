package entity

import "time"

type Users struct {
	Id             int32          `json:"id" gorm:"primaryKey"`
	Username       string         `json:"username" gorm:"type:varchar(255)"`
	NamaLengkap    string         `json:"nama_lengkap" gorm:"type:varchar(255)"`
	Password       string         `json:"password" gorm:"type:varchar(255)"`
	NomorTelepon   string         `json:"nomor_telepon" gorm:"type:varchar(255)"`
	JenisKelamin   int            `json:"jenis_kelamin" gorm:"type:varchar(1);default:0"`
	Email          string         `json:"email" gorm:"type:varchar(255)"`
	Alamat         string         `json:"alamat" gorm:"type:text"`
	KoleksiPribadi KoleksiPribadi `json:"koleksi_pribadi" gorm:"foreignKey:UserId"`
	UlasanBuku     UlasanBuku     `json:"ulasan_buku" gorm:"foreignKey:UserId"`
	UserRoles      UserRoles      `json:"user_roles" gorm:"foreignKey:UserId"`
	CreateAt       time.Time      `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt       time.Time      `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (Users) TableName() string {
	return "users"
}

type Roles struct {
	Id        int32     `json:"id" gorm:"primaryKey"`
	NameRole  int32     `json:"name_role" gorm:"type:varchar(255)"`
	UserRoles UserRoles `json:"user_roles" gorm:"foreignKey:RolesId"`
	CreateAt  time.Time `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (Roles) TableName() string {
	return "roles"
}

type UserRoles struct {
	Id       int32     `json:"id" gorm:"primaryKey"`
	UserId   int32     `json:"user_id" gorm:"foreignKey:UserId"`
	RolesId  int32     `json:"roles_id" gorm:"foreignKey:RolesId"`
	CreateAt time.Time `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (UserRoles) TableName() string {
	return "user_role"
}

type Buku struct {
	Id             int32          `json:"id" gorm:"primaryKey"`
	Judul          string         `json:"judul" gorm:"type:varchar(255)"`
	Penulis        string         `json:"penulis" gorm:"type:varchar(255)"`
	Penerbit       string         `json:"penerbit" gorm:"type:varchar(255)"`
	Cover          string         `json:"cover" gorm:"type:varchar(255)"`
	BackCover      string         `json:"back_cover" gorm:"type:varchar(255) NULL"`
	JumlahHalaman  int            `json:"jumlah_halaman" gorm:"type:int(11)"`
	TahunTerbit    int            `json:"tahun_terbit" gorm:"type:int(11)"`
	CreateAt       time.Time      `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt       time.Time      `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	KategoriBuku   KategoriBuku   `json:"kategori_buku" gorm:"foreignKey:BukuId"`
	KoleksiPribadi KoleksiPribadi `json:"koleksi_pribadi" gorm:"foreignKey:BukuId"`
	UlasanBuku     UlasanBuku     `json:"ulasan_buku" gorm:"foreignKey:BukuId"`
}

func (Buku) TableName() string {
	return "buku"
}

type Kategori struct {
	Id           int32          `json:"id" gorm:"primaryKey"`
	Kategori     string         `json:"kategori" gorm:"type:varchar(255)"`
	KategoriBuku []KategoriBuku `json:"kategori_buku" gorm:"foreignKey:KategoriId"`
	CreateAt     time.Time      `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt     time.Time      `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (Kategori) TableName() string {
	return "kategori"
}

type KategoriBuku struct {
	Id         int32     `json:"id" gorm:"primaryKey"`
	BukuId     int32     `json:"buku_id"`
	KategoriId int32     `json:"kategori_id"`
	CreateAt   time.Time `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt   time.Time `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (KategoriBuku) TableName() string {
	return "kategori_buku"
}

type KoleksiPribadi struct {
	Id       int32     `json:"id" gorm:"primaryKey"`
	BukuId   int32     `json:"buku_id"`
	UserId   int32     `json:"user_id"`
	CreateAt time.Time `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (KoleksiPribadi) TableName() string {
	return "koleksi_pribadi"
}

type UlasanBuku struct {
	Id       int32     `json:"id" gorm:"primaryKey"`
	BukuId   int32     `json:"buku_id"`
	UserId   int32     `json:"user_id"`
	Ulasan   string    `json:"ulasan" grom:"type:text"`
	Rating   string    `json:"rating" grom:"type:int(11)"`
	CreateAt time.Time `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (UlasanBuku) TableName() string {
	return "ulasan_buku"
}

type Peminjaman struct {
	Id       int32     `json:"id" gorm:"primaryKey"`
	BukuId   int32     `json:"buku_id"`
	UserId   int32     `json:"user_id"`
	Ulasan   string    `json:"ulasan" grom:"type:text"`
	Rating   string    `json:"rating" grom:"type:int(11)"`
	CreateAt time.Time `json:"create_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (Peminjaman) TableName() string {
	return "peminjaman"
}
