/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 */

package jenkinsbee

import (
    "github.com/muesli/beehive/modules"
)

type JenkinsBeeFactory struct {
}

func (factory *JenkinsBeeFactory) New (name, description string, options modules.BeeOptions) modules.ModuleInterface {
    bee := JenkinsBee {
        url: options.GetValue("url").(string),
    }
    bee.Module = modules.Module{name, factory.Name(), description}
    return &bee
}

func (factory *JenkinsBeeFactory) Name() string {
    return "jenkinsbee"
}

func (factory *JenkinsBeeFactory) Description() string {
    return "A bee that triggers and reads info from Jenkins-Builds"
}

func (factory *JenkinsBeeFactory) Options() []modules.BeeOptionDescriptor {
    opts := []modules.BeeOptionDescriptor{
        modules.BeeOptionDescriptor{
            Name: "url",
            Description: "The url the jenkins-installation is reachable at",
            Type: "string",
        },
    }
    return opts
}

func (factory *JenkinsBeeFactory) Events() []modules.EventDescriptor {
    events := []modules.EventDescriptor{
        modules.EventDescriptor{
            Namespace: factory.Name(),
            Name:   "statuschange",
            Description: "the status of a job has changed",
            Options: []modules.PlaceholderDescriptor {
                modules.PlaceholderDescriptor {
                    Name: "name",
                    Description: "Name of the job",
                    Type: "string",
                },
                modules.PlaceholderDescriptor {
                    Name: "status",
                    Description: "Current status of the job ('red' or 'blue')",
                    Type: "string",
                },
                modules.PlaceholderDescriptor {
                    Name: "url",
                    Description: "URL of the affected job",
                    Type : "string",
                },
            },
        },
    }
    return events
}

func (factory *JenkinsBeeFactory) Actions() []modules.ActionDescriptor {
    actions := []modules.ActionDescriptor{
        modules.ActionDescriptor{
            Namespace: factory.Name(),
            Name: "trigger",
            Description: "Trigger a build on this jenkins machine",
            Options: []modules.PlaceholderDescriptor{
                modules.PlaceholderDescriptor{
                    Name: "job",
                    Description: "Name of the job on which to trigger a build",
                    Type: "string",
                },
            },
        },
    }
    return actions
}

func init () {
    f := JenkinsBeeFactory{}
    modules.RegisterFactory(&f)
}
