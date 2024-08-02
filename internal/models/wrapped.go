package models

import "github.com/spf13/pflag"

type Flags struct {
	Root []string
	Leaf []string
}

func (f *Flags) ExtractFlags(root *pflag.FlagSet, leaf *pflag.FlagSet) *pflag.FlagSet {
	var toReturn pflag.FlagSet
	for _, flag := range f.Root {
		if root.Lookup(flag) != nil {
			toReturn.AddFlag(root.Lookup(flag))
		}
	}
	for _, flag := range f.Leaf {
		if leaf.Lookup(flag) != nil {
			toReturn.AddFlag(leaf.Lookup(flag))
		}
	}
	return &toReturn
}

type ComposeFlags map[string]*Flags

func CreateComposeFlags(root []string, leaf []string) *Flags {
	return &Flags{
		Root: root,
		Leaf: leaf,
	}
}

func (c *ComposeFlags) Contains(command string) bool {
	_, ok := (*c)[command]
	return ok
}

func (c *ComposeFlags) Get(command string) []*Flags {
	return []*Flags{(*c)[command]}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
