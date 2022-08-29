import json
import jwt
import web

def decide(service, claims):
    try:
        rulesFile = open("rules/rules.json", "r")
        rules = json.load(rulesFile)
        decision = rules[service]
        if decision == "allowed":
            return True
        elif decision == "moderator":
            return claims["moderator"] == True
        else:
            return False
    except (FileNotFoundError, KeyError) as err:
        return err
 
def auth_GET(data, jwtKey):
    response = {
        "allowed": False
    }

    if 'token' in data and 'service' in data:
        token = data.token
        service = data.service

        try:
            claims = jwt.decode(token, jwtKey, algorithms=["HS256"], options={"verify_signature": True})

            try:
                response["allowed"] = decide(service, claims)
            except FileNotFoundError:
                web.webapi.internalerror()
                print("No rules.json found. Did you convert the m4 file to json using m4 rules.m4 > rules.json?")
            except KeyError:
                web.webapi.badrequest()
                response["allowed"] = False
        except jwt.InvalidTokenError:
            response["allowed"] = False
    else:
        web.webapi.badrequest()
    
    web.webapi.header("Access-Control-Allow-Origin", "*")
    web.webapi.header("Content-Type", "application/json") 
    return json.dumps(response) 
