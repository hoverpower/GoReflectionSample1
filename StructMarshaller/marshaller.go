package StructMarshaller

import (
	"fmt"
	"reflect"
)

// allowedFieldKinds содерит список допустимых Kinds
// для прочих типов пришлось бы писать более сложную логику проверки и присвоения
var allowedFieldKinds = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
	reflect.String,
}

var ErrArgumentNotStructure = fmt.Errorf("argument is not a struct")
var ErrFieldUnsupportedKind = fmt.Errorf("field of unsupported king present")
var ErrArgumentNotPointer = fmt.Errorf("argument is not pointer to struct")
var ErrFieldTypeMismatch = fmt.Errorf("field type not match value type")

// StructMarshall принимает структуру, содержащую только логические, числовые и строковые поля,
// и возвращает мапу, в которой ключи - это имена полей, а значения - соответствующие значения полей структуры
// ошибка ErrArgumentNotStructure возникнет, если в функцию передать не структуру,
// ошибка ErrFieldUnsupportedKind возникнет, если у структуры есть поле неподдерживаемого Kind
func StructMarshall(v any) (map[string]any, error) {
	reflVal := reflect.ValueOf(v)

	//Мы работаем только со структурами
	if reflVal.Kind() != reflect.Struct {
		return nil, ErrArgumentNotStructure
	}

	numFields := reflVal.NumField()
	m := make(map[string]any, numFields)

	// проходимся по всем полям структуры
	for i := 0; i < numFields; i++ {
		fieldVal := reflVal.Field(i)

		// Проверяем, работает ли с такими Kind наша функция
		if !checkKind(fieldVal.Kind()) {
			return nil, ErrFieldUnsupportedKind
		}

		// проверяем, экспортируемый ли тип поля
		if !reflVal.Type().Field(i).IsExported() {
			continue
		}
		// если все проверки прошли, то в мапу вписываем ключ: имя поля в структуре,
		// значение: значение поля, обернутое в interface{}
		m[reflVal.Type().Field(i).Name] = fieldVal.Interface()
	}

	return m, nil
}

// StructUnmarshall принимает указатель на структуру и мапу, содержащую строковые имена полей и их значения,
// и заполняет соответствующие поля структуры их значениями. Если поля с нужным именем нет, значение мапы пропускается.
// если тип поля не соответствует типу значения возвращается ErrFieldTypeMismatch. Для полей неподдерживаемого Kind
// возникает ошибка ErrFieldUnsupportedKind.
func StructUnmarshall(v any, fieldsData map[string]any) error {
	reflVal := reflect.ValueOf(v)

	// проверяем, передали ли нам значение по указателю
	if reflVal.Kind() != reflect.Pointer {
		return ErrArgumentNotPointer
	}

	//переходим от указателя к значению
	reflVal = reflVal.Elem()

	// проверяем, вяляется ли значение - структурой
	if reflVal.Kind() != reflect.Struct {
		return ErrArgumentNotStructure
	}

	for key, val := range fieldsData {
		argField := reflVal.FieldByName(key)

		// Проверяем, есть ли поле в структуре, если поля нет - пропускаем это значение в мапе
		if argField.Equal(reflect.Value{}) {
			continue
		}

		// Проверяем, совпадают ли типы у поля и у значения в мапе
		if argField.Type() != reflect.ValueOf(val).Type() {
			return ErrFieldTypeMismatch
		}

		// Проверяем, работает ли с такими Kind наша функция
		if !checkKind(argField.Type().Kind()) {
			return ErrFieldUnsupportedKind
		}

		// Проверяем, можно ли присвоить значение этому полю
		if !argField.CanSet() {
			continue
		}

		//Если все проверки прошли, то присваиваем значение
		argField.Set(reflect.ValueOf(val))
	}
	return nil
}

// checkKind проверяет, входит ли данный Kind в список поддерживаемых пакетом,
// в go 1.21 и выше можно было бы использовать функцию slices.Contains
func checkKind(k reflect.Kind) bool {
	for _, v := range allowedFieldKinds {
		if k == v {
			return true
		}
	}
	return false
}
