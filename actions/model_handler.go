package actions

import (
	. "buf1/models"

	"github.com/gobuffalo/buffalo"
	// 	"fmt"
	// 	"os"
)

// HomeHandler is a default handler to serve up
// a home page.
func ModelListHandler(c buffalo.Context) error {
	var u []User
	//tx := c.Value("tx").(*pop.Connection)
	err := DB.Eager().All(&u)
	//fmt.Fprintf(os.Stderr, "user: %v", u)
	if err == nil {
		c.Set("users", u)
		return c.Render(200, r.HTML("userlist.html"))
	} else {
		c.Set("error", err.Error())
		c.Set("DBURL", DB.URL())
		return c.Render(403, r.HTML("error.html"))
	}
}
