package tutamen
import "errors"

func GetSecretSuperEasy(collection, secret string) (string, error) {

	dir, err := GetConfigDir()
	if err != nil {
		return "", errors.New("Unable to find config dir: " + err.Error())
	}

	cfg, err := GetConfig(dir)
	if err != nil {
		return "", errors.New("Error parsing config file(s): " + err.Error())
	}

	cli, err := NewClientV1(cfg.CertPath, cfg.KeyPath, cfg.ACUrl, cfg.SSUrl)
	if err != nil {
		return "", errors.New("Error creating client: " + err.Error())
	}

	return cli.GetSecretEasy(collection, secret)
}

func (cli *client) GetSecretEasy(collection, secret string) (string, error) {


	auth, err := cli.GetAuthorization(collection)
	if err != nil {
		return "", errors.New("Error getting authorization: " + err.Error())
	}

	token, err := cli.WaitForToken(auth)
	if err != nil {
		return "", errors.New("Error getting token: " + err.Error())
	}

	val, err := cli.GetSecret(collection, secret, token)
	if err != nil {
		return "", errors.New("Error getting secret: " + err.Error())
	}

	return val, nil
}
