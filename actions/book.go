package actions

import (
	"buf1/models"

	. "buf1/actions/funclib"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

// BookAdd default implementation.
func BookAdd(c buffalo.Context) error {
	if SessionValid(c) == nil {
		c.Set("book", models.Book{})
		return c.Render(200, r.HTML("book/add.html"))
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}

// BookRemove default implementation.
func BookRemove(c buffalo.Context) error {
	if SessionValid(c) == nil {
		if err := c.Request().ParseForm(); err != nil {
			return err
		}
		//logger.Infof("Request ID: %#v", c.Request().Form["ID"])
		id := c.Request().Form["ID"][0]
		b := models.Book{}

		tx := c.Value("tx").(*pop.Connection)

		// find a book
		if err := tx.Where("id = ?", id).First(&b); err != nil {
			return errors.New("Book not found." + err.Error())
		}

		c.Set("book", b)
		return c.Render(200, r.HTML("book/remove.html"))
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}

// BookEdit default implementation.
func BookEdit(c buffalo.Context) error {
	if SessionValid(c) == nil {
		if err := c.Request().ParseForm(); err != nil {
			return err
		}
		id := c.Request().Form["ID"][0]
		c.Set("current_user_id", c.Session().Get("current_user_id"))
		b := models.Book{}

		tx := c.Value("tx").(*pop.Connection)

		// find a book
		if err := tx.Where("id = ?", id).First(&b); err != nil {
			return errors.New("Book not found." + err.Error())
		}

		c.Set("book", b)
		return c.Render(200, r.HTML("book/edit.html"))
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}

// BookEditcommit default implementation.
func BookEditcommit(c buffalo.Context) error {
	if SessionValid(c) == nil {
		if err := c.Request().ParseForm(); err != nil {
			return err
		}
		id := c.Request().Form["ID"][0]

		b := models.Book{}

		tx := c.Value("tx").(*pop.Connection)

		// find a book
		if err := tx.Where("id = ?", id).First(&b); err != nil {
			return errors.New("Book not found." + err.Error())
		}

		if err := c.Bind(&b); err != nil {
			return err
		}

		if err := tx.Update(&b); err != nil {
			return err
		}
		logger.Infof("Book saved:%#v", b)
		c.Flash().Add("success", "Book successfully saved.")
		c.Set("book", b)
		return c.Redirect(302, "/book/list")
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}

// BookAddcommit default implementation.
func BookAddcommit(c buffalo.Context) error {
	if SessionValid(c) == nil {
		b := models.Book{}

		tx := c.Value("tx").(*pop.Connection)

		if err := c.Bind(&b); err != nil {
			return err
		}

		b.UserId = c.Session().Get("current_user_id").(uuid.UUID)

		if err := tx.Save(&b); err != nil {
			return err
		}

		c.Flash().Add("success", "Book successfully saved.")
		return c.Redirect(302, "/book/list")
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}

// BookRemovecommit default implementation.
func BookRemovecommit(c buffalo.Context) error {
	if SessionValid(c) == nil {
		if err := c.Request().ParseForm(); err != nil {
			return err
		}
		//logger.Infof("Request ID: %#v", c.Request().Form["ID"][0])
		id := c.Request().Form["ID"][0]

		b := models.Book{}
		tx := c.Value("tx").(*pop.Connection)

		if err := tx.Where("id = ?", id).First(&b); err != nil {
			return errors.New("Book not found." + err.Error())
		}

		if err := tx.Destroy(&b); err != nil {
			return errors.New("Book could not be deleted." + err.Error())
		}

		c.Flash().Add("success", "Book successfully removed!")
		return c.Redirect(302, "/book/list")
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}

// BookList default implementation.
func BookList(c buffalo.Context) error {
	if SessionValid(c) == nil {
		c.Set("current_user_id", c.Session().Get("current_user_id"))
		b := []models.Book{}

		tx := c.Value("tx").(*pop.Connection)

		// find a book
		if err := tx.Where("user_id = ?", c.Session().Get("current_user_id").(uuid.UUID)).All(&b); err != nil {
			return errors.New("Books not found." + err.Error())
		}

		c.Set("books", b)
		return c.Render(200, r.HTML("book/list.html"))
	} else {
		return errors.New("No valid session exists! Please log in first!")
	}
}
