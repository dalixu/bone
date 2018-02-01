package bone

import (
	"fmt"
	"reflect"
)

//Injector 依赖注入
type Injector struct {
	Parent *Injector

	Provider []interface{}
}

//Call 注意这里ij可以为nil 只要fv的参数为0 Panic if t is not kind of Func
func (ij *Injector) Call(fv reflect.Value) ([]reflect.Value, error) {
	ft := fv.Type()
	var in = make([]reflect.Value, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		argType := ft.In(i)
		val := ij.Get(argType)
		if !val.IsValid() {
			return nil, fmt.Errorf("Value not found for type %v", argType)
		}

		in[i] = val
	}
	return fv.Call(in), nil
}

//Get 依赖注入 获取指定类型对应的Value
func (ij *Injector) Get(t reflect.Type) reflect.Value {
	for _, v := range ij.Provider {

		vt := reflect.TypeOf(v)
		if vt == t || (t.Kind() == reflect.Interface && vt.Implements(t)) {
			return reflect.ValueOf(v)
		}
	}
	//找不到 往父查找
	if ij.Parent != nil {
		return ij.Parent.Get(t)
	}
	return reflect.Value{}
}
