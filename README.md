[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_go-cloud-k8s-user-group&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_go-cloud-k8s-user-group)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_go-cloud-k8s-user-group&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_go-cloud-k8s-user-group)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_go-cloud-k8s-user-group&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_go-cloud-k8s-user-group)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_go-cloud-k8s-user-group&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_go-cloud-k8s-user-group)

# go-cloud-k8s-user-group
go-cloud-k8s-user-group is a user and group microservice written in Golang 
that allows authentication with 2FA and sends a JWT. 

_This project showcases how to build a container image with nerdctl, in a secured way (scan of CVE done with Trivy) and how to deploy it on Kubernetes_

## Dependencies
[Echo: high performance, extensible, minimalist Go web framework](https://echo.labstack.com/)

[deepmap/oapi-codegen: OpenAPI Client and Server Code Generator](https://github.com/deepmap/oapi-codegen)

[pgx: PostgreSQL Driver and Toolkit](https://pkg.go.dev/github.com/jackc/pgx)

[Json Web Token for Go (RFC 7519)](https://github.com/cristalhq/jwt)


## Project Layout and conventions
This project uses the Standard Go Project Layout : https://github.com/golang-standards/project-layout
