// Package swagger Code generated by swaggo/swag. DO NOT EDIT
package swagger

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
        "/": {
            "get": {
                "description": "Get metric all from storage",
                "tags": [
                    "V1 API"
                ],
                "summary": "List metrics",
                "operationId": "ListHandler",
                "responses": {
                    "200": {
                        "description": "Metrics list",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Inernal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/update/": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2 API"
                ],
                "summary": "Update metrics",
                "operationId": "UpdateHandlerV2",
                "parameters": [
                    {
                        "description": "Metric Name",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.MetricsV2"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.MetricsV2"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Inernal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/update/{type}/{name}/": {
            "get": {
                "description": "Get metric from storage",
                "tags": [
                    "V1 API"
                ],
                "summary": "Get metric",
                "operationId": "GetHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Metric name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Metric type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Inernal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/update/{type}/{name}/{value}/": {
            "post": {
                "tags": [
                    "V1 API"
                ],
                "summary": "Update metrics",
                "operationId": "UpdateHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Metric name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Metric type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Metric value ",
                        "name": "value",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Inernal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/updates/": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2 API"
                ],
                "summary": "Batch update",
                "operationId": "BatchUpdateHandler",
                "parameters": [
                    {
                        "description": "Metrics request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.MetricsV2"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.MetricsV2"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Inernal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/value/": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2 API"
                ],
                "summary": "Get metrics",
                "operationId": "GetHandlerV2",
                "parameters": [
                    {
                        "description": "Metric request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.MetricsV2"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.MetricsV2"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Inernal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.MetricType": {
            "type": "string",
            "enum": [
                "gauge",
                "counter"
            ],
            "x-enum-varnames": [
                "GaugeType",
                "CounterType"
            ]
        },
        "model.MetricsV2": {
            "type": "object",
            "properties": {
                "delta": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/model.MetricType"
                },
                "value": {
                    "type": "number"
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
