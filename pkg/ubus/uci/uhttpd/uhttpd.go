/*
Copyright 2025 Daimonas Labs.

Licensed under the GNU General Public License, Version 3 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package uhttpd

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

const (
	// the name of this config
	Config = "uhttpd"

	// these are static values for the uci.StaticSectionOptions.Type field
	UHTTPd = "uhttpd"
	Cert   = "cert"
)

var (
	Sections []string
)

func init() {
	Sections = []string{Cert, UHTTPd}
}

type CertSection struct {
	uci.StaticSectionOptions `json:",inline"`
	CertSectionOptions       `json:",inline"`
}

func (in *CertSection) DeepCopyInto(out *CertSection) {
	*out = *in
}

type CertSectionOptions struct {
	// Size of the generated RSA key in bits.
	Bits *uci.Int `json:"bits,omitempty"`
	// Common name covered by the certificate.
	CommonName *string `json:"commonname,omitempty"`
	// ISO country code of the certificate issuer.
	Country *string `json:"country,omitempty"`
	// Validity time of the generated certificates in days.
	Days    *uci.Int `json:"days,omitempty"`
	ECCurve *string  `json:"ec_curve,omitempty"`
	KeyType *string  `json:"key_type,omitempty"`
	// Location/city of the certificate issuer.
	Location *string `json:"location,omitempty"`
	// Organization name covered by the certificate.
	Organization *string `json:"organization,omitempty"`
	// State of the certificate issuer.
	State *string `json:"state,omitempty"`
}

func (CertSectionOptions) IsConfigSectionOptions() {}

type UHTTPdSection struct {
	uci.StaticSectionOptions `json:",inline"`
	UHTTPdSectionOptions     `json:",inline"`
}

func (in *UHTTPdSection) DeepCopyInto(out *UHTTPdSection) {
	*out = *in
}

type UHTTPdSectionOptions struct {
	// ASN.1/DER or PEM certificate used to serve HTTPS connections.
	Cert *string `json:"cert,omitempty"`
	// Defines the prefix for CGI scripts, relative to the document root.
	CGIPrefix *string `json:"cgi_prefix,omitempty"`
	// Config file in Busybox httpd format for additional settings.
	Config *string `json:"config,omitempty"`
	// Virtual URL of file or CGI script to handle 404 request.
	ErrorPage *string `json:"error_page,omitempty"`
	// Defines the server document root.
	Home *string `json:"home,omitempty"`
	// Connection reuse: HTTP keepalive.
	HTTPKeepAlive *uci.Int `json:"http_keepalive,omitempty"`
	// Index file(s) to use for directories.
	IndexFile *string `json:"index_file,omitempty"`
	// Index file to use for directories (alternative form).
	IndexPage *string `json:"index_page,omitempty"`
	// ASN.1/DER or PEM private key used to serve HTTPS connections.
	Key *string `json:"key,omitempty"`
	// Specifies the ports and addresses to listen on for plain HTTP access.
	ListenHTTP *uci.List `json:"listen_http,omitempty"`
	// Specifies the ports and addresses to listen on for encrypted HTTPS access.
	ListenHTTPS *uci.List `json:"listen_https,omitempty"`
	// Lua handler script used to initialize the Lua runtime on server start.
	LuaHandler *string `json:"lua_handler,omitempty"`
	// Defines the prefix for dispatching requests to the embedded Lua interpreter.
	LuaPrefix *uci.List `json:"lua_prefix,omitempty"`
	// Maximum number of concurrent connections.
	MaxConnections *uci.Int `json:"max_connections,omitempty"`
	// Maximum number of concurrent requests.
	MaxRequests *uci.Int `json:"max_requests,omitempty"`
	// Maximum wait time for network activity in seconds.
	NetworkTimeout *uci.Int `json:"network_timeout,omitempty"`
	// Do not generate directory listings if enabled.
	NoDirlists *uci.Bool `json:"no_dirlists,omitempty"`
	// Do not follow symbolic links if enabled.
	NoSymlinks *uci.Bool `json:"no_symlinks,omitempty"`
	// Basic authentication realm when prompting the client for credentials.
	Realm         *string   `json:"realm,omitempty"`
	RedirectHTTPS *uci.Bool `json:"redirect_https,omitempty"`
	// Reject requests from RFC1918 IPs to public IPs (DNS rebinding protection).
	RFC1918Filter *uci.Bool `json:"rfc1918_filter,omitempty"`
	// Maximum wait time for CGI or Lua requests in seconds.
	ScriptTimeout *uci.Int `json:"script_timeout,omitempty"`
	TCPKeepAlive  *uci.Int `json:"tcp_keepalive,omitempty"`
	// Enable CORS HTTP headers on JSON-RPC API.
	UbusCors *uci.Bool `json:"ubus_cors,omitempty"`
	// Do not authenticate JSON-RPC requests against UBUS session API.
	UbusNoauth *uci.Bool `json:"ubus_noauth,omitempty"`
	// URL prefix for UBUS via JSON-RPC handler.
	UbusPrefix *string `json:"ubus_prefix,omitempty"`
	// Override ubus socket path.
	UbusSocket *string `json:"ubus_socket,omitempty"`
}

func (UHTTPdSectionOptions) IsConfigSectionOptions() {}
