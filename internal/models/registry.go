package models

import "fmt"

type RegistryInfo struct {
	Registry string `json:"registry"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegistryInfo) ValidateRegistry() error {
	if r.Registry == "" {
		return fmt.Errorf("missing registry")
	}
	if r.Username == "" {
		return fmt.Errorf("missing username")
	}
	if r.Password == "" {
		return fmt.Errorf("missing password")
	}
	return nil
}
