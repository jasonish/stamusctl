package models

import (
	"log"
	"os/user"
	"strings"

	"stamus-ctl/internal/app"
)

type Release struct {
	Name      string // the name given to the release
	User      string // userid
	Group     string // groupid
	Location  string // path to the install (eg: /home/.../config, the place where the compose.yaml file will be placed)
	IsUpgrade bool   // see helm
	IsInstall bool   // see helm
	Service   string
}

func NewRelease(name, location string, isUpgrade, isInstall bool) *Release {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return &Release{
		User:      currentUser.Uid,
		Group:     currentUser.Gid,
		Name:      name,
		Location:  location,
		IsUpgrade: isUpgrade,
		IsInstall: isInstall,
		Service:   app.StamusAppName + ":" + app.Version,
	}
}

func (s *Release) AsMap() map[string]interface{} {
	prefix := "Release"
	return map[string]interface{}{
		prefix + ".name":      s.Name,
		prefix + ".user":      s.User,
		prefix + ".group":     s.Group,
		prefix + ".location":  s.Location,
		prefix + ".isUpgrade": s.IsUpgrade,
		prefix + ".isInstall": s.IsInstall,
		prefix + ".service":   s.Service,
	}
}

func (s *Release) SetName(name string) *Release {
	s.Name = name
	return s
}
func (s *Release) SetLocation(location string) *Release {
	s.Location = location
	return s
}

func (s *Release) SetIsUpgrade(isUpgrade bool) *Release {
	s.IsUpgrade = isUpgrade
	return s
}

func (s *Release) SetIsInstall(isInstall bool) *Release {
	s.IsInstall = isInstall
	return s
}

func (s *Release) SetService(service string) *Release {
	s.Service = service
	return s
}

type Template struct {
	templateName    string
	templateVersion string
}

func NewTemplate(name string) *Template {
	splitted := strings.Split(name, "/")
	version := strings.Join(splitted[len(splitted)-2:], "/")
	return &Template{
		templateName:    name,
		templateVersion: version,
	}
}

func (t *Template) AsMap() map[string]interface{} {
	prefix := "Template"
	return map[string]interface{}{
		prefix + ".name":    t.templateName,
		prefix + ".version": t.templateVersion,
	}
}
