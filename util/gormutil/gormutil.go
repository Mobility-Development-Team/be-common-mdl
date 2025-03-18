package gu

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// BuildOption Setup Preload option
type BuildOption struct {
	PreloadName      string
	OnlySelectFields []string
	WhereMap         map[string]interface{}
}
type GormScopesOption struct {
	Obj            interface{}
	PreloadOptions []PreloadOption
	IsDebug        bool
}

type PreloadOption struct {
	BuildOption
	ExtraAnonymousStructs []string
}

type GormScopesEasilyOption struct {
	Obj                       interface{}
	PreloadNames              []string
	OnlySelectFields          []string
	WhereMap                  map[string]map[string]interface{}
	EachExtraAnonymousStructs []string
	IsDebug                   bool
}

const ConstTagGormutil = "gu"
const ConstTagExtractField = "extractField"

const ConstFirstParentKey = "Obj"
const ConstAnonymousParent = "AnonymousEmbed"

// Deprecated: 切片指针获取不出来，有问题
func ReflectStruct(s interface{}, parent string, result map[string][]string) (string, error) {
	valueOf := reflect.ValueOf(s)
	if valueOf.Kind() != reflect.Struct {
		return "", fmt.Errorf("input is not a struct type: %+v", valueOf.Kind())
	}
	typeOf := reflect.TypeOf(s)

	numFields := valueOf.NumField()
	for i := 0; i < numFields; i++ {
		fieldValue := valueOf.Field(i)
		fieldType := typeOf.Field(i)
		isAnonymousStruct := fieldType.Anonymous
		tagValue := fieldType.Tag.Get(ConstTagGormutil)

		nextParent := parent + "." + fieldType.Name

		// 如果是指针类型，则判断实际指向是什么类型
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				// 创建一个指针指向类型的实例
				elemType := fieldValue.Type().Elem()
				fieldValue = reflect.New(elemType).Elem()
				// newInstance := fieldValue.Interface()
				// if reflect.ValueOf(newInstance).Kind() == reflect.Struct {
				// 	_, err := ReflectStruct(newInstance, nextParent, result)
				// 	if err != nil {
				// 		return "", err
				// 	}
				// }
			}
		}
		var t1 *time.Time
		switch {
		case fieldValue.Kind() == reflect.Struct && !isAnonymousStruct &&
			fieldType.Type != reflect.TypeOf(time.Time{}) &&
			fieldType.Type != reflect.TypeOf(t1):
			_, err := ReflectStruct(fieldValue.Interface(), nextParent, result)
			if err != nil {
				return "", err
			}
		case IsSliceOfStruct(fieldValue):
			if fieldValue.Len() == 0 || fieldValue.IsNil() {
				// 创建一个新的实例并传递给递归调用
				sliceElemType := fieldValue.Type().Elem()
				newInstance := reflect.New(sliceElemType).Elem().Interface()
				if reflect.ValueOf(newInstance).Kind() == reflect.Ptr && reflect.ValueOf(newInstance).Elem().Kind() == reflect.Struct {
					newInstance = reflect.ValueOf(newInstance).Elem().Interface()
				}
				_, err := ReflectStruct(newInstance, nextParent, result)
				if err != nil {
					return "", err
				}
			} else {
				_, err := ReflectStruct(fieldValue.Index(0).Interface(), nextParent, result)
				if err != nil {
					return "", err
				}
			}

		case IsSliceOfPrt(fieldValue):
			if fieldValue.Len() > 0 && !fieldValue.IsNil() {
				// 只处理非空切片指针
				_, err := ReflectStruct(fieldValue.Index(0).Interface(), nextParent, result)
				if err != nil {
					return "", err
				}
			}
		case isAnonymousStruct:
			// 是匿名结构体，需要把匿名结构体Fields提取出来放到它的嵌套的节点Fields里
			aEmbedNewResult := map[string][]string{}
			_, err := ReflectStruct(fieldValue.Interface(), ConstAnonymousParent, aEmbedNewResult)
			if err != nil {
				return "", err
			}
			result[parent] = append(result[parent], aEmbedNewResult[ConstAnonymousParent]...)
		default:
			result[parent] = append(result[parent], fieldType.Name)
		}

		if tagValue == ConstTagExtractField {
			//如果是extractField,需额外提交到所在的结构体的fields里
			result[parent] = append(result[parent], fieldType.Name)
		}
	}
	return parent, nil
}

// IsSliceOfStruct 判断给定的 reflect.Value 是否为包含结构体元素的切片 or 是否包含指针结构体
func IsSliceOfStruct(fieldValue reflect.Value) bool {
	if fieldValue.Kind() == reflect.Slice {
		sliceElemType := fieldValue.Type().Elem()
		if sliceElemType.Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

// IsSliceOfPrt 判断给定的 reflect.Value 是否为包含指针的切片
func IsSliceOfPrt(fieldValue reflect.Value) bool {
	if fieldValue.Kind() == reflect.Slice {
		sliceElemType := fieldValue.Type().Elem()

		// 如果是指针类型，则判断实际指向是什么类型
		if sliceElemType.Kind() == reflect.Ptr {
			return true
		}
	}
	return false
}

/*
Deprecated: 切片指针获取不出来，有问题
  - @description: 提取obj struct 字段，自动生成preload select 等语句；eg:
    -- []func(*gorm.DB) *gorm.DB{
    SelectBuild(db, []string{"name", "id"}),
    PreloadBuild(db, BuildOption{
    OnlySelectFields: []string{"age", "user_id", "id"},
    preloadName:      "Profile",
    }),
    PreloadBuild(db, BuildOption{
    OnlySelectFields: []string{"name", "no", "user_id", "id"},
    preloadName:      "Company",
    }),
    PreloadBuild(db, BuildOption{
    OnlySelectFields: []string{"name", "no", "company_id", "id"},
    preloadName:      "Company.Address",
    })
    }

- @param {*gorm.DB} db

- @param {GormScopesOption} GormScopesOption 后续需要其他sql语法，可以在这里扩展即可，如order by之类的
-- {struct{}} Obj 必填;需要提取Fields的结构体
-- {[]PreloadOption} PreloadOptions 需要提取Fields的结构体
--- {string} PreloadName 可空，空时只会查主表;需要preload的表名，注意这是结构体的key，不是类型 eg:Companies|Companies.Address|Profile
--- {[]string} OnlySelectFields 可空;是空数组就自动提取Obj的Fileds查询；有字符串切片，只会用字符串切片查询
--- {map[string]interface{}} WhereMap 可空;PreloadName的where条件，eg:map[string]interface{}{"age":1}
--- {[]string} ExtraAnonymousStructs 可空;匿名嵌套结构体，需要在PreloadName查出的Fields里额外增加这些匿名结构体的字段，eg:传入[]string{"Model"}，会把Model里面的所有一层字段赋值给PreloadName的Fields里

- @return {[]func(*gorm.DB) *gorm.DB} scopesFunc

- @return {error} err
*/
func GetGormScopes(db *gorm.DB, option GormScopesOption) (scopesFunc []func(*gorm.DB) *gorm.DB, err error) {
	structFieldsMap := make(map[string][]string)
	structName, err := ReflectStruct(option.Obj, ConstFirstParentKey, structFieldsMap)
	if err != nil {
		return nil, err
	}
	if option.IsDebug {
		fmt.Printf("[GetGromScopes]option:\n%+v;\nstructFieldsMap:\n%+v\n", option, structFieldsMap)
	}
	// select处理
	if selectFields, exist := structFieldsMap[structName]; exist {
		scopesFunc = append(scopesFunc, SelectBuild(db, selectFields))
	}
	// preload处理
	for _, preloadOption := range option.PreloadOptions {
		key := structName + "." + preloadOption.PreloadName
		currentFields := preloadOption.OnlySelectFields

		if len(currentFields) == 0 {
			if fields, exist := structFieldsMap[key]; exist {
				currentFields = fields
			} else {
				continue
			}
		}

		for _, extraStruct := range preloadOption.ExtraAnonymousStructs {
			if extraFields, exist := structFieldsMap[key+"."+extraStruct]; exist {
				currentFields = append(currentFields, extraFields...)
			}
		}

		preloadOption.OnlySelectFields = currentFields
		if option.IsDebug {
			fmt.Printf("[GetGromScopes]preloadName:%+v; currentFields:%+v\n", preloadOption.PreloadName, currentFields)
		}

		scopesFunc = append(scopesFunc, PreloadBuild(db, preloadOption.BuildOption))
	}

	return scopesFunc, nil
}

// Deprecated: 切片指针获取不出来，有问题
func GetGormScopesEasily(db *gorm.DB, option GormScopesEasilyOption) (scopesFunc []func(*gorm.DB) *gorm.DB, err error) {

	preloadOptions := []PreloadOption{}
	for _, pn := range option.PreloadNames {
		preloadOptionsTemp := PreloadOption{}
		preloadOptionsTemp.PreloadName = pn
		// 这里每一层都是去找同个名字field字段
		preloadOptionsTemp.ExtraAnonymousStructs = option.EachExtraAnonymousStructs

		if len(option.WhereMap) > 0 {
			wm, exist := option.WhereMap[pn]
			if exist {
				preloadOptionsTemp.BuildOption.WhereMap = wm
			}
		}
		preloadOptions = append(preloadOptions, preloadOptionsTemp)
	}

	return GetGormScopes(db,
		GormScopesOption{
			Obj:            option.Obj,
			PreloadOptions: preloadOptions,
			IsDebug:        option.IsDebug,
		},
	)
}

func PreloadBuild(db *gorm.DB, option BuildOption) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(option.PreloadName, WhereBuild(db, option.WhereMap), SelectBuild(db, option.OnlySelectFields))
	}
}

func PreloadGet(db *gorm.DB, preloadName string, onlySelectFields []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(preloadName, SelectBuild(db, onlySelectFields))
	}
}

func WhereBuild(db *gorm.DB, whereMap map[string]interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(whereMap) == 0 {
			return db
		}
		return db.Where(whereMap)
	}
}

func SelectBuild(db *gorm.DB, onlySelectFields []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(onlySelectFields) == 0 {
			return db
		}
		return db.Select(onlySelectFields)
	}
}
