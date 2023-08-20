<h1 align="center">
    dnsaudit
  <br>
</h1>

<h4 align="center">A command-line utility for auditing DNS configuration using Zonemaster API</h4>


<p align="center">
  <a href="#install">🏗️ Install</a>
  <a href="#usage">⛏️ Usage</a>
  <br>
</p>


![dnsaudit](https://github.com/devanshbatham/dnsaudit/blob/main/static/dnsaudit.png?raw=true)

# Install
To install dnsaudit, follow these steps:


```
go install github.com/devanshbatham/dnsaudit@latest
```

# Usage
dnsaudit is a tool for auditing DNS configurations. Here are some examples of how to use the tool:

- Audit a domain's DNS configuration:
  ```sh
  dnsaudit -domain example.com
  ```

- Update the audit results for a domain:
  ```sh
  dnsaudit -domain example.com -update
  ```

Here are the available command-line flags:

| Flag        | Description                                                        | Example                    |
|-------------|--------------------------------------------------------------------|----------------------------|
| `-domain`   | Specify the domain name to audit.                                  | `dnsaudit -domain example.com` |
| `-update`   | Update the audit results for the specified domain.                | `dnsaudit -domain example.com -update` |

The tool will provide information about the DNS configuration of the specified domain, including issues and warnings.

