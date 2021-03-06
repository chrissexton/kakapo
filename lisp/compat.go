package lisp

import (
	"fmt"
	"path"
	"reflect"
)

// The map of available imports
var _go_imports = map[string]map[string]interface{}{}

func ExposeImport(name string, pkg map[string]interface{}) {
	_go_imports[name] = pkg
}

// Expose an identifier globally.
func ExposeGlobal(id string, x interface{}) {
	global.define(sym(id), wrapGo(x))
}

func builtinImport(sc *scope, ss []sexpr) sexpr {
	if len(ss) != 1 {
		panic("Invalid number of arguments")
	}

	pkgPath, ok := ss[0].(string)
	if !ok {
		panic("Invalid argument")
	}

	pkgName := path.Base(pkgPath)

	// find the package in _go_imports
	pkg, found := _go_imports[pkgPath]
	if !found {
		panic("Package not found")
	}

	// import each item
	for name, _go := range pkg {
		if name[0] != '\x00' {
			sc.define(sym(pkgName+"."+name), wrapGo(_go))
		} else {
			// import all methods on this object
			importMethods(sc, pkgName, name[1:], _go)
			// and for the pointer version as well
			t := _go.(reflect.Type)
			importMethods(sc, pkgName, name[1:], reflect.PtrTo(t))
		}
	}
	return Nil
}

func importMethods(sc *scope, pkgName, name string, r interface{}) {
	t := r.(reflect.Type)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		n := fmt.Sprintf("%s.%s.%s", pkgName, name, m.Name)
		mName := m.Name
		sc.define(sym(n), function(func(sc *scope, ss []sexpr) sexpr {
			if len(ss) == 0 {
				panic("Invalid number of arguments")
			}
			v := reflect.ValueOf(ss[0])
			if v.Type() != t {
				panic("Invalid argument")
			}
			fun := v.MethodByName(mName)
			t := fun.Type()
			ni := t.NumIn()
			ss = ss[1:]
			if ni != len(ss) && !t.IsVariadic() {
				panic("Invalid number of arguments")
			}

			vs := make([]reflect.Value, len(ss))
			for i, s := range ss {
				// get argument type
				var at reflect.Type
				if t.IsVariadic() && i >= ni-1 {
					st := t.In(ni - 1)
					at = st.Elem()
				} else {
					at = t.In(i)
				}
				// TODO convert any cons and function arguments
				vs[i] = forGo(s, at)
			}
			r := fun.Call(vs)
			if len(r) == 0 {
				return Nil
			}
			return wrapGoval(r[0])
		}))
	}
}

// Wrap the given value as a usable s-expression.
func wrapGo(_go interface{}) sexpr {
	return wrapGoval(reflect.ValueOf(_go))
}

func wrapGoval(r reflect.Value) sexpr {
	typ := r.Type()
	kind := typ.Kind()
	switch kind {
	case reflect.Bool:
		b := r.Bool()
		if b {
			return float64(1)
		} else {
			return Nil
		}
	case reflect.Int:
		return float64(r.Int())
	case reflect.Int8:
		return float64(r.Int())
	case reflect.Int16:
		return float64(r.Int())
	case reflect.Int32:
		return float64(r.Int())
	case reflect.Int64:
		return float64(r.Int())
	case reflect.Uint:
		return float64(r.Uint())
	case reflect.Uint8:
		return float64(r.Uint())
	case reflect.Uint16:
		return float64(r.Uint())
	case reflect.Uint32:
		return float64(r.Uint())
	case reflect.Uint64:
		return float64(r.Uint())
	case reflect.Uintptr:
		return native(r.Interface()) // TODO
	case reflect.Float32:
		return float64(r.Float())
	case reflect.Float64:
		return float64(r.Float())
	case reflect.Complex64:
		return native(r.Interface()) // TODO
	case reflect.Complex128:
		return native(r.Interface()) // TODO
	case reflect.Array:
		return native(r.Interface()) // TODO
	case reflect.Chan:
		return native(r.Interface()) // TODO
	case reflect.Func:
		return wrapFunc(r.Interface())
	case reflect.Interface:
		return native(r.Interface()) // TODO
	case reflect.Map:
		return native(r.Interface()) // TODO
	case reflect.Ptr:
		return native(r.Interface()) // TODO
	case reflect.Slice:
		return native(r.Interface()) // TODO
	case reflect.String:
		return r.String()
	case reflect.Struct:
		return native(r.Interface()) // TODO
	case reflect.UnsafePointer:
		return native(r.Interface()) // TODO
	}
	return Nil
}

func forGo(v sexpr, typ reflect.Type) reflect.Value {
	kind := typ.Kind()
	switch kind {
	case reflect.Bool:
		return reflect.ValueOf(v != Nil)
	case reflect.Int:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(int(f))
	case reflect.Int8:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(int8(f))
	case reflect.Int16:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(int16(f))
	case reflect.Int32:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(int32(f))
	case reflect.Int64:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(int64(f))
	case reflect.Uint:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(uint(f))
	case reflect.Uint8:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(uint8(f))
	case reflect.Uint16:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(uint16(f))
	case reflect.Uint32:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(uint32(f))
	case reflect.Uint64:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(uint64(f))
	case reflect.Uintptr:
		panic("Invalid argument") // TODO
	case reflect.Float32:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(float32(f))
	case reflect.Float64:
		f, ok := v.(float64)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(f)
	case reflect.Complex64:
		panic("Invalid argument") // TODO
	case reflect.Complex128:
		panic("Invalid argument") // TODO
	case reflect.Array:
		panic("Invalid argument") // TODO
	case reflect.Chan:
		panic("Invalid argument") // TODO
	case reflect.Func:
		// a reflect.Func is expected of our sexpr
		panic("Cannot do raw callbacks yet, sorry") // XXX TODO
	case reflect.Interface:
		// TODO do some checks
		reflect.ValueOf(v)
	case reflect.Map:
		panic("Invalid argument") // TODO
	case reflect.Ptr:
		panic("Invalid argument") // TODO
	case reflect.Slice:
		panic("Invalid argument") // TODO
	case reflect.String:
		s, ok := v.(string)
		if !ok {
			panic("Invalid argument")
		}
		return reflect.ValueOf(s)
	case reflect.Struct:
		panic("Invalid argument") // TODO
	case reflect.UnsafePointer:
		panic("Invalid argument") // can't handle this
	}
	return reflect.ValueOf(v)
}

func wrapFunc(f interface{}) function {
	// TODO patch reflect so we can do type compatibility-checking
	return func(sc *scope, ss []sexpr) sexpr {
		fun := reflect.ValueOf(f)

		t := fun.Type()
		ni := t.NumIn()
		if ni != len(ss) && !t.IsVariadic() {
			panic("Invalid number of arguments")
		}

		vs := make([]reflect.Value, len(ss))
		for i, s := range ss {
			// get argument type
			var at reflect.Type
			if t.IsVariadic() && i >= ni-1 {
				st := t.In(ni - 1)
				at = st.Elem()
			} else {
				at = t.In(i)
			}
			// TODO convert any cons and function arguments
			vs[i] = forGo(s, at)
		}
		r := fun.Call(vs)
		if len(r) == 0 {
			return Nil
		}
		return wrapGoval(r[0])
	}
}
