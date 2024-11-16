package controller

import (
	"errors"

	"github.com/gocroot/model"
)

func GetMemberByAttributeInProject(project model.Project, attribute string, value string) (*model.MenuItem, error) {
	for _, member := range project.Menu {
		switch attribute {
		case "email":
			if member.Name == value {
				return &member, nil
			}
		case "githubusername":
			if member.ID == value {
				return &member, nil
			}
		default:
			return nil, errors.New("unknown attribute")
		}
	}
	return nil, errors.New("member not found")
}
