# Beehive Docker Container

You can install beehive with docker if you'd prefer not to set up a working Go environment.

Make sure you have docker installed. See the [install instructions](https://docs.docker.com/get-started/).

## Building your own container image

Make sure to change into the directory that contains the `Dockerfile`.
It is sufficient to copy the file to an empty directory on your local filesystem.
You can also clone the repository with git and cd into the `docker` directory.

To start building the image, execute the following command:

    docker build -t beehive .

## Running beehive in a container

Once the image build process was finished, you can spin up the container with this command:

    docker run --name beehive -d -p 8181:8181 beehive

The `--name` parameter will give your container a name so that you can easily reference it with future commands.  
The `-d` flag instructs docker to detatch from the container once it was started.  
The `-p` parameter tells docker to map port 8181 on the host machine to port 8181 in the container.

By default, the beehive daemon will bind on port 8181 and expects to be accessed via `http://localhost:8181`.
You may want to specify your own parameters to the daemon by just appending them like so:

    docker run --name beehive -d -p 8181:8181 beehive -bind 0.0.0.0:8181 -canonicalurl https://mydomain.tld

You could then use a reverse proxy to dispatch requests to `https://mydomain.tld` to your container.

## Upgrading to a newer beehive version

To update to the latest beehive version (i.e. the `HEAD` this repository) you can simply 1) stop and remove your existing container and 2) re-build the image rebuild the image via

    docker build --no-cache -t beehive .

After the new image was built, you can then again execute the `docker run` command as explained above and profit from new features.
In order to not loose your configuration (bees & chains) each time you upgrade the container, you can mount a local beehive.conf file to the container which will then be used to persist any configuration on your local filesystem. To do that, simply adapt the `docker run` command as follows.

    docker run --name beehive -d -p 8181:8181 -v /path/to/beehive.conf:/app/beehive.conf beehive

### Useful commands

| Intend | Command |
| ------ | ------- |
| Stop the running container | `docker stop beehive` |
| Start the container (e.g. after reboot) | `docker start beehive` |
| Follow log output of the beehive container | `docker logs -f beehive` |
| Remove the container (e.g. to re-create a container with different parameters) | `docker rm beehive` |
| Remove the container and any persisted configuration (i.e. your bees and chains) | `docker rm beehive -v` |
| Remove old images that are not referenced by a (stopped or running) container | `docker image prune` |
