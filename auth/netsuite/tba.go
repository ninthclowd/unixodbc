package netsuite

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ninthclowd/unixodbc"

	"time"
)

var _ unixodbc.ConnectionStringFactory = (*TBA)(nil)

type TBA struct {
	Host           string
	Port           int
	DataSource     string
	RoleId         int
	ConsumerKey    string
	ConsumerSecret string
	TokenKey       string
	TokenSecret    string
	AccountId      string
}

func (s *TBA) ConnectionString() (string, error) {
	n, err := s.nonce()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("DRIVER=/opt/netsuite/odbcclient/lib64/ivoa27.so;Host=%s;Port=%d;Encrypted=1;AllowSinglePacketLogout=1;Truststore=/opt/netsuite/odbcclient/cert/ca3.cer;SDSN=%s;UID=TBA;PWD=%s;CustomProperties=AccountID=%s;RoleID=%d",
		s.Host,
		s.Port,
		s.DataSource,
		s.token(n, time.Now()),
		s.AccountId,
		s.RoleId), nil
}

func (s *TBA) nonce() (string, error) {
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(nonceBytes), nil
}

func (s *TBA) token(nonce string, timestamp time.Time) string {
	base := fmt.Sprintf("%s&%s&%s&%s&%d",
		s.AccountId,
		s.ConsumerKey,
		s.TokenKey,
		nonce,
		timestamp.Unix())
	secret := []byte(s.ConsumerSecret + "&" + s.TokenSecret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(base))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return base + "&" + signature + "&HMAC-SHA256"
}
