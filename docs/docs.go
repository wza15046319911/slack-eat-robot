// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/sd/cpu": {
            "get": {
                "description": "CPUCheck",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sd"
                ],
                "summary": "CPUCheck checks the cpu usage.",
                "responses": {
                    "200": {
                        "description": "OK - Load average: xx, xx, xx | Cores: x",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/sd/disk": {
            "get": {
                "description": "DiskCheck",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sd"
                ],
                "summary": "DiskCheck checks the disk usage.",
                "responses": {
                    "200": {
                        "description": "OK - Free space: xxxMB (xxGB) / xxxMB (xxGB) | Used: xx%",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/sd/health": {
            "get": {
                "description": "HealthCheck",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "sd"
                ],
                "summary": "HealthCheck shows ` + "`" + `OK` + "`" + ` as the ping-pong result.",
                "responses": {
                    "200": {
                        "description": "OK ",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/sd/ram": {
            "get": {
                "description": "RAMCheck",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sd"
                ],
                "summary": "RAMCheck checks the disk usage.",
                "responses": {
                    "200": {
                        "description": "OK - Free space: xxMB (xxGB) / xxMB (xxGB) | Used: xx%",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/sd/version": {
            "get": {
                "description": "versionCheck",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "sd"
                ],
                "summary": "VersionCheck show the version info as Running status",
                "responses": {
                    "200": {
                        "description": "service version info ",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/slack": {
            "post": {
                "description": "Handle APP_MENTION event or other call back event sent from slack",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "event"
                ],
                "summary": "api",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
