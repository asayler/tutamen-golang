package tutamen
import "testing"

const (
	TLS_CRT    = "/tut.crt"
	TLS_KEY    = "/tut.key"
	AC_SERVER  = "ac.tutamen-test.bdr1.volaticus.net"
	SS_SERVER  = "ss.tutamen-test.bdr1.volaticus.net"
	COLLECTION = "ebcdb067-469d-44af-b52f-1925e68645b9"
	SECRET     = "3828262f-3f0b-490f-bab3-399efe5897ab"
)

func Test1(t *testing.T) {

	cli, err := NewClientV1(TLS_CRT, TLS_KEY, AC_SERVER, SS_SERVER)
	if err != nil {
		t.Fatal("Error creating client: " + err.Error())
	}

	auth, err := cli.GetAuthorization(COLLECTION)
	if err != nil {
		t.Fatal("Error getting authorization: " + err.Error())
	}

	token, err := cli.WaitForToken(auth)
	if err != nil {
		t.Fatal("Error getting token: " + err.Error())
	}

	secret, err := cli.GetSecret(COLLECTION, SECRET, token)
	if err != nil {
		t.Fatal("Error getting secret: " + err.Error())
	}

	if secret != "open" {
		t.Fatalf("Got secret '%s', expected 'open'", secret)
	}
}

func TestEasy(t *testing.T) {

	cli, err := NewClientV1(TLS_CRT, TLS_KEY, AC_SERVER, SS_SERVER)
	if err != nil {
		t.Fatal("Error creating client: " + err.Error())
	}

	secret, err := cli.GetSecretEasy(COLLECTION, SECRET)
	if err != nil {
		t.Fatal("Error getting secret: " + err.Error())
	}

	if secret != "open" {
		t.Fatalf("Got secret '%s', expected 'open'", secret)
	}
}
