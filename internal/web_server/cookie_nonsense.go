package server

import (
	"encoding/json"
	"net/http"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/npmaile/focusbot/internal/models"
	"github.com/npmaile/focusbot/pkg/logerooni"
)

func createCookie(key string, value any, w http.ResponseWriter) error {
	token, err := paseto.MakeToken(map[string]interface{}{key: value}, []byte{})
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	encrypted := token.V4Encrypt(pasetoKey, nil)

	realCookie := http.Cookie{
		Name:    "n8authCookie",
		Value:   encrypted,
		Expires: time.Now().Add(2 * time.Hour),
	}

	http.SetCookie(w, &realCookie)
	return err

}

type megaGuildConfig struct {
	*models.GuildConfig
	ServerName string
}

func getCookie(r *http.Request) (servers []megaGuildConfig, error error) {
	parser := paseto.NewParser()
	cookie, err := r.Cookie("n8authCookie")
	if err != nil {
		logerooni.Errorf("it's broken for this reason in particular: %s", err.Error())
	}
	tok, err := parser.ParseV4Local(pasetoKey, cookie.Value, nil)
	if err != nil {
		logerooni.Errorf("unable to parse token: %s", err.Error())
	}
	maybeServersRaw := tok.Claims()["servers"]
	ret := []megaGuildConfig{}
	err = json.Unmarshal([]byte(maybeServersRaw.(string)), &ret)
	return ret, err
}
