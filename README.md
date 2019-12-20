# GHORGS
(= GitHub ORGanizationS) is a small tool initially used to list
GitHub projects and users belonging to SonyMobile organizational account
at github.com/SonyMobile, but is extended to work with any
organizational account.

It uses GitHub's GraphQL API to query for repositories and users belonging to
an organization and then writes the output into a csv files.

## Set up
* The project is written in go, so you have to set up go environment.
* The bare minimum is to have the latest installation of go for your environment and
  set GOPATH to the root of your go projects and GOBIN to the path where you want
  the executables to end up after `go install`.

### config

#### config.yaml
* url: URL to GitHub API, should be https://api.github.com/graphql
* token: String security token used on Github. Required GitHub scopes covered by token are:
  * user,
  * public_repo,
  * repo,
  * repo_deployment,
  * repo:status,
  * read:repo_hook,
  * read:org,
  * read:public_key,
  * read:gpg_key
* organization: Organizational account which is being analyzed
* per_page: Integer denoting the number of items listed in paged output
* time_out: Seconds until connection is abandoned

#### gql and json
* gql files are GraphQL queries used for testing in Explore mode on GitHub
* json files are GraphQL queries translated to json format required by
  GitHub's API:
```
{
    "query": "<literal GraphQL query from gql file>"
}
```

### Dependencies
The only external dependency for now is to yaml.v3.
Run `go get gopkg.in/yaml.v3` to make sure it's downloaded.

## Build
`go install`

## Run
`ghorgs [-h] [-v] [-t Token] [-o Organization]`

where:

* h = prints help
* v = enables verbose prints to stdout
* t = overrides config token
* o = overrides config organization

## LICENSE
Currently the tool is proprietary.
