package inventory

import "os/exec"
import "bytes"
import "strings"

func readGoDeps(packagePath string) ([]Project, error) {
	var returned []Project
	// Read the names of all Go dependencies.
	deps, err := goListDeps(packagePath)
	if err != nil {
		return nil, err
	}
	goStandardPackages := loadGoStandardPackageList()
	// Iterate the package names.
	for _, dep := range deps {
		// Skip packages in the Go standard library.
		if isKnownStandardGoPackage(dep, goStandardPackages) {
			continue
		}
		// Run `go list` for package information.
		info, err := goListPackageInfo(dep)
		if err != nil {
			continue
		}
		// Skip packages in the Go standard library.
		if info.Standard {
			continue
		}
		// Try to read licensezero.json in the package's path.
		projects, err := ReadLicenseZeroJSON(info.Dir)
		if err != nil {
			continue
		}
		for _, project := range projects {
			projectID := project.Envelope.Manifest.ProjectID
			if alreadyHaveProject(returned, projectID) {
				continue
			}
			project.Type = "go"
			project.Name = info.Name
		}
	}
	return returned, nil
}

func goListDeps(packagePath string) ([]string, error) {
	list := exec.Command("go", "list", "-f", "{{ join .Deps \"\\n\" }}")
	list.Dir = packagePath
	var stdout bytes.Buffer
	list.Stdout = &stdout
	err := list.Run()
	if err != nil {
		return nil, err
	}
	deps := strings.Split(string(stdout.Bytes()), "\n")
	// Remove empty string after final newline.
	if len(deps) != 0 {
		deps = deps[0 : len(deps)-1]
	}
	return deps, nil
}

type goPackageInfo struct {
	Name       string
	Dir        string
	ImportPath string
	Standard   bool
}

func goListPackageInfo(name string) (*goPackageInfo, error) {
	list := exec.Command("go", "list", "-f", "{{ .Name }}\n{{ .Dir }}\n{{ .ImportPath }}\n{{ .Standard }}\n", name)
	var stdout bytes.Buffer
	list.Stdout = &stdout
	err := list.Run()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(stdout.Bytes()), "\n")
	return &goPackageInfo{
		Name:       lines[0],
		Dir:        lines[1],
		ImportPath: lines[2],
		Standard:   lines[3] == "true",
	}, nil
}

func loadGoStandardPackageList() []string {
	return []string{
		"archive",
		"archive/tar",
		"archive/zip",
		"bufio",
		"builtin",
		"bytes",
		"compress",
		"compress/bzip2",
		"compress/flate",
		"compress/gzip",
		"compress/lzw",
		"compress/zlib",
		"container",
		"container/heap",
		"container/list",
		"container/ring",
		"context",
		"crypto",
		"crypto/aes",
		"crypto/cipher",
		"crypto/des",
		"crypto/dsa",
		"crypto/ecdsa",
		"crypto/elliptic",
		"crypto/hmac",
		"crypto/md5",
		"crypto/rand",
		"crypto/rc4",
		"crypto/rsa",
		"crypto/sha1",
		"crypto/sha256",
		"crypto/sha512",
		"crypto/subtle",
		"crypto/tls",
		"crypto/x509",
		"crypto/x509/pkix",
		"database",
		"database/sql",
		"database/sql/driver",
		"debug",
		"debug/dwarf",
		"debug/elf",
		"debug/gosym",
		"debug/macho",
		"debug/pe",
		"debug/plan9obj",
		"encoding",
		"encoding/ascii85",
		"encoding/asn1",
		"encoding/base32",
		"encoding/base64",
		"encoding/binary",
		"encoding/csv",
		"encoding/gob",
		"encoding/hex",
		"encoding/json",
		"encoding/pem",
		"encoding/xml",
		"errors",
		"expvar",
		"flag",
		"fmt",
		"go",
		"go/ast",
		"go/build",
		"go/constant",
		"go/doc",
		"go/format",
		"go/importer",
		"go/parser",
		"go/printer",
		"go/scanner",
		"go/token",
		"go/types",
		"hash",
		"hash/adler32",
		"hash/crc32",
		"hash/crc64",
		"hash/fnv",
		"html",
		"html/template",
		"image",
		"image/color",
		"image/color/palette",
		"image/draw",
		"image/gif",
		"image/jpeg",
		"image/png",
		"index",
		"index/suffixarray",
		"io",
		"io/ioutil",
		"log",
		"log/syslog",
		"math",
		"math/big",
		"math/bits",
		"math/cmplx",
		"math/rand",
		"mime",
		"mime/multipart",
		"mime/quotedprintable",
		"net",
		"net/http",
		"net/http/internal",
		"net/http/cgi",
		"net/http/cookiejar",
		"net/http/fcgi",
		"net/http/httptest",
		"net/http/httptrace",
		"net/http/httputil",
		"net/http/pprof",
		"net/mail",
		"net/rpc",
		"net/rpc/jsonrpc",
		"net/smtp",
		"net/textproto",
		"net/url",
		"os",
		"os/exec",
		"os/signal",
		"os/user",
		"path",
		"path/filepath",
		"plugin",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/debug",
		"runtime/msan",
		"runtime/pprof",
		"runtime/race",
		"runtime/trace",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"testing",
		"testing/iotest",
		"testing/quick",
		"text",
		"text/scanner",
		"text/tabwriter",
		"text/template",
		"text/template/parse",
		"time",
		"unicode",
		"unicode/utf16",
		"unicode/utf8",
		"unsafe",
	}
}

func isKnownStandardGoPackage(name string, list []string) bool {
	if strings.HasPrefix(name, "internal/") {
		return true
	}
	if strings.HasPrefix(name, "runtime/internal/") {
		return true
	}
	if strings.HasPrefix(name, "crypto/internal/") {
		return true
	}
	for _, item := range list {
		if item == name {
			return true
		}
	}
	return false
}
