config
======

CLI tool to fetch secret variables from AWS at service runtime.

## Usage

Download the latest binary from the releases page and run like so,

```
$ ./config <SSM parameter name>
```

**This will print the decrypted value to stdout.** Be aware of this when using it at service runtime,
and always opt for using this in conjunction with an environment variable, set at runtime.

```
ENV_VAR=$(./config <SSM parameter name>)
```

### Authorization

When including this tool in a service, the service must have the following IAM statements
included in its role policy.

```
...
{
    "Effect": "Allow",
    "Action": "ssm:GetParameters",
    "Resource": [
        "arn:aws:ssm:*:*:parameter/<SSM parameter name>"
        ( add other parameters as necessary )
    ]
},
{
    "Effect": "Allow",
    "Action": "kms:Decrypt",
    "Resource": "arn:aws:kms:*:*:key/<SSM key id>"
}
...
```

where `<SSM key id>` is the UUID of the encryption key you chose when pushing the value to SSM.
