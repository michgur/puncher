package design

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/a-h/templ"
)

func formatString(s string, args interface{}) string {
	println("------------------------")
	println("------------------------")
	println("------------------------")
	println("------------------------")
	println("------------------------")
	// print type of args
	fmt.Println(reflect.TypeOf(args))
	println(args)
	j, err := json.Marshal(&args)
	if err != nil {
		panic(err)
	}
	// convert json to map
	m := make(map[string]interface{})
	err = json.Unmarshal(j, &m)
	if err != nil {
		panic(err)
	}
	// replace each "{ key }" with "value"
	for k, v := range m {
		fmt.Println(k, v)
		s = strings.ReplaceAll(s, fmt.Sprintf("{ %s }", k), fmt.Sprintf("%v", v))
	}
	return s
}

// generate class lists with dynamic args
func Class(s string, args interface{}) templ.Attributes {
	return templ.Attributes{
		"class": formatString(s, args),
	}
}

func Style(s string, args interface{}) templ.Attributes {
	return templ.Attributes{
		"style": formatString(s, args),
	}
}
