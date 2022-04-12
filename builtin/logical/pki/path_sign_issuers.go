package pki

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssuerSignIntermediate(b *backend) *framework.Path {
	pattern := "issuers/" + framework.GenericNameRegex(issuerRefParam) + "/sign-intermediate"
	return pathIssuerSignIntermediateRaw(b, pattern)
}

func pathSignIntermediate(b *backend) *framework.Path {
	pattern := "root/sign-intermediate"
	return pathIssuerSignIntermediateRaw(b, pattern)
}

func pathIssuerSignIntermediateRaw(b *backend, pattern string) *framework.Path {
	path := &framework.Path{
		Pattern: pattern,
		Fields: map[string]*framework.FieldSchema{
			issuerRefParam: {
				Type:        framework.TypeString,
				Description: `Reference to issuer; either "default" for the configured default issuer, an identifier of an issuer, or the name assigned to the issuer.`,
				Default:     "default",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathIssuerSignIntermediate,
		},

		HelpSynopsis:    pathIssuerSignIntermediateHelpSyn,
		HelpDescription: pathIssuerSignIntermediateHelpDesc,
	}

	path.Fields = addCACommonFields(path.Fields)
	path.Fields = addCAIssueFields(path.Fields)

	path.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     "",
		Description: `PEM-format CSR to be signed.`,
	}

	path.Fields["use_csr_values"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: false,
		Description: `If true, then:
1) Subject information, including names and alternate
names, will be preserved from the CSR rather than
using values provided in the other parameters to
this path;
2) Any key usages requested in the CSR will be
added to the basic set of key usages used for CA
certs signed by this path; for instance,
the non-repudiation flag;
3) Extensions requested in the CSR will be copied
into the issued certificate.`,
	}

	return path
}

const (
	pathIssuerSignIntermediateHelpSyn  = `Issue an intermediate CA certificate based on the provided CSR.`
	pathIssuerSignIntermediateHelpDesc = `
This API endpoint allows for signing the specified CSR, adding to it a basic
constraint for IsCA=True. This allows the issued certificate to issue its own
leaf certificates.

Note that the resulting certificate is not imported as an issuer in this PKI
mount. This means that you can use the resulting certificate in another Vault
PKI mount point or to issue an external intermediate (e.g., for use with
another X.509 CA).

See the API documentation for more information about required parameters.
`
)

func pathIssuerSignSelfIssued(b *backend) *framework.Path {
	pattern := "issuers/" + framework.GenericNameRegex(issuerRefParam) + "/sign-self-issued"
	return buildPathIssuerSignSelfIssued(b, pattern)
}

func pathSignSelfIssued(b *backend) *framework.Path {
	pattern := "root/sign-self-issued"
	return buildPathIssuerSignSelfIssued(b, pattern)
}

func buildPathIssuerSignSelfIssued(b *backend, pattern string) *framework.Path {
	path := &framework.Path{
		Pattern: pattern,
		Fields: map[string]*framework.FieldSchema{
			issuerRefParam: {
				Type:        framework.TypeString,
				Description: `Reference to issuer; either "default" for the configured default issuer, an identifier of an issuer, or the name assigned to the issuer.`,
				Default:     "default",
			},
			"certificate": {
				Type:        framework.TypeString,
				Description: `PEM-format self-issued certificate to be signed.`,
			},
			"require_matching_certificate_algorithms": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: `If true, require the public key algorithm of the signer to match that of the self issued certificate.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathIssuerSignSelfIssued,
		},

		HelpSynopsis:    pathIssuerSignSelfIssuedHelpSyn,
		HelpDescription: pathIssuerSignSelfIssuedHelpDesc,
	}

	return path
}

const (
	pathIssuerSignSelfIssuedHelpSyn  = `Re-issue a self-signed certificate based on the provided certificate.`
	pathIssuerSignSelfIssuedHelpDesc = `
This API endpoint allows for signing the specified self-signed certificate,
effectively allowing cross-signing of external root CAs. This allows for an
alternative validation path, chaining back through this PKI mount. This
endpoint is also useful in a rolling-root scenario, allowing devices to trust
and validate later (or earlier) root certificates and their issued leaves.

Usually the sign-intermediate operation is preferred to this operation.

Note that this is a very privileged operation and should be extremely
restricted in terms of who is allowed to use it. All values will be taken
directly from the incoming certificate and only verification that it is
self-issued will be performed.

Configured URLs for CRLs/OCSP/etc. will be copied over and the issuer will
be this mount's CA cert. Other than that, all other values will be used
verbatim from the given certificate.

See the API documentation for more information about required parameters.
`
)