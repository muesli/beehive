# Configuration Encryption

Beehive's supports encrypting the configuration file using AES+GCM.

## Usage

To encrypt the configuration for the first time, simply start Beehive using a `crypto` URL for the configuration:

```
./beehive --config crypto://mysecret@$HOME/.config/beehive/beehive.conf`
```

You could also use the `BEEHIVE_CONFIG_PASSWORD` environment variable to define the password:

```
BEEHIVE_CONFIG_PASSWORD=mysecret ./beehive --config crypto://$HOME/.config/beehive/beehive.conf`
```

This will use the key `mysecret` to encrypt/decrypt the configuration file.

Once the configuration has been encrypted, it's no longer necessary to use a `crypto:` URL, Beehive will automatically detect it's encrypted.
That is, something like:

```
BEEHIVE_CONFIG_PASSWORD=mysecret beehive --config /path/to/config
```

Will happily detect and load an encrypted configuration file.

## Using user keyrings to store the password

You can also set the `BEEHIVE_CONFIG_PASSWORD_COMMAND` environment variable to automatically retrieve the password from an external command.
For example this environment setting will retrieve the password using the Secret Service API (Gnome Keyring):

```
export BEEHIVE_CONFIG_PASSWORD_COMMAND="secret-tool lookup user behive"
```

Something similar could be written to do it on macOS using Keychain and its `security(1)` CLI.

## Decrypting the configuration

Use `--decrypt` with a valid password:

```
beehive --decrypt --config crypto://mysecret@/path/to/config/file
```

or using an environment variable:

```
BEEHIVE_CONFIG_PASSWORD=mysecret beehive --decrypt --config crypto:///path/to/config/file
```

You can also use omit `--config` when using the default configuration path:

```
BEEHIVE_CONFIG_PASSWORD=mysecret beehive --decrypt
```

## Troubleshooting

```
FATA[0000] Error loading user config file /home/rubiojr/.config/beehive/beehive.conf. err: cipher: message authentication failed
```

Means the password used to decrypt the configuration file is not valid.

## Notes

The encrypted configuration file includes a 12 bytes header (`beehiveconf+`) that makes it possible to identify the file as an encrypted configuration file:

```
head -c 12 beehive-encrypted.conf
beehiveconf+
```
