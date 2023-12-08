package common

import "reflect"

func FindItemInArray[T any, B any](objectiveArray []T, item B, fieldName string) (int, *T) {
	itemValue := reflect.ValueOf(item)
	itemField, ok := extractField(itemValue, fieldName)
	if !ok {
		return -1, nil
	}

	for idx, currentItem := range objectiveArray {
		currentValue := reflect.ValueOf(currentItem)
		currentField, ok := extractField(currentValue, fieldName)
		if ok && reflect.DeepEqual(currentField, itemField) {
			return idx, &objectiveArray[idx]
		}
	}
	return -1, nil
}

// extractField intenta extraer un campo con el nombre dado del valor dado.
// Devuelve el valor del campo y un booleano que indica si la extracción fue exitosa.
func extractField(v reflect.Value, fieldName string) (interface{}, bool) {
	if v.Kind() == reflect.Struct {
		field := v.FieldByName(fieldName)
		if field.IsValid() {
			return field.Interface(), true
		}
	}
	return nil, false
}

func ItemHasChanges[T comparable](newValue *T, originalValue *T) bool {
	if newValue == nil && originalValue == nil {
		return false
	}
	if (newValue == nil && originalValue != nil) || (newValue != nil && originalValue == nil) {
		return true
	}

	// Obtener los valores reflectidos
	newValueReflect := reflect.ValueOf(newValue)
	originalValueReflect := reflect.ValueOf(originalValue)

	// Dereferenciar hasta obtener el valor subyacente real
	for newValueReflect.Kind() == reflect.Ptr {
		newValueReflect = newValueReflect.Elem()
	}
	for originalValueReflect.Kind() == reflect.Ptr {
		originalValueReflect = originalValueReflect.Elem()
	}

	// Usar reflect.DeepEqual para comparar los valores
	return !reflect.DeepEqual(newValueReflect.Interface(), originalValueReflect.Interface())
}
