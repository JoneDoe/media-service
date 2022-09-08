package middleware

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"istorage/logger"
	"istorage/utils"
	"istorage/wsse"
)

var (
	wsseProfile   = `WSSE profile="UsernameToken"`
	context       *gin.Context
	rxSplitHeader = regexp.MustCompile(`\s*,\s*`)
	rxKeyValue    = regexp.MustCompile(`^(\w+)="(.+)"$`)
)

func AuthGuard(c *gin.Context) {
	context = c

	endpoint, _ := os.LookupEnv("FIREWALL_ENDPOINT")
	if endpoint == "" {
		return
	}

	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	var err error

	if authorization == wsseProfile {
		err = WsseGuard()
	} else {
		err = BearerGuard(endpoint, authorization)
	}

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	c.Next()
}

func WsseGuard() error {
	kv := splitWsseHeader(context.GetHeader("X-WSSE"))
	wsseUser, _ := os.LookupEnv("WSSE_INTERNAL_USERNAME")

	if kv["Username"] != wsseUser {
		return errors.New("unauthorized")
	}

	wssePassword, _ := os.LookupEnv("WSSE_INTERNAL_SECRET")
	passwordDigest := wsse.CreatePasswordDigest(kv["Nonce"], kv["Created"], wssePassword)

	if passwordDigest != kv["PasswordDigest"] {
		return errors.New("unauthorized")
	}

	if validateLifetime(kv["Created"]) == false {
		return errors.New("unauthorized")
	}

	return nil
}

func BearerGuard(endpoint, authorization string) error {
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		utils.Response{context}.Error(http.StatusUnauthorized, "Authorization service unavailable")

		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("unauthorized")
	}

	return nil
}

func splitWsseHeader(wsseHeader string) map[string]string {
	kv := map[string]string{}

	wsseHeader = strings.TrimPrefix(wsseHeader, "UsernameToken ")

	parts := rxSplitHeader.Split(wsseHeader, -1)
	for _, part := range parts {
		m := rxKeyValue.FindStringSubmatch(part)
		if m != nil {
			kv[m[1]] = m[2]
		}
	}

	return kv
}

func validateLifetime(wsseCreated string) bool {
	now := time.Now()
	created, _ := time.Parse(time.RFC3339, wsseCreated)

	if created.After(now) {
		return false
	}

	wsseLifetime, _ := os.LookupEnv("WSSE_LIFETIME")
	lifetime, _ := time.ParseDuration(wsseLifetime)

	if created = created.Add(lifetime); created.Before(now) {
		return false
	}

	return true
}
