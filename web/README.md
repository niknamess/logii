### Usage

```
$ go run main.go --help
usage: main [<flags>] [<dir>...]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=8080          Port number to host the server
  -r, --restrict           Enforce PAM authentication (single level)
  -a, --acl=ACL            enable Access Control List with users in the provided file
  -t, --cron="0h"          configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]
  -s, --secure             Run Server with TLS
  -c, --cert="server.crt"  Server Certificate
  -k, --key="server.key"   Server Key File



#### Examples
```
./webtail
```
This will run the server on port `8080` and look for files in the current Directory
```
./webtail --port 15000 /var/log/tomcat /tmp/
```
This will run the server on port 15000 and recursively look for files in `/var/log/tomcat` and `/tmp` directories (provided the permissions)

```
./webtail /var/log/tomcat /tmp/ --restrict
```
This will add an authentication layer over it. Once you navigate to the home page, it will redirect to the `/login` page and ask for username and password. Since this is supposed to be as generic as possible, hence it uses PAM authentication to authenticate the user. You need to provide the credentials that you would use to login to the host on which the server is hosted. Right now it would only authenticate via PAM if only a single step is required.



./webtail /var/log/tomcat --restrict --cron 5h
```
This will make the server re-index files every 5 hours. And the **New** files will be served only after a page refresh, once it has been indexed.




### TLS Server

To run the server with TLS enabled use `--secure` flag. It will search for `server.crt` and `server.key` files in the current directory, if not will fail.

```
./webtail /var/log/tomcat --secure --cert /path/to/server.crt --key /path/to/server.key --port 8443 --restrict
```

If you are running it on `--restrict` mode then it is recommended to use `--secure` flag as well to protect the login credentials on the wire.

Server Accepts connections only on `TLSv1.1` and above    
List of CiphersSuites supported by the server:
```
tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
tls.TLS_RSA_WITH_AES_256_GCM_SHA384
tls.TLS_RSA_WITH_AES_256_CBC_SHA
```

to connect to a particular cipher you can use       
**TLSv1.1**
```
echo "" | openssl s_client -connect localhost:8443 -cipher AES256-SHA -tls1_1 -quiet 2>/dev/null
```
**TLSv1.2**
```
 echo "" | openssl s_client -connect localhost:8443 -cipher ECDHE-RSA-AES128-SHA256 -tls1_2 -quiet 2>/dev/null
```

**Note**: It requires the private key to be in the plain text format. i.e. it should not be passphrase protected. Can be done via
```
openssl rsa -in [file1.key] -out [file2.key]
```

### Cron

Cron option supports only 2 formats of time: `days` and `hours`.      
You can say something like `5h` or `1d` or `100h` or `4d`. Zero prefixed time intervals are not allowed and will fail.    
**Note**: By default cron is not enabled and will not re-index files.
