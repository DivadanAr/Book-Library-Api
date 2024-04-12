package authcontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main.go/config/database"
	"main.go/helpers"
	"main.go/models/request"
	"main.go/models/response"
	"strconv"
)

func RolesGet(c *fiber.Ctx) error {
	var (
		role    []response.Roles
		rowrole response.Roles
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

	roleQry, err := db.QueryContext(ctx, `
	SELECT Id, NameRole 
	FROM roles 
	`+cond+
		helpers.Limit(c.Query("Limit"))+" "+helpers.Offset(c.Query("Offset"))+`;`)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer roleQry.Close()
	for roleQry.Next() {
		err := roleQry.Scan(&rowrole.Id, &rowrole.NameRole)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		role = append(role, rowrole)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(roles.Id) 
		FROM roles
		`).Scan(&count)
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Roles": role,
		"Total": count,
	}, nil)
	return c.JSON(res)
}

func RolesDetail(c *fiber.Ctx) error {
	var (
		role    []response.Roles
		rowrole response.Roles
	)

	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	roleQry, err := db.QueryContext(ctx, `
	SELECT Id, NameRole 
	FROM roles WHERE Id = ?
	`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	defer roleQry.Close()
	for roleQry.Next() {
		err := roleQry.Scan(&rowrole.Id, &rowrole.NameRole)
		if err != nil {
			res := helpers.GetResponse(500, nil, err)
			return c.Status(res.Status).JSON(res)
		}

		role = append(role, rowrole)
	}

	res := helpers.GetResponse(200, fiber.Map{
		"Roles": role,
	}, nil)
	return c.JSON(res)
}

func RolesPost(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	role := request.Roles{}
	if err := c.BodyParser(&role); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	qry, err := db.ExecContext(ctx, `
	INSERT INTO roles (NameRole)
	VALUES (?)
`, role.NameRole)

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	role.Id, _ = qry.LastInsertId()

	res := helpers.GetResponse(200, role, nil)
	return c.JSON(res)
}

func RolesPut(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	role := request.Roles{}
	if err := c.BodyParser(&role); err != nil {
		return c.JSON(helpers.GetResponse(500, nil, err))
	}

	_, err := db.ExecContext(ctx, `
	UPDATE roles SET NameRole = ? WHERE Id = ?`, role.NameRole, c.Params("id"))

	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	role.Id, _ = strconv.ParseInt(c.Params("id"), 10, 64)

	res := helpers.GetResponse(200, role, nil)
	return c.JSON(res)
}

func RolesDelete(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
	DELETE FROM roles WHERE Id = ?`, c.Params("id"))
	if err != nil {
		res := helpers.GetResponse(500, nil, err)
		return c.Status(res.Status).JSON(res)
	}

	res := helpers.GetResponse(200, "Delete Success", nil)
	return c.JSON(res)
}
