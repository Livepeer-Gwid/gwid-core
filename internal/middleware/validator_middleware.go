package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gwid.io/gwid-core/internal/utils"
)

var defaultValidationMessages = map[string]string{
	// Basic validations
	"required":         "{0} is required",
	"required_if":      "{0} is required when {1}",
	"required_unless":  "{0} is required unless {1}",
	"required_with":    "{0} is required when {1} is present",
	"required_without": "{0} is required when {1} is not present",

	// String validations
	"min":             "{0} must be at least {1} characters",
	"max":             "{0} must be at most {1} characters",
	"len":             "{0} must be exactly {1} characters",
	"alpha":           "{0} must contain only alphabetic characters",
	"alphanum":        "{0} must contain only alphanumeric characters",
	"alphaunicode":    "{0} must contain only unicode alphabetic characters",
	"alphanumunicode": "{0} must contain only unicode alphanumeric characters",
	"numeric":         "{0} must be a numeric value",
	"number":          "{0} must be a valid number",
	"hexadecimal":     "{0} must be a hexadecimal string",
	"hexcolor":        "{0} must be a valid HEX color",
	"lowercase":       "{0} must be in lowercase",
	"uppercase":       "{0} must be in uppercase",
	"uuid":            "{0} must be a valid UUID",
	"uuid3":           "{0} must be a valid UUID v3",
	"uuid4":           "{0} must be a valid UUID v4",
	"uuid5":           "{0} must be a valid UUID v5",

	// Format validations
	"email":     "{0} must be a valid email address",
	"url":       "{0} must be a valid URL",
	"uri":       "{0} must be a valid URI",
	"base64":    "{0} must be a valid Base64 string",
	"base64url": "{0} must be a valid Base64URL string",
	"json":      "{0} must be a valid JSON string",
	"jwt":       "{0} must be a valid JWT token",
	"hostname":  "{0} must be a valid hostname",
	"fqdn":      "{0} must be a fully qualified domain name",
	"ip":        "{0} must be a valid IP address",
	"ipv4":      "{0} must be a valid IPv4 address",
	"ipv6":      "{0} must be a valid IPv6 address",
	"datetime":  "{0} must be a valid datetime in format {1}",
	"timezone":  "{0} must be a valid timezone",

	// Comparison validations
	"eq":       "{0} must be equal to {1}",
	"ne":       "{0} must not be equal to {1}",
	"gt":       "{0} must be greater than {1}",
	"gte":      "{0} must be greater than or equal to {1}",
	"lt":       "{0} must be less than {1}",
	"lte":      "{0} must be less than or equal to {1}",
	"eqfield":  "{0} must be equal to field {1}",
	"nefield":  "{0} must not be equal to field {1}",
	"gtfield":  "{0} must be greater than field {1}",
	"gtefield": "{0} must be greater than or equal to field {1}",
	"ltfield":  "{0} must be less than field {1}",
	"ltefield": "{0} must be less than or equal to field {1}",

	// Collection validations
	"unique":   "{0} must contain unique values",
	"oneof":    "{0} must be one of {1}",
	"contains": "{0} must contain {1}",
	"excludes": "{0} must not contain {1}",

	// File validations
	"file":  "{0} must be a valid file",
	"dir":   "{0} must be a valid directory path",
	"path":  "{0} must be a valid filesystem path",
	"image": "{0} must be a valid image file",

	// Country/region validations
	"iso3166_1_alpha2":        "{0} must be a valid ISO 3166-1 alpha-2 country code",
	"iso3166_1_alpha3":        "{0} must be a valid ISO 3166-1 alpha-3 country code",
	"iso3166_1_alpha_numeric": "{0} must be a valid ISO 3166-1 alpha-numeric country code",
	"iso3166_2":               "{0} must be a valid ISO 3166-2 subdivision code",
	"iso4217":                 "{0} must be a valid ISO 4217 currency code",

	// Other common validations
	"boolean":    "{0} must be a boolean value",
	"rgb":        "{0} must be a valid RGB color",
	"rgba":       "{0} must be a valid RGBA color",
	"hsl":        "{0} must be a valid HSL color",
	"hsla":       "{0} must be a valid HSLA color",
	"ssn":        "{0} must be a valid Social Security Number",
	"creditcard": "{0} must be a valid credit card number",
	"ean":        "{0} must be a valid EAN barcode",
	"isbn":       "{0} must be a valid ISBN",
	"issn":       "{0} must be a valid ISSN",
	"postcode":   "{0} must be a valid postal code for locale {1}",
	"latitude":   "{0} must be a valid latitude coordinate",
	"longitude":  "{0} must be a valid longitude coordinate",
}

func ValidateRequestMiddleware[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input T

		if err := c.ShouldBindJSON(&input); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {

				errors := make(map[string]string)

				for _, fieldError := range validationErrors {
					errors[utils.ToSnakeCase(fieldError.Field())] = getValidationMessage(fieldError)
				}

				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"success": false,
					"errors":  errors,
				})

				return
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid request data",
			})
			return
		}

		c.Set("validatedInput", input)
		c.Next()
	}
}

func getValidationMessage(fieldError validator.FieldError) string {
	var field reflect.StructField
	var found bool

	val := fieldError.Value()
	rt := reflect.TypeOf(val)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Struct {
		field, found = rt.FieldByName(fieldError.Field())
	}

	if found {
		if message := field.Tag.Get("message"); message != "" {
			return message
		}
	}

	if template, ok := defaultValidationMessages[fieldError.Tag()]; ok {
		result := strings.Replace(template, "{0}", utils.ToSnakeCase(fieldError.Field()), 1)
		if fieldError.Param() != "" {
			result = strings.Replace(result, "{1}", fieldError.Param(), 1)
		}
		return result
	}

	return fmt.Sprintf("%s failed %s validation", utils.ToSnakeCase(fieldError.Field()), fieldError.Tag())
}
