package entityresolution

// Resolutions is part of an Slot
// An Resolutions object represents the results of resolving the words captured
// from the user’s utterance.
//
// For example:
// User can define synonyms that should resolve to specific slot value
// with the following sample json in their interaction model:
// {
//      "types": [
//          {
//          "name": "WeatherTypes",
//          "values": [
//              {
//              "id" : "RAIN",
//              "name" : {
//                  "value": "rain",
//                  "synonyms": ["shower", "storm", "rainstorm"]
//              }
//              }
//          ]
//          }
//      ]
// }
//
// In the above example, "shower" is defined as synonym for "rain".
// When user says "shower", a response containing the following slot sample may
// be returned to the skill:
// {
//     "slots": {
//          "WeatherType": {
//              "name": "WeatherType",
//              "value": "shower",
//              "resolutions": {
//                  "resolutionsPerAuthority": [{
//                      "authority": "{authority-url}",
//                      "status": {
//                          "code": "ER_SUCCESS_MATCH"
//                      },
//                      "values": [{
//                          "value": {
//                              "name": "rain",
//                              "id": "RAIN"
//                              }
//                           }]
//                     }]
//                 }
//          }
// }
//
// In the above response json, "shower" is still passed as value to ensure
// backward-compatibility. But "rain" appears as resolution value
// Resolutions is included for slots that use a custom slot type
//  or a built-in slot type that have been extended with custom values.
// Note that resolutions is not included for built-in slot types that you have
// not extended.

type StatusCode string

// Indication of the results of attempting to resolve the user utterance against
// the defined slot types.
// This can be one of the following:
// ER_SUCCESS_MATCH: The spoken value matched a value or synonym explicitly
// defined in your custom slot type.
// ER_SUCCESS_NO_MATCH: The spoken value did not match any values or synonyms
// explicitly defined in your custom slot type.
// ER_ERROR_TIMEOUT: An error occurred due to a timeout.
// ER_ERROR_EXCEPTION: An error occurred due to an exception during processing.
const (
	ER_SUCCESS_MATCH    StatusCode = "ER_SUCCESS_MATCH"
	ER_SUCCESS_NO_MATCH StatusCode = "ER_SUCCESS_NO_MATCH"
	ER_ERROR_TIMEOUT    StatusCode = "ER_ERROR_TIMEOUT"
	ER_ERROR_EXCEPTION  StatusCode = "ER_ERROR_EXCEPTION"
)

type Status struct {
	Code StatusCode `json:"code"`
}

type Value struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ValueWrapper struct {
	Value
}

type Resolution struct {
	Authority string         `json:"authority"`
	Status    Status         `json:"status"`
	Values    []ValueWrapper `json:"values"`
}

type Resolutions struct {
	ResolutionsPerAuthority []Resolution `json:"resolutionsPerAuthority"`
}
