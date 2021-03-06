package tutamen
import "testing"

const (
	COLLECTION = "ebcdb067-469d-44af-b52f-1925e68645b9"
	SECRET     = "3828262f-3f0b-490f-bab3-399efe5897ab"
)

func TestConfig(t *testing.T) {

	cfg, err := GetConfig()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log("%+v\n", cfg)
}

func TestAll(t *testing.T) {

	cfg, err := GetConfig()
	if err != nil {
		t.Fatal(err.Error())
	}

	cli, err := NewClientV1(cfg.CertPath, cfg.KeyPath, cfg.ACUrl, cfg.SSUrl)
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

func TestSuperEasy(t *testing.T) {

	secret, err := GetSecretSuperEasy(COLLECTION, SECRET)
	if err != nil {
		t.Fatal("Error getting secret: " + err.Error())
	}

	if secret != "open" {
		t.Fatalf("Got secret '%s', expected 'open'", secret)
	}
}
