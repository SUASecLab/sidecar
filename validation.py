import json
import jwt
import web

def validation_GET(data, jwtKey):
    result = {
        "error": None,
        "valid": False
    }
    
    if 'token' in data:
        token = data.token
        try:
            jwt.decode(token, jwtKey, algorithms=["HS256"], options={"verify_signature": True})
            result["valid"] = True
        except jwt.InvalidTokenError:
            result["error"] = "The provided token is invalid"
            result["valid"] = False
    else:
        web.webapi.badrequest()
        result["error"] = "No token provided"
    
    web.webapi.header("Access-Control-Allow-Origin", "*")
    web.webapi.header("Content-Type", "application/json")
    return json.dumps(result)
