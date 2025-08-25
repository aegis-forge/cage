<p align="center">
  <img width="100" src="assets/branding/logo.svg" alt="cage logo"> <br><br>
  <img src="https://img.shields.io/badge/go-^v1.23.0-blue" alt="Go version">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

# cage-go

## Installing

To use CAGE, first install the latest version of it using `go get`.

```sh
go get -u github.com/aegis-forge/cage@latest
```

Then, you can include it in your application as such:

```go
import "github.com/aegis-forge/cage"
```

## Usage

To check if a software package is vulnerable or not, you can use the following working example. In this case, we are checking if the GitHub Action `tj-actions/branch-names@v7` is vulnerable or not.

```go
package main

import (
	"log"
	"time"
	
	cage "github.com/aegis-forge/cage-go"
)

func main() {
	advisories := cage.Github{}
	sources := []cage.Source{advisories}
	
	semver, err := cage.NewSemver("7")
	
	if err != nil {
		log.Fatal(err)
	}
	
	packg, err := cage.NewPackage("tj-actions", "branch-names", time.Now(), semver)
	
	if err != nil {
		log.Fatal(err)
	}
	
	vulns, err := packg.IsVulnerable(sources)
	
	if err != nil {
		log.Fatal(err)
	}
	
	parsed, err := json.MarshalIndent(vulns, "", "  ")
	
	log.Print(string(parsed))
}
```

> [!NOTE]
> If you need to check a large number of Actions, you can add your own personal GitHub token. To do so, after initializing `cage.Github{}`, simply call the `SetToken()` method on the source and pass it your token

<details>
    <summary>Output</summary>

By running the code above, you will get the following JSON-formatted output (as of `2025-08-25 11:38:58`):

```json
[
  {
    "cve": "CVE-2025-54416",
    "cwes": [
      "CWE-77"
    ],
    "cvss": 9.1,
    "published": "2025-07-25T19:28:22Z",
    "vulnerable_ranges": [
      {
        "start": "v0.0.0",
        "end": "v8.2.1",
        "left": true,
        "right": true
      }
    ],
    "patched_ranges": [
      {
        "start": "v9.0.0",
        "end": "",
        "left": true,
        "right": false
      }
    ]
  },
  {
    "cve": "CVE-2023-49291",
    "cwes": [
      "CWE-20"
    ],
    "cvss": 9.3,
    "published": "2023-12-05T23:30:10Z",
    "vulnerable_ranges": [
      {
        "start": "v0.0.0",
        "end": "v7.0.7",
        "left": true,
        "right": false
      }
    ],
    "patched_ranges": [
      {
        "start": "v7.0.7",
        "end": "",
        "left": true,
        "right": false
      }
    ]
  }
]
```
</details>

## Extend

CAGE can be extended by adding custom `Source`s. To do so, the new struct must include the methods included in the `Source` interface.

```go
type Source interface {
	GetVulnerabilities(Package) ([]Vulnerability, error)
	CompareVulnerabilities([]Vulnerability, Package) ([]Vulnerability, error)
}
```

In the case of the `GetVulnerabilities()` method, it will be invoked by the `Package` in its method `Package.IsVulnerable()`. Given the name of the package passed, it retrieves all the vulnerabilities for that package from the source's database.

On the other hand, `CompareVulnerabilities()` is basically a custom rule to compare the `Vulnerability` objects to the `Package.version` field. It returns all detected vulnerabilities.

## Vulnerability Sources

- [X] [GitHub Advisory Database](https://github.com/advisories): `cage.Github{}`
- [ ] [NIST National Vulnerability Database](https://nvd.nist.gov/)
