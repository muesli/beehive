# Configuration Encryption

Beehive's supports encrypting the configuration file using AES+GCM.

## Usage

To encrypt the configuration for the first time, simply start Beehive using a `crypto` URL for the configuration:

```
./beehive --config crypto://x:mysecret@$HOME/.config/beehive/beehive.conf`
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
