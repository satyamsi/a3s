// Code generated by elegen. DO NOT EDIT.
// Source: go.aporeto.io/elemental (templates/model.gotpl)

package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueLDAP represents the model of a issueldap
type IssueLDAP struct {
	// The password for the user.
	Password string `json:"password" msgpack:"password" bson:"-" mapstructure:"password,omitempty"`

	// The LDAP username.
	Username string `json:"username" msgpack:"username" bson:"-" mapstructure:"username,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueLDAP returns a new *IssueLDAP
func NewIssueLDAP() *IssueLDAP {

	return &IssueLDAP{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueLDAP) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueLDAP{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueLDAP) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueLDAP{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueLDAP) BleveType() string {

	return "issueldap"
}

// DeepCopy returns a deep copy if the IssueLDAP.
func (o *IssueLDAP) DeepCopy() *IssueLDAP {

	if o == nil {
		return nil
	}

	out := &IssueLDAP{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueLDAP.
func (o *IssueLDAP) DeepCopyInto(out *IssueLDAP) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueLDAP: %s", err))
	}

	*out = *target.(*IssueLDAP)
}

// Validate valides the current information stored into the structure.
func (o *IssueLDAP) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("password", o.Password); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("username", o.Username); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if len(requiredErrors) > 0 {
		return requiredErrors
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// SpecificationForAttribute returns the AttributeSpecification for the given attribute name key.
func (*IssueLDAP) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := IssueLDAPAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return IssueLDAPLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*IssueLDAP) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return IssueLDAPAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *IssueLDAP) ValueForAttribute(name string) interface{} {

	switch name {
	case "password":
		return o.Password
	case "username":
		return o.Username
	}

	return nil
}

// IssueLDAPAttributesMap represents the map of attribute for IssueLDAP.
var IssueLDAPAttributesMap = map[string]elemental.AttributeSpecification{
	"Password": {
		AllowedChoices: []string{},
		ConvertedName:  "Password",
		Description:    `The password for the user.`,
		Exposed:        true,
		Name:           "password",
		Required:       true,
		Type:           "string",
	},
	"Username": {
		AllowedChoices: []string{},
		ConvertedName:  "Username",
		Description:    `The LDAP username.`,
		Exposed:        true,
		Name:           "username",
		Required:       true,
		Type:           "string",
	},
}

// IssueLDAPLowerCaseAttributesMap represents the map of attribute for IssueLDAP.
var IssueLDAPLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"password": {
		AllowedChoices: []string{},
		ConvertedName:  "Password",
		Description:    `The password for the user.`,
		Exposed:        true,
		Name:           "password",
		Required:       true,
		Type:           "string",
	},
	"username": {
		AllowedChoices: []string{},
		ConvertedName:  "Username",
		Description:    `The LDAP username.`,
		Exposed:        true,
		Name:           "username",
		Required:       true,
		Type:           "string",
	},
}

type mongoAttributesIssueLDAP struct {
}
