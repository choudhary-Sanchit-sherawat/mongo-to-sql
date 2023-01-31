package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//first we has to take input of mogo json query and Unmarshal([]byte(s), &m) into map

// create a func for the conver the mogo into sql

// in func we we take a map interface input
// and after we take into map condiion key and val
// and in the loop we take switch case ==  key
// first hendle the and condition we take into array "key =  "$and" " then we take in loop
func main() {

	// m := make(map[string]string)

	// m["k1"] = "2"
	// m["k2"] = "3"

	// k := convetintoint(m)
	// fmt.Println(k)

	// s := `{ "$and": [ { "age" : { "$gte" : 20 } },  { "$or" : [  { "postalCode" : 1000 },  { "$and" : [  { "lastname" : "johnson" }, { "firstname" : "ohn"} ] }  ] }  ] } `

	s := `{ "delete" :{"id": 1000} }`

	// s := `{
	// 	"$and": [
	// 	       { "x": { "$gt": 0 } },
	// 		   { "$or": { "$lt": 0 } },
	// 	 	  {"x": { "$eq":{"$or":["sancit","2","3"]}}},
	// 		   {"x":{"$in":["sancit","2","3"]}},
	// 	 	   {"$or": [
	// 		   { "x": { "$ne": 0 } },
	// 		   { "x": { "$ne": 0 } }
	// 	 	   ]}
	// 		]

	// 	  }  `

	// s := `{
	// 	"$and": [
	// 	   { "x": { "$gt": 0 } },
	// 	   { "x": { "$lt": 0 } },
	// 	  {"x": { "$eq": 3 }},
	// 	   {"x":{"$in":["1","2","3"]}},
	// 	   {"$or": [
	// 	   { "x": { "$ne": 0 } },
	// 	   { "x": { "$ne": 0 } }
	// 	   ]}
	// 	]

	//  }  `
	// s := `{

	// 	"$and": [{
	// 	" (x ":{ "$gt" :  0  }
	// 	 },{
	// 	"$and": [{
	// 	"x ":{ "$lt" :  0  }
	// 	 },{
	// 	"$and": [{
	// 	"x " : 3
	// 	 },{
	// 	"$and": [{
	// 	"x": {"$in":  [1,2,3] }
	// 	 },{
	// 	"$or": [{
	// 	 "x " : { "$ne":  0 }
	//   },{ "x " : { "$ne":  0)) }
	// 	}]
	// 	}]
	// 	}]
	// 	}]
	// 	}]
	//  }`

	// s := `{ "lastname" : "johnson" }`
	// s := `{

	// 	"$and": [{
	// 	"age":{ "$gte" :  20  }
	// 	 },{
	// 	"$or": [{
	// 	 "postalCode " :  1000
	//   },{
	// 	"$and": [{
	// 	"lastname " :  "johnson"
	//   },{ " firstname " :  "ohn"
	// 	}]
	// 	}]
	// 	}]
	//  }`

	// fmt.Println(s)

	// s := `{ "$and": [ { "age" : { "$gte" : 20 } },  { "$or" : [  { "postalCode" : 1000 },  { "$and" : [  { "lastname" : "johnson" }, { "firstname" : "ohn"} ] }  ] }  ] } `

	// s := `{
	// 	"$and": [
	// 		{ "$or": [ { "qty": { "$lt" : 10 } }, { "qty" : { "$gt": 50 } } ] },
	// 		{ "$or": [ { "sale": true }, { "price" : { "$lt" : 5 } } ] }
	// 	]
	// }`

	r, err := Unmarshal(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)

}

type FilterOp string

const (
	eq  FilterOp = "$eq"
	gt  FilterOp = "$gt"
	gte FilterOp = "$gte"
	in  FilterOp = "$in"
	lt  FilterOp = "$lt"
	lte FilterOp = "$lte"
	ne  FilterOp = "$ne"

	and FilterOp = "$and"
	or  FilterOp = "$or"
	nin FilterOp = "$nin"
)

func Unmarshal(s string) (string, error) {

	var m map[string]any

	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return "", err
	}
	K, err := MongoToSql(m)
	if err != nil {
		return K, err
	}
	return K, nil

}

func MongoToSql(mongoQuery map[string]any) (string, error) {
	// fmt.Println(mongoQuery)
	// sqlQuery for the stored ek jgah rakh ne ke leye
	var err error
	sqlQuery := []string{}
	for k, v := range mongoQuery {
		switch FilterOp(k) {
		case or:
			x := "OR"
			sqlQuery, err = MapInterface(v, sqlQuery, x)
			if err != nil {
				return k, err
			}
		case and:
			x := "AND"
			sqlQuery, err = MapInterface(v, sqlQuery, x)
			if err != nil {
				return k, err
			}

		case in:
			x := "IN"
			sqlQuery, err = Interface(v, sqlQuery, x)
			if err != nil {
				return k, err
			}
		case nin:
			x := "NOT IN"
			sqlQuery, err = Interface(v, sqlQuery, x)
			if err != nil {
				return k, err
			}

		case ne:
			x := "!="
			sqlQuery, err = Convert(sqlQuery, x, v)
			if err != nil {
				return k, err
			}
		case eq:
			x := "="
			sqlQuery, err = Convert(sqlQuery, x, v)
			if err != nil {
				return k, err
			}
		case gt:
			x := ">"
			sqlQuery, err = Convert(sqlQuery, x, v)
			if err != nil {
				return k, err
			}
		case lt:
			x := "<"
			sqlQuery, err = Convert(sqlQuery, x, v)
			if err != nil {
				return k, err
			}
		case gte:
			x := ">="
			sqlQuery, err = Convert(sqlQuery, x, v)
			if err != nil {
				return k, err
			}
		case lte:
			x := "<="
			sqlQuery, err = Convert(sqlQuery, x, v)
			if err != nil {
				return k, err
			}
		default:
			sqlQuery, err = Default(v, sqlQuery, k)
			if err != nil {
				return "", nil
			}

		}
	}
	return strings.Join(sqlQuery, " AND "), nil
}

func Convert(sqlQuery []string, x string, v any) ([]string, error) {
	_, maps := v.(map[string]any)
	_, inter := v.([]any)
	if maps || inter {
		log.Fatal("pls enter the correct value is not int or float, String")
	}
	sqlQuery = append(sqlQuery, fmt.Sprintf("%v %v", x, v))
	return sqlQuery, nil
}

func Interface(v any, sqlQuery []string, x string) ([]string, error) {
	var conditions []string

	value, ok := v.([]any)
	if !ok {
		return nil, SetErr(value, ok)
	}
	for _, condition := range value {
		conditions = append(conditions, fmt.Sprintf("%v", condition))
	}

	sqlQuery = append(sqlQuery, ""+x+" ("+strings.Join(conditions, ",")+")")
	return sqlQuery, nil
}

// mapInterface for the AND & OR
func MapInterface(v any, sqlQuery []string, x string) ([]string, error) {
	// ager or condition hai to loop chale or hr loop ke baD or lgega
	conditions := []string{}
	value, ok := v.([]any)
	if !ok {
		return nil, SetErr(value, ok)
	}

	for _, condition := range value {
		value, ok := condition.(map[string]any)
		if !ok {
			return nil, SetErr(value, ok)
		}

		m, err := MongoToSql(value)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, m)
	}
	// sqlQuery = append(sqlQuery, strings.Join(conditions, " OR ya AND "))
	sqlQuery = append(sqlQuery, "("+strings.Join(conditions, " "+x+" ")+")")
	return sqlQuery, nil
}

func Default(v any, sqlQuery []string, k string) ([]string, error) {
	var condition string
	var err error
	if _, ok := v.(map[string]any); ok {

		condition, err = MongoToSql(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		sqlQuery = append(sqlQuery, fmt.Sprintf("%v %v", k, condition))

	} else if _, ok := v.(bool); ok {

		sqlQuery = append(sqlQuery, fmt.Sprintf(`%v = %v`, k, v))

	} else {

		if value, ok := v.([]any); ok {
			return nil, SetErr(value, ok)
		}
		sqlQuery = append(sqlQuery, fmt.Sprintf(`%v = "%v"`, k, v))
	}
	return sqlQuery, nil
}

func SetErr(value any, ok bool) error {
	return fmt.Errorf("pls enter the correct value: %v, typeStatus : %v,", value, ok)
}

//
//
//
//

//
//

//

//
// /
//
//

// func mongoToSql(mongoQuery map[string]interface{}) string {
// 	// sqlQuery for the stored ek jgah rakh ne ke leye
// 	sqlQuery := []string{}
// 	for k, v := range mongoQuery {
// 		switch k {
// 		case "$or":
// 			conditions := []string{}
// 			// ager or condition hai to loop chale or hr loop ke baD or lgega
// 			for _, condition := range v.([]interface{}) {
// 				conditions = append(conditions, mongoToSql(condition.(map[string]interface{})))
// 			}
// 			// sqlQuery = append(sqlQuery, strings.Join(conditions, " OR "))
// 			sqlQuery = append(sqlQuery, "("+strings.Join(conditions, " OR ")+")")
// 		case "$and":
// 			var conditions []string
// 			for _, condition := range v.([]interface{}) {
// 				conditions = append(conditions, mongoToSql(condition.(map[string]interface{})))
// 			}
// 			sqlQuery = append(sqlQuery, "("+strings.Join(conditions, " AND ")+")")
// 		case "$in":
// 			var conditions []string
// 			for _, condition := range v.([]interface{}) {
// 				conditions = append(conditions, fmt.Sprintf("%v", condition))
// 			}
// 			sqlQuery = append(sqlQuery, "IN ("+strings.Join(conditions, ",")+")")
// 		case "$nin":
// 			var conditions []string
// 			for _, condition := range v.([]interface{}) {
// 				conditions = append(conditions, fmt.Sprintf("%v", condition))
// 			}
// 			sqlQuery = append(sqlQuery, " NOT IN ("+strings.Join(conditions, ",")+")")
// 		case "$expr":
// 			// ager Expr hai to isme ki se statement bhi hoge then ye chlega
// 			condition := Expr(v.(map[string]interface{}))
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("(%v)", condition))
// 		case "$divide":
// 			// join them with the / operator
// 			operands := v.([]interface{})
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("(%v / %v)", operands[0], operands[1]))
// 		case "$ne":
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("!= %v", v))
// 		case "$eq":
// 			// operands := v.([]interface{})
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("= %v", v))
// 		case "$gt":
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("> %v", v))
// 		case "$lt":
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("< %v", v))
// 		case "$gte":
// 			sqlQuery = append(sqlQuery, fmt.Sprintf(">= %v", v))
// 		case "$lte":
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("<= %v", v))
// 		default:
// 			var condition string
// 			if _, ok := v.(map[string]interface{}); ok {

// 				condition = mongoToSql(v.(map[string]interface{}))
// 				sqlQuery = append(sqlQuery, fmt.Sprintf("%v %v", k, condition))
// 			} else {
// 				sqlQuery = append(sqlQuery, fmt.Sprintf(`%v = "%v"`, k, v))

// 			}

// 		}
// 	}
// 	return strings.Join(sqlQuery, " AND ")
// }

// // func isType(b interface{}) bool {
// // 	var A map[string]interface{}
// // 	return fmt.Sprintf("%T", b) == fmt.Sprintf("%T", A)

// // }

// func Expr(mongoQuery map[string]interface{}) string {
// 	var sqlQuery []string
// 	// Iterate over the MongoDB query and build the SQL query
// 	for k, v := range mongoQuery {
// 		// fmt.Println(k, v)
// 		switch k {
// 		case "$ne":
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("<> %v", v))
// 		case "$divide":
// 			operands := v.([]interface{})
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("(%v / %v)", operands[0], operands[1]))
// 		case "$eq":
// 			// fmt.Println(mongoQuery, "ffff")
// 			operands := v.([]interface{})
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("%v = %v", mongoToSql(operands[0].(map[string]interface{})), operands[1]))
// 		default:
// 			condition := mongoToSql(v.(map[string]interface{}))
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("%v %v", k, condition))
// 		}
// 	}
// 	return strings.Join(sqlQuery, " AND ")
// }

// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============
// ===================  !!!!!}!}!}!}! =========================                  ============

// =================================================================================================================
// ===================================================================================================================
// ===================================================================================================================
// ===================================================================================================================
// ===================================================================================================================
// ===================================================================================================================
// ===================================================================================================================

// func Default(mongoQuery map[string]interface{}) string {
// 	var sqlQuery []string
// 	// Iterate over the MongoDB query and build the SQL query
// 	for k, v := range mongoQuery {
// 		// fmt.Println(k, v)
// 		switch v.(type) {
// 		// case map[string]interface{}:
// 		// 	fmt.Println(reflect.ValueOf(v).Kind())
// 		// 	condition := mongoToSql(v.(map[string]interface{}))
// 		// 	sqlQuery = append(sqlQuery, fmt.Sprintf("(%v %v)", k, condition))
// 		// case []interface{}:
// 		// 	operands := v.([]interface{})
// 		// 	sqlQuery = append(sqlQuery, fmt.Sprintf("(%v / %v)", operands[0], operands[1]))
// 		case float64:
// 			// fmt.Println(mongoQuery, "ffff")
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("<> %v", v))
// 		default:
// 			// For all other cases, simply add the field and value to the SQL query
// 			condition := mongoToSql(v.(map[string]interface{}))
// 			sqlQuery = append(sqlQuery, fmt.Sprintf("%v %v", k, condition))
// 		}
// 	}
// 	return strings.Join(sqlQuery, " AND ")

// }

// func parseMap(aMap map[string]interface{}) {
// 	for key, val := range aMap {
// 		// fmt.Println(val.(type))
// 		switch concreteVal := val.(type) {
// 		case map[string]interface{}:
// 			fmt.Println(key, "ppp")
// 			parseMap(val.(map[string]interface{}))
// 		case []interface{}:
// 			fmt.Println(key, "123")
// 			if strings.Contains(key, "$divide") {
// 				parseArray(val.([]interface{}), "$divide")

// 			} else {
// 				parseArray(val.([]interface{}), "")
// 			}
// 			// fmt.Println(full)

// 		default:
// 			// fmt.Println(regexp.Compile(key,""))
// 			if strings.Contains(key, "$ne") {
// 				fmt.Println(key, "!=", concreteVal)
// 			} else {
// 				fmt.Println(key, "=", concreteVal)
// 			}

// 		}
// 	}
// }

// func parseArray(anArray []interface{}, c string) {
// 	// var concreteVal interface{}
// 	// var i int
// 	for i, val := range anArray {
// 		switch concreteVal := val.(type) {
// 		case map[string]interface{}:
// 			fmt.Println("Index:", i)
// 			parseMap(val.(map[string]interface{}))
// 			// return i, concreteVal
// 		case []interface{}:
// 			fmt.Println("Index:", i)
// 			parseArray(val.([]interface{}), "")
// 			// return i, concreteVal
// 		default:
// 			if c == "$divide" {
// 				if i == 0 {
// 					full = fmt.Sprint(val, "/")
// 				}
// 				if i == 1 {
// 					full = full + fmt.Sprint(val)
// 				}

// 			}
// 			fmt.Println(full)
// 			fmt.Println("Index", i, ":", concreteVal)
// 			// return i, concreteVal

// 		}

// 	}

// }

// func Loop(k string, v interface{}, m map[string]interface{}) (interface{}, string) {
// 	for k, v := range m {
// 		fmt.Println(k, "dfghj", v)
// 		// if v!=nil{
// 		// 	for
// 		// }
// 	}
// 	return v, k
// }
// case "$in":

// 	// fmt.Println(mongoQuery, "ffff")
// 	operands := v.([]interface{})
// 	// fmt.Println(operands, "operandseq")
// 	// condition := mongoToSql(v.([]interface{}.(map[string]interface{})))
// 	// fmt.Println(condition)
// 	// var l =
// 	sqlQuery = append(sqlQuery, fmt.Sprintf("%v = %v", operands...))
// 	// mongoToSql(v.(map[string]interface{}))
// 	// fmt.Println(sqlQuery)

// var conditions string
// // fmt.Println("$or")
// // ager or condition hai to loop chale or hr loop ke baD or lgega
// for _, condition := range v.([]interface{}) {
// 	fmt.Println(condition, "aas")
// 	// if v.(type) != map[string]interface{}{

// 	// }
// 	switch c := condition.(type) {
// 	// fmt.Println(c)
// 	case float64:
// 		fmt.Println(conditions, c)
// 		// conditions = append(conditions, fmt.Sprintf("%v", c))

// 	default:
// 		// conditions = mongoToSql(v.(map[string]interface{}))
// 		fmt.Println(conditions)

// 	}
// }

// if reflect.ValueOf(v).Kind() != (float64){

// }
// fmt.Println(operands, "operandseq")
// condition := mongoToSql(v.([]interface{}.(map[string]interface{})))
// fmt.Println(condition)
// var l =

// type User struct {
// 	Massages []Message
// }

// type Message struct {
// 	Name       string
// 	Version    int64
// 	DirectionS Direction
// }

// type Direction int

// const (
// 	AND Direction = iota
// 	OR
// 	Eq
// 	Lt
// )

//	type Con struct{
//		AND Con
//		OR Con
//	}
// var full string

// fmt.Println(operands, "operandseq")
// condition := mongoToSql(v.([]interface{}.(map[string]interface{})))
// fmt.Println(condition)
// var l =

// Join the individual parts of the SQL query with the AND keyword and return the result
// fmt.Println(strings.Join(sqlQuery, " AND "))

// fmt.Println("kkk")
// condition := mongoToSql(v.(map[string]interface{}))
// fmt.Println(condition, "condition")
// condition := mongoToSql(v.(map[string]interface{})["$eq"].([]interface{})[0].(map[string]interface{}))
// fmt.Println(condition, "condition")

// for _, condition := range v.([]interface{}) {
// 	fmt.Println(condition)
// 	conditions = append(conditions, fmt.Sprintf(" %v ", condition))
// }

// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/
// ]][[][][][][][][][][][][][][]]|||||||||||||||||||||//////////|\\\\\\\\\\\\//////////\\\\\\\\/\/\//\/\///\/\/\/\//\\/\//\/\/\/\\//\/\/\\/\/\/\/\/\/

// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// DDDDDDDD[][][][][]]][][]{}{}{}{}{}()()()()()hello word i  am hear where  are you :::::::::: new think alway looking for me......[][][][][][
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
// 0001234468/62166740305.0m'hvv  ge50gqr\frf+*f+r9V9+6+1F++++++*-*-*-*-*--*-*-*-*-*-*-**-*-**-*-*-*-*-*-*-***-*-*-*-**-**-*-*-*-*-*-*-****-*-*-
// *-*-*-*-*---*--*-*-*-**-*-*-*---*-*-*-*-*-*--**-*-*-*-*-*-*--*-**-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*
// +*-+-+-+-+-+-+*---**-*--*-*-*-*-*-*-**-*----*--*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-***-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-**-*-*-*-*-*-*-*
