package e_tests

import (
	"testing"
	sentinel "github.com/elcengine/elemental/plugins/sentinel"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestValidator(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Basic validations", t, func() {

		type CreateUserDTO struct {
			Name string `augmented_validate:"unique=users" json:"name"`
			Age  int    `validate:"max=150,min=18" json:"age"`
		}

		Convey("Inherited validations", func() {
			request := CreateUserDTO{
				Name: e_mocks.Eredin.Name,
				Age:  10,
			}
			err := sentinel.Legitimize(request)
			So(err.Error(), ShouldEqual, "Key: 'CreateUserDTO.Age' Error:Field validation for 'Age' failed on the 'min' tag")
		})

		Convey("Unique document validation", func() {
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
	})
}
