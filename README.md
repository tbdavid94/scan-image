Targets

- Container Image
- Filesystem
- Git Repository (remote)
- Virtual Machine Image
- Kubernetes
- AWS

Scanners (what scanner can find there):

- OS packages and software dependencies in use (SBOM)
- Known vulnerabilities (CVEs)
- IaC issues and misconfigurations
- Sensitive information and secrets
- Software licenses

```shell
curl -X POST -H "Content-Type: application/json" -d '{"image_name": "your_image_name"}' http://localhost:8080/enqueue
```
