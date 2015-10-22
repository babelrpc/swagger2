package swagger2

import (
	"errors"
	"mime"
	"net/mail"
	"net/url"
	"strings"
)

// ErrorList is a slice of errors from validation
type ErrorList []error

// String converts the errors into text
func (e ErrorList) String() string {
	text := ""
	for _, s := range e {
		if text != "" {
			text += "\n"
		}
		text += s.Error()
	}
	return text
}

// Indent converts the error list to text, using the provided indentation
func (e ErrorList) Indent(ind string) string {
	text := ""
	for _, s := range e {
		if text != "" {
			text += "\n"
		}
		text += ind + s.Error()
	}
	return text
}

// isUrl returns true if the string is a properly formed url
func isUrl(surl string) bool {
	_, err := url.ParseRequestURI(surl)
	if err != nil {
		return false
	}
	return true
}

// isEmail returns true if the string is a properly formed email
func isEmail(email string) bool {
	m, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	if m.Name != "" {
		return false
	}
	return true
}

// isMime returns true if the string is a properly formed mime-type
func isMime(smime string) bool {
	_, _, err := mime.ParseMediaType(smime)
	if err != nil {
		return false
	}
	return true
}

// isHost returns true if the string is a properly formed host name with optional port
func isHost(host string) bool {
	if len(host) == 0 {
		return false
	}
	u, err := url.Parse("http://" + host)
	if err != nil {
		return false
	}
	if u.Opaque != "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return false
	}
	return true
}

// isPath returns true if the string is a properly formed path starting with /
func isPath(path string) bool {
	if len(path) == 0 || path[0] != '/' {
		return false
	}
	u, err := url.Parse(path)
	if err != nil {
		return false
	}
	if u.Scheme != "" || u.Opaque != "" || u.User != nil || u.Host != "" || u.RawQuery != "" || u.Fragment != "" {
		return false
	}
	return true
}

// Validate confirms that the node is set up correctly
func (s *Swagger) Validate() []error {
	errs := make([]error, 0)
	// Required. Specifies the Swagger Specification version being used. It can be used by the Swagger UI and other clients to interpret the API listing. The value MUST be "2.0".
	if strings.TrimSpace(s.Swagger) == "" {
		errs = append(errs, errors.New("swagger element is required"))
	}
	if s.Swagger != "2.0" {
		errs = append(errs, errors.New("swagger version must be 2.0"))
	}
	// Required. Provides metadata about the API. The metadata can be used by the clients if needed.
	errs = append(errs, s.Info.Validate()...)
	// The host (name or ip) serving the API. This MUST be the host only and does not include the scheme nor sub-paths. It MAY include a port. If the host is not included, the host serving the documentation is to be used (including the port). The host does not support path templating.
	if s.Host != "" {
		if !isHost(s.Host) {
			errs = append(errs, errors.New(s.Host+" is not a valid host name"))
		}
	}
	// The base path on which the API is served, which is relative to the host. If it is not included, the API is served directly under the host. The value MUST start with a leading slash (/). The basePath does not support path templating.
	if s.BasePath != "" {
		if !isPath(s.BasePath) {
			errs = append(errs, errors.New(s.BasePath+" is not a valid path"))
		}
	}
	// The transfer protocol of the API. Values MUST be from the list: "http", "https", "ws", "wss". If the schemes is not included, the default scheme to be used is the one used to access the specification.
	if s.Schemes != nil {
		for _, t := range s.Schemes {
			if t != "http" && t != "https" && t != "ws" && t != "wss" {
				errs = append(errs, errors.New("scheme not supported: "+t))
			}
		}
	}
	// A list of MIME types the APIs can consume. This is global to all APIs but can be overridden on specific API calls. Value MUST be as described under Mime Types.
	if s.Consumes != nil {
		for _, t := range s.Consumes {
			if !isMime(t) {
				errs = append(errs, errors.New("consumes: "+t+" is not a valid mime type"))
			}
		}
	}
	// A list of MIME types the APIs can produce. This is global to all APIs but can be overridden on specific API calls. Value MUST be as described under Mime Types.
	if s.Produces != nil {
		for _, t := range s.Produces {
			if !isMime(t) {
				errs = append(errs, errors.New("produces: "+t+" is not a valid mime type"))
			}
		}
	}
	// Required. The available paths and operations for the API.
	if s.Paths == nil || len(s.Paths) == 0 {
		errs = append(errs, errors.New("paths are required"))
	}
	if s.Paths != nil {
		for _, t := range s.Paths {
			errs = append(errs, t.Validate()...)
		}
	}
	// An object to hold data types produced and consumed by operations.
	if s.Definitions != nil {
		for _, t := range s.Definitions {
			errs = append(errs, t.Validate()...)
		}
	}
	// An object to hold parameters that can be used across operations. This property does not define global parameters for all operations.
	if s.Parameters != nil {
		for _, t := range s.Parameters {
			errs = append(errs, t.Validate()...)
		}
	}
	// An object to hold responses that can be used across operations. This property does not define global responses for all operations.
	if s.Responses != nil {
		for _, t := range s.Responses {
			errs = append(errs, t.Validate()...)
		}
	}
	// Security scheme definitions that can be used across the specification.
	if s.SecurityDefinitions != nil {
		for n, t := range s.SecurityDefinitions {
			if strings.TrimSpace(n) == "" {
				errs = append(errs, errors.New("security defintions must be named"))
			}
			errs = append(errs, t.Validate()...)
		}
	}
	// A declaration of which security schemes are applied for the API as a whole. The list of values describes alternative security schemes that can be used (that is, there is a logical OR between the security requirements). Individual operations can override this definition.
	// s.Security - figure out what to validate later
	// A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared may be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
	if s.Tags != nil {
		for _, t := range s.Tags {
			errs = append(errs, t.Validate()...)
		}
	}
	// Additional external documentation.
	if s.ExternalDocs != nil {
		errs = append(errs, s.ExternalDocs.Validate()...)
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Info) Validate() []error {
	errs := make([]error, 0)
	// Required. The title of the application.
	if strings.TrimSpace(s.Title) == "" {
		errs = append(errs, errors.New("info.title is required"))
	}
	// A short description of the application. GFM syntax can be used for rich text representation.
	// s.Description - not required
	// The Terms of Service for the API.
	// s.TermsOfService - not required
	// The contact information for the exposed API.
	if s.Contact != nil {
		errs = append(errs, s.Contact.Validate()...)
	}
	// The license information for the exposed API.
	if s.License != nil {
		errs = append(errs, s.License.Validate()...)
	}
	// Required Provides the version of the application API (not to be confused by the specification version).
	if strings.TrimSpace(s.Version) == "" {
		errs = append(errs, errors.New("info.version is required"))
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Contact) Validate() []error {
	errs := make([]error, 0)
	// The identifying name of the contact person/organization.
	// s.Name - not required
	// The URL pointing to the contact information. MUST be in the format of a URL.
	if s.Url != "" && !isUrl(s.Url) {
		errs = append(errs, errors.New("info.contact.url: "+s.Url+" is not a valid URL"))
	}
	// The email address of the contact person/organization. MUST be in the format of an email address.
	if s.Email != "" && !isEmail(s.Email) {
		errs = append(errs, errors.New("info.contact.email: "+s.Email+" is not a valid email address"))
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *License) Validate() []error {
	errs := make([]error, 0)
	// Required. The license name used for the API.
	if s.Url != "" && strings.TrimSpace(s.Name) == "" {
		errs = append(errs, errors.New("info.license.name is required"))
	}
	// A URL to the license used for the API. MUST be in the format of a URL.
	if s.Url != "" && !isUrl(s.Url) {
		errs = append(errs, errors.New("info.license.url: "+s.Url+" is not a valid url"))
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *PathItem) Validate() []error {
	errs := make([]error, 0)
	// Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. If there are conflicts between the referenced definition and this Path Item's definition, the behavior is undefined.
	if strings.TrimSpace(s.Ref) != "" && (s.Get != nil || s.Put != nil || s.Post != nil || s.Delete != nil || s.Options != nil || s.Head != nil || s.Patch != nil || s.Parameters != nil || len(s.Parameters) > 0) {
		errs = append(errs, errors.New("paths.pathitem.ref specificed but other elements are present too"))
	}
	// A definition of a GET operation on this path.
	if s.Get != nil {
		errs = append(errs, s.Get.Validate()...)
	}
	// A definition of a PUT operation on this path.
	if s.Put != nil {
		errs = append(errs, s.Put.Validate()...)
	}
	// A definition of a POST operation on this path.
	if s.Post != nil {
		errs = append(errs, s.Post.Validate()...)
	}
	// A definition of a DELETE operation on this path.
	if s.Delete != nil {
		errs = append(errs, s.Delete.Validate()...)
	}
	// A definition of a OPTIONS operation on this path.
	if s.Options != nil {
		errs = append(errs, s.Options.Validate()...)
	}
	// A definition of a HEAD operation on this path.
	if s.Head != nil {
		errs = append(errs, s.Head.Validate()...)
	}
	// A definition of a PATCH operation on this path.
	if s.Patch != nil {
		errs = append(errs, s.Patch.Validate()...)
	}
	// A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the Swagger Object's parameters. There can be one "body" parameter at most.
	if s.Parameters != nil {
		for _, t := range s.Parameters {
			errs = append(errs, t.Validate()...)
		}
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Operation) Validate() []error {
	errs := make([]error, 0)
	// A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	// s.Tags - nothing to validate
	// A short summary of what the operation does. For maximum readability in the swagger-ui, this field SHOULD be less than 120 characters.
	// s.Summary - not required
	// A verbose explanation of the operation behavior. GFM syntax can be used for rich text representation.
	// s.Description - not required
	// Additional external documentation for this operation.
	if s.ExternalDocs != nil {
		errs = append(errs, s.ExternalDocs.Validate()...)
	}
	// A friendly name for the operation. The id MUST be unique among all operations described in the API. Tools and libraries MAY use the operation id to uniquely identify an operation.
	// s.OperationId - not required, must be unique
	// A list of MIME types the operation can consume. This overrides the consumes definition at the Swagger Object. An empty value MAY be used to clear the global definition. Value MUST be as described under Mime Types.
	if s.Consumes != nil {
		for _, t := range s.Consumes {
			if !isMime(t) {
				errs = append(errs, errors.New("operation.consumes: "+t+" is not a valid mime type"))
			}
		}
	}
	// A list of MIME types the operation can produce. This overrides the produces definition at the Swagger Object. An empty value MAY be used to clear the global definition. Value MUST be as described under Mime Types.
	if s.Produces != nil {
		for _, t := range s.Produces {
			if !isMime(t) {
				errs = append(errs, errors.New("operation.produces: "+t+" is not a valid mime type"))
			}
		}
	}
	// A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it, but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the Swagger Object's parameters. There can be one "body" parameter at most.
	if s.Parameters != nil {
		for _, t := range s.Parameters {
			errs = append(errs, t.Validate()...)
		}
	}
	// Required. The list of possible responses as they are returned from executing this operation.
	if s.Responses == nil || len(s.Responses) == 0 {
		errs = append(errs, errors.New("operation.responses is required"))
	}
	if s.Responses != nil {
		for _, t := range s.Responses {
			errs = append(errs, t.Validate()...)
		}
	}
	// The transfer protocol for the operation. Values MUST be from the list: "http", "https", "ws", "wss". The value overrides the Swagger Object schemes definition.
	if s.Schemes != nil {
		for _, t := range s.Schemes {
			if t != "http" && t != "https" && t != "ws" && t != "wss" {
				errs = append(errs, errors.New("operation scheme not supported: "+t))
			}
		}
	}
	// Declares this operation to be deprecated. Usage of the declared operation should be refrained. Default value is false.
	// s.Deprecated - nothing to validate
	// A declaration of which security schemes are applied for this operation. The list of values describes alternative security schemes that can be used (that is, there is a logical OR between the security requirements). This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
	// s.Security - figure out what to validate later
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Documentation) Validate() []error {
	errs := make([]error, 0)
	// A short description of the target documentation. GFM syntax can be used for rich text representation.
	// s. Description - not required
	// Required. The URL for the target documentation. Value MUST be in the format of a URL.
	if s.Url == "" {
		errs = append(errs, errors.New("externalDocs: url is required"))
	} else if !isUrl(s.Url) {
		errs = append(errs, errors.New("externalDocs: "+s.Url+" is not a valid url"))
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Parameter) Validate() []error {
	errs := make([]error, 0)
	// Required. The name of the parameter. Parameter names are case sensitive. If in is "path", the name field MUST correspond to the associated path segment from the path field in the Paths Object. See Path Templating for further information. For all other cases, the name corresponds to the parameter name used based on the in property.
	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, errors.New("parameter: name is required"))
	}
	// Required. The location of the parameter. Possible values are "query", "header", "path", "formData" or "body".
	if strings.TrimSpace(s.In) == "" {
		errs = append(errs, errors.New("parameter: in is required"))
	} else if s.In != "query" && s.In != "header" && s.In != "path" && s.In != "formData" && s.In != "body" {
		errs = append(errs, errors.New("parameter: "+s.In+"is not a valid value for in"))
	}
	// A brief description of the parameter. This could contain examples of use. GFM syntax can be used for rich text representation.
	// s.Description - not required
	// Determines whether this parameter is mandatory. If the parameter is in "path", this property is required and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	// s.Required - not required
	// (for in=body) Required. The schema defining the type used for the body parameter.
	if s.Schema != nil {
		if s.In != "body" {
			errs = append(errs, errors.New("parameter: in must be \"body\" when using a schema"))
		}
		errs = append(errs, s.Schema.Validate()...)
	}
	// Other fields
	if s.In != "body" {
		errs = append(errs, s.ItemsDef.Validate()...)
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *ItemsDef) Validate() []error {
	errs := make([]error, 0)
	// NOTE: Not going to work hard here, the rules for JSON schema are complex
	// Required. The reference string.
	// s.Ref
	// Required. The type of the parameter. Since the parameter is not located at the request body, it is limited to simple types (that is, not an object). The value MUST be one of "string", "number", "integer", "boolean", "array" or "file". If type is "file", the consumes MUST be either "multipart/form-data" or " application/x-www-form-urlencoded" and the parameter MUST be in "formData".
	// s.Type
	// The extending format for the previously mentioned type. See Data Type Formats for further details.
	// s.Format
	// Required if type is "array". Describes the type of items in the array.
	// s.Items
	if s.Items != nil {
		errs = append(errs, s.Items.Validate()...)
	}
	// Determines the format of the array if type array is used. Possible values are: csv - comma separated values foo,bar. ssv - space separated values foo bar. tsv - tab separated values foo\tbar. pipes - pipe separated values foo|bar. multi - corresponds to multiple parameter instances instead of multiple values for a single instance foo=bar&foo=baz. This is valid only for parameters in "query" or "formData". Default value is csv.
	// s.CollectionFormat
	// Sets a default value to the parameter. The type of the value depends on the defined type. See http://json-schema.org/latest/json-schema-validation.html#anchor101.
	// s.Default
	// See http://json-schema.org/latest/json-schema-validation.html#anchor17.
	// s.Maximum
	// See http://json-schema.org/latest/json-schema-validation.html#anchor17.
	// s.ExclusiveMaximum
	// See http://json-schema.org/latest/json-schema-validation.html#anchor21.
	// s.Minimum
	// See http://json-schema.org/latest/json-schema-validation.html#anchor21.
	// s.ExclusiveMinimum
	// See http://json-schema.org/latest/json-schema-validation.html#anchor26.
	// s.MaxLength
	// See http://json-schema.org/latest/json-schema-validation.html#anchor29.
	// s.MinLength
	// See http://json-schema.org/latest/json-schema-validation.html#anchor33.
	// s.Pattern
	// See http://json-schema.org/latest/json-schema-validation.html#anchor42.
	// s.MaxItems
	// See http://json-schema.org/latest/json-schema-validation.html#anchor45.
	// s.MinItems
	// See http://json-schema.org/latest/json-schema-validation.html#anchor49.
	// s.UniqueItems
	// See http://json-schema.org/latest/json-schema-validation.html#anchor76.
	// s.Enum
	// See http://json-schema.org/latest/json-schema-validation.html#anchor14.
	// s.MultipleOf
	// Used for maps
	if s.AdditionalProperties != nil {
		errs = append(errs, s.AdditionalProperties.Validate()...)
	}
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Response) Validate() []error {
	errs := make([]error, 0)
	// Required. A short description of the response. GFM syntax can be used for rich text representation.
	if strings.TrimSpace(s.Description) == "" {
		errs = append(errs, errors.New("response.description is required"))
	}
	// A definition of the response structure. It can be a primitive, an array or an object. If this field does not exist, it means no content is returned as part of the response. As an extension to the Schema Object, its root type value may also be "file". This SHOULD be accompanied by a relevant produces mime-type.
	if s.Schema != nil {
		errs = append(errs, s.Schema.Validate()...)
	}
	// A list of headers that are sent with the response.
	if s.Headers != nil {
		for _, t := range s.Headers {
			errs = append(errs, t.Validate()...)
		}
	}
	// An example of the response message.
	// s.Example - could be anything
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Header) Validate() []error {
	errs := make([]error, 0)
	// A short description of the header.
	// s.Description - nothing to validate
	// Other fields
	errs = append(errs, s.ItemsDef.Validate()...)
	return errs
}

// Validate confirms that the node is set up correctly
func (s *Tag) Validate() []error {
	errs := make([]error, 0)
	// Required. The name of the tag.
	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, errors.New("tag.name is required"))
	}
	// A short description for the tag. GFM syntax can be used for rich text representation.
	// s.Description - not required
	// Additional external documentation for this tag.
	if s.ExternalDocs != nil {
		errs = append(errs, s.ExternalDocs.Validate()...)
	}

	return errs
}

// Validate confirms that the node is set up correctly
func (s *Schema) Validate() []error {
	errs := make([]error, 0)
	// NOTE: Not sure how to validate some of these
	// s.Title
	// s.Description
	// s.MaxProperties
	// s.MinProperties
	// s.Required
	// s.AllOf
	// s.Properties
	if s.Properties != nil {
		for _, t := range s.Properties {
			errs = append(errs, t.Validate()...)
		}
	}
	// Adds support for polymorphism. The discriminator is the schema property name that is used to differentiate between other schema that inherit this schema. The property name used MUST be defined at this schema and it MUST be in the required property list. When used, the value MUST be the name of this schema or any schema that inherits it.
	// s.Discriminator
	// Relevant only for Schema "properties" definitions. Declares the property as "read only". This means that it MAY be sent as part of a response but MUST NOT be sent as part of the request. Properties marked as readOnly being true SHOULD NOT be in the required list of the defined schema. Default value is false.
	// s.ReadOnly
	// This MAY be used only on properties schemas. It has no effect on root schemas. Adds Additional metadata to describe the XML representation format of this property.
	if s.Xml != nil {
		errs = append(errs, s.Xml.Validate()...)
	}
	// Additional external documentation for this schema.
	if s.ExternalDocs != nil {
		errs = append(errs, s.ExternalDocs.Validate()...)
	}
	// A free-form property to include a an example of an instance for this schema.
	// s.Example

	// Other fields
	errs = append(errs, s.ItemsDef.Validate()...)

	return errs
}

// Validate confirms that the node is set up correctly
func (s *Xml) Validate() []error {
	errs := make([]error, 0)
	// NOTE: not clear the best way to validate here
	// Replaces the name of the element/attribute used for the described schema property. When defined within the Items Object (items), it will affect the name of the individual XML elements within the list. When defined alongside type being array (outside the items), it will affect the wrapping element and only if wrapped is true. If wrapped is false, it will be ignored.
	// s.Name
	// The URL of the namespace definition. Value SHOULD be in the form of a URL.
	// s.Namespace
	// The prefix to be used for the name.
	// s.Prefix
	// Declares whether the property definition translates to an attribute instead of an element. Default value is false.
	// s.Attribute
	// MAY be used only for an array definition. Signifies whether the array is wrapped (for example, <books><book/><book/></books>) or unwrapped (<book/><book/>). Default value is false. The definition takes effect only when defined alongside type being array (outside the items).
	// s.Wrapped
	return errs
}

// Validate confirms that the node is set up correctly
func (s *SecurityDefinition) Validate() []error {
	errs := make([]error, 0)
	// Required. The type of the security scheme. Valid values are "basic", "apiKey" or "oauth2".
	if s.Type != "basic" && s.Type != "apiKey" && s.Type != "oauth2" {
		errs = append(errs, errors.New("securityDefinition.type must be \"basic\", \"apiKey\", or \"oauth2\""))
	}
	// A short description for security scheme.
	// s.Description - not required
	// Required. The name of the header or query parameter to be used.
	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, errors.New("securityDefinition.name is required"))
	}
	// Required The location of the API key. Valid values are "query" or "header".
	if s.In != "query" && s.In != "header" {
		errs = append(errs, errors.New("securityDefinition.in must be \"query\" or \"header\""))
	}
	// Required. The flow used by the OAuth2 security scheme. Valid values are "implicit", "password", "application" or "accessCode".
	if s.Flow != "implicit" && s.Flow != "password" && s.Flow != "application" && s.Flow != "accessCode" {
		errs = append(errs, errors.New("securityDefinition.flow must be \"implicit\", \"password\", \"application\", or \"accessCode\""))
	}
	// Required. The authorization URL to be used for this flow. This SHOULD be in the form of a URL.
	if s.AuthorizationUrl == "" || !isUrl(s.AuthorizationUrl) {
		errs = append(errs, errors.New("securityDefinition.authorizationUrl: "+s.AuthorizationUrl+" is not a valid URL"))
	}
	// Required. The token URL to be used for this flow. This SHOULD be in the form of a URL.
	if s.TokenUrl == "" || !isUrl(s.TokenUrl) {
		errs = append(errs, errors.New("securityDefinition.tokenUrl: "+s.TokenUrl+" is not a valid URL"))
	}
	// Required. The available scopes for the OAuth2 security scheme.
	// s.Scopes - see how to validate this later
	return errs
}
