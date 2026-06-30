package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	UserAge struct {
		Age int `validate:"min:18|max:50"`
	}
	Role struct {
		Role UserRole `validate:"in:admin,stuff"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	AppInvalidTag struct {
		Version string `validate:"len5"`
	}
	AppInvalidTag2 struct {
		Version string `validate:"len:5y"`
	}

	Nums struct {
		Version string `validate:"regexp:\\d+|len:20"`
	}

	UserEmail struct {
		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func makeVe(field string, err error) ValidationErrors {
	ve := make(ValidationErrors, 1)
	ve[0] = ValidationError{field, err}
	return ve
}

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in:          AppInvalidTag{Version: "1234"},
		expectedErr: makeVe("Version", ErrUnknownTag),
	},
	{
		in:          AppInvalidTag2{Version: "1234"},
		expectedErr: makeVe("Version", ErrBadParameter),
	},
	{
		in:          App{Version: "12345"},
		expectedErr: nil,
	},

	{
		in:          App{Version: "1234"},
		expectedErr: makeVe("Version", ErrLen),
	},
	{
		in:          UserAge{Age: 10},
		expectedErr: makeVe("Age", ErrMin),
	},

	{
		in:          UserAge{Age: 100},
		expectedErr: makeVe("Age", ErrMax),
	},

	{
		in:          Response{Code: 100},
		expectedErr: makeVe("Code", ErrIn),
	},
	{
		in:          Response{Code: 200},
		expectedErr: nil,
	},
	{
		in:          Role{Role: "none"},
		expectedErr: makeVe("Role", ErrIn),
	},
	{
		in:          Role{Role: "admin"},
		expectedErr: nil,
	},
	{
		in:          UserEmail{Email: "ggg*mail.su"},
		expectedErr: makeVe("Email", ErrRegexp),
	},
	{
		in:          UserEmail{Email: "ggg@mail.su"},
		expectedErr: nil,
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			/*tt := tt
			t.Parallel()*/
			err := Validate(tt.in)

			if errors.Is(err, tt.expectedErr) {
				err = tt.expectedErr
			}
			if tt.expectedErr == nil {
				require.NoError(t, err, tt.expectedErr)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
