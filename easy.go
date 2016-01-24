package tutamen
import "errors"

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
