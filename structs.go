package swagger2

import ()

// This is the root document object for the API specification. It combines what previously was the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
type Swagger struct {
	Swagger             string              `yaml:"swagger" json:"swagger"`                                             // Required. Specifies the Swagger Specification version being used. It can be used by the Swagger UI and other clients to interpret the API listing. The value MUST be "2.0".
	Info                Info                `yaml:"info" json:"info"`                                                   // Required. Provides metadata about the API. The metadata can be used by the clients if needed.
	Host                string              `yaml:"host,omitempty" json:"host,omitempty"`                               // The host (name or ip) serving the API. This MUST be the host only and does not include the scheme nor sub-paths. It MAY include a port. If the host is not included, the host serving the documentation is to be used (including the port). The host does not support path templating.
	BasePath            string              `yaml:"basePath,omitempty" json:"basePath,omitempty"`                       // The base path on which the API is served, which is relative to the host. If it is not included, the API is served directly under the host. The value MUST start with a leading slash (/). The basePath does not support path templating.
	Schemes             []string            `yaml:"schemes,omitempty" json:"schemes,omitempty"`                         // The transfer protocol of the API. Values MUST be from the list: "http", "https", "ws", "wss". If the schemes is not included, the default scheme to be used is the one used to access the specification.
	Consumes            []string            `yaml:"consumes,omitempty" json:"consumes,omitempty"`                       // A list of MIME types the APIs can consume. This is global to all APIs but can be overridden on specific API calls. Value MUST be as described under Mime Types.
	Produces            []string            `yaml:"produces,omitempty" json:"produces,omitempty"`                       // A list of MIME types the APIs can produce. This is global to all APIs but can be overridden on specific API calls. Value MUST be as described under Mime Types.
	Paths               Paths               `yaml:"paths" json:"paths"`                                                 // Required. The available paths and operations for the API.
	Definitions         Definitions         `yaml:"definitions,omitempty" json:"definitions,omitempty"`                 // An object to hold data types produced and consumed by operations.
	Parameters          Parameters          `yaml:"parameters,omitempty" json:"parameters,omitempty"`                   // An object to hold parameters that can be used across operations. This property does not define global parameters for all operations.
	Responses           Responses           `yaml:"responses,omitempty" json:"responses,omitempty"`                     // An object to hold responses that can be used across operations. This property does not define global responses for all operations.
	SecurityDefinitions SecurityDefinitions `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"` // Security scheme definitions that can be used across the specification.
	Security            []Security          `yaml:"security,omitempty" json:"security,omitempty"`                       // A declaration of which security schemes are applied for the API as a whole. The list of values describes alternative security schemes that can be used (that is, there is a logical OR between the security requirements). Individual operations can override this definition.
	Tags                []Tag               `yaml:"tags,omitempty" json:"tags,omitempty"`                               // A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared may be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
	ExternalDocs        *Documentation      `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`               // Additional external documentation.
}

// The object provides metadata about the API. The metadata can be used by the clients if needed, and can be presented in the Swagger-UI for convenience.
type Info struct {
	Title          string   `yaml:"title" json:"title"`                                       // Required. The title of the application.
	Description    string   `yaml:"description,omitempty" json:"description,omitempty"`       // A short description of the application. GFM syntax can be used for rich text representation.
	TermsOfService string   `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"` // The Terms of Service for the API.
	Contact        *Contact `yaml:"contact,omitempty" json:"contact,omitempty"`               // The contact information for the exposed API.
	License        *License `yaml:"license,omitempty" json:"license,omitempty"`               // The license information for the exposed API.
	Version        string   `yaml:"version" json:"version"`                                   // Required Provides the version of the application API (not to be confused by the specification version).

	// Vendor extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `yaml:"name,omitempty" json:"name,omitempty"`   // The identifying name of the contact person/organization.
	Url   string `yaml:"url,omitempty" json:"url,omitempty"`     // The URL pointing to the contact information. MUST be in the format of a URL.
	Email string `yaml:"email,omitempty" json:"email,omitempty"` // The email address of the contact person/organization. MUST be in the format of an email address.
}

// License information for the exposed API.
type License struct {
	Name string `yaml:"name" json:"name"`                   // Required. The license name used for the API.
	Url  string `yaml:"url,omitempty" json:"url,omitempty"` // A URL to the license used for the API. MUST be in the format of a URL.
}

// Holds the relative paths to the individual endpoints. The path is appended to the basePath in order to construct the full URL. The Paths may be empty, due to ACL constraints.
type Paths map[string]PathItem

// Note Path Extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.

// Describes the operations available on a single path. A Path Item may be empty, due to ACL constraints. The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available.
type PathItem struct {
	Ref        string      `yaml:"$ref,omitempty" json:"$ref,omitempty"`             // Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. If there are conflicts between the referenced definition and this Path Item's definition, the behavior is undefined.
	Get        *Operation  `yaml:"get,omitempty" json:"get,omitempty"`               // A definition of a GET operation on this path.
	Put        *Operation  `yaml:"put,omitempty" json:"put,omitempty"`               // A definition of a PUT operation on this path.
	Post       *Operation  `yaml:"post,omitempty" json:"post,omitempty"`             // A definition of a POST operation on this path.
	Delete     *Operation  `yaml:"delete,omitempty" json:"delete,omitempty"`         // A definition of a DELETE operation on this path.
	Options    *Operation  `yaml:"options,omitempty" json:"options,omitempty"`       // A definition of a OPTIONS operation on this path.
	Head       *Operation  `yaml:"head,omitempty" json:"head,omitempty"`             // A definition of a HEAD operation on this path.
	Patch      *Operation  `yaml:"patch,omitempty" json:"patch,omitempty"`           // A definition of a PATCH operation on this path.
	Parameters []Parameter `yaml:"parameters,omitempty" json:"parameters,omitempty"` // A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the Swagger Object's parameters. There can be one "body" parameter at most.

	// Vendor extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.
}

// Describes a single API operation on a path.
type Operation struct {
	Tags         []string       `yaml:"tags,omitempty" json:"tags,omitempty"`                 // A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Summary      string         `yaml:"summary,omitempty" json:"summary,omitempty"`           // A short summary of what the operation does. For maximum readability in the swagger-ui, this field SHOULD be less than 120 characters.
	Description  string         `yaml:"description,omitempty" json:"description,omitempty"`   // A verbose explanation of the operation behavior. GFM syntax can be used for rich text representation.
	ExternalDocs *Documentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"` // Additional external documentation for this operation.
	OperationId  string         `yaml:"operationId,omitempty" json:"operationId,omitempty"`   // A friendly name for the operation. The id MUST be unique among all operations described in the API. Tools and libraries MAY use the operation id to uniquely identify an operation.
	Consumes     []string       `yaml:"consumes,omitempty" json:"consumes,omitempty"`         // A list of MIME types the operation can consume. This overrides the consumes definition at the Swagger Object. An empty value MAY be used to clear the global definition. Value MUST be as described under Mime Types.
	Produces     []string       `yaml:"produces,omitempty" json:"produces,omitempty"`         // A list of MIME types the operation can produce. This overrides the produces definition at the Swagger Object. An empty value MAY be used to clear the global definition. Value MUST be as described under Mime Types.
	Parameters   []Parameter    `yaml:"parameters,omitempty" json:"parameters,omitempty"`     // A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it, but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the Swagger Object's parameters. There can be one "body" parameter at most.
	Responses    Responses      `yaml:"responses,omitempty" json:"responses,omitempty"`       // Required. The list of possible responses as they are returned from executing this operation.
	Schemes      []string       `yaml:"schemes,omitempty" json:"schemes,omitempty"`           // The transfer protocol for the operation. Values MUST be from the list: "http", "https", "ws", "wss". The value overrides the Swagger Object schemes definition.
	Deprecated   bool           `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`     // Declares this operation to be deprecated. Usage of the declared operation should be refrained. Default value is false.
	Security     []Security     `yaml:"security,omitempty" json:"security,omitempty"`         // A declaration of which security schemes are applied for this operation. The list of values describes alternative security schemes that can be used (that is, there is a logical OR between the security requirements). This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.

	// Vendor extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.
}

// Allows referencing an external resource for extended documentation.
type Documentation struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"` // A short description of the target documentation. GFM syntax can be used for rich text representation.
	Url         string `yaml:"url" json:"url"`                                     // Required. The URL for the target documentation. Value MUST be in the format of a URL.
}

// Describes a single operation parameter.
//
// A unique parameter is defined by a combination of a name and location.
//
// There are five possible parameter types.
//
// Path - Used together with Path Templating, where the parameter value is actually part of the operation's URL. This does not include the host or base path of the API. For example, in /items/{itemId}, the path parameter is itemId.
// Query - Parameters that are appended to the URL. For example, in /items?id=###, the query parameter is id.
// Header - Custom headers that are expected as part of the request.
// Body - The payload that's appended to the HTTP request. Since there can only be one payload, there can only be one body parameter. The name of the body parameter has no effect on the parameter itself and is used for documentation purposes only. Since Form parameters are also in the payload, body and form parameters cannot exist together for the same operation.
// Form - Used to describe the payload of an HTTP request when either application/x-www-form-urlencoded or multipart/form-data are used as the content type of the request (in Swagger's definition, the consumes property of an operation). This is the only parameter type that can be used to send files, thus supporting the file type. Since form parameters are sent in the payload, they cannot be declared together with a body parameter for the same operation. Form parameters have a different format based on the content-type used (for further details, consult http://www.w3.org/TR/html401/interact/forms.html#h-17.13.4):
//   application/x-www-form-urlencoded - Similar to the format of Query parameters but as a payload. For example, foo=1&bar=swagger - both foo and bar are form parameters. This is normally used for simple parameters that are being transferred.
//   multipart/form-data - each parameter takes a section in the payload with an internal header. For example, for the header Content-Disposition: form-data; name="submit-name" the name of the parameter is submit-name. This type of form parameters is more commonly used for file transfers.
type Parameter struct {
	Name        string `yaml:"name" json:"name"`                                   // Required. The name of the parameter. Parameter names are case sensitive. If in is "path", the name field MUST correspond to the associated path segment from the path field in the Paths Object. See Path Templating for further information. For all other cases, the name corresponds to the parameter name used based on the in property.
	In          string `yaml:"in" json:"in"`                                       // Required. The location of the parameter. Possible values are "query", "header", "path", "formData" or "body".
	Description string `yaml:"description,omitempty" json:"description,omitempty"` // A brief description of the parameter. This could contain examples of use. GFM syntax can be used for rich text representation.
	Required    *bool  `yaml:"required,omitempty" json:"required,omitempty"`       // Determines whether this parameter is mandatory. If the parameter is in "path", this property is required and its value MUST be true. Otherwise, the property MAY be included and its default value is false.

	// If in is "body":
	Schema *Schema `yaml:"schema,omitempty" json:"schema,omitempty"` // Required. The schema defining the type used for the body parameter.

	// If in is any value other than "body":
	ItemsDef `yaml:",omitempty,inline"`

	// Vendor extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.
}

/*
Data Type Formats

Common Name   type      format      Comments
-----------   -------   ---------   ---------------------------------
integer       integer   int32       signed 32 bits
long          integer   int64       signed 64 bits
float         number    float
double        number    double
string        string
byte          string    byte
boolean       boolean
date          string    date        As defined by full-date - RFC3339
dateTime      string    date-time   As defined by date-time - RFC3339
*/

// An limited subset of JSON-Schema's items object. It is used by parameter definitions that are not located in "body".
type ItemsDef struct {
	Ref                  string        `yaml:"$ref,omitempty" json:"$ref,omitempty"`                                 // Required. The reference string.
	Type                 string        `yaml:"type,omitempty" json:"type,omitempty"`                                 // Required. The type of the parameter. Since the parameter is not located at the request body, it is limited to simple types (that is, not an object). The value MUST be one of "string", "number", "integer", "boolean", "array" or "file". If type is "file", the consumes MUST be either "multipart/form-data" or " application/x-www-form-urlencoded" and the parameter MUST be in "formData".
	Format               string        `yaml:"format,omitempty" json:"format,omitempty"`                             // The extending format for the previously mentioned type. See Data Type Formats for further details.
	Items                *ItemsDef     `yaml:"items,omitempty" json:"items,omitempty"`                               // Required if type is "array". Describes the type of items in the array.
	CollectionFormat     string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`         // Determines the format of the array if type array is used. Possible values are: csv - comma separated values foo,bar. ssv - space separated values foo bar. tsv - tab separated values foo\tbar. pipes - pipe separated values foo|bar. multi - corresponds to multiple parameter instances instead of multiple values for a single instance foo=bar&foo=baz. This is valid only for parameters in "query" or "formData". Default value is csv.
	Default              interface{}   `yaml:"default,omitempty" json:"default,omitempty"`                           // Sets a default value to the parameter. The type of the value depends on the defined type. See http://json-schema.org/latest/json-schema-validation.html#anchor101.
	Maximum              *float64      `yaml:"maximum,omitempty" json:"maximum,omitempty"`                           // See http://json-schema.org/latest/json-schema-validation.html#anchor17.
	ExclusiveMaximum     *bool         `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`         // See http://json-schema.org/latest/json-schema-validation.html#anchor17.
	Minimum              *float64      `yaml:"minimum,omitempty" json:"minimum,omitempty"`                           // See http://json-schema.org/latest/json-schema-validation.html#anchor21.
	ExclusiveMinimum     *bool         `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`         // See http://json-schema.org/latest/json-schema-validation.html#anchor21.
	MaxLength            *int          `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`                       // See http://json-schema.org/latest/json-schema-validation.html#anchor26.
	MinLength            *int          `yaml:"minLength,omitempty" json:"minLength,omitempty"`                       // See http://json-schema.org/latest/json-schema-validation.html#anchor29.
	Pattern              *string       `yaml:"pattern,omitempty" json:"pattern,omitempty"`                           // See http://json-schema.org/latest/json-schema-validation.html#anchor33.
	MaxItems             *int          `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`                         // See http://json-schema.org/latest/json-schema-validation.html#anchor42.
	MinItems             *int          `yaml:"minItems,omitempty" json:"minItems,omitempty"`                         // See http://json-schema.org/latest/json-schema-validation.html#anchor45.
	UniqueItems          *bool         `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`                   // See http://json-schema.org/latest/json-schema-validation.html#anchor49.
	Enum                 []interface{} `yaml:"enum,omitempty" json:"enum,omitempty"`                                 // See http://json-schema.org/latest/json-schema-validation.html#anchor76.
	MultipleOf           *float64      `yaml:"multipleOf,omitempty" json:"multipleOf,omitempty"`                     // See http://json-schema.org/latest/json-schema-validation.html#anchor14.
	AdditionalProperties *ItemsDef     `yaml:"additionalProperties,omitempty" json:"additionalProperties,omitempty"` // Used for maps
}

// A container for the expected responses of an operation. The container maps a HTTP response code to the expected response. It is not expected from the documentation to necessarily cover all possible HTTP response codes, since they may not be known in advance. However, it is expected from the documentation to cover a successful operation response and any known errors.
// The default can be used a default response object for all HTTP codes that are not covered individually by the specification.
// The Responses Object MUST contain at least one response code, and it SHOULD be the response for a successful operation call.
type Responses map[string]Response

// "default": The documentation of responses other than the ones declared for specific HTTP response codes. It can be used to cover undeclared responses. Reference Object can be used to link to a response that is defined at the Swagger Object's responses section.
// {HTTP code}: Any HTTP status code can be used as the property name (one property per HTTP status code). Describes the expected response for that HTTP status code. Reference Object can be used to link to a response that is defined at the Swagger Object's responses section.
// Extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.

// Describes a single response from an API Operation.
type Response struct {
	Description string  `yaml:"description" json:"description"`               // Required. A short description of the response. GFM syntax can be used for rich text representation.
	Schema      *Schema `yaml:"schema,omitempty" json:"schema,omitempty"`     // A definition of the response structure. It can be a primitive, an array or an object. If this field does not exist, it means no content is returned as part of the response. As an extension to the Schema Object, its root type value may also be "file". This SHOULD be accompanied by a relevant produces mime-type.
	Headers     Headers `yaml:"headers,omitempty" json:"headers,omitempty"`   // A list of headers that are sent with the response.
	Examples    Example `yaml:"examples,omitempty" json:"examples,omitempty"` // An example of the response message.
}

// Lists the headers that can be sent as part of a response.
type Headers map[string]Header

// Header object
type Header struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"` // A short description of the header.
	ItemsDef    `yaml:",omitempty,inline"`
}

// Allows sharing examples for operation responses. keys must be a mime type
type Example map[string]interface{}

// Allows adding meta data to a single tag that is used by the Operation Object. It is not mandatory to have a Tag Object per tag used there.
type Tag struct {
	Name         string         `yaml:"name" json:"name"`                                     // Required. The name of the tag.
	Description  string         `yaml:"description,omitempty" json:"description,omitempty"`   // A short description for the tag. GFM syntax can be used for rich text representation.
	ExternalDocs *Documentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"` // Additional external documentation for this tag.

	// Vendor extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.
}

// JSON schema
type Schema struct {
	ItemsDef      `yaml:",omitempty,inline"`
	Title         string            `yaml:"title,omitempty" json:"title,omitempty"`
	Description   string            `yaml:"description,omitempty" json:"description,omitempty"`
	MaxProperties *int              `yaml:"maxProperties,omitempty" json:"maxProperties,omitempty"`
	MinProperties *int              `yaml:"minProperties,omitempty" json:"minProperties,omitempty"`
	Required      []string          `yaml:"required,omitempty" json:"required,omitempty"`
	AllOf         []Schema          `yaml:"allOf,omitempty" json:"allOf,omitempty"`
	Properties    map[string]Schema `yaml:"properties,omitempty" json:"properties,omitempty"`
	Discriminator string            `yaml:"discriminator,omitempty" json:"discriminator,omitempty"` // Adds support for polymorphism. The discriminator is the schema property name that is used to differentiate between other schema that inherit this schema. The property name used MUST be defined at this schema and it MUST be in the required property list. When used, the value MUST be the name of this schema or any schema that inherits it.
	ReadOnly      *bool             `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`           // Relevant only for Schema "properties" definitions. Declares the property as "read only". This means that it MAY be sent as part of a response but MUST NOT be sent as part of the request. Properties marked as readOnly being true SHOULD NOT be in the required list of the defined schema. Default value is false.
	Xml           *Xml              `yaml:"xml,omitempty" json:"xml,omitempty"`                     // This MAY be used only on properties schemas. It has no effect on root schemas. Adds Additional metadata to describe the XML representation format of this property.
	ExternalDocs  *Documentation    `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`   // Additional external documentation for this schema.
	Example       interface{}       `yaml:"example,omitempty" json:"example,omitempty"`             // A free-form property to include a an example of an instance for this schema.
}

// The xml property allows extra definitions when translating the JSON definition to XML. The XML Object contains additional information about the available options.
type Xml struct {
	Name      string `yaml:"name,omitempty" json:"name,omitempty"`           // Replaces the name of the element/attribute used for the described schema property. When defined within the Items Object (items), it will affect the name of the individual XML elements within the list. When defined alongside type being array (outside the items), it will affect the wrapping element and only if wrapped is true. If wrapped is false, it will be ignored.
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"` // The URL of the namespace definition. Value SHOULD be in the form of a URL.
	Prefix    string `yaml:"prefix,omitempty" json:"prefix,omitempty"`       // The prefix to be used for the name.
	Attribute *bool  `yaml:"attribute,omitempty" json:"attribute,omitempty"` // Declares whether the property definition translates to an attribute instead of an element. Default value is false.
	Wrapped   *bool  `yaml:"wrapped,omitempty" json:"wrapped,omitempty"`     // MAY be used only for an array definition. Signifies whether the array is wrapped (for example, <books><book/><book/></books>) or unwrapped (<book/><book/>). Default value is false. The definition takes effect only when defined alongside type being array (outside the items).
}

// An object to hold data types that can be consumed and produced by operations. These data types can be primitives, arrays or models.
type Definitions map[string]Schema

// An object to hold parameters to be reused across operations. Parameter definitions can be referenced to the ones defined here.
type Parameters map[string]Parameter

// A declaration of the security schemes available to be used in the specification. This does not enforce the security schemes on the operations and only serves to provide the relevant details for each scheme.
type SecurityDefinitions map[string]SecurityDefinition

// Allows the definition of a security scheme that can be used by the operations. Supported schemes are basic authentication, an API key (either as a header or as a query parameter) and OAuth2's common flows (implicit, password, application and access code).
type SecurityDefinition struct {
	Type             string `yaml:"type" json:"type"`                                   // Required. The type of the security scheme. Valid values are "basic", "apiKey" or "oauth2".
	Description      string `yaml:"description,omitempty" json:"description,omitempty"` // A short description for security scheme.
	Name             string `yaml:"name" json:"name"`                                   // Required. The name of the header or query parameter to be used.
	In               string `yaml:"in" json:"in"`                                       // Required The location of the API key. Valid values are "query" or "header".
	Flow             string `yaml:"flow" json:"flow"`                                   // Required. The flow used by the OAuth2 security scheme. Valid values are "implicit", "password", "application" or "accessCode".
	AuthorizationUrl string `yaml:"authorizationUrl" json:"authorizationUrl"`           // Required. The authorization URL to be used for this flow. This SHOULD be in the form of a URL.
	TokenUrl         string `yaml:"tokenUrl" json:"tokenUrl"`                           // Required. The token URL to be used for this flow. This SHOULD be in the form of a URL.
	Scopes           Scopes `yaml:"scopes" json:"scopes"`                               // Required. The available scopes for the OAuth2 security scheme.

	// Vendor extensions: Allows extensions to the Swagger Schema. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object. See Vendor Extensions for further details.
}

// Lists the available scopes for an OAuth2 security scheme.
type Scopes map[string]string

// Lists the required security schemes to execute this operation. The object can have multiple security schemes declared in it which are all required (that is, there is a logical AND between the schemes).
type Security map[string][]string
