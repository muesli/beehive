/*
 *    Copyright (C) 2017      Sergio Rubio
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Sergio Rubio <sergio@rubio.im>
 */

package s3bee

import (
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go"
	"github.com/muesli/beehive/bees"
)

type S3Bee struct {
	bees.Bee
	client *minio.Client
}

// Action triggers the action passed to it.
func (bee *S3Bee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "upload":
		bucket := ""
		action.Options.Bind("bucket", &bucket)

		path := ""
		action.Options.Bind("path", &path)

		objectPath := ""
		action.Options.Bind("object_path", &objectPath)

		if objectPath == "" {
			objectPath = filepath.Base(path)
		}

		_, err := bee.client.FPutObject(bucket, objectPath, path, minio.PutObjectOptions{ContentType: mime.TypeByExtension(filepath.Ext(path))})
		if err != nil {
			bee.LogFatal(err)
		}
	default:
		panic("Unknown action triggered in " + bee.Name() + ": " + action.Name)
	}

	return outs
}

func (bee *S3Bee) ReloadOptions(options bees.BeeOptions) {
	bee.SetOptions(options)

	endpoint := getConfigValue("endpoint", &options)

	var useSSL bool
	options.Bind("use_ssl", &useSSL)

	accessKeyID := getConfigValue("access_key_id", &options)
	secretAccessKey := getConfigValue("secret_access_key", &options)

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		panic(err)
		return
	}

	bee.client = client
}

func getConfigValue(key string, options *bees.BeeOptions) string {
	var value string
	options.Bind(key, &value)

	if strings.HasPrefix(value, "env://") {
		buf := strings.TrimPrefix(value, "env://")
		value = os.Getenv(string(buf))
	}

	return strings.TrimSpace(value)
}
