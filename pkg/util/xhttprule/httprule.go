package xhttprule

import (
	"regexp"

	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// HasHTTPRules returns true when the given method descriptor is annotated with
// a google.api.http option.
func HasHTTPRules(m protoreflect.MethodDescriptor) bool {
	got := proto.GetExtension(m.Options(), apb.E_Http).(*apb.HttpRule)
	return got != nil
}

// GetHTTPRules returns a slice of HTTP rules for a given method descriptor.
//
// Note: This returns a slice -- it takes the google.api.http annotation,
// and then flattens the values in `additional_bindings`.
// This allows rule authors to simply range over all of the HTTP rules,
// since the common case is to want to apply the checks to all of them.
func GetHTTPRules(m protoreflect.MethodDescriptor) []*HTTPRule {
	rules := []*HTTPRule{}

	// Get the method options.
	opts := m.Options()

	// Get the "primary" rule (the direct google.api.http annotation).
	if x := proto.GetExtension(opts, apb.E_Http); x != nil {
		httpRule := x.(*apb.HttpRule)
		if parsedRule := parseRule(httpRule); parsedRule != nil {
			rules = append(rules, parsedRule)

			// Add any additional bindings and flatten them into `rules`.
			for _, binding := range httpRule.GetAdditionalBindings() {
				rules = append(rules, parseRule(binding))
			}
		}
	}

	// Done; return the rules.
	return rules
}

func parseRule(rule *apb.HttpRule) *HTTPRule {
	oneof := map[string]string{
		"GET":    rule.GetGet(),
		"POST":   rule.GetPost(),
		"PUT":    rule.GetPut(),
		"PATCH":  rule.GetPatch(),
		"DELETE": rule.GetDelete(),
	}
	if custom := rule.GetCustom(); custom != nil {
		oneof[custom.GetKind()] = custom.GetPath()
	}

	for method, uri := range oneof {
		if uri != "" {
			return &HTTPRule{
				Method:       method,
				URI:          uri,
				Body:         rule.GetBody(),
				ResponseBody: rule.GetResponseBody(),
			}
		}
	}
	return nil
}

// HTTPRule defines a parsed, easier-to-query equivalent to `apb.HttpRule`.
type HTTPRule struct {
	// The HTTP method. Guaranteed to be in all caps.
	// This is set to "CUSTOM" if the Custom property is set.
	Method string

	// The HTTP URI (the value corresponding to the selected HTTP method).
	URI string

	// The `body` value forwarded from the generated proto's HttpRule.
	Body string

	// The `response_body` value forwarded from the generated proto's HttpRule.
	ResponseBody string
}

// GetVariables returns the variable segments in a URI as a map.
//
// For a given variable, the key is the variable's field path. The value is the
// variable's template, which will match segment(s) of the URL.
//
// For more details on the path template syntax, see
// https://github.com/googleapis/googleapis/blob/6e1a5a066659794f26091674e3668229e7750052/google/api/http.proto#L224.
func (h *HTTPRule) GetVariables() map[string]string {
	vars := map[string]string{}

	// Replace the version template variable with "v".
	uri := VersionedSegment.ReplaceAllString(h.URI, "v")
	for _, match := range plainVar.FindAllStringSubmatch(uri, -1) {
		vars[match[1]] = "*"
	}
	for _, match := range varSegment.FindAllStringSubmatch(uri, -1) {
		vars[match[1]] = match[2]
	}
	return vars
}

// GetPlainURI returns the URI with variable segment information removed.
func (h *HTTPRule) GetPlainURI() string {
	return plainVar.ReplaceAllString(
		varSegment.ReplaceAllString(
			VersionedSegment.ReplaceAllString(h.URI, "v"),
			"$2"),
		"*")
}

var (
	plainVar         = regexp.MustCompile(`\{([^}=]+)\}`)
	varSegment       = regexp.MustCompile(`\{([^}=]+)=([^}]+)\}`)
	VersionedSegment = regexp.MustCompile(`\{\$api_version\}`)
)
