package decode

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestHookTranslateKeys(t *testing.T) {
	var testcases = []struct {
		name     string
		data     interface{}
		expected interface{}
	}{
		{
			name: "target of type struct, with struct receiver",
			data: map[string]interface{}{
				"S": map[string]interface{}{
					"None":   "no translation",
					"OldOne": "value1",
					"oldtwo": "value2",
				},
			},
			expected: Config{
				S: TypeStruct{
					One:  "value1",
					Two:  "value2",
					None: "no translation",
				},
			},
		},
		{
			name: "target of type ptr, with struct receiver",
			data: map[string]interface{}{
				"PS": map[string]interface{}{
					"None":   "no translation",
					"OldOne": "value1",
					"oldtwo": "value2",
				},
			},
			expected: Config{
				PS: &TypeStruct{
					One:  "value1",
					Two:  "value2",
					None: "no translation",
				},
			},
		},
		{
			name: "target of type ptr, with ptr receiver",
			data: map[string]interface{}{
				"PTR": map[string]interface{}{
					"None":      "no translation",
					"old_THREE": "value3",
					"old_four":  "value4",
				},
			},
			expected: Config{
				PTR: &TypePtrToStruct{
					Three: "value3",
					Four:  "value4",
					None:  "no translation",
				},
			},
		},
		{
			name: "target of type ptr, with struct receiver",
			data: map[string]interface{}{
				"PTRS": map[string]interface{}{
					"None":      "no translation",
					"old_THREE": "value3",
					"old_four":  "value4",
				},
			},
			expected: Config{
				PTRS: TypePtrToStruct{
					Three: "value3",
					Four:  "value4",
					None:  "no translation",
				},
			},
		},
		{
			name: "target of type map",
			data: map[string]interface{}{
				"Blob": map[string]interface{}{
					"one": 1,
					"two": 2,
				},
			},
			expected: Config{
				Blob: map[string]interface{}{
					"one": 1,
					"two": 2,
				},
			},
		},
		{
			name: "value already exists for canonical key",
			data: map[string]interface{}{
				"PS": map[string]interface{}{
					"OldOne": "value1",
					"One":    "original1",
					"oldTWO": "value2",
					"two":    "original2",
				},
			},
			expected: Config{
				PS: &TypeStruct{
					One: "original1",
					Two: "original2",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{}
			md := new(mapstructure.Metadata)
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				DecodeHook: HookTranslateKeys,
				Metadata:   md,
				Result:     &cfg,
			})
			require.NoError(t, err)

			require.NoError(t, decoder.Decode(tc.data))
			require.Equal(t, cfg, tc.expected, "decode metadata: %#v", md)
		})
	}
}

type Config struct {
	S    TypeStruct
	PS   *TypeStruct
	PTR  *TypePtrToStruct
	PTRS TypePtrToStruct
	Blob map[string]interface{}
}

type TypeStruct struct {
	One  string
	Two  string
	None string
}

func (m TypeStruct) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"oldone": "One",
		"oldtwo": "two",
	}
}

type TypePtrToStruct struct {
	Three string
	Four  string
	None  string
}

func (m *TypePtrToStruct) DecodeKeyMapping() map[string]string {
	return map[string]string{
		"old_three": "Three",
		"old_four":  "four",
		"oldfour":   "four",
	}
}

func TestHookTranslateKeys_TargetStructHasPointerReceiver(t *testing.T) {
	target := &TypePtrToStruct{}
	md := new(mapstructure.Metadata)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: HookTranslateKeys,
		Metadata:   md,
		Result:     target,
	})
	require.NoError(t, err)

	data := map[string]interface{}{
		"None":      "no translation",
		"Old_Three": "value3",
		"OldFour":   "value4",
	}
	expected := &TypePtrToStruct{
		None:  "no translation",
		Three: "value3",
		Four:  "value4",
	}
	require.NoError(t, decoder.Decode(data))
	require.Equal(t, target, expected, "decode metadata: %#v", md)
}
