#!/bin/env python3
from os import getenv
import web

# Modules
from auth import *
from issuance import *
from userinfo import *
from validation import *

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

# Declaration of routes

# Authorize requests
class auth:
    def GET(self):
        return auth_GET(web.input(), jwtKey)

# Token issuance
class issuance:
    def GET(self):
        return issuance_GET(web.input(), jwtKey, dbUser, dbPassword, dbHost, dbName, jitsiIssuer, jitsiKey)

# User information
class userinfo:
    def GET(self):
        return userinfo_GET(web.input(), dbUser, dbPassword, dbHost, dbName)

# Token validation
class validate:
    def GET(self):
        return validation_GET(web.input(), jwtKey)

# Start the program
if __name__ == "__main__":    
    if dbUser == None or dbPassword == None or dbHost == None or dbName == None:
        print("Database information incomplete")
    if jitsiIssuer == None or jitsiKey == None:
        print("Jitsi information incomplete")
    
    print("Starting up sidecar")

    #web.config.debug = False
    app = web.application(urls, globals())
    app.run()
