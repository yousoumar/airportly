// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/metadata/airports": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metadata"
                ],
                "summary": "Get all available airport IDs",
                "operationId": "getAvailableAirportIds",
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/{airportIATA}/available-metrics": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metadata"
                ],
                "summary": "Get available metrics for a specific airport",
                "operationId": "getAvailableMetrics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The IATA code of the airport",
                        "name": "airportIATA",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/{airportIATA}/metric/{metric}": {
            "get": {
                "description": "Get data for a specific metric at an airport between two times",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data"
                ],
                "summary": "Get data between two times",
                "operationId": "getDataBetweenTwoTimes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The IATA code of the airport",
                        "name": "airportIATA",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The type of metric (e.g., pressure, temperature, wind-speed)",
                        "name": "metric",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The start time in RFC3339 format",
                        "name": "startTime",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The end time in RFC3339 format",
                        "name": "endTime",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/data_type.DataType"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/{airportIATA}/metric/{metric}/average": {
            "get": {
                "description": "Get the average value of a specific metric at an airport for a given date",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Average"
                ],
                "summary": "Get average value of a metric in a day",
                "operationId": "getAverageForSingleTypeInDay",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The IATA code of the airport",
                        "name": "airportIATA",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The type of metric (e.g., pressure, temperature, wind-speed)",
                        "name": "metric",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The date in RFC3339 format",
                        "name": "date",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/main.SuccessfulAverageResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/{airportIATA}/metric/{metric}/date-range": {
            "get": {
                "description": "Get the date interval for a specific metric at an airport",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metadata"
                ],
                "summary": "Get date interval of a specific metric",
                "operationId": "getDateInterval",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The IATA code of the airport",
                        "name": "airportIATA",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The type of metric (e.g., pressure, temperature, wind-speed)",
                        "name": "metric",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/main.DateIntervalResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/{airportIATA}/metrics/average": {
            "get": {
                "description": "Get the average value of all metrics at an airport for a given date",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Average"
                ],
                "summary": "Get average value of all metrics in a day",
                "operationId": "getAverageForAllTypesInDay",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The IATA code of the airport",
                        "name": "airportIATA",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The date in RFC3339 format",
                        "name": "date",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/main.AverageAllResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "data_type.DataType": {
            "type": "object",
            "properties": {
                "airportId": {
                    "type": "string"
                },
                "sensorId": {
                    "type": "integer"
                },
                "sensorType": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "main.AverageAllResponse": {
            "type": "object",
            "properties": {
                "pressure": {
                    "type": "number"
                },
                "temperature": {
                    "type": "number"
                },
                "wind-speed": {
                    "type": "number"
                }
            }
        },
        "main.DateIntervalResponse": {
            "type": "object",
            "properties": {
                "endTime": {
                    "type": "string"
                },
                "startTime": {
                    "type": "string"
                }
            }
        },
        "main.SuccessfulAverageResponse": {
            "type": "object",
            "properties": {
                "average": {
                    "type": "number"
                },
                "unit": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
