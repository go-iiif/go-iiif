package server

// https://github.com/FiloSottile/mkcert/issues/45

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const MKCERT string = "mkcert"

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "mkcert", NewMkCertServer)
}

func NewMkCertServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	root := q.Get("root")

	if root == "" {

		root = os.TempDir()

	} else {

		abs_root, err := filepath.Abs(root)

		if err != nil {
			return nil, err
		}

		info, err := os.Stat(abs_root)

		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			return nil, errors.New("Invalid root")
		}

		root = abs_root
	}

	server_uri := fmt.Sprintf("https://%s", u.Host)
	server_u, err := url.Parse(server_uri)

	if err != nil {
		return nil, err
	}

	tls_cert, tls_key, err := mkCert(server_u, root)

	if err != nil {
		return nil, err
	}

	server_params := url.Values{}
	server_params.Set("cert", tls_cert)
	server_params.Set("key", tls_key)

	server_u.RawQuery = server_params.Encode()

	https_uri := server_u.String()
	return NewServer(ctx, https_uri)
}

func mkCert(u *url.URL, root string) (string, string, error) {

	err := mkCertInstall()

	if err != nil {
		return "", "", err
	}

	parts := strings.Split(u.Host, ":")
	host := parts[0]

	cert_fname := fmt.Sprintf("%s-cert.pem", host)
	key_fname := fmt.Sprintf("%s-key.pem", host)

	cert_path := filepath.Join(root, cert_fname)
	key_path := filepath.Join(root, key_fname)

	args := []string{
		"-cert-file",
		cert_path,
		"-key-file",
		key_path,
		host,
	}

	cmd := exec.Command(MKCERT, args...)
	err = cmd.Run()

	if err != nil {
		return "", "", err
	}

	return cert_path, key_path, nil
}

func mkCertInstall() error {

	// unfortunately there is no way from the CLI tool to check whether
	// mkcert is installed and some of the built-in methods for testing
	// state (in mkcert.go) are marked as private so there's no way to
	// access them... (20200420/thisisaaronland)

	log.Println("Checking whether mkcert is installed. If it is not you may be prompted for your password (in order to install certificate files")

	cmd := exec.Command(MKCERT, "-install")
	return cmd.Run()
}
