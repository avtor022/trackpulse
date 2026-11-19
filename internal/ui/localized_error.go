package ui

import (
	"errors"
	"fmt"
	"regexp"

	"trackpulse/internal/locale"
)

// localizedError converts expected service errors to messages for the user.
// Unexpected technical details remain available after a localized explanation.
func localizedError(err error) error {
	if err == nil {
		return nil
	}
	keys := map[string]string{
		"full name is required":                           "error.required.competitor_name",
		"competitor number must be positive":              "error.invalid.number",
		"competitor is required":                          "error.required.competitor",
		"RC model is required":                            "error.required.model",
		"transponder number is required":                  "error.required.transponder_number",
		"brand is required":                               "error.required.brand",
		"model name is required":                          "error.required.name",
		"scale is required":                               "error.required.scale",
		"model type is required":                          "error.required.model_type",
		"competition title is required":                   "error.required.title",
		"competition type is required":                    "error.required.type",
		"invalid year":                                    "error.invalid.year",
		"competitor not found":                            "error.competitor_not_found",
		"RC model not found":                              "error.model_not_found",
		"competitor model not found":                      "error.transponder_not_found",
		"competition not found":                           "error.competition_not_found",
		"participant already exists for this competition": "error.participant_exists",
	}
	if key, ok := keys[err.Error()]; ok {
		return errors.New(locale.T(key))
	}
	patterns := []struct{ pattern, key string }{
		{`^competitor number (\d+) already exists$`, "error.competitor_number_exists"},
		{`^transponder number '(.+)' already exists$`, "error.transponder_exists"},
		{`^brand '(.+)' already exists$`, "error.brand_exists"},
		{`^scale '(.+)' already exists$`, "error.scale_exists"},
		{`^model type '(.+)' already exists$`, "error.model_type_exists"},
		{`^cannot delete brand '(.+)': it is used by model '(.+)'$`, "error.brand_in_use"},
		{`^cannot delete scale '(.+)': it is used by model '(.+)'$`, "error.scale_in_use"},
		{`^cannot delete model type '(.+)': it is used by model '(.+)'$`, "error.model_type_in_use"},
		{`^participant with model (.+) already exists for this competition$`, "error.participant_model_exists"},
		{`^competitor model (.+) not found$`, "error.transponder_id_not_found"},
	}
	for _, entry := range patterns {
		matches := regexp.MustCompile(entry.pattern).FindStringSubmatch(err.Error())
		if matches != nil {
			args := make([]interface{}, len(matches)-1)
			for i, value := range matches[1:] {
				args[i] = value
			}
			return fmt.Errorf(locale.T(entry.key), args...)
		}
	}
	return fmt.Errorf("%s: %w", locale.T("error.operation_failed"), err)
}
