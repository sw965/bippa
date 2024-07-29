package bippa

import (
	omwjson "github.com/sw965/omw/json"
)

type TypeData map[Type]float64

func (t TypeData) ToEasyRead() EasyReadTypeData {
	e := EasyReadTypeData{}
	for k, v := range t {
		e[k.ToString()] = v
	}
	return e
}

type EasyReadTypeData map[string]float64

func(e EasyReadTypeData) From() (TypeData, error) {
	d := TypeData{}
	for k, v := range e {
		t, err := StringToType(k)
		if err != nil {
			return TypeData{}, err
		}
		d[t] = v
	}
	return d, nil
}

type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	e, err := omwjson.Load[EasyReadTypedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	d, err := e.From()
	if err != nil {
		panic(err)
	}
	return d
}()

func (t Typedex) EffectivenessValue(atk Type, def Types) float64 {
	v := 1.0
	for _, e := range def {
		v *= t[atk][e]
	}
	return v
}

func (t Typedex) ToEasyRead() EasyReadTypedex {
	e := EasyReadTypedex{}
	for k, v := range t {
		e[k.ToString()] = v.ToEasyRead()
	}
	return e
}

type EasyReadTypedex map[string]EasyReadTypeData

func (e EasyReadTypedex) From() (Typedex, error) {
	d := Typedex{}
	for k, v := range e {
		t, err := StringToType(k)
		if err != nil {
			return Typedex{}, err
		}
		d[t], err = v.From()
		if err != nil {
			return Typedex{}, err
		}
	}
	return d, nil
}