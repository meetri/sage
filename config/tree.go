package config

import (
	"errors"
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	//"reflect"
	//"regexp"
	//"strings"
)

type Tree struct {
	Data Map
	Env  Map
	Sel  Map
}

func (t *Tree) SmartLoad(fn string) (err error) {

	t.Data = Load(fn)
	t.Env = t.Data.Find("_env").(Map)
	delete(t.Data, "_env")

	return
}

func (t *Tree) Select(fn string) (err error) {

	var env Map
	if t.Sel, env = t.Data.Select(fn, t.Env); t.Sel == nil {
		err = errors.New("what's the deal yo")
	}
	_ = env

	return

}
