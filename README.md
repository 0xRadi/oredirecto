# oredirecto
a simple open redirect scanner in Go, that takes URLs from stdin

# Installation
```bash
go install github.com/0xRadi/oredirecto@latest
```

# Usage
```bash
cat open_redirect_urls.txt | oredirecto
```

```bash
echo https://httpbin.domain.io/redirect-to\?url\=sample | oredirecto
[Found] [Header] [http://payload.com] https://httpbin.domain.io/redirect-to?url=http%3A%2F%2Fpayload.com
```


