package tutamen

import (
	"os"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
)

type Config struct {
	ACName    string
	ACUrl     string
	SSName    string
	SSUrl     string
	Account   string
	Client    string
	CertPath  string
	KeyPath   string
}

func GetConfigDir() (string, error) {

	var dir string

	xdg_dir := os.Getenv("XDG_CONFIG_HOME")
	if xdg_dir != "" {
		dir = xdg_dir + "/tutamen"
		stat, err := os.Stat(dir)
		if err == nil && stat.IsDir() {
			return dir, nil
		}
		dir = xdg_dir + "/pytutamen_client"
		stat, err = os.Stat(dir)
		if err == nil && stat.IsDir() {
			return dir, nil
		}
	}

	home_dir := os.Getenv("HOME")
	if home_dir != "" {
		dir = home_dir + "/.config/tutamen"
		stat, err := os.Stat(dir)
		if err == nil && stat.IsDir() {
			return dir, nil
		}
		dir = home_dir + "/.config/pytutamen_client"
		stat, err = os.Stat(dir)
		if err == nil && stat.IsDir() {
			return dir, nil
		}
	}

	return "", errors.New("No configuration directory found")
}

func GetConfig(dir string) (*Config, error) {

	cfg := new(Config)

	inicfg, err := ini.Load(
		dir + "/core.conf",
		dir + "/srv_ac.conf",
		dir + "/srv_storage.conf")
	if err != nil {
		return nil, err
	}

	cfg.ACName = inicfg.Section("defaults").Key("ac_server").String()
	if cfg.ACName == "" {
		return nil, errors.New("No default AC server defined")
	}

	cfg.ACUrl = inicfg.Section(cfg.ACName).Key("url").String()
	if cfg.ACUrl == "" {
		return nil, errors.New("No URL given for AC server '"+cfg.ACName+"'")
	}

	cfg.SSName = inicfg.Section("defaults").Key("storage_server").String()
	if cfg.SSName == "" {
		return nil, errors.New("No default storage server defined")
	}

	cfg.SSUrl = inicfg.Section(cfg.SSName).Key("url").String()
	if cfg.SSUrl == "" {
		return nil, errors.New("No URL given for storage server '"+cfg.SSName+"'")
	}

	cfg.Account = inicfg.Section("defaults").Key("account").String()
	if cfg.Account == "" {
		return nil, errors.New("No default account defined")
	}

	cfg.Client = inicfg.Section("defaults").Key("client").String()
	if cfg.Client == "" {
		return nil, errors.New("No default client defined")
	}

	cfg.CertPath = fmt.Sprintf("%s/accounts/%s/clients/%s/%s_crt.pem",
		dir, cfg.Account, cfg.Client, cfg.ACName)

	cfg.KeyPath = fmt.Sprintf("%s/accounts/%s/clients/%s/key.pem",
		dir, cfg.Account, cfg.Client)

	return cfg, nil
}
