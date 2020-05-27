# GHORGS
(= GitHub ORGanizationS) is a small tool initially used to list
GitHub projects and users belonging to SonyMobile organizational account
at github.com/SonyMobile, but is extended to work with any
organizational account.

It uses GitHub's GraphQL API (v4) to query for repositories and users belonging to
an organization and then writes the output into a csv files.
It uses GitHub's REST-like API (v3) to delete repositories and remove users.

## Set up
* The project is written in go, so you have to set up go environment.
* The bare minimum is to have the latest installation of go for your environment and
  set GOPATH to the root of your go projects and GOBIN to the path where you want
  the executables to end up after `go install`.

### config

#### config.yaml
* url: URL to GitHub API, should be https://api.github.com/
  * That way both v3 and v4 API are internally differentiated.
* user: username required for `git clone` in `ghorgs archive` command
* token: String security token used on Github. Required GitHub scopes covered by token are:
  * user,
  * public_repo,
  * delete_repo
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
  GitHub's v4 API:
```
{
    "query": "<literal GraphQL query from gql file>"
}
```

### Dependencies
Current dependencies are to `cobra` (https://github.com/spf13/cobra) and
`viper` (https://github.com/spf13/viper).
Make sure you run:
`go get github.com/spf13/cobra` and `go get github.com/spf13/viper`
respectively to get the dependcies.
Also, upon modification of module, re-init the module
with `go mod init`, which will recreate the go.mod and go.sum files.

## Build
`go install`

## Run
Usage:
```
  ghorgs [command]

  Available Commands:
    archive     Archive GitHub repositories according to given criteria.
    dump        Dumps the requested entities into a csv file.
    help        Help about any command
    remove      Remove GitHub users according to given criteria.
    version     prints version of ghorgs tool

  Flags:
    -d, --dry-run               Perform a dry run of the command without actually executing it in the end.
    -h, --help                  help for ghorgs
    -o, --organization string   Organizational account on GitHub analyzed.
    -t, --token string          Security token used on Github. Overrides the token from configuration file.
                                Required GitHub scopes covered by a single token in the config file are:
                                  - user,
                                  - delete_repo,
                                  - public_repo,
                                  - repo,
                                  - repo_deployment,
                                  - repo:status,
                                  - read:repo_hook,
                                  - read:org,
                                  - read:public_key,
                                  - read:gpg_key.
                               Individual commands don't require all the scopes, so different tokens can be
                               used in the command line for different commands.
    -u, --user string           User name of the owner of token. (Needed with 'git clone'.)
    -v, --verbose               Toggle debug printouts.
```
Use "ghorgs [command] --help" for more information about a command.

### Dump command
Dumps the requested entities into a csv file.

Usage:
```
  ghorgs dump [flags]

  Flags:
    -b, --by string         Name of the entity field to use for sorting the result of the dump.
        If empty, default sort on GitHub is creation date.
    -e, --entities string   'all' for full dump or comma separated list of one or more of:
        users, repos, teams. (default "all")
    -h, --help              help for dump

  Global Flags:
    -d, --dry-run               Perform a dry run of the command without actually executing it in the end.
    -o, --organization string   Organizational account on GitHub analyzed.
    -t, --token string          Security token used on Github. Overrides the token from configuration file.
                                Required GitHub scopes covered by a single token in the config file are:
                                  - user,
                                  - delete_repo,
                                  - public_repo,
                                  - repo,
                                  - repo_deployment,
                                  - repo:status,
                                  - read:repo_hook,
                                  - read:org,
                                  - read:public_key,
                                  - read:gpg_key.
                               Individual commands don't require all the scopes, so different tokens can be
                               used in the command line for different commands.
    -u, --user string           User name of the owner of token. (Needed with 'git clone'.)
    -v, --verbose               Toggle debug printouts.
```

### Archive command
Remove GitHub repositories according to given criteria and archive to a given folder.
Uses v4 API for caching, v3 API for 'delete repository' operation.

Usage:
```
  ghorgs archive [flags]

  Flags:
    -b, --backup         Only backup the repositories. DO NOT REMOVE them.
    -h, --help           help for archive
    -n, --n int          Number of repositories to archive.
        * If --n is used together with --since, then the result is:
          "the number --n of repositories to archive --since point in time - whichever comes first."
        * If used alone, then the result is:
          "the least active number of repositories to archive".
        NOTE: It will be ignored if used with --repos.
        (default 1)
    -O, --out string     Output folder where archives of repositories are recorded. (default ".")
    -q, --quiet          DO NOT ask user for confirmation.(Use with care, e.g. in scripts where interaction is minimal or impossible.)
    -r, --repos string   Comma separated list of repositories to archive.
        * Name can contain alphanumeric and special characters '_', '.' and '-'.
        * If --repos is used with --since, then the result is:
          "archive the repositories from --repos list if they have been inactive --since this point in time.
        NOTE: --n will be ignored if used with --repos.
    -s, --since string   Remove repositories inactive since this date (YYYY-MM-DD).
        * If --since is used together with --n, then the result is:
          "the number --n of repositories to archive --since point in time - whichever comes first."
        * If --since is used together with --repos, then the result is:
          "archive the repositories from --repos list if they have been inactive --since this point in time".

  Global Flags:
    -d, --dry-run               Perform a dry run of the command without actually executing it in the end.
    -o, --organization string   Organizational account on GitHub analyzed.
    -t, --token string          Security token used on Github. Overrides the token from configuration file.
                                Required GitHub scopes covered by a single token in the config file are:
                                  - user,
                                  - delete_repo,
                                  - public_repo,
                                  - repo,
                                  - repo_deployment,
                                  - repo:status,
                                  - read:repo_hook,
                                  - read:org,
                                  - read:public_key,
                                  - read:gpg_key.
                               Individual commands don't require all the scopes, so different tokens can be
                               used in the command line for different commands.
    -u, --user string           User name of the owner of token. (Needed with 'git clone'.)
    -v, --verbose               Toggle debug printouts.
```

### Remove command
Remove GitHub users according to given criteria.
Uses v4 API for caching, v3 API for 'remove user' operation.

Usage:
```
  ghorgs remove [flags]

  Flags:
    -m, --MFA            Remove users without MFA set up.
    -a, --access         Remove users without access to any repository owned by the organization.
    -c, --company        Remove users without company affiliation.
    -h, --help           help for remove
    -q, --quiet          DO NOT ask user for confirmation. (Use with care, e.g. in scripts where interaction is minimal or impossible.)
    -r, --users string   Comma separated list of users to remove. Name can contain alphanumeric and special characters '_', '.' and '-'.

  Global Flags:
    -d, --dry-run               Perform a dry run of the command without actually executing it in the end.
    -o, --organization string   Organizational account on GitHub analyzed.
    -t, --token string          Security token used on Github. Overrides the token from configuration file.
                                Required GitHub scopes covered by a single token in the config file are:
                                  - user,
                                  - delete_repo,
                                  - public_repo,
                                  - repo,
                                  - repo_deployment,
                                  - repo:status,
                                  - read:repo_hook,
                                  - read:org,
                                  - read:public_key,
                                  - read:gpg_key.
                               Individual commands don't require all the scopes, so different tokens can be
                               used in the command line for different commands.
    -u, --user string           User name of the owner of token. (Needed with 'git clone'.)
    -v, --verbose               Toggle debug printouts.
```

## LICENSE
Currently the tool is proprietary.
