package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rwxrob/klogin/internal/clusters"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// ReqOIDCPass grants an authentication token suitable for logging
// into Kubernetes by making an OIDC Resource Owner Password Flow
// request to the OIDC issuer URL (iurl). The token is returned along
// with the other response data as a map[string]any.
//
// Issuer URL (iurl) is the OIDC URL at which the discovery URL
// .well-known/openid-configuration endpoint will be
// queried to discover the token URL and other data.
//
// Client ID (cid) is passed as client_id to issuer URL (iurl)
// identifiying an OIDC supporting application.
//
// Client secret (csec) is only required if the OIDC client
// authentication type is set to "confidential". When set to "public"
// the ClientSecret is required even though---depending on the auth
// flow---it is not necessarily a "secret" and can be safely embedded in
// application source code.
func ReqOIDCPass(user, pass, iurl, cid, csec string) (map[string]any, error) {

	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("client_id", cid)
	params.Add("client_secret", csec)
	params.Add("scope", "openid")
	params.Add("username", user)
	params.Add("password", pass)
	preader := strings.NewReader(params.Encode())

	client := &http.Client{}
	ctx := oidc.ClientContext(context.Background(), client)
	prov, err := oidc.NewProvider(ctx, iurl)
	if err != nil {
		return nil, err
	}
	tokurl := prov.Endpoint().TokenURL

	req, err := http.NewRequest(`POST`, tokurl, preader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200: // do nothing
	case 403, 401:
		return nil, fmt.Errorf("%v authorization denied for server for %q", resp.StatusCode, user)
	default:
		return nil, fmt.Errorf("%v unexpected status code from server for %q", resp.StatusCode, user)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OIDC password token response for %q", user)
	}

	data := new(map[string]any)
	if err := json.Unmarshal(buf, data); err != nil {
		return nil, err
	}

	return *data, nil
}

// ParseTarget takes either an unqualified cluster name or a qualified
// name combined with a username prefix separated by an at sign (ex:
// user@prod).
func ParseTarget(in string) (uname, clname string) {
	it := strings.SplitN(in, `@`, 2)
	switch len(it) {
	case 2:
		clname = it[1]
		uname = it[0]
	case 1:
		clname = it[0]
	}
	return
}

// LoginROPC logs a user inferred from conf.Contexts[conf.CurrentContext]
// .AuthInfo into that cluster using the clusters.Map to provide the
// authentication data required. The CurrentContext must therefore
// always be set and the AuthInfo to which it points must always be of
// the form user@cluster. An error is returned if any of these
// requirements is not met.
//
// The OIDC/OAuth2 Resource Owner Password Credentials flow
// (grant_type=password, known as "direct access" in Keycloak) is used
// to authenticate and grant a JWT token. Returns an error if that
// cluster is not in clusters.Names.  Saves the id_token (not
// access_token, per Kubernetes official documentation) from the JWT
// returned by the OIDC issuer into the users/credentials/AuthInfo
// section of the appropriate KUBECONFIG file by passing an official
// api.Config struct to clientcmd.ModifyConf (the only official
// supported way to update configuration files consistent with kubectl
// behavior). The reserved context (named after the reserved cluster
// name) is always created, updated, and persisted as well.
func LoginROPC(conf *api.Config, pass string) error {

	ctx, has := conf.Contexts[conf.CurrentContext]
	if !has {
		return fmt.Errorf(`unable to infer target user and cluster`)
	}

	uname, clname := ParseTarget(ctx.AuthInfo)
	cl, has := clusters.Map[clname]
	if !has {
		return fmt.Errorf(`unsupported cluster: %v`, clname)
	}

	// first ensure we get a successful login before persisting anything

	grant, err := ReqOIDCPass(
		uname, pass, cl.OIDCIssuerURL, cl.ClientID, cl.ClientSecret,
	)
	if err != nil {
		log.Fatal(err)
	}
	token, isstring := grant[`id_token`].(string)
	if !isstring {
		log.Fatal(`id_token not found or is not a string`)
	}

	// add/update the cluster api.Config entry

	cluster := api.NewCluster()
	cluster.Server = cl.APIServerURL
	cluster.CertificateAuthorityData = cl.CA
	conf.Clusters[cl.Name] = cluster

	// update/add the user@cluster user/credential/AuthInfo entry

	delete(conf.AuthInfos, cl.Name) // cleanup old, unqualified entries
	authinfo := api.NewAuthInfo()
	authinfo.Token = token
	authinfoid := strings.Join([]string{uname, cl.Name}, `@`)
	conf.AuthInfos[authinfoid] = authinfo

	// save modified configuration

	o := clientcmd.NewDefaultPathOptions()
	return clientcmd.ModifyConfig(o, *conf, true)
}

// LoginAuth saves the token into the kubeconfig file using the same
// methods as the kubectl program itself. The token should be obtained
// by interactively prompting the user to paste or otherwise provide the
// token obtained from a web browser that has completed the standard
// Oauth2 standard user authentication flow. The conf.CurrentContext
// must be set and the AuthInfo to which it points must always be of the
// form user@cluster. An error is returned if any of these requirements
// is not met.
func LoginAuth(conf *api.Config, token string) error {

	ctx, has := conf.Contexts[conf.CurrentContext]
	if !has {
		return fmt.Errorf(`unable to infer target user and cluster`)
	}

	uname, clname := ParseTarget(ctx.AuthInfo)
	cl, has := clusters.Map[clname]
	if !has {
		return fmt.Errorf(`unsupported cluster: %v`, clname)
	}

	// add/update the cluster api.Config entry

	cluster := api.NewCluster()
	cluster.Server = cl.APIServerURL
	cluster.CertificateAuthorityData = cl.CA
	conf.Clusters[cl.Name] = cluster

	// update/add the user@cluster user/credential/AuthInfo entry

	delete(conf.AuthInfos, cl.Name) // cleanup old, unqualified entries
	authinfo := api.NewAuthInfo()
	authinfo.Token = token
	authinfoid := strings.Join([]string{uname, cl.Name}, `@`)
	conf.AuthInfos[authinfoid] = authinfo

	// save modified configuration

	o := clientcmd.NewDefaultPathOptions()
	return clientcmd.ModifyConfig(o, *conf, true)
}
