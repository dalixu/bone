package bone

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//Scope 数据绑定 上下文
type Scope interface {
	Parent() Scope
	FieldByName(key string, tp reflect.Type) (interface{}, error)
	MethodByName(key string, tp reflect.Type) (interface{}, error)
	InterfaceByName(keys []string, tp reflect.Type, find func(reflect.Value, string, bool) reflect.Value) (interface{}, error)
}

//数据绑定过程中需要通过Scope树来查找绑定的变量
type scope struct {
	data   interface{} //上下文
	key    string      //用于支持dom 扩展指令 for
	parent Scope
}

//CreateScope 创建上下文
func CreateScope(data interface{}, key string, parent Scope) Scope {
	return scope{
		data:   data,
		key:    key,
		parent: parent,
	}
}

func (se scope) Parent() Scope {
	return se.parent
}

func (se scope) FieldByName(key string, tp reflect.Type) (interface{}, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("empty method or field name")
	}
	return se.InterfaceByName(strings.Split(key, "."), tp, func(src reflect.Value, key string, last bool) reflect.Value {
		if src.Kind() == reflect.Ptr {
			src = src.Elem()
		}
		return src.FieldByName(key)
	})
}

func (se scope) MethodByName(key string, tp reflect.Type) (interface{}, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("empty method or field name")
	}
	return se.InterfaceByName(strings.Split(key, "."), tp, func(src reflect.Value, key string, last bool) reflect.Value {
		if last {
			return src.MethodByName(key)
		}
		if src.Kind() == reflect.Ptr {
			src = src.Elem()
		}
		return src.FieldByName(key)

	})
}

func (se scope) InterfaceByName(keys []string, tp reflect.Type, find func(host reflect.Value, key string, last bool) reflect.Value) (interface{}, error) {
	if len(keys) <= 0 {
		return nil, errors.New("empty method or field name")
	}
	if se.key != "" {
		if se.key != keys[0] {
			if se.parent != nil {
				return se.parent.InterfaceByName(keys, tp, find)
			}
			return nil, fmt.Errorf("can't find method or field %s", keys[0])
		}
		keys = keys[1:]
	}

	i := 0
	src := reflect.ValueOf(se.data)
	for i < len(keys) {
		value := find(src, keys[i], i == len(keys)-1)
		if !value.IsValid() || (tp != nil && value.Type() != tp) {
			if se.parent != nil {
				return se.parent.InterfaceByName(keys, tp, find)
			}
			return nil, fmt.Errorf("can't find method or field %s", keys[i])
		}
		src = value
		i++
	}
	return src.Interface(), nil

}
