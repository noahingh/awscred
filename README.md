# AWSCRED

---

AWSCRED is a tool to generate a AWS session token and manage it easily. The core concept of this tool is that **it reflects on a new credentials file, not aws credentials**, it doesn't intrude the aws credentials file.


## Install

### Source code

```shell
$ git clone git@github.com:hanjunlee/awscred.git
$ mv awscred
$ go install
```

### Go 

```shell
$ go get -u github.com/hanjunlee/awscred
```

### Brew

TBD

## Usage

1. Run a new daemon - It creates a new daemon which reflect a session token on new credentials.

```shell
$ awscred run
```

2. Set up the profile - It configures the serial number and the duration. These values are used as options to generate a session token.

```shell
# set the profile enabled
$ awscred on PROFILE

# set up the configuration
$ awscred set --serial arn:aws:iam::XXXXXXXX:mfa/USER PROFILE

$ awscred info
NAME         ON       SERIAL                                    DURATION    EXPIRED
PROFILE      true     arn:aws:iam::XXXXXXXX:mfa/USER PROFILE    43200       0001-01-01T00:00:00Z
...
```

3. Generate a new session token 

```shell
# generate a new session token
$ awscred gen --code XXXXXX PROFILE

$ awscred info
NAME         ON       SERIAL                                    DURATION    EXPIRED
PROFILE      true     arn:aws:iam::XXXXXXXX:mfa/USER PROFILE    43200       2020-08-22T23:43:50Z (10.9h)
...
```

4. Modify the location of shared credentials file - By changing [the location of shared credentials file](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html#the-shared-credentials-file), `aws` command use the new credentials file. 

```shell
export AWS_SHARED_CREDENTIALS_FILE="~/.awscred/credentials"
```

5. Terminate 

```shell
# terminate the daemon
$ awscred terminate 

$ unset AWS_SHARED_CREDENTIALS_FILE
```

## LICENSE

[MIT License](./LICENSE)

## CHANGELOG

[Changelogs](./CHANGELOG)
