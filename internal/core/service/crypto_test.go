package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	pubkeyData := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxscinIin68E0Dn+AQunE
/GTkHvqTOSN63PK/693Sap638lMwXnhnnZkr+Ts/uqgEwqefUW05DNPH9+s+CaZ4
0ZyyaOudZ6G3sgVpwxAsqUoIPBdnL/XhYMsqjZy8eQ+h2k3m7hP5iDkWxRV/YH52
WL7vHPU3LLzyNv30lG5szYHvStcGDsOB6TXVOYNpC0BveBwL2E45BDeAMlLoLMOC
6C2jMhjfZBwyKz3xEoJXSgjh4vjCHPTRyMyBhWOKHTWa4LeMAt6bbnYFKJB4eQyc
xY0wjXc7V57ZFic7LbLjxzh/Do/ZzJE7UsBpZYoy9ZB36ajMb5nPRm8Y1/l+/mEe
g3UBufr0yNw7xy/hfIXP75vUH1wODZOR+qbrGSrILWT2IgmQmXLxtu9CTqkqgT3X
yl28GGdN9T7AnA2KWSWHyXijV/WeYFJUxKQ1++IlIZY07j6D1IwHIX3zW67j6hly
rzUyShrRCPomvkZmNr81BDMu4Ua7TQFaze66+aK5t1dZv/Bf/obv3YBedhnq7DyI
w7OgXbNenZXKcLRm0Aw9zVx+WbdojQmmhEGMuhJXsqqePIPHcLZobzVmp5QV1lRr
CO5yygE0w2K5lC5zS0V4aBRdiWLDxlEDOyMQaJGXeF9a5S4r2aGsmxdxUSbSc+Ic
69Dj1PkIkrjipibF6HXSKksCAwEAAQ==
-----END PUBLIC KEY-----
`

	privkeyData := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
	expected := `
	[{"delta":4,"id":"PollCount","type":"counter"},
	{"value":0.5773257469936748,"id":"RandomValue","type":"gauge"},
	{"value":942416,"id":"Alloc","type":"gauge"},
	{"value":2990,"id":"BuckHashSys","type":"gauge"},{"value":121,"id":"Frees","type":"gauge"},
	{"value":0,"id":"GCCPUFraction","type":"gauge"},{"value":1811440,"id":"GCSys","type":"gauge"},
	{"value":942416,"id":"HeapAlloc","type":"gauge"},{"value":1712128,"id":"HeapIdle","type":"gauge"},
	{"value":1990656,"id":"HeapInuse","type":"gauge"},{"value":1972,"id":"HeapObjects","type":"gauge"},
	{"value":1695744,"id":"HeapReleased","type":"gauge"},{"value":3702784,"id":"HeapSys","type":"gauge"},
	{"value":0,"id":"LastGC","type":"gauge"},{"value":0,"id":"Lookups","type":"gauge"},
	{"value":9600,"id":"MCacheInuse","type":"gauge"},{"value":15600,"id":"MCacheSys","type":"gauge"},
	{"value":45120,"id":"MSpanInuse","type":"gauge"},{"value":48960,"id":"MSpanSys","type":"gauge"},
	{"value":2093,"id":"Mallocs","type":"gauge"},{"value":4194304,"id":"NextGC","type":"gauge"},
	{"value":0,"id":"NumForcedGC","type":"gauge"},{"value":0,"id":"NumGC","type":"gauge"},
	{"value":1093698,"id":"OtherSys","type":"gauge"},{"value":0,"id":"PauseTotalNs","type":"gauge"},
	{"value":491520,"id":"StackInuse","type":"gauge"},{"value":491520,"id":"StackSys","type":"gauge"},
	{"value":7166992,"id":"Sys","type":"gauge"},{"value":942416,"id":"TotalAlloc","type":"gauge"},
	{"value":16544485376,"id":"TotalMemory","type":"gauge"},
	{"value":4914249728,"id":"FreeMemory","type":"gauge"},
	{"value":0.9999999999768078,"id":"CPUutilization1","type":"gauge"},
	{"value":0,"id":"CPUutilization2","type":"gauge"},
	{"value":0,"id":"CPUutilization3","type":"gauge"},
	{"value":0,"id":"CPUutilization4","type":"gauge"},
	{"value":1.010101010077347,"id":"CPUutilization5","type":"gauge"},
	{"value":2.9702970297221394,"id":"CPUutilization6","type":"gauge"},
	{"value":7.999999999992724,"id":"CPUutilization7","type":"gauge"},
	{"value":1.9801980198006428,"id":"CPUutilization8","type":"gauge"}]
`

	pubkeyFile := "/tmp/public.pem"
	f, err := os.OpenFile(pubkeyFile, os.O_CREATE|os.O_WRONLY, 0777)
	require.NoError(t, err)
	defer os.Remove(pubkeyFile)

	f.Write([]byte(pubkeyData))
	f.Close()

	privkeyFile := "/tmp/private.pem"
	f, err = os.OpenFile(privkeyFile, os.O_CREATE|os.O_WRONLY, 0777)
	require.NoError(t, err)
	defer os.Remove(privkeyFile)

	f.Write([]byte(privkeyData))
	f.Close()

	pubKey, err := NewPublicKey(pubkeyFile)
	require.NoError(t, err)

	enc, err := Encrypt(pubKey, []byte(expected))
	require.NoError(t, err)

	privKey, err := NewPrivateKey(privkeyFile)
	require.NoError(t, err)

	dec, err := Decrypt(privKey, enc)
	require.NoError(t, err)
	assert.Equal(t, expected, string(dec))
}
