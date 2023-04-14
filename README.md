# multiwhois

`multiwhois` is a Go package for querying domain information from WHOIS servers. It uses the `github.com/domainr/whois` package to fetch the WHOIS information and parses the response into a struct containing the domain owner, expiration date, and availability status.

## Supported TLDs
We currently support:
- .com
- .info
- .be
- .nl

## Installation

To install the package, use the following command:

```
go get github.com/melvinsh/multiwhois
```

## Usage

``` go
package main

import (
	"fmt"

	"github.com/melvinsh/multiwhois"
)

func main() {
	info, err := multiwhois.QueryDomainInfo("google.com")
	if err != nil {
		panic(err)
	}

	fmt.Println("Domain:", info.Domain)
    fmt.Println("Available:", info.IsAvailable)
	fmt.Println("Expiration:", info.Expiration)
	fmt.Println("Full WHOIS response:", info.FullWhois)
}
```

The `QueryDomainInfo` function takes a domain name as its argument and returns a `*DomainInfo` struct and an error. The `*DomainInfo` struct contains the following fields:

- `Domain` (string): the domain name
- `Expiration` (time.Time): the expiration date of the domain
- `IsAvailable` (bool): whether the domain is available or not
- `FullWhois` (string): the full WHOIS response

## Testing

To run the tests, use the following command:

```
go test github.com/melvinsh/multiwhois
```

## Configuration

`multiwhois` uses YAML files to configure the regex patterns for parsing WHOIS responses. The YAML files are located in the `tlds` directory and are named after the TLD they represent (e.g. `com.yml` for the `.com` TLD).

The YAML file must contain the following fields:

- `availability` (string): a regex pattern for determining the availability of the domain
- `expiration` (string): a regex pattern for extracting the expiration date from the WHOIS response

Here is an example YAML file for the `.com` TLD:

``` yaml
expiration: 'Registry Expiry Date: ([^\r\n]+)'
expiration_format: '2006-01-02T15:04:05Z'
availability: 'No match for'
```