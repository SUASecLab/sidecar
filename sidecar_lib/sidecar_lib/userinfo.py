import json
import jwt
from pymongo import MongoClient
import time
import web

def userinfo_GET(data, dbUser, dbPassword, dbHost, dbName):
    result = {
        "exists": False,
        "isAdmin": False
    }
    
    if 'uuid' in data:
        uuid = data.uuid
        client = MongoClient("mongodb://" + dbUser + ":" + dbPassword +
                            "@" + dbHost + ":27017/" + dbName)
        db = client.workadventure
        user = db.users.find_one({'uuid': uuid})
        
        # User not found
        if user == None:
            result["exists"] = False
            result["isAdmin"] = False
        else:
            result["exists"] = True
            result["isAdmin"] = "admin" in user["tags"]
        
        client.close()
    else:
        web.webapi.badrequest()
        result["error"] = "No UUID provided"
    
    web.webapi.header("Access-Control-Allow-Origin", "*")
    web.webapi.header("Content-Type", "application/json")
    return json.dumps(result)
