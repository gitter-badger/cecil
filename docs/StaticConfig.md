


# Customize MailGun Settings 

```
$ export MAILERDOMAIN=mg.yourdomain.co
$ export MAILERAPIKEY=key-<fill in here>
$ export MAILERPUBLICAPIKEY=pubkey-<fill in here>
```

You can find the mailer (Mailgun) API keys at [mailgun.com/app/account/security](https://mailgun.com/app/account/security)  For `MAILERAPIKEY` use the value in `Active API Key` and for `MAILERPUBLICAPIKEY` use `Email Validation Key`

# Customize keypair for signing JWT tokens 

Cecil uses JWT tokens in a few places to verify the authenticity of links sent to users via email.  In order for this to work, it needs an RSA keypair.

If not provided, it will generate a keypair on it's own and use it, and emit it on the console.  However, if you want to restart the `cecil` process and re-use the generated keypair, check the logs from the first run and capture the emitted private key into an environment variable named `CECIL_RSA_PRIVATE`:

```
$ export CECIL_RSA_PRIVATE='-----BEGIN RSA PRIVATE KEY----- MIIEowIBAAKCAQEAt ... -----END RSA PRIVATE KEY-----
```

# Event recording

To enable recording of all SQS events, set the `EVENTLOGDIR` environment variable and point it to a directory that is accessible to the process.

```
$ export EVENTLOGDIR=/path/to/directory
```