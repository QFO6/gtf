package gtf

import (
	"encoding/json"
	"fmt"
	"html/template"
	htmlTemplate "html/template"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	textTemplate "text/template"
	"time"

	humanize "github.com/dustin/go-humanize"
	blackfriday "github.com/russross/blackfriday/v2"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/xeonx/timeago"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var striptagsRegexp = regexp.MustCompile("<[^>]*?>")

// recovery will silently swallow all unexpected panics.
func recovery() {
	recover()
}

var GtfTextFuncMap = textTemplate.FuncMap{
	"toValue": func(value interface{}) interface{} {
		defer recover()
		// convert primitive.ObjectID to string
		switch value.(type) {
		case primitive.ObjectID:
			return value.(primitive.ObjectID).Hex()
		default:
			return value
		}
	},
	"timeIn": func(t time.Time, locName string) string {
		defer recovery()
		if t.IsZero() {
			return ""
		}

		if locName == "" {
			locName = "Asia/Shanghai"
		}

		loc, err := time.LoadLocation(locName)
		if err == nil && loc != nil {
			t = t.In(loc)
		}
		return t.Format("2006-01-02 15:04 -07")
	},
	"funcMap": func(v ...interface{}) []interface{} {
		defer recovery()
		return v
	},

	"asQuery": func(query string) string {
		defer recovery()
		return url.QueryEscape(query)
	},

	"asURL": func(query string) template.URL {
		defer recovery()
		return template.URL(query)
	},
	"isChecked": func(values interface{}, option interface{}) string {
		defer recovery()
		list := []string{}
		switch values.(type) {
		case string:
			return ""
		case []string:
			list = values.([]string)
		}
		//for checkbox mark check value
		for _, item := range list {
			if item == option.(string) {
				return "checked"
			}
		}
		return ""
	},
	"objectId": func(value interface{}) string {
		defer recovery()
		switch value.(type) {
		case primitive.ObjectID:
			return value.(primitive.ObjectID).Hex()
		default:
			return value.(string)
		}
	},
	"parseUrl": func(path string, r *http.Request) string {
		defer recovery()
		link, _ := url.ParseRequestURI(r.RequestURI)
		link.Path = path
		return link.String()
	},
	"duration": func(start, stop time.Time) float64 {
		defer recovery()
		loading := stop.Sub(start).Seconds()
		return loading
	},
	"existobjectid": func(values []primitive.ObjectID, id string) bool {
		defer recovery()
		for _, value := range values {
			if value.Hex() == id {
				return true
			}
		}
		return false
	},
	"sameobjectid": func(value interface{}, id string) bool {
		defer recovery()
		switch value.(type) {
		case primitive.ObjectID:
			return value.(primitive.ObjectID).Hex() == id
		default:
			return fmt.Sprintf("%v", value) == id
		}
	},
	"minus": func(value interface{}, i int) int {
		defer recovery()
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return value.(int) - i
		default:
			return 0
		}
	},

	"delQuery": func(r *http.Request, k string) string {
		defer recovery()
		link, _ := url.ParseRequestURI(r.RequestURI)
		values := link.Query()

		values.Del(k)
		link.RawQuery = values.Encode()
		return link.String()
	},
	"setQuery": func(r *http.Request, k, v string) string {
		defer recovery()
		link, _ := url.ParseRequestURI(r.RequestURI)
		values := link.Query()

		values.Set(k, v)
		link.RawQuery = values.Encode()
		return link.String()
	},
	"getQuery": func(r *http.Request, k string) string {
		defer recovery()
		link, _ := url.ParseRequestURI(r.RequestURI)
		values := link.Query()
		return values.Get(k)
	},
	"repeat": func(count int, str string) string {
		defer recovery()
		return strings.Repeat(str, count)
	},
	"getInt": func(value interface{}) int {
		defer recovery()
		if value == nil {
			return 0
		}
		switch value.(type) {
		case *int:
			return *value.(*int)
		default:
			return value.(int)
		}
	},
	"istrue": func(value interface{}) bool {
		defer recovery()
		//for *bool type
		v := value.(*bool)
		if v == nil || *v == false {
			return false
		}
		return true
	},
	"humanizeSize": func(size interface{}) string {
		defer recovery()
		switch v := size.(type) {
		case float64:
			out := uint64(int64(v))
			return humanize.Bytes(out)
		case int64:
			return humanize.Bytes(uint64(v))
		}
		return "NA"
	},
	"isblank": func(value string) bool {
		defer recovery()
		s := strings.TrimSpace(value)
		if s == "" {
			return true
		}
		return false
	},
	"renderTime": func(value interface{}) string {
		defer recovery()
		switch value.(type) {
		case time.Time:
			startTime := value.(time.Time)
			loading := time.Now().Sub(startTime).Seconds() * 1000
			return fmt.Sprintf("%.2fms", loading)
		default:
			return ""
		}
	},
	"markdown2": func(value string) template.HTML {
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs
		parser := parser.NewWithExtensions(extensions)

		md := []byte(value)
		output := markdown.ToHTML(md, parser, nil)
		return template.HTML(output)
	},
	"markdown": func(value string) template.HTML {
		defer recovery()

		md := []byte(value)
		output := blackfriday.Run(md)

		return template.HTML(string(output))
	},
	"timeago": func(value time.Time) string {
		defer recovery()
		return timeago.English.Format(value)
	},
	"asHTMLAttr": func(value string) template.HTMLAttr {
		defer recovery()
		return template.HTMLAttr(value)
	},
	"asCSS": func(value string) template.CSS {
		defer recovery()
		return template.CSS(value)
	},
	"asHTML": func(value string) template.HTML {
		defer recovery()
		return template.HTML(value)
	},
	"asJS": func(value string) template.JS {
		defer recovery()
		return template.JS(value)
	},
	"existin": func(list interface{}, value string) bool {
		defer recovery()
		if list == nil {
			return false
		}
		var nlist []string
		switch list.(type) {
		case *[]string:
			nlist = *list.(*[]string)
		default:
			nlist = list.([]string)
		}

		for _, item := range nlist {
			if item == value {
				return true
			}
		}
		return false
	},
	"tojson": func(value interface{}) template.JS {
		defer recovery()
		out, _ := json.Marshal(value)
		return template.JS(string(out))
	},
	"gettitle": func(value string) string {
		defer recovery()
		list := strings.Split(value, ".")
		return list[len(list)-1]
	},
	"replace": func(s1 string, s2 string) string {
		defer recovery()

		return strings.Replace(s2, s1, "", -1)
	},
	"findreplace": func(s1 string, s2 string, s3 string) string {
		defer recovery()

		return strings.Replace(s3, s1, s2, -1)
	},
	"title": func(s string) string {
		defer recovery()
		return strings.Title(s)
	},
	"default": func(arg interface{}, value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			if v.Len() == 0 {
				return arg
			}
		case reflect.Bool:
			if !v.Bool() {
				return arg
			}
		default:
			return value
		}

		return value
	},
	"length": func(value interface{}) int {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return v.Len()
		case reflect.String:
			return len([]rune(v.String()))
		case reflect.Ptr:
			return v.Elem().Len()
		}

		return 0
	},
	"lower": func(s string) string {
		defer recovery()

		return strings.ToLower(s)
	},
	"upper": func(s string) string {
		defer recovery()

		return strings.ToUpper(s)
	},
	"truncatechars": func(n int, s string) string {
		defer recovery()

		if n < 0 {
			return s
		}

		r := []rune(s)
		rLength := len(r)

		if n >= rLength {
			return s
		}

		if n > 3 && rLength > 3 {
			return string(r[:n-3]) + "..."
		}

		return string(r[:n])
	},
	"urlencode": func(s string) string {
		defer recovery()

		return url.QueryEscape(s)
	},
	"wordcount": func(s string) int {
		defer recovery()

		return len(strings.Fields(s))
	},
	"divisibleby": func(arg interface{}, value interface{}) bool {
		defer recovery()

		var v float64
		switch value.(type) {
		case int, int8, int16, int32, int64:
			v = float64(reflect.ValueOf(value).Int())
		case uint, uint8, uint16, uint32, uint64:
			v = float64(reflect.ValueOf(value).Uint())
		case float32, float64:
			v = reflect.ValueOf(value).Float()
		default:
			return false
		}

		var a float64
		switch arg.(type) {
		case int, int8, int16, int32, int64:
			a = float64(reflect.ValueOf(arg).Int())
		case uint, uint8, uint16, uint32, uint64:
			a = float64(reflect.ValueOf(arg).Uint())
		case float32, float64:
			a = reflect.ValueOf(arg).Float()
		default:
			return false
		}

		return math.Mod(v, a) == 0
	},
	"lengthis": func(arg int, value interface{}) bool {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return v.Len() == arg
		case reflect.String:
			return len([]rune(v.String())) == arg
		}

		return false
	},
	"trim": func(s string) string {
		defer recovery()

		return strings.TrimSpace(s)
	},
	"capfirst": func(s string) string {
		defer recovery()

		return strings.ToUpper(string(s[0])) + s[1:]
	},
	"pluralize": func(arg string, value interface{}) string {
		defer recovery()

		flag := false
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			flag = v.Int() == 1
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			flag = v.Uint() == 1
		default:
			return ""
		}

		if !strings.Contains(arg, ",") {
			arg = "," + arg
		}

		bits := strings.Split(arg, ",")

		if len(bits) > 2 {
			return ""
		}

		if flag {
			return bits[0]
		}

		return bits[1]
	},
	"yesno": func(yes string, no string, value bool) string {
		defer recovery()

		if value {
			return yes
		}

		return no
	},
	"rjust": func(arg int, value string) string {
		defer recovery()

		n := arg - len([]rune(value))

		if n > 0 {
			value = strings.Repeat(" ", n) + value
		}

		return value
	},
	"ljust": func(arg int, value string) string {
		defer recovery()

		n := arg - len([]rune(value))

		if n > 0 {
			value = value + strings.Repeat(" ", n)
		}

		return value
	},
	"center": func(arg int, value string) string {
		defer recovery()

		n := arg - len([]rune(value))

		if n > 0 {
			left := n / 2
			right := n - left
			value = strings.Repeat(" ", left) + value + strings.Repeat(" ", right)
		}

		return value
	},
	"filesizeformat": func(value interface{}) string {
		defer recovery()

		var size float64

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			size = float64(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			size = float64(v.Uint())
		case reflect.Float32, reflect.Float64:
			size = v.Float()
		default:
			return ""
		}

		var KB float64 = 1 << 10
		var MB float64 = 1 << 20
		var GB float64 = 1 << 30
		var TB float64 = 1 << 40
		var PB float64 = 1 << 50

		filesizeFormat := func(filesize float64, suffix string) string {
			return strings.Replace(fmt.Sprintf("%.1f %s", filesize, suffix), ".0", "", -1)
		}

		var result string
		if size < KB {
			result = filesizeFormat(size, "bytes")
		} else if size < MB {
			result = filesizeFormat(size/KB, "KB")
		} else if size < GB {
			result = filesizeFormat(size/MB, "MB")
		} else if size < TB {
			result = filesizeFormat(size/GB, "GB")
		} else if size < PB {
			result = filesizeFormat(size/TB, "TB")
		} else {
			result = filesizeFormat(size/PB, "PB")
		}

		return result
	},
	"apnumber": func(value interface{}) interface{} {
		defer recovery()

		name := [10]string{"one", "two", "three", "four", "five",
			"six", "seven", "eight", "nine"}

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() < 10 {
				return name[v.Int()-1]
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v.Uint() < 10 {
				return name[v.Uint()-1]
			}
		}

		return value
	},
	"intcomma": func(value interface{}) string {
		defer recovery()

		v := reflect.ValueOf(value)

		var x uint
		minus := false
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() < 0 {
				minus = true
				x = uint(-v.Int())
			} else {
				x = uint(v.Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x = uint(v.Uint())
		default:
			return ""
		}

		var result string
		for x >= 1000 {
			result = fmt.Sprintf(",%03d%s", x%1000, result)
			x /= 1000
		}
		result = fmt.Sprintf("%d%s", x, result)

		if minus {
			result = "-" + result
		}

		return result
	},
	"ordinal": func(value interface{}) string {
		defer recovery()

		v := reflect.ValueOf(value)

		var x uint
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() < 0 {
				return ""
			}
			x = uint(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x = uint(v.Uint())
		default:
			return ""
		}

		suffixes := [10]string{"th", "st", "nd", "rd", "th", "th", "th", "th", "th", "th"}

		switch x % 100 {
		case 11, 12, 13:
			return fmt.Sprintf("%d%s", x, suffixes[0])
		}

		return fmt.Sprintf("%d%s", x, suffixes[x%10])
	},
	"first": func(value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			return string([]rune(v.String())[0])
		case reflect.Slice, reflect.Array:
			return v.Index(0).Interface()
		}

		return ""
	},
	"last": func(value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			str := []rune(v.String())
			return string(str[len(str)-1])
		case reflect.Slice, reflect.Array:
			return v.Index(v.Len() - 1).Interface()
		}

		return ""
	},
	"join": func(arg string, value []string) string {
		defer recovery()

		return strings.Join(value, arg)
	},
	"slice": func(start int, end int, value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)

		if start < 0 {
			start = 0
		}

		switch v.Kind() {
		case reflect.String:
			str := []rune(v.String())

			if end > len(str) {
				end = len(str)
			}

			return string(str[start:end])
		case reflect.Slice:
			return v.Slice(start, end).Interface()
		}
		return ""
	},
	"random": func(value interface{}) interface{} {
		defer recovery()

		rand.Seed(time.Now().UTC().UnixNano())

		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			str := []rune(v.String())
			return string(str[rand.Intn(len(str))])
		case reflect.Slice, reflect.Array:
			return v.Index(rand.Intn(v.Len())).Interface()
		}

		return ""
	},
	"randomintrange": func(min, max int, value interface{}) int {
		defer recovery()

		rand.Seed(time.Now().UTC().UnixNano())
		return rand.Intn(max-min) + min
	},
	"striptags": func(s string) string {
		return strings.TrimSpace(striptagsRegexp.ReplaceAllString(s, ""))
	},
}

var GtfFuncMap = htmlTemplate.FuncMap(GtfTextFuncMap)

// gtf.New is a wrapper function of template.New(https://golang.org/pkg/html/template/#New).
// It automatically adds the gtf functions to the template's function map
// and returns template.Template(http://golang.org/pkg/html/template/#Template).
func New(name string) *htmlTemplate.Template {
	return htmlTemplate.New(name).Funcs(GtfFuncMap)
}

// gtf.Inject injects gtf functions into the passed FuncMap.
// It does not overwrite the original function which have same name as a gtf function.
func Inject(funcs map[string]interface{}) {
	for k, v := range GtfFuncMap {
		if _, ok := funcs[k]; !ok {
			funcs[k] = v
		}
	}
}

// gtf.ForceInject injects gtf functions into the passed FuncMap.
// It overwrites the original function which have same name as a gtf function.
func ForceInject(funcs map[string]interface{}) {
	for k, v := range GtfFuncMap {
		funcs[k] = v
	}
}

// gtf.Inject injects gtf functions into the passed FuncMap.
// It prefixes the gtf functions with the specified prefix.
// If there are many function which have same names as the gtf functions,
// you can use this function to prefix the gtf functions.
func InjectWithPrefix(funcs map[string]interface{}, prefix string) {
	for k, v := range GtfFuncMap {
		funcs[prefix+k] = v
	}
}
