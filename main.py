#!/bin/env python3
from apispec import APISpec
from os import getenv
import web

# Modules
from sidecar_lib.auth import *
from sidecar_lib.issuance import *
from sidecar_lib.userinfo import *
from sidecar_lib.validation import *

# Environment variables
dbUser = getenv("DB_USER")
dbPassword = getenv("DB_PASSWORD")
dbHost = getenv("DB_HOST")
dbName = getenv("DB_NAME")
jwtKey = getenv("JWT_KEY")
jitsiIssuer = getenv("JITSI_ISS")
jitsiKey = getenv("SECRET_JITSI_KEY")

# Definition of routes
urls = (
    '/auth', 'auth',
    '/issuance', 'issuance',
    '/userinfo', 'userinfo',
    '/validate', 'validate'
)

# OpenAPI
spec = APISpec(
    title = "SUASecLab Sidecar",
    version = "0.0.1",
    openapi_version="3.0.3",
    info=dict(
        description="This is the API documentation for the SUSecLab's sidecar. The sidecar is responsible for all security related operations in the SUASecLab, e.g. generating tokens for the users or checking if a user has the permission to do something.",
        license=dict(
            name="GPL-3.0",
            url="https://www.gnu.org/licenses/gpl-3.0.en.html"
        ),
        contact=dict(
            email="t.tefke@stud.fh-sm.de"
        ),
    ),
)

# Declaration of routes

# Authorize requests
spec.components.schema(
    "AuthDecision",
    {
        "type": "object",
        "properties": {
            "allowed": {
                "type": "boolean",
                "example": True,
            },
        },
    },
)

spec.components.parameter(
    "AuthRequestService",
    "query",
     {
        "name": "service",
        "in": "query",
        "description": "The service or task for which the user should be authenticated.",
        "required": True,
        "schema": {"type": "string"},
        "example": "updateComponents"
    }
)

spec.components.parameter(
    "AuthRequestToken",
    "query",
     {
        "name": "token",
        "in": "query",
        "description": "The JWT token of the user.",
        "required": True,
        "schema": {"type": "string", "format": "JWT"},
    }
)

spec.path(
    path="/auth",
    summary="Make auth decision",
    description="API call to the sidebar in order to make an auth decision.",
    parameters=list(
        {
            "AuthRequestService": "AuthRequestService",
            "AuthRequestToken": "AuthRequestToken",
        }
    ),
    operations=dict(
        get=dict(
            summary="Get auth decision",
            description="Get auth decision from the sidecar. The decision tells whether the user is allowed to perform a specific operation.",
            responses = {
                "200": {
                    "description": "Decision could be made and is being returned.",
                    "content": {"application/json": {"schema": "AuthDecision"}}
                },
                "400": {
                    "description": "Decision could not be made, required parameters are (partly) missing."
                }
            }
        )
    ),
)

class auth:
    def GET(self):
        return auth_GET(web.input(), jwtKey)

# Token issuance
spec.components.schema(
    "TokenIssuanceResult",
    {
        "type": "object",
        "properties": {
            "error": {
                "type": "string",
                "example": None,
            },
            "token": {
                "type": "string",
                "format": "JWT"
            }
        },
    },
)

spec.components.parameter(
    "IssuanceRequestUUID",
    "query",
     {
        "name": "uuid",
        "in": "query",
        "description": "The UUID to specify a user account. This must be handed over when issuing a token for a SUASecLab user.",
        "schema": {"type": "string", "format": "UUID"},
    }
)

spec.components.parameter(
    "IssuanceRequestName",
    "query",
     {
        "name": "name",
        "in": "query",
        "description": "The name of a user. Must be handed over together with a token when issuing a token for a Jitsi Meet room.",
        "schema": {"type": "string"},
    }
)

spec.components.parameter(
    "IssuanceRequestToken",
    "query",
     {
        "name": "token",
        "in": "query",
        "description": "The SUASecLab token of a user. Must be handed over together with a name when issuing a token for a Jitsi Meet room.",
        "schema": {"type": "string", "format": "JWT"},
    }
)

spec.path(
    path="/issuance",
    summary="Issue JWT token",
    description="This route can be used for generating a user token for the SUASecLab or for receiving a token for a Jitsi Meet instance. If a token for the SUASecLab should be generated, a UUID must be handed over. Otherwise, if a Jitsi Meet token should be generated, an existing SUASecLab user token and a name to be displayed in Jitsi must be provided.",
    parameters=list(
        {
            "IssuanceRequestUUID": "IssuanceRequestUUID",
            "IssuanceRequestName": "IssuanceRequestName",
            "IssuanceRequestToken": "IssuanceRequestToken",
        }
    ),
    operations=dict(
        get=dict(
            summary="Get JWT token",
            description="Get a JWT token for the user to authenticate her on the SUASecLab or a Jitsi Meet room.",
            responses = {
                "200": {
                    "description": "Token was issued and is being returned.",
                    "content": {"application/json": {"schema": "TokenIssuanceResult"}}
                },
                "400": {
                    "description": "The request was invalid, not enough parameters specified."
                },
                "403": {
                    "description": "Token could not be issued, because there exists no user having the specified UUID."
                },
                "409": {
                    "description": "The provided JWT token is invalid."
                }
            }
        )
    ),
)

class issuance:
    def GET(self):
        return issuance_GET(web.input(), jwtKey, dbUser, dbPassword, dbHost, dbName, jitsiIssuer, jitsiKey)

# User information
spec.components.schema(
    "UserInfoResult",
    {
        "type": "object",
        "properties": {
            "exists": {
                "type": "boolean",
                "example": True,
            },
            "isAdmin": {
                "type": "boolean",
                "example": True
            }
        },
    },
)

spec.components.parameter(
    "UserInfoUUID",
    "query",
     {
        "name": "uuid",
        "in": "query",
        "required": True,
        "description": "The UUID to specify tje user account. This must be handed over for receiving user information.",
        "schema": {"type": "string", "format": "UUID"},
    }
)

spec.path(
    path="/userinfo",
    summary="Get user information",
    description="Find out whether an account exists and if the user is an admin.",
    parameters=list({"UserInfoUUID":"UserInfoUUID"}),
    operations=dict(
        get=dict(
            summary="Get user information",
            description="Find out whether an account exists and if the user is an admin.",
            responses = {
                "200": {
                    "description": "The request was processed, information is being returned.",
                    "content": {"application/json": {"schema": "UserInfoResult"}}
                },
                "400": {
                    "description": "No user UUID was provided."
                }
            }
        )
    )
)

class userinfo:
    def GET(self):
        return userinfo_GET(web.input(), dbUser, dbPassword, dbHost, dbName)

# Token validation
spec.components.schema(
    "ValidationResult",
    {
        "type": "object",
        "properties": {
            "error": {
                "type": "string",
                "example": None,
            },
            "valid": {
                "type": "boolean",
                "example": True
            }
        },
    },
)

spec.components.parameter(
    "ValidateToken",
    "query",
     {
        "name": "token",
        "in": "query",
        "required": True,
        "description": "The SUASecLab user token of a user.",
        "schema": {"type": "string", "format": "JWT"},
    }
)

spec.path(
    path="/validate",
    summary="Validate JWT tokens",
    description="Validate user JWT tokens for the SUASecLab.",
    parameters=list({"ValidateToken":"ValidateToken"}),
    operations=dict(
        get=dict(
            summary="Validate JWT tokens",
            description="Validate user JWT tokens for the SUASecLab.",
            responses = {
                "200": {
                    "description": "The request was processed, information about the state of the token is being returned.",
                    "content": {"application/json": {"schema": "ValidationResult"}}
                },
                "400": {
                    "description": "No token was provided."
                }
            }
        )
    )
)
class validate:
    def GET(self):
        return validation_GET(web.input(), jwtKey)

# Start the program  
if dbUser == None or dbPassword == None or dbHost == None or dbName == None:
    print("Database information incomplete")
if jitsiIssuer == None or jitsiKey == None:
    print("Jitsi information incomplete")

with open("openapi.yaml", "w") as api_spec:
    api_spec.write(spec.to_yaml())
    print("API specification written to file")
    
print("Starting up sidecar")

web.config.debug = False
application = web.application(urls, globals()).wsgifunc()
