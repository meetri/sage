package config

/*
	Parse YAML Config - ( supports Docker Compose Format )
*/

import (
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strings"
)

type Map map[interface{}]interface{}

//Load RAW YAML data
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
		panic("problem parsing yaml data")
	}
	return

}

//Export YAML
func (mm Map) Save() []byte {
	out, _ := yaml.Marshal(mm)
	spew.Printf("%s", out)
	return out
}

//Expand config's extended files and services
func (mm Map) Select(path string) Map {

	elem := mm.Find(path).(Map)

	extend := elem.Find("extends/file")
	service := elem.Find("extends/service")

	for extend != nil {
		es := Load(extend.(string)).Find(service.(string))
		if es != nil {
			eservice := es.(Map)
			for k, v := range eservice {
				mergeKey(false, eservice, elem, k, v)
			}

			extend = eservice.Find("extends/file")
			service = eservice.Find("extends/service")
		} else {
			panic("failed loading " + extend.(string) + "::" + service.(string))
		}
	}

	return mm.Find(path).(Map)
}

// grab element specified by searchPath ie. '/elem1/subelem2/elem3'
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
func mergeKey(override bool, srcMap interface{}, dstMap interface{}, key interface{}, val interface{}) {

	if _, ok := dstMap.(Map)[key]; !ok {
		dstMap.(Map)[key] = val
	} else {
		valType := reflect.TypeOf(val).Kind()

		switch valType {
		case reflect.Map:
			nm := srcMap.(Map)[key]
			if "map" == reflect.TypeOf(nm).Kind().String() {
				for k, v := range nm.(Map) {
					mergeKey(override, srcMap.(Map)[key], dstMap.(Map)[key], k, v)
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
