# S3 bee

Uploads files to S3 compatible storage:

- AWS Signature Version 4

  - Amazon S3
  - Minio

- AWS Signature Version 2

  - Google Cloud Storage (Compatibility Mode)
  - Openstack Swift + Swift3 middleware
  - Ceph Object Gateway
  - Riak CS

See <https://github.com/minio/minio-go>

## Configuration

- endpoint: S3 host (Path to monitor (file or directory).
- access_key_id: S3 access key. Prefix the value with `env://` to retrieve the key from the environment. Example: `env://AWS_ACCESS_KEY_ID` or `ASDFWERSDF123WERF`
- secret_access_key: S3 secret access key. Prefix the value with `env://` to retrieve the key from the environment. Example: `env://AWS_SECRET_ACCESS_KEY` or `vASiudxSHReo4elkajsdklfu827389234sdfsdf`
- use_ssl: Defaults to `true`.
- region: AWS region. Defaults to `us-east-1`.

## Actions

### upload

- bucket: target bucket to upload the file.
- path: Source file path.
- object_path: Destination file path. Defaults to the source file name.

## Credits

AWS logo: <https://hu.wikipedia.org/wiki/Amazon_S3#/media/File:AWS_Simple_Icons_AWS_Cloud.svg>
