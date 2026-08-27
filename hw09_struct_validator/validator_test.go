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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
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

	InvalidStruct struct {
		Code int    `validate:"in:200,404,500"`
		Body []byte `validate:"len:6"`
	}

	InvalidRule1 struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty" validate:"regexp:"`
	}

	InvalidRule2 struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty" validate:"minmax:1-100"`
	}

	InvalidRule3 struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty" validate:"regexp:["`
	}
)

var (
	stuff UserRole = "stuff"
	admin UserRole = "admin"

	IDFN     = "ID"
	AgeFN    = "Age"
	EmailFN  = "Email"
	PhonesFN = "Phones"
	RoleFN   = "Role"

	tests = []struct {
		in              interface{}
		expectedErr     error
		testDescription string
	}{
		{ // 0
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Таня",
				Age:    25,
				Email:  "tania@mail.ru",
				Role:   stuff,
				Phones: []string{"79161234560", "79101234561"},
				meta:   nil,
			},
			expectedErr:     nil,
			testDescription: "Различные вадидаторы",
		},
		{ // 1
			in: &User{
				ID:     "123456789012345678901234567890123457",
				Name:   "Ваня",
				Age:    40,
				Email:  "tania2@mail.ru",
				Role:   admin,
				Phones: []string{"79161234562", "79101234563"},
				meta:   nil,
			},
			expectedErr:     nil,
			testDescription: "Проверка работы с сылкой на правильную структуру",
		},
		{ // 2
			in: App{
				Version: "1.2.3",
			},
			expectedErr:     nil,
			testDescription: "Один валидатор",
		},

		{ // 3
			in: Token{
				Header:    []byte("Auth"),
				Payload:   []byte(""),
				Signature: []byte(""),
			},
			expectedErr:     nil,
			testDescription: "Проверка работы, когда не заданы правила валидации",
		},
		{ // 4
			in: Response{
				Code: 404,
				Body: "sdlfdkfdlsmvfld",
			},
			expectedErr:     nil,
			testDescription: "Один валидатор",
		},
		{ // 5
			in: []User{
				{
					ID:     "123456789012345678901234567890123458",
					Name:   "Таня",
					Age:    25,
					Email:  "tania3@mail.ru",
					Role:   stuff,
					Phones: []string{"79161234564", "79101234565"},
					meta:   nil,
				},
				{
					ID:     "123456789012345678901234567890123459",
					Name:   "Ваня",
					Age:    40,
					Email:  "tania4@mail.ru",
					Role:   stuff,
					Phones: []string{"79161234566", "79101234567"},
					meta:   nil,
				},
			},
			expectedErr:     ErrUnsupportedType,
			testDescription: "Проверка ошибки при передаче массива структур",
		},

		{ // 6
			in:              []int{5},
			expectedErr:     ErrUnsupportedType,
			testDescription: "Проверка работы с сылкой на простой тип",
		},

		{ // 7
			in: User{
				ID:     "1",
				Name:   "Матвейка",
				Age:    1,
				Email:  "-",
				Role:   "baby",
				Phones: []string{"1", "2", "3"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrLenValidation,
				},
				ValidationError{
					Field: AgeFN,
					Err:   ErrMinValidation,
				},
				ValidationError{
					Field: EmailFN,
					Err:   ErrRegexpValidation,
				},
				ValidationError{
					Field: RoleFN,
					Err:   ErrEnumValidation,
				},
				ValidationError{
					Field: PhonesFN,
					Err:   ErrLenValidation,
				},
				ValidationError{
					Field: PhonesFN,
					Err:   ErrLenValidation,
				},
				ValidationError{
					Field: PhonesFN,
					Err:   ErrLenValidation,
				},
			},
			testDescription: "Различные ошибки при валидации",
		},

		{ // 8
			in: User{
				ID:     "1234567890123456789012345678901234567",
				Name:   "Петр Иванович",
				Age:    60,
				Email:  "pi_mail.ru",
				Role:   "gf",
				Phones: []string{},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrLenValidation,
				},
				ValidationError{
					Field: AgeFN,
					Err:   ErrMaxValidation,
				},
				ValidationError{
					Field: EmailFN,
					Err:   ErrRegexpValidation,
				},
				ValidationError{
					Field: RoleFN,
					Err:   ErrEnumValidation,
				},
			},
			testDescription: "Различные ошибки при валидации некоторых полей",
		},

		{ // 9
			in: App{
				Version: "1.2.3.4",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrLenValidation,
				},
			},
			testDescription: "Один несрабатывающий валидатор",
		},

		{ // 10
			in: Response{
				Code: 403,
				Body: "sdlfdkfdlsmvfld",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrEnumValidation,
				},
			},
			testDescription: "Один неверный валидатор",
		},
		{ // 11
			in: InvalidStruct{
				Code: 403,
				Body: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
			},
			expectedErr:     ErrUnsupportedField,
			testDescription: "Один неверный валидатор",
		},

		{ // 12
			in: InvalidRule1{
				Code: 200,
				Body: "fdsfqwrfesadfvsawfcd",
			},
			expectedErr:     ErrInvalidRuleValue,
			testDescription: "Пустое правило в валидаторе",
		},
		{ // 13
			in: InvalidRule2{
				Code: 500,
				Body: "sdfdfsedgwqeqe",
			},
			expectedErr:     ErrInvalidRule,
			testDescription: "Неверное правило в валидаторе",
		},
		{ // 14
			in: InvalidRule3{
				Code: 400,
				Body: "weewetasdawe3q",
			},
			expectedErr:     ErrInvalidRuleValue,
			testDescription: "Неправильное регуляное выражение в правиле валидатора",
		},
	}
)

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr != nil {
				if errors.As(err, &ValidationErrors{}) {
					require.EqualError(t, err, tt.expectedErr.Error())
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				require.Nil(t, err)
			}
			_ = tt
		})
	}
}
