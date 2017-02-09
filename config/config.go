/*
	Load nested YAML app configuration
*/
package config

import (
	//"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

type Map map[interface{}]interface{}

//type Map map[interface{}]interface{}

//Export YAML
func (mm Map) Save() []byte {
	out, _ := yaml.Marshal(mm)
	spew.Printf("%s", out)
	return out
}

//LoadRAW YAML data
func LoadRaw(data []byte) (mm Map, err error) {
	err = yaml.Unmarshal([]byte(data), &mm)
	return
}

//Load YAML from file
func Load(fn string) (mm Map) {

	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		panic("problem reading from " + fn)
	}
	mm, err = LoadRaw(raw)
	if err != nil {
		panic("YAML Parse Error:" + err.Error())
	}
	return

}

//Templatize a NestedMaps object which replaces "${key}" with the value from env[key]
func (mm Map) Templatize(env Map) {
	reg := regexp.MustCompile("\\$\\{([^\\}]+)\\}")
	for k, v := range mm {
		if reflect.TypeOf(v).Kind() == reflect.String {
			matches := reg.FindAllString(v.(string), -1)
			for _, match := range matches {

				envKey := strings.Replace(
					strings.Replace(match, "}", "", 1), "${", "", 1)

				if env[envKey] != nil {
					newVal := env[envKey].(string)
					mm[k] = strings.Replace(mm[k].(string), match, newVal, -1)
				}
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Map {
			v.(Map).Templatize(env)
		}
	}
}

func (mm Map) mergeExtendedPath(file, service, pathkey string, parentEnv Map) Map {

	var env Map

	elem := mm.Find(file).(Map)

	ex := Load(file)
	es := ex.Find(service)

	if g := ex.Find("_env"); g == nil {
		env = parentEnv
	} else {
		env = g.(Map)
		for k, v := range parentEnv {
			env[k] = v
		}
	}

	var elempath Map
	if pathkey == "" {
		elempath = elem
	} else {
		elem[pathkey] = make(Map)
		elempath = elem[pathkey].(Map)
	}

	if es != nil {
		eservice := es.(Map)
		eservice.Templatize(env)
		for k, v := range eservice {
			mergeMaps(false, eservice, elempath, k, v)
		}
		mm.Templatize(env)

		return eservice

	}

	return nil
}

//Select and Expand config's extended files
func (mm Map) Select(path string, parentEnv Map) (selmap Map, env Map) {

	elem := mm.Find(path).(Map)

	extend := elem.Find("extends/file")
	service := elem.Find("extends/service")
	pathkey := elem.Find("extends/path")

	if extend != nil {

		elem.mergeExtendedPath(extend.(string), service.(string), pathkey.(string), parentEnv)

	}

	env = make(Map)
	for extend != nil {

		ex := Load(extend.(string))
		es := ex.Find(service.(string))
		if g := ex.Find("_env"); g == nil {
			env = parentEnv
		} else {
			env = g.(Map)
			for k, v := range parentEnv {
				env[k] = v
			}
		}

		var elempath Map
		if pathkey == nil {
			elempath = elem
		} else {
			elem[pathkey] = make(Map)
			elempath = elem[pathkey].(Map)
		}

		if es != nil {
			eservice := es.(Map)
			eservice.Templatize(env)
			for k, v := range eservice {
				mergeMaps(false, eservice, elempath, k, v)
			}
			mm.Templatize(env)

			extend = eservice.Find("extends/file")
			service = eservice.Find("extends/service")
			pathkey = eservice.Find("extends/path")

		} else {
			panic("failed loading " + extend.(string) + "::" + service.(string))
		}
	}

	selmap = mm.Find(path).(Map)
	return
}

// grab element specified by searchPath in nested Map ie. '/elem1/subelem2/elem3'
func (mm Map) Find(searchPath string) interface{} {

	pathArr := strings.Split(searchPath, "/")
	node := mm[pathArr[0]]
	if node != nil {
		for _, path := range pathArr[1:] {
			if reflect.Map == reflect.TypeOf(node).Kind() {
				node, _ = node.(Map)[path]
			} else {
				panic("can't find " + searchPath)
			}
		}
	}
	return node
}

// merge map trees, choose to override or not from src
func mergeMaps(override bool, srcMap interface{}, dstMap interface{}, key interface{}, val interface{}) {

	if _, ok := dstMap.(Map)[key]; !ok {
		dstMap.(Map)[key] = val
	} else {
		valType := reflect.TypeOf(val).Kind()

		switch valType {
		case reflect.Map:
			nm := srcMap.(Map)[key]
			if "map" == reflect.TypeOf(nm).Kind().String() {
				for k, v := range nm.(Map) {
					mergeMaps(override, srcMap.(Map)[key], dstMap.(Map)[key], k, v)
				}
			}
			break
		case reflect.Slice:
			//merge slices
			if reflect.TypeOf(dstMap.(Map)[key]).Kind() == reflect.Slice {
				pval := dstMap.(Map)[key].([]interface{})
				for _, v := range val.([]interface{}) {
					skip := false
					for _, iv := range pval {
						if iv == v {
							skip = true
						}
					}
					if !skip {
						pval = append(pval, v)
					}
				}
				dstMap.(Map)[key] = pval
			} else {
				dstMap.(Map)[key] = val
			}

			break
		case reflect.String:
			if override && dstMap.(Map)[key] != val {
				dstMap.(Map)[key] = val
			}
		}

	}
}
