---
layout: default
type: api
navgroup: docs
shortname: grpc/credentials
title: grpc/credentials
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: grpc/credentials
---
# credentials
--
    import "github.com/antha-lang/antha/internal/google.golang.org/grpc/credentials"

Package credentials implements various credentials supported by gRPC library,
which encapsulate all the state needed by a client to authenticate with a server
and make various assertions, e.g., about the client's identity, role, or whether
it is authorized to make a particular call.

## Usage

#### func  NewContext

```go
func NewContext(ctx context.Context, authInfo AuthInfo) context.Context
```
NewContext creates a new context with authInfo attached.

#### type AuthInfo

```go
type AuthInfo interface {
	AuthType() string
}
```

AuthInfo defines the common interface for the auth information the users are
interested in.

#### func  FromContext

```go
func FromContext(ctx context.Context) (authInfo AuthInfo, ok bool)
```
FromContext returns the authInfo in ctx if it exists.

#### type Credentials

```go
type Credentials interface {
	// GetRequestMetadata gets the current request metadata, refreshing
	// tokens if required. This should be called by the transport layer on
	// each request, and the data should be populated in headers or other
	// context. uri is the URI of the entry point for the request. When
	// supported by the underlying implementation, ctx can be used for
	// timeout and cancellation.
	// TODO(zhaoq): Define the set of the qualified keys instead of leaving
	// it as an arbitrary string.
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	// RequireTransportSecurity indicates whether the credentails requires
	// transport security.
	RequireTransportSecurity() bool
}
```

Credentials defines the common interface all supported credentials must
implement.

#### type ProtocolInfo

```go
type ProtocolInfo struct {
	// ProtocolVersion is the gRPC wire protocol version.
	ProtocolVersion string
	// SecurityProtocol is the security protocol in use.
	SecurityProtocol string
	// SecurityVersion is the security protocol version.
	SecurityVersion string
}
```

ProtocolInfo provides information regarding the gRPC wire protocol version,
security protocol, security protocol version in use, etc.

#### type TLSInfo

```go
type TLSInfo struct {
	State tls.ConnectionState
}
```

TLSInfo contains the auth information for a TLS authenticated connection. It
implements the AuthInfo interface.

#### func (TLSInfo) AuthType

```go
func (t TLSInfo) AuthType() string
```

#### type TransportAuthenticator

```go
type TransportAuthenticator interface {
	// ClientHandshake does the authentication handshake specified by the corresponding
	// authentication protocol on rawConn for clients. It returns the authenticated
	// connection and the corresponding auth information about the connection.
	ClientHandshake(addr string, rawConn net.Conn, timeout time.Duration) (net.Conn, AuthInfo, error)
	// ServerHandshake does the authentication handshake for servers. It returns
	// the authenticated connection and the corresponding auth information about
	// the connection.
	ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error)
	// Info provides the ProtocolInfo of this TransportAuthenticator.
	Info() ProtocolInfo
	Credentials
}
```

TransportAuthenticator defines the common interface for all the live gRPC wire
protocols and supported transport security protocols (e.g., TLS, SSL).

#### func  NewClientTLSFromCert

```go
func NewClientTLSFromCert(cp *x509.CertPool, serverName string) TransportAuthenticator
```
NewClientTLSFromCert constructs a TLS from the input certificate for client.

#### func  NewClientTLSFromFile

```go
func NewClientTLSFromFile(certFile, serverName string) (TransportAuthenticator, error)
```
NewClientTLSFromFile constructs a TLS from the input certificate file for
client.

#### func  NewServerTLSFromCert

```go
func NewServerTLSFromCert(cert *tls.Certificate) TransportAuthenticator
```
NewServerTLSFromCert constructs a TLS from the input certificate for server.

#### func  NewServerTLSFromFile

```go
func NewServerTLSFromFile(certFile, keyFile string) (TransportAuthenticator, error)
```
NewServerTLSFromFile constructs a TLS from the input certificate file and key
file for server.

#### func  NewTLS

```go
func NewTLS(c *tls.Config) TransportAuthenticator
```
NewTLS uses c to construct a TransportAuthenticator based on TLS.
