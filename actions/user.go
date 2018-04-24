package actions

import (
	"database/sql"
	"strings"

	crypt "github.com/amoghe/go-crypt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	//"golang.org/x/crypto/bcrypt"
	. "buf1/actions/funclib"
	"buf1/models"

	"github.com/pkg/errors"
)

// UserAdd default implementation.
func UserAdd(c buffalo.Context) error {
	if SessionValid(c) != nil {
		c.Set("user", models.User{})
		return c.Render(200, r.HTML("user/add.html"))
	} else {
		return errors.New("A valid session exists! Please log out first!")
	}
}

func UserAddcommit(c buffalo.Context) error {
	if SessionValid(c) != nil {
		u := models.User{}

		tx := c.Value("tx").(*pop.Connection)

		if err := c.Bind(&u); err != nil {
			return err
		}

		if u.Password != "" && u.PasswordConfirmation != "" && u.Password == u.PasswordConfirmation {
			var crypt_error error
			if salt, err_salt := CreateSalt("SHA-512"); err_salt != nil {
				return errors.New("Password salt failed:" + err_salt.Error())
			} else {
				if u.PasswordHash, crypt_error = crypt.Crypt(u.Password, salt); crypt_error != nil {
					return errors.New("Password crypt failed:" + crypt_error.Error())
				}
			}
		}

		logger.Infof("User: %v", u)

		if err := tx.Save(&u); err != nil {
			return err
		}
		c.Session().Clear()
		c.Flash().Add("success", "User successfully saved.")
		return c.Redirect(302, "/user/login")
	} else {
		return errors.New("A valid session exists! Please log out first!")
	}
}

// UserGet default implementation.
func UserGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("user/get.html"))
}

// UserRemove default implementation.
func UserRemove(c buffalo.Context) error {
	if SessionValid(c) == nil {
		c.Set("current_user_id", c.Session().Get("current_user_id"))
		u := models.User{}

		tx := c.Value("tx").(*pop.Connection)

		// find a user with the email
		if err := tx.Where("id = ?", c.Session().Get("current_user_id").(uuid.UUID)).First(&u); err != nil {
			return errors.New("User not found." + err.Error())
		}

		c.Set("user", u)
		return c.Render(200, r.HTML("user/remove.html"))
	} else {
		return errors.New("No valid session")
	}

}

func UserRemoveconfirm(c buffalo.Context) error {
	if SessionValid(c) == nil {
		u := models.User{}
		tx := c.Value("tx").(*pop.Connection)

		if err := tx.Where("id = ?", c.Session().Get("current_user_id").(uuid.UUID)).First(&u); err != nil {
			return errors.New("User not found." + err.Error())
		}

		if err := tx.Destroy(&u); err != nil {
			return errors.New("User could not be deleted." + err.Error())
		}

		c.Session().Clear()
		c.Flash().Add("success", "Account successfully removed! Session destroyed!")
		return c.Redirect(302, "/")
	} else {
		return errors.New("No valid session")
	}
}

// UserEdit default implementation.
func UserEdit(c buffalo.Context) error {
	if SessionValid(c) == nil {
		c.Set("current_user_id", c.Session().Get("current_user_id"))
		u := models.User{}

		tx := c.Value("tx").(*pop.Connection)

		// find a user with the email
		if err := tx.Where("id = ?", c.Session().Get("current_user_id").(uuid.UUID)).First(&u); err != nil {
			return errors.New("User not found." + err.Error())
		}

		c.Set("user", u)
		return c.Render(200, r.HTML("user/edit.html"))
	} else {
		return errors.New("No valid session")
	}
}

func UserEditcommit(c buffalo.Context) error {
	if SessionValid(c) == nil {
		u := models.User{}

		tx := c.Value("tx").(*pop.Connection)
		if err := tx.Where("id = ?", c.Session().Get("current_user_id").(uuid.UUID)).First(&u); err != nil {
			return errors.New("User not found." + err.Error())
		}

		if err := c.Bind(&u); err != nil {
			return err
		}

		if u.Password != "" && u.PasswordConfirmation != "" && u.Password == u.PasswordConfirmation {
			var crypt_error error
			if salt, err_salt := CreateSalt("SHA-512"); err_salt != nil {
				return errors.New("Password salt failed:" + err_salt.Error())
			} else {
				if u.PasswordHash, crypt_error = crypt.Crypt(u.Password, salt); crypt_error != nil {
					return errors.New("Password crypt failed:" + crypt_error.Error())
				}
			}
		}

		logger.Infof("User: %v", u)

		if err := tx.Update(&u); err != nil {
			return err
		}
		c.Flash().Add("success", "User successfully saved.")
		return c.Redirect(302, "/user/edit")
	} else {
		return errors.New("No valid session")
	}
}

// UserLogin default implementation.
func UserLogin(c buffalo.Context) error {
	c.Set("errors", nil)
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("user/login.html"))
}

// UserLogincommit default implementation.
func UserLogincommit(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("user/login.html"))
	}

	// find a user with the email or username
	//if err := tx.Where("(email = ? or username = ?)", strings.ToLower(u.Email), strings.ToLower(u.Username)).First(u); err != nil {
	if err := tx.Where("(email = ? or username = ?)", strings.ToLower(u.Email), u.Username).First(u); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	//err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))

	//logger.Infof("PasswordHash: %s, Password: %s", u.PasswordHash, u.Password)
	if err := CryptCheck(u.PasswordHash, u.Password); err != nil {
		logger.Infof("Error: %s", err.Error())
		return bad()
	}
	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "You have been logged in!")
	return c.Redirect(302, "/")
}

// UserLogout default implementation.
func UserLogout(c buffalo.Context) error {
	if SessionValid(c) == nil {
		c.Session().Clear()
		c.Flash().Add("success", "You have been logged out!")
		return c.Redirect(302, "/")
	} else {
		c.Session().Clear()
		c.Flash().Add("success", "Error while logging out!")
		//return errors.New("Error while logging out")
		return c.Redirect(302, "/")
	}
}
