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
  "paths": {
    "/api/address/geocode": {
      "post": {
        "description": "Search by geocode",
        "operationId": "Geocode",
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
        "operationId": "Search",
        "responses": {
          "200": {
            "$ref": "#/responses/apiResponseSearch"
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "description": "Login",
        "operationId": "Login",
        "responses": {
          "200": {
            "$ref": "#/responses/apiLoginResponse"
          }
        }
      }
    },
    "/auth/register": {
      "post": {
        "description": "Register",
        "operationId": "Register",
        "responses": {
          "200": {
            "$ref": "#/responses/apiRegisterResponse"
          }
        }
      }
    },
    "/user/list": {
      "post": {
        "description": "List",
        "operationId": "List",
        "responses": {
          "200": {
            "$ref": "#/responses/apiListResponse"
          }
        }
      }
    },
    "/user/profile": {
      "post": {
        "description": "Profile",
        "operationId": "Profile",
        "responses": {
          "200": {
            "$ref": "#/responses/apiProfileResponse"
          }
        }
      }
    }
  },
  "definitions": {
    "Geo": {
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
      "x-go-package": "proxy/internal/modules/controller"
    },
    "GeoResponse": {
      "type": "object",
      "properties": {
        "addresses": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Geo"
          },
          "x-go-name": "Addresses"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    },
    "ListUser": {
      "type": "object",
      "properties": {
        "users": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/UserFromRpc"
          },
          "x-go-name": "Users"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    },
    "ProfileRequest": {
      "description": "User endponit types",
      "type": "object",
      "properties": {
        "email": {
          "type": "string",
          "x-go-name": "Email"
        },
        "phone": {
          "type": "string",
          "x-go-name": "Phone"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    },
    "ProfileResponse": {
      "type": "object",
      "properties": {
        "email": {
          "type": "string",
          "x-go-name": "Email"
        },
        "password": {
          "type": "string",
          "x-go-name": "Password"
        },
        "phone": {
          "type": "string",
          "x-go-name": "Phone"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    },
    "RequestAuth": {
      "description": "Auth register and login type",
      "type": "object",
      "properties": {
        "email": {
          "type": "string",
          "x-go-name": "Email"
        },
        "password": {
          "type": "string",
          "x-go-name": "Password"
        },
        "phone": {
          "type": "string",
          "x-go-name": "Phone"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    },
    "RequestGeoGeo": {
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
      "x-go-package": "proxy/internal/modules/controller"
    },
    "RequestGeoSearch": {
      "description": "Geo search and geocode type",
      "type": "object",
      "properties": {
        "query": {
          "type": "string",
          "x-go-name": "Query"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    },
    "UserFromRpc": {
      "type": "object",
      "properties": {
        "email": {
          "type": "string",
          "x-go-name": "Email"
        },
        "password": {
          "type": "string",
          "x-go-name": "Password"
        },
        "phone": {
          "type": "string",
          "x-go-name": "Phone"
        }
      },
      "x-go-package": "proxy/internal/modules/controller"
    }
  },
  "responses": {
    "apiGeocodeResponse": {
      "schema": {
        "$ref": "#/definitions/GeoResponse"
      }
    },
    "apiListResponse": {
      "schema": {
        "$ref": "#/definitions/ListUser"
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
    "apiProfileResponse": {
      "schema": {
        "$ref": "#/definitions/ProfileResponse"
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
        "$ref": "#/definitions/GeoResponse"
      }
    }
  }
}
`
