package helpers

import (
	"crypto/rand"
	"fmt"
	// "math/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() []byte {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	return salt
}
func HashPassword(password string, salt []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func RandStr(n int) (string, error) {
	char := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	result := make([]rune, n)
	for i, v := range b {
		result[i] = char[int(v)%len(char)]
	}

	return string(result), nil
}

//func SingleFileUpload(c *fiber.Ctx, inputName string) (string, string, string, string) {
//	var fileName, name, ext, typeFile string
//	file, _ := c.FormFile(inputName)
//	if file != nil {
//		name = file.Filename
//		dateNow := time.Now().Unix()
//		ext = filepath.Ext(strings.ToLower(name))
//		randName, _ := RandStr(16)
//		fileName = strings.ToLower(randName) + strconv.FormatInt(dateNow, 10) + "." + ext[1:]
//
//		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
//			src, _ := file.Open()
//			img, _ := imaging.Decode(src)
//			width, _ := strconv.Atoi(os.Getenv("WIDTH_IMAGE"))
//			height, _ := strconv.Atoi(os.Getenv("HEIGHT_IMAGE"))
//			quality, _ := strconv.Atoi(os.Getenv("QUALITY_IMAGE"))
//			compressed := imaging.Resize(img, width, height, imaging.Lanczos)
//			imaging.Save(compressed, fmt.Sprintf("assets/%s", fileName), imaging.JPEGQuality(quality))
//		} else {
//			if err := c.SaveFile(file, fmt.Sprintf("assets/%s", fileName)); err != nil {
//				log.Println(err)
//			}
//		}
//
//		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
//			typeFile = "Gambar"
//		} else if ext == ".mp4" || ext == ".mkv" {
//			typeFile = "Video"
//		} else {
//			typeFile = "Dokumen"
//		}
//
//	} else {
//		fileName = ""
//		name = ""
//		ext = ""
//		typeFile = ""
//	}
//	return fileName, name, ext, typeFile
//}

func Csql(cond string, field string, value interface{}) string {
	var result string
	switch v := value.(type) {
	case string:
		if value != "" {
			result = fmt.Sprintf("%s%s%s%s%v%s", cond, " ", field, " '", v, "'")
		} else {
			result = ""
		}
	case int:
		result = fmt.Sprintf("%s%s%s%s%v", cond, " ", field, " ", v)
	case float64:
		result = fmt.Sprintf("%s%s%s%s%v", cond, " ", field, " ", v)
	case bool:
		result = fmt.Sprintf("%s%s%s%s%v", cond, " ", field, " ", v)
	default:
		result = ""
	}

	return result
}

func Limit(n interface{}) string {
	var result string
	switch v := n.(type) {
	case string:
		if n != "" {
			result = fmt.Sprintf("%s%s", "LIMIT ", v)
		} else {
			result = fmt.Sprintf("%s%s", " ", " ")
		}
	default:
		result = fmt.Sprintf("%s%s", " ", " ")
	}

	return result
}

func Offset(n interface{}) string {
	var result string
	switch v := n.(type) {
	case string:
		if n != "" {
			result = fmt.Sprintf("%s%s", "OFFSET ", v)
		} else {
			result = fmt.Sprintf("%s%s", " ", " ")
		}
	default:
		result = fmt.Sprintf("%s%s", " ", " ")
	}

	return result
}
