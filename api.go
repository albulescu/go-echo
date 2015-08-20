package main

import (
	"encoding/json"
	"fmt"
	"github.com/dvsekhvalnov/jose2go"
	"time"
)

/**
 * Provide informations about the connections
 * by accessing the /info route on the bind port
 *
 * @route /api/info
 */
func ApiInfo(req Request, res Response) {
	res.Error("Not implemented", 0, 500)
}

func ApiMe(req Request, res Response) {
	res.Success(req.GetUser(), 200)
}

/**
 * Authentication
 *
 * @param {[type]} request  Request  [description]
 * @param {[type]} response Response [description]
 *
 * @route /api/auth
 */
func ApiAuth(req Request, res Response) {

	var UserId string = "111"

	//TODO: Authenticate user and if everything is ok set user in UserId var

	now := time.Now()

	expire, err := time.ParseDuration(configApi.Expire)

	if err != nil {
		res.InternalError(INTERNAL_ERROR_CONFIG_INVALID)
		return
	}

	data := JWTPayload{
		Id:  UserId,
		Iat: now.Unix(),
		Exp: now.Add(expire).Unix(),
	}

	jsonString, _ := json.Marshal(data)

	token, err := jose.Encrypt(string(jsonString), jose.DIR, jose.A128GCM, []byte(configApi.SecretKey), jose.Zip(jose.DEF))

	if err != nil {
		res.Error(fmt.Sprintf("Fail to encode jwt. %s", err), 0, 500)
		return
	}

	res.Success(map[string]string{
		"token": token,
	}, 200)
}
