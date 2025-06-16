package tests

import (
	"fmt"
	"regexp"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
	"github.com/samber/lo"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCoreSchemaOptions(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Schema variations", t, func() {
		Convey(fmt.Sprintf("Should use the default database of %s", t.Name()), func() {
			So(UserModel.Collection().Database().Name(), ShouldEqual, t.Name())
		})
		Convey("Should automatically create a collection with given name", func() {
			So(UserModel.Collection().Name(), ShouldEqual, "users")
		})
		Convey("Collection should be a plural of the model name if not specified", func() {
			CastleModel := elemental.NewModel[Castle]("Castle", elemental.NewSchema(map[string]elemental.Field{
				"Name": {
					Type:     elemental.String,
					Required: true,
				},
			})).SetDatabase(t.Name())
			CastleModel.Create(Castle{Name: "Kaer Morhen"}).Exec()
			So(CastleModel.Collection().Name(), ShouldEqual, "castles")
		})
		Convey("Should create a capped collection", func() {
			collectionOptions := options.CreateCollectionOptions{}
			collectionOptions.SetCapped(true)
			collectionOptions.SetSizeInBytes(1024)
			KingdomModel := elemental.NewModel[Kingdom]("Kingdom-Temporary", elemental.NewSchema(map[string]elemental.Field{
				"Name": {
					Type:     elemental.String,
					Required: true,
				},
			}, elemental.SchemaOptions{
				Database:          t.Name(),
				CollectionOptions: collectionOptions,
			}))
			KingdomModel.Create(Kingdom{Name: "Nilfgaard"}).Exec()
			So(KingdomModel.IsCapped(), ShouldBeTrue)
		})
		Convey("Should use the specified database", func() {
			DATABASE := fmt.Sprintf("%s_%s", t.Name(), "temporary_1")
			MonsterModel := elemental.NewModel[Monster]("Monster-Temporary", elemental.NewSchema(map[string]elemental.Field{
				"Name": {
					Type:     elemental.String,
					Required: true,
				},
			}, elemental.SchemaOptions{
				Database: DATABASE,
			}))
			MonsterModel.Create(Monster{Name: "Nekker"}).Exec()
			So(MonsterModel.Collection().Database().Name(), ShouldEqual, DATABASE)
		})
		Convey("Should validate a document against the schema", func() {
			Convey("Required field", func() {
				So(func() {
					UserModel.Validate(User{})
				}, ShouldPanicWith, fmt.Errorf("field Name is required"))
				So(func() {
					UserModel.Validate(User{Name: "Geralt"})
				}, ShouldNotPanic)
			})
			Convey("Required field with default", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type:     elemental.String,
						Required: true,
						Default:  "Placeholder",
					},
				}))
				So(func() {
					Model.Validate(User{})
				}, ShouldPanicWith, fmt.Errorf("field Name is required"))
				So(func() {
					Model.Validate(User{Name: "Geralt"})
				}, ShouldNotPanic)
			})
			Convey("Pointer type check", func() {
				type User struct {
					Name *string `json:"name"`
				}
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type: elemental.String,
					},
				}))
				So(func() {
					Model.Validate(User{Name: lo.ToPtr("Geralt")})
				}, ShouldNotPanic)
			})
			Convey("Type check", func() {
				type InvalidUser struct {
					Name int64 `json:"name"`
				}
				Model := elemental.NewModel[InvalidUser](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type: elemental.String,
					},
				}))
				So(func() {
					Model.Validate(InvalidUser{Name: 12345})
				}, ShouldPanicWith, fmt.Errorf("field Name has an invalid type. It must be of type string"))
			})
			Convey("Min check", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Age": {
						Type: elemental.Int,
						Min:  10,
					},
				}))
				So(func() {
					Model.Validate(User{Age: 5})
				}, ShouldPanicWith, fmt.Errorf("field Age must be greater than or equal to 10"))
				So(func() {
					Model.Validate(User{Age: 15})
				}, ShouldNotPanic)
				So(func() {
					Model.Validate(User{Age: 10})
				}, ShouldNotPanic)
			})
			Convey("Max check", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Age": {
						Type: elemental.Int,
						Max:  120,
					},
				}))
				So(func() {
					Model.Validate(User{Age: 121})
				}, ShouldPanicWith, fmt.Errorf("field Age must be less than or equal to 120"))
				So(func() {
					Model.Validate(User{Age: 50})
				}, ShouldNotPanic)
				So(func() {
					Model.Validate(User{Age: 120})
				}, ShouldNotPanic)
			})
			Convey("Length check", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type:   elemental.String,
						Length: 10,
					},
				}))
				So(func() {
					Model.Validate(User{Name: "Geralt"})
				}, ShouldNotPanic)
				So(func() {
					Model.Validate(User{Name: "Geralt of Rivia"})
				}, ShouldPanicWith, fmt.Errorf("field Name must be less than or equal to 10 characters"))
			})
			Convey("Regex check", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type:  elemental.String,
						Regex: regexp.MustCompile("^[A-Z]+$"),
					},
				}))
				So(func() {
					Model.Validate(User{Name: "G1Cc"})
				}, ShouldPanicWith, fmt.Errorf("field Name must match the regex pattern ^[A-Z]+$"))
				So(func() {
					Model.Validate(User{Name: "GERALT"})
				}, ShouldNotPanic)
			})
			Convey("Ignore non existing definitions", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type: elemental.String,
					},
					"NonExistingField": {
						Type: elemental.String,
					},
				}))
				So(func() {
					Model.Validate(User{Name: "Geralt"})
				}, ShouldNotPanic)
			})
		})
		Convey("Should use default values when provided", func() {
			Convey("Default value of a primitive", func() {
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type: elemental.String,
					},
					"Age": {
						Type:    elemental.Int,
						Default: 30,
					},
				}))
				user := Model.Create(User{Name: uuid.NewString()}).ExecT()
				So(user.Age, ShouldEqual, 30)
			})
			Convey("Default value of a struct", func() {
				type UserPreferences struct {
					Language string `json:"language" bson:"language"`
					Theme    string `json:"theme" bson:"theme"`
				}
				type User struct {
					Name        string           `json:"name" bson:"name"`
					Preferences *UserPreferences `json:"preferences" bson:"preferences"`
				}
				Model := elemental.NewModel[User](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
					"Name": {
						Type: elemental.String,
					},
					"Preferences": {
						Type: elemental.Struct,
						Default: UserPreferences{
							Language: "en",
							Theme:    "light",
						},
					},
				}))
				user := Model.Create(User{Name: uuid.NewString()}).ExecT()
				So(user.Preferences.Language, ShouldEqual, "en")
				So(user.Preferences.Theme, ShouldEqual, "light")
			})
		})
	})
}
