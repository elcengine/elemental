package e_tests

import (
	"testing"

	elemental "github.com/elcengine/elemental/core"
	sentinel "github.com/elcengine/elemental/plugins/sentinel"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestValidator(t *testing.T) {

	e_test_setup.SeededConnection()

	elemental.NativeModel.SetCollection("occupations").InsertMany([]map[string]any{
		{
			"occupation":     "Witcher",
			"minimum_income": 100,
		},
		{
			"occupation":     "Mage",
			"minimum_income": 200,
		},
	}).Exec()

	defer e_test_setup.Teardown()

	Convey("Basic validations", t, func() {

		Convey("Inherited validations", func() {
			type CreateUserDTO struct {
				Name string `augmented_validate:"unique=users" json:"name"`
				Age  int    `validate:"max=150,min=18" json:"age"`
			}
			request := CreateUserDTO{
				Name: e_mocks.Eredin.Name,
				Age:  10,
			}
			err := sentinel.Legitimize(request)
			So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Age' Error:Field validation for 'Age' failed on the 'min' tag")
		})

		Convey("Unique document validation", func() {
			type CreateUserDTO struct {
				Name string `augmented_validate:"unique=users" json:"name"`
				Age  int    `validate:"max=150,min=18" json:"age"`
			}
			Convey("Should return error if document already exists", func() {
				request := CreateUserDTO{
					Name: e_mocks.Caranthir.Name,
					Age:  100,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Name' Error:Field validation for 'Name' failed on the 'unique' tag")
			})
			Convey("Should not return error if document does not exist", func() {
				request := CreateUserDTO{
					Name: "Foltest",
					Age:  30,
				}
				err := sentinel.Legitimize(request)
				So(err, ShouldBeNil)
			})
			Convey("Should return error if document already exists - DTO specifying model", func() {
				type CreateUserDTOWithModel struct {
					Name string `augmented_validate:"unique=User" json:"name"`
					Age  int    `validate:"max=150,min=18" json:"age"`
				}
				request := CreateUserDTOWithModel{
					Name: e_mocks.Eredin.Name,
					Age:  100,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTOWithModel.Name' Error:Field validation for 'Name' failed on the 'unique' tag")
			})
			Convey("Should return error if document already exists - DTO specifying custom field", func() {
				type CreateUserDTOWithCustomField struct {
					Name string `augmented_validate:"unique=User->name"`
					Age  int    `validate:"max=150,min=18"`
				}
				request := CreateUserDTOWithCustomField{
					Name: e_mocks.Imlerith.Name,
					Age:  100,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTOWithCustomField.Name' Error:Field validation for 'Name' failed on the 'unique' tag")
			})
			Convey("Should return error if document already exists - DTO specifying custom database name", func() {
				type CreateUserDTOWithCustomDatabase struct {
					Name string `augmented_validate:"unique=User" database:"elemental_secondary" json:"name"`
					Age  int    `validate:"max=150,min=18" json:"age"`
				}
				request := CreateUserDTOWithCustomDatabase{
					Name: "Radovid",
					Age:  100,
				}
				UserModel.SetDatabase(e_mocks.SECONDARY_DB).Create(User{
					Name: "Radovid",
					Age:  100,
				}).Exec()
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTOWithCustomDatabase.Name' Error:Field validation for 'Name' failed on the 'unique' tag")
			})
		})

		Convey("Exists document validation", func() {
			type CreateUserDTO struct {
				Name       string `json:"name"`
				Age        int    `validate:"max=150,min=18" json:"age"`
				Occupation string `augmented_validate:"exists=occupations" json:"occupation"`
			}
			Convey("Should return error if document doesn't already exist", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Caranthir.Name,
					Age:        100,
					Occupation: "Druid",
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Occupation' Error:Field validation for 'Occupation' failed on the 'exists' tag")
			})
			Convey("Should not return error if document exists", func() {
				request := CreateUserDTO{
					Name:       "Letho",
					Age:        30,
					Occupation: "Witcher",
				}
				err := sentinel.Legitimize(request)
				So(err, ShouldBeNil)
			})
		})

		Convey("Greater than document validation", func() {
			type CreateUserDTO struct {
				Name       string `json:"name"`
				Age        int    `validate:"max=150,min=18" json:"age"`
				Occupation string `json:"occupation"`
				Income     int    `augmented_validate:"greater_than=occupations->minimum_income" ref:"occupation" json:"income"`
			}
			Convey("Should return error if document is not greater than the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     50,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Income' Error:Field validation for 'Income' failed on the 'greater_than' tag")
			})
			Convey("Should not return error if document is greater than the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     150,
				}
				err := sentinel.Legitimize(request)
				So(err, ShouldBeNil)
			})
		})

		Convey("Greater than or equal to document validation", func() {
			type CreateUserDTO struct {
				Name       string `json:"name"`
				Age        int    `validate:"max=150,min=18" json:"age"`
				Occupation string `json:"occupation"`
				Income     int    `augmented_validate:"greater_than_or_equal_to=occupations->minimum_income" ref:"occupation" json:"income"`
			}
			Convey("Should return error if document is not greater than or equal to the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     99,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Income' Error:Field validation for 'Income' failed on the 'greater_than_or_equal_to' tag")
			})
			Convey("Should not return error if document is greater than or equal to the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     100,
				}
				err := sentinel.Legitimize(request)
				So(err, ShouldBeNil)
			})
		})

		Convey("Less than document validation", func() {
			type CreateUserDTO struct {
				Name       string `json:"name"`
				Age        int    `validate:"max=150,min=18" json:"age"`
				Occupation string `json:"occupation"`
				Income     int    `augmented_validate:"less_than=occupations->minimum_income" ref:"occupation" json:"income"`
			}
			Convey("Should return error if document is not less than the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     150,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Income' Error:Field validation for 'Income' failed on the 'less_than' tag")
			})
			Convey("Should not return error if document is less than the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     90,
				}
				err := sentinel.Legitimize(request)
				So(err, ShouldBeNil)
			})
		})

		Convey("Less than or equal to document validation", func() {
			type CreateUserDTO struct {
				Name       string `json:"name"`
				Age        int    `validate:"max=150,min=18" json:"age"`
				Occupation string `json:"occupation"`
				Income     int    `augmented_validate:"less_than_or_equal_to=occupations->minimum_income" ref:"occupation" json:"income"`
			}
			Convey("Should return error if document is not less than or equal to the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     151,
				}
				err := sentinel.Legitimize(request)
				So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Income' Error:Field validation for 'Income' failed on the 'less_than_or_equal_to' tag")
			})
			Convey("Should not return error if document is less than or equal to the specified field", func() {
				request := CreateUserDTO{
					Name:       e_mocks.Geralt.Name,
					Age:        e_mocks.Geralt.Age,
					Occupation: e_mocks.Geralt.Occupation,
					Income:     100,
				}
				err := sentinel.Legitimize(request)
				So(err, ShouldBeNil)
			})
		})
	})
}
