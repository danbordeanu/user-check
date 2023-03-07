package model

import (
	"fmt"
)

type UserCheck struct {
	Request
	Isid string `json:"isid" example:"bordeanu"`
}

func (r *UserCheck) Validate() error {
	if r.Isid == "" {
		return fmt.Errorf("user name is a required parameter")
	}
	return nil
}
