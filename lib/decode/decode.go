/*Package decode provides tools for customizing the decoding of configuration
blogs, into structures using mapstructure.
*/
package decode

import (
	"fmt"
	"reflect"
	"strings"
)

// KeyTranslator identifies a structure as one which should have its keys
// translated by HookTranslateKeys.
type KeyTranslator interface {
	// DecodeKeyMapping must return a mapping of:
	//     Lower case deprecated key -> canonical key
	// If the data contains a key which matches the deprecated key, the value
	// for that key is moved to the canonical key. The deprecated key is matched
	// case-insensitive.
	// If the canonical key already exists in the data, the deprecated key is
	// not modified.
	DecodeKeyMapping() map[string]string
}

var typeOfKeyTranslator = reflect.TypeOf((*KeyTranslator)(nil)).Elem()

// HookTranslateKeys is a mapstructure decode hook which translates keys in a
// map to their canonical value. Implement KeyTranslator on a struct to
// have HookTranslateKeys translate keys. See KeyTranslator for more details.
//
// This hook should ONLY be used to maintain backwards compatibility with
// deprecated keys. For new structures use struct tags to set the desired
// serialization key.
//
// TODO: If field aliases were identified as struct tags, instead of a map
// returned from a method on the struct, we could avoid creating instances of
// the type.
func HookTranslateKeys(_, to reflect.Type, data interface{}) (interface{}, error) {
	// When the target is a pointer mapstructure will call the hook again with
	// the struct. Return immediately if target is a pointer to avoid doing the
	// work twice.
	if to.Kind() == reflect.Ptr {
		return data, nil
	}

	// Avoid doing any work if data is not a map
	target, ok := data.(map[string]interface{})
	if !ok {
		return data, nil
	}

	var translator KeyTranslator
	switch {
	case to.Implements(typeOfKeyTranslator):
		translator = reflect.Zero(to).Interface().(KeyTranslator)
	case reflect.PtrTo(to).Implements(typeOfKeyTranslator):
		translator = reflect.New(to).Interface().(KeyTranslator)
	default:
		return data, nil
	}

	rules := translator.DecodeKeyMapping()
	for k, v := range target {
		lowerK := strings.ToLower(k)
		canonKey, ok := rules[lowerK]
		if !ok {
			continue
		}

		// if there is a value for the canonical key then keep it
		if _, ok := target[canonKey]; ok {
			continue
		}

		// otherwise translate to the canonical key
		delete(target, k)
		target[canonKey] = v
	}
	return target, nil
}

// HookNormalizeHCLNestedBlocks looks for []map[string]interface{} is the source
// data. If the target is not a slice or an array it attempts to unpack 1 item
// out of the slice. If there are more items an error is returned.
//
// The HCL language allows one to repeat map keys which forces it to
// store nested structs as []map[string]interface{} instead of
// map[string]interface{}. This is an ambiguity which makes the
// generated structures incompatible with a corresponding JSON
// struct.
//
// TODO: test mapstructure errors properly when source type is incompatible with
// target.
func HookNormalizeHCLNestedBlocks(from, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.Slice && (to.Kind() == reflect.Slice || to.Kind() == reflect.Array) {
		return data, nil
	}

	switch d := data.(type) {
	case []map[string]interface{}:
		switch {
		case len(d) == 0:
			return nil, nil
		case len(d) == 1:
			return d[0], nil
		default:
			// TODO: is it safe to include data in the error? Can it contain secrets?
			return data, fmt.Errorf("HCL nested block with more than one element is not supported: %s", data)
		}
	default:
		return data, nil
	}
}
