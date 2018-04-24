package funclib

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

func SessionValid(c buffalo.Context) error {
	if c.Session().Get("current_user_id") != nil {
		return nil
	} else {
		return errors.New("No Session ID found!")
	}
}
