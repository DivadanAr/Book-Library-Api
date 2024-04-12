package authcontroller

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"net/mail"
	"os"
	"strings"
	"time"
)

type Account struct {
	NamaLengkap    string `json:"nama_lengkap"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Alamat         string `json:"alamat"`
	NomorTelepon   string `json:"nomor_telepon"`
	JenisKelamin   string `json:"jenis_kelamin"`
	ProfilePicture string `json:"profile_picture"`
}

type Users struct {
	Id             int32   `json:"id"`
	Username       *string `json:"username"`
	NamaLengkap    *string `json:"nama_lengkap"`
	Password       *string `json:"password"`
	NomorTelepon   *string `json:"nomor_telepon"`
	Email          *string `json:"email"`
	JenisKelamin   *string `json:"jenis_kelamin"`
	Alamat         *string `json:"alamat"`
	ProfilePicture *string `json:"profile_picture"`
	NameRole       *string `json:"name_role"`
}

func UserGet(c *fiber.Ctx) error {
	var (
		user    = []Users{}
		rowuser = Users{}
		cond    string
		count   int
		search  string
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	search = c.Query("Search")
	if search != "" {
		cond += `WHERE (NameRole LIKE '%` + search + `%')`
	}

	userQry, err := db.QueryContext(ctx, `
	SELECT users.Id, users.Username, users.Email, users.NamaLengkap, users.NomorTelepon, users.JenisKelamin, users.Alamat, users.ProfilePicture, roles.NameRole FROM user_role JOIN users ON users.Id = user_role.UserId JOIN roles ON roles.Id = user_role.RolesId
	`+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer userQry.Close()
	for userQry.Next() {
		err := userQry.Scan(&rowuser.Id, &rowuser.Username, &rowuser.Email, &rowuser.NamaLengkap, &rowuser.NomorTelepon, &rowuser.JenisKelamin, &rowuser.Alamat, &rowuser.ProfilePicture, &rowuser.NameRole)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		user = append(user, rowuser)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(users.Id) 
		FROM users
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"User":  user,
		"Total": count,
	}, nil)
	return c.JSON(res)
}

func UserDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, errrole := db.ExecContext(ctx, `
	DELETE FROM user_role WHERE UserId = ?`, c.Params("id"))
	if errrole != nil {
		res := helpers.GetResponse(500, nil, errrole)
		return c.Status(res.Status).JSON(res)
	} else {

		_, errdelete := db.ExecContext(ctx, `
	DELETE FROM users WHERE Id = ?`, c.Params("id"))
		if errdelete != nil {
			res := helpers.GetResponse(500, nil, errdelete)
			return c.Status(res.Status).JSON(res)
		}
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}

func Register(c *fiber.Ctx) error {
	var users response.Users
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	account := Account{}
	user := request.RoleUser{}

	if err := c.BodyParser(&account); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := mail.ParseAddress(account.Email)
	if err != nil {
		res := helpers.GetResponse(401, nil, errors.New("pastikan email sudah benar"))
		return c.JSON(res)
	}

	var usernameChecker string
	db.QueryRowContext(ctx, `
    SELECT Username
    FROM users
    WHERE Username = ?;
    `, account.Username).Scan(&usernameChecker)

	if usernameChecker != "" {
		res := helpers.GetResponse(409, nil, errors.New("username telah digunakan"))
		return c.Status(res.Status).JSON(res)
	}

	var emailChecker string
	db.QueryRowContext(ctx, `
    SELECT Email
    FROM users
    WHERE Email = ?;
    `, account.Email).Scan(&emailChecker)

	if emailChecker != "" {
		log.Println("email udah ada")
		res := helpers.GetResponse(409, nil, errors.New("email telah digunakan"))
		return c.Status(res.Status).JSON(res)
	}

	password, err := helpers.HashPassword(account.Password, helpers.GenerateSalt())
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	var query string
	query = `
		INSERT INTO users(Email, Username, Password, Alamat, NamaLengkap)
		VALUES(?,?,?,?,?)`

	auth, err := db.ExecContext(ctx, query, account.Email, account.Username, password, account.Alamat, account.NamaLengkap)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}
	userid, _ := auth.LastInsertId()

	qry, err := db.ExecContext(ctx, `
	INSERT INTO user_role (RolesId, UserId)
	VALUES (?, ?)
`, 3, userid)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	user.Id, _ = qry.LastInsertId()
	if errlog := c.BodyParser(&account); errlog != nil {
		return c.JSON(helpers.GetResponse(500, nil, errlog))
	}

	errlog := db.QueryRowContext(ctx, `
        SELECT id, username, password, email
        FROM users
        WHERE email = ? OR username = ?;
        `, account.Email, account.Username).Scan(&users.Id, &users.Username, &users.Password, &users.Email)

	if errlog == sql.ErrNoRows {
		res := helpers.GetResponse(401, nil, errors.New("unauthorized"))
		return c.Status(res.Status).JSON(res)
	} else if errlog != nil {
		res := helpers.GetResponse(500, nil, errlog)
		return c.Status(res.Status).JSON(res)
	}

	if errlog := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(account.Password)); errlog != nil {
		res := helpers.GetResponse(fiber.StatusUnauthorized, nil, errors.New("unauthorized"))
		return c.Status(res.Status).JSON(res)
	}

	claims := jwt.MapClaims{
		"userid":   users.Id,
		"username": users.Username,
		"email":    users.Email,
		"exp":      time.Now().Add(time.Hour * 6).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, errlog := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if errlog != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"UserId":   users.Id,
		"Username": users.Username,
		"Email":    users.Email,
		"Token":    t,
	}, nil)
	return c.Status(res.Status).JSON(res)
}

func SetupAddress(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	user := Account{}
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := db.ExecContext(ctx, `
	UPDATE users SET Alamat = ? WHERE Id = ?`, user.Alamat, c.Locals("UserId"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"alamat":  user.Alamat,
		"user_id": c.Locals("UserId").(float64),
	}, nil)
	return c.JSON(res)
}

func SetupProfile(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	user := Account{}
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	//file, errFile := c.FormFile("profile_picture")
	//
	//if errFile != nil {
	//	fmt.Printf("gaada foto")
	//	//return errFile
	//}
	//
	//fileName := file.Filename
	//extension := fileName[strings.LastIndex(fileName, "."):]
	//
	//timestampProfile := time.Now().Format("20060102150405")
	//profileName := "profile-" + user.NamaLengkap + "-" + user.Username + timestampProfile + extension
	//
	//c.SaveFile(file, "public/uploads/profile/"+profileName)

	_, err := db.ExecContext(ctx, `
	UPDATE users SET Username = ?, NamaLengkap = ?, Email = ?, NomorTelepon = ?, JenisKelamin = ? WHERE Id = ?`, user.Username, user.NamaLengkap, user.Email, user.NomorTelepon, user.JenisKelamin, c.Locals("UserId"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, user, nil)
	return c.JSON(res)
}

func SetupPhotoProfile(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	file, errFile := c.FormFile("profile_picture")

	if errFile != nil {
		return errFile
	}

	fileName := file.Filename
	extension := fileName[strings.LastIndex(fileName, "."):]

	timestampProfile := time.Now().Format("20060102150405")
	profileName := "profile-" + c.Locals("Username").(string) + "-" + c.Locals("Email").(string) + timestampProfile + extension

	c.SaveFile(file, "public/uploads/profile/"+profileName)

	_, err := db.ExecContext(ctx, `
	UPDATE users SET ProfilePicture = ? WHERE Id = ?`, profileName, c.Locals("UserId"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Photo Profile": profileName,
	}, nil)
	return c.JSON(res)
}

//func SetupProfileMerge(c *fiber.Ctx) error {
//	db := database.ConnectDB()
//	defer db.Close()
//	ctx := context.Background()
//
//	user := Account{}
//	if err := c.BodyParser(&user); err != nil {
//		return c.JSON(helpers.GetResponse(500, nil, err))
//	}
//
//	file, errFile := c.FormFile("profile_picture")
//
//	if errFile != nil {
//		//fmt.Printf("gaada foto")
//		return errFile
//	}
//
//	fileName := file.Filename
//	extension := fileName[strings.LastIndex(fileName, "."):]
//
//	timestampProfile := time.Now().Format("20060102150405")
//	profileName := "profile-" + user.NamaLengkap + "-" + user.Username + timestampProfile + extension
//
//	c.SaveFile(file, "public/uploads/profile/"+profileName)
//
//	_, err := db.ExecContext(ctx, `
//	UPDATE users SET Username = ?, NamaLengkap = ?, Email = ?, NomorTelepon = ?, JenisKelamin = ?, ProfilePicture = ? WHERE Id = ?`, user.Username, user.NamaLengkap, user.Email, user.NomorTelepon, user.JenisKelamin, profileName, c.Locals("UserId"))
//
//	if err != nil {
//		res := helpers.GetResponse(500, nil, err)
//		return c.Status(res.Status).JSON(res)
//	}
//
//	res := helpers.GetResponse(200, user, nil)
//	return c.JSON(res)
//}

func Login(c *fiber.Ctx) error {
	var users response.Users
	var role response.Roles
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	account := Account{}
	if errlog := c.BodyParser(&account); errlog != nil {
		return c.JSON(helpers.GetResponse(500, nil, errlog))
	}

	errlog := db.QueryRowContext(ctx, `
        SELECT users.id, users.username, users.password, users.email, roles.id, roles.NameRole
        FROM users JOIN user_role ON users.Id = user_role.UserId JOIN roles ON user_role.RolesId = roles.Id
        WHERE email = ? OR username = ?;
        `, account.Email, account.Username).Scan(&users.Id, &users.Username, &users.Password, &users.Email, &role.Id, &role.NameRole)

	if errlog == sql.ErrNoRows {
		res := helpers.GetResponse(401, nil, errors.New("unauthorized"))
		return c.Status(res.Status).JSON(res)
	} else if errlog != nil {
		res := helpers.GetResponse(500, nil, errlog)
		return c.Status(res.Status).JSON(res)
	}

	if errlog := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(account.Password)); errlog != nil {
		res := helpers.GetResponse(fiber.StatusUnauthorized, nil, errors.New("unauthorized"))
		return c.Status(res.Status).JSON(res)
	}

	claims := jwt.MapClaims{
		"userid":   users.Id,
		"username": users.Username,
		"email":    users.Email,
		"roleid":   role.Id,
		"role":     role.NameRole,
		"exp":      time.Now().Add(time.Hour * 6).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, errlog := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if errlog != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"UserId":   users.Id,
		"Username": users.Username,
		"Email":    users.Email,
		"RoleId":   role.Id,
		"Role":     role.NameRole,
		"Token":    t,
	}, nil)
	return c.Status(res.Status).JSON(res)
}
