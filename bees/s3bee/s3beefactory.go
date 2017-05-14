/*
 *    Copyright (C) 2017 Sergio Rubio
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
	"github.com/muesli/beehive/bees"
)

// S3BeeFactory is a factory for S3Bees.
type S3BeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *S3BeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := S3Bee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *S3BeeFactory) ID() string {
	return "s3bee"
}

// Name returns the name of this Bee.
func (factory *S3BeeFactory) Name() string {
	return "S3Bee"
}

// Description returns the description of this Bee.
func (factory *S3BeeFactory) Description() string {
	return "Upload files to S3 compatible storage"
}

// Image returns the filename of an image for this Bee.
func (factory *S3BeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *S3BeeFactory) LogoColor() string {
	return "#4b4b4b"
}

// Options returns the options available to configure this Bee.
func (factory *S3BeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "endpoint",
			Description: "S3 compatible endpoint",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "access_key_id",
			Description: "Access Key ID",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "secret_access_key",
			Description: "Secret Access Key",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "use_ssl",
			Description: "Use SSL",
			Type:        "bool",
			Mandatory:   true,
			Default:     true,
		},
		{
			Name:        "region",
			Description: "S3 region",
			Type:        "string",
			Mandatory:   true,
			Default:     "us-east-1",
		},
	}
	return opts
}

func (factory *S3BeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "upload",
			Description: "Uploads a file to S3 compatible storage",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "bucket",
					Description: "Which bucket to upload the file to",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "path",
					Description: "Path to the file to upload",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "object_path",
					Description: "Target object path (defaults to the source file name)",
					Type:        "string",
					Mandatory:   false,
				},
			},
		},
	}
	return actions
}

func init() {
	f := S3BeeFactory{}
	bees.RegisterFactory(&f)
}
