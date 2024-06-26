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
        "/admin/deductions/k-receipt": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Admin set k-receipt deduction",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Admin set k-receipt deduction",
                "parameters": [
                    {
                        "description": "Amount to set personal deduction",
                        "name": "amount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.Deduction"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/admin.KReceiptDeduction"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/admin.Err"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/admin.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/admin.Err"
                        }
                    }
                }
            }
        },
        "/admin/deductions/personal": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Admin set personal deduction",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Admin set personal deduction",
                "parameters": [
                    {
                        "description": "Amount to set personal deduction",
                        "name": "amount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.Deduction"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/admin.PersonalDeduction"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/admin.Err"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/admin.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/admin.Err"
                        }
                    }
                }
            }
        },
        "/tax/calculations": {
            "post": {
                "description": "Calculate tax",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "summary": "Calculate tax",
                "parameters": [
                    {
                        "description": "Amount to calculate tax",
                        "name": "amount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tax.TaxInformation"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tax.TaxResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    }
                }
            }
        },
        "/tax/calculations/upload-csv": {
            "post": {
                "description": "Upload csv file and calculate tax",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "summary": "Upload csv file and calculate tax",
                "parameters": [
                    {
                        "type": "file",
                        "description": "this is a test file",
                        "name": "taxFile",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tax.CsvTaxResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "admin.Deduction": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "minimum": 0
                }
            }
        },
        "admin.Err": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "admin.KReceiptDeduction": {
            "type": "object",
            "properties": {
                "kReceipt": {
                    "type": "number"
                }
            }
        },
        "admin.PersonalDeduction": {
            "type": "object",
            "properties": {
                "personalDeduction": {
                    "type": "number"
                }
            }
        },
        "tax.Allowance": {
            "type": "object",
            "properties": {
                "allowanceType": {
                    "$ref": "#/definitions/tax.AllowanceType"
                },
                "amount": {
                    "type": "number",
                    "minimum": 0
                }
            }
        },
        "tax.AllowanceType": {
            "type": "string",
            "enum": [
                "donation",
                "k-receipt"
            ],
            "x-enum-varnames": [
                "AllowanceTypeDonation",
                "AllowanceTypeKReceipt"
            ]
        },
        "tax.CsvTaxRecord": {
            "type": "object",
            "properties": {
                "tax": {
                    "type": "number"
                },
                "taxRefund": {
                    "type": "number"
                },
                "totalIncome": {
                    "type": "number"
                }
            }
        },
        "tax.CsvTaxResponse": {
            "type": "object",
            "properties": {
                "taxes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.CsvTaxRecord"
                    }
                }
            }
        },
        "tax.Err": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "tax.TaxInformation": {
            "type": "object",
            "required": [
                "totalIncome"
            ],
            "properties": {
                "allowances": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.Allowance"
                    }
                },
                "totalIncome": {
                    "type": "number",
                    "minimum": 0
                },
                "wht": {
                    "type": "number",
                    "minimum": 0
                }
            }
        },
        "tax.TaxLevel": {
            "type": "object",
            "properties": {
                "level": {
                    "type": "string"
                },
                "tax": {
                    "type": "number"
                }
            }
        },
        "tax.TaxResult": {
            "type": "object",
            "properties": {
                "tax": {
                    "type": "number"
                },
                "taxLevel": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.TaxLevel"
                    }
                },
                "taxRefund": {
                    "type": "number"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "K-Tax API",
	Description:      "This is an API for K-Tax.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
