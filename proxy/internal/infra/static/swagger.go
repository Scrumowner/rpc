package static

const (
	SwaggerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
   <meta charset="UTF-8">
   <meta name="viewport" content="width=device-width, initial-scale=1.0">
   <meta http-equiv="X-UA-Compatible" content="ie=edge">
   <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
   <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-standalone-preset.js"></script> -->
   <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
   <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-bundle.js"></script> -->
   <link rel="stylesheet" href="//unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
   <!-- <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui.css" /> -->
	<style>
		body {
			margin: 0;
		}
	</style>
   <title>Swagger</title>
</head>
<body>
   <div id="swagger-ui"></div>
   <script>
       window.onload = function() {
         SwaggerUIBundle({
           url: "swagger/swagger",
           dom_id: '#swagger-ui',
           presets: [
             SwaggerUIBundle.presets.apis,
             SwaggerUIStandalonePreset
           ],
           layout: "StandaloneLayout"
         })
       }
   </script>
</body>
</html>
`
)

var SwagJson = `
{
  "swagger": "2.0",
  "info": {},
  "securityDefinitions": {
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "Bearer token for authentication"
    }
  },
  "paths": {
    "/api/address/geocode": {
      "post": {
        "description": "Search by geocode",
        "operationId": "apiGeocodeRequest",
        "parameters": [
          {
            "name": "Body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/GeocodeRequest"
            }
          },
          {
            "type": "string",
            "x-go-name": "Token",
            "name": "Authorization",
            "in": "header"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/apiGeocodeResponse"
          }
        }
      }
    },
    "/api/address/search": {
      "post": {
        "description": "Search by address",
        "operationId": "apiRequestSearch",
        "parameters": [
          {
            "name": "Body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/SearchRequest"
            }
          },
          {
            "type": "string",
            "x-go-name": "Token",
            "name": "Authorization",
            "in": "header"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/apiResponseSearch"
          }
        }
      }
    },
    "/api/login": {
      "post": {
        "description": "Login",
        "operationId": "apiLoginRequest",
        "parameters": [
          {
            "name": "Body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/User"
            }
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/apiLoginResponse"
          }
        }
      }
    },
    "/api/register": {
      "post": {
        "description": "Register",
        "operationId": "apiRegisterRequest",
        "parameters": [
          {
            "name": "Body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          {
            "type": "string",
            "x-go-name": "Token",
            "name": "Authorization",
            "in": "header"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/apiRegisterResponse"
          }
        }
      }
    }
  },
  "definitions": {
    "GeocodeRequest": {
      "type": "object",
      "properties": {
        "lat": {
          "type": "string",
          "x-go-name": "Lat"
        },
        "lng": {
          "type": "string",
          "x-go-name": "Lng"
        }
      },
      "x-go-package": "hugoproxy-main/middleware/handler"
    },
    "ReworkedSearch": {
      "type": "object",
      "properties": {
        "lat": {
          "type": "string",
          "x-go-name": "GeoLat"
        },
        "lon": {
          "type": "string",
          "x-go-name": "GeoLon"
        },
        "result": {
          "type": "string",
          "x-go-name": "Result"
        }
      },
      "x-go-package": "hugoproxy-main/middleware/handler"
    },
    "ReworkedSearchResponse": {
      "type": "object",
      "properties": {
        "addresses": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ReworkedSearch"
          },
          "x-go-name": "Addresses"
        }
      },
      "x-go-package": "hugoproxy-main/middleware/handler"
    },
    "SearchRequest": {
      "type": "object",
      "properties": {
        "query": {
          "type": "string",
          "x-go-name": "Query"
        }
      },
      "x-go-package": "hugoproxy-main/middleware/handler"
    },
    "User": {
      "type": "object",
      "properties": {
        "password": {
          "type": "string",
          "x-go-name": "Password"
        },
        "username": {
          "type": "string",
          "x-go-name": "Username"
        }
      },
      "x-go-package": "hugoproxy-main/middleware/handler"
    }
  },
  "responses": {
    "apiGeocodeResponse": {
      "schema": {
        "$ref": "#/definitions/ReworkedSearchResponse"
      }
    },
    "apiLoginResponse": {
      "schema": {
        "type": "object",
        "properties": {
          "token": {
            "type": "string",
            "x-go-name": "Token"
          }
        }
      }
    },
    "apiRegisterResponse": {
      "schema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string",
            "x-go-name": "Message"
          }
        }
      }
    },
    "apiResponseSearch": {
      "schema": {
        "$ref": "#/definitions/ReworkedSearch"
      }
    }
  }
}
`
