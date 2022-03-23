package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllUsers(c echo.Context) error {
	db := connect()
	defer db.Close()

	query := "SELECT ID, name, age, address, email FROM users"

	ID, _ := strconv.Atoi(c.QueryParam("id"))
	name := c.QueryParam("name")
	age, _ := strconv.Atoi(c.QueryParam("age"))
	address := c.QueryParam("address")
	email := c.QueryParam("email")

	if ID != 0 {
		query += " WHERE id = " + strconv.Itoa(ID)
	} else if name != "" {
		if ID != 0 {
			query += ", name = '" + name + "'"
		} else {
			query += " WHERE name = '" + name + "'"
		}
	} else if age != 0 {
		if ID != 0 || name != "" {
			query += ", age = " + strconv.Itoa(age)
		} else {
			query += " WHERE age = " + strconv.Itoa(age)
		}
	} else if address != "" {
		if ID != 0 || name != "" || age != 0 {
			query += ", address = '" + strconv.Itoa(age) + "'"
		} else {
			query += " WHERE address = '" + address + "'"
		}
	} else if email != "" {
		if ID != 0 || name != "" || age != 0 || address != "" {
			query += ", email = '" + strconv.Itoa(age) + "'"
		} else {
			query += " WHERE email = '" + email + "'"
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		return sendNotFoundResponse(c, "Table Not Found")
	}

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			return sendBadRequestResponse(c, "Error Field Undefined")
		} else {
			users = append(users, user)
		}
	}

	if len(users) != 0 {
		return sendUserSuccessResponse(c, users, "Get Success")
	} else {
		return sendBadRequestResponse(c, "Error Array Size Not Correct")
	}
}

func InsertUser(c echo.Context) error {
	db := connect()
	defer db.Close()

	name := c.FormValue("name")
	age, _ := strconv.Atoi(c.FormValue("age"))
	address := c.FormValue("address")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if name == "" || age == 0 || address == "" || email == "" || password == "" {
		return sendNotFoundResponse(c, "Value Not Found")
	}

	password = GetMD5hash(password)

	result, errQuery := db.Exec("INSERT INTO users(name, age, address, email, password) VALUES (?, ?, ?, ?, ?)", name, age, address, email, password)

	userID, _ := result.LastInsertId()
	rows, _ := db.Query("SELECT id, name, age, address, email FROM users WHERE id=?", userID)

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			return sendBadRequestResponse(c, "Error Field Undefined")
		} else {
			users = append(users, user)
		}
	}

	if errQuery == nil {
		return sendUserSuccessResponse(c, users, "Insert Success")
	} else {
		return sendBadRequestResponse(c, "Error Can Not Insert")
	}
}

func UpdateUser(c echo.Context) error {
	db := connect()
	defer db.Close()

	userID := c.Param("userID")
	name := c.FormValue("name")
	age, _ := strconv.Atoi(c.FormValue("age"))
	address := c.FormValue("address")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if name == "" || age == 0 || address == "" || email == "" || password == "" {
		return sendNotFoundResponse(c, "Value Not Found")
	}

	password = GetMD5hash(password)

	result, errQuery := db.Exec("UPDATE users SET name=?, age=?, address=?, email=?, password=? WHERE id=?", name, age, address, email, password, userID)
	rows, _ := db.Query("SELECT id, name, age, address, email FROM users WHERE id=?", userID)

	num, _ := result.RowsAffected()

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			sendBadRequestResponse(c, "Error Field Undefined")
		} else {
			users = append(users, user)
		}
	}

	if errQuery == nil {
		if num == 0 {
			return sendBadRequestResponse(c, "Error 0 Rows Affected")
		} else {
			return sendUserSuccessResponse(c, users, "Update Success")
		}
	} else {
		return sendBadRequestResponse(c, "Error Can Not Update")
	}
}

func DeleteUser(c echo.Context) error {
	db := connect()
	defer db.Close()

	userID := c.Param("userID")

	rows, _ := db.Query("SELECT id, name, age, address, email FROM users WHERE id=?", userID)

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			return sendBadRequestResponse(c, "Error Field Undefined")
		} else {
			users = append(users, user)
		}
	}

	_, errQuery := db.Exec("DELETE FROM transactions WHERE userID=?", userID)

	if errQuery == nil {
		result2, errQuery2 := db.Exec("DELETE FROM users WHERE id=?", userID)

		num2, _ := result2.RowsAffected()

		if errQuery2 == nil {
			if num2 == 0 {
				return sendBadRequestResponse(c, "Error 0 Rows Affected")
			} else {
				return sendUserSuccessResponse(c, users, "Delete Success")
			}
		} else {
			return sendBadRequestResponse(c, "Error Can Not Delete")
		}
	} else {
		return sendBadRequestResponse(c, "Error Can Not Delete")
	}
}

func Login(c echo.Context) error {
	db := connect()
	defer db.Close()

	email := c.FormValue("email")
	password := c.FormValue("password")

	password = GetMD5hash(password)

	rows, err := db.Query("SELECT ID, name, age, address, email FROM users WHERE email=? AND password=?", email, password)
	if err != nil {
		return sendNotFoundResponse(c, "Table Not Found")
	}

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			return sendBadRequestResponse(c, "Error Field Undefined")
		} else {
			users = append(users, user)
		}
	}

	if len(users) != 0 {
		generateToken(c, user.ID, user.Name, 0)
		return sendLoginLogoutSuccessResponse(c, true)
	} else {
		return sendNotFoundResponse(c, "User Not Found")
	}
}

func Logout(c echo.Context) error {
	resetUserToken(c)
	return sendLoginLogoutSuccessResponse(c, false)
}

func GetMD5hash(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
