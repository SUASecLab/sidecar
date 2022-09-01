import json
import jwt
from pymongo import MongoClient
import time
import web

from sidecar_lib.auth import decide

def issuance_GET(data, jwtKey, dbUser, dbPassword, dbHost, dbName, jitsiIssuer, jitsiKey):
    result = {
        "error": None,
        "token": ""
    }
    currentTime = round(time.time())
    
    if 'uuid' in data:
        uuid = data.uuid
        client = MongoClient("mongodb://" + dbUser + ":" + dbPassword +
                            "@" + dbHost + ":27017/" + dbName)
        db = client.workadventure
        user = db.users.find_one({'uuid': uuid})
        
        # User not found
        if user == None:
            web.webapi.forbidden()
            result["error"] = "User does not exist"
        else:
            tags = user["tags"]
            payload = {
                "uuid": uuid,
                "tags": tags,
                "moderator": "admin" in tags,
                "exp": currentTime + 60 * 60 * 24,
                "nbf": currentTime - 10,
                "iat": currentTime 
            }
            result["token"] = jwt.encode(payload, jwtKey, 'HS256')
        client.close()
    elif 'name' in data and 'token' in data:
        # generate Jitsi token
        name = data.name
        token = data.token
        try:
            claims = jwt.decode(token, jwtKey, algorithms=["HS256"], options={"verify_signature": True})
            isModerator = decide("jitsiModerator", claims)
            payload = {
                "context": {
                    "user": {
                        "name": name
                    }
                },
                "nbf": currentTime - 10,
                "aud": "jitsi",
                "iss": jitsiIssuer,
                "room": "*",
                "moderator": isModerator,
                "iat": currentTime,
                "exp": currentTime + 60
            }
            result["token"] = jwt.encode(payload, jitsiKey, 'HS256')
        except jwt.InvalidTokenError:
            result["error"] = "The provided token is invalid"
    else:
        web.webapi.badrequest()
        result["error"] = "No UUID provided"
    
    web.webapi.header("Access-Control-Allow-Origin", "*")
    web.webapi.header("Content-Type", "application/json")
    return json.dumps(result)
