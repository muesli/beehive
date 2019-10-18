Beehive
=======

## Docker Installation

You can install beehive with docker if you'd prefer not to set up a working Go environment.

Make sure you have docker installed. See the [install instructions](https://docs.docker.com/engine/getstarted/step_one/).

### Using a prebuilt container image

The simplest way to set up beehive with docker is to simply pull a prebuilt image.

    docker pull fribbledom/beehive

### Building your own container image (skip if using a prebuilt image)

Make sure you're currently in the docker directory of the repository.
You can simply clone the repository with git and cd into the directory.

    git clone --recursive https://github.com/muesli/beehive.git
    cd beehive

Alternatively if you have the package installed with `go get` you can navigate
to `$GOPATH/src/github.com/muesli/beehive`

Either way once you're there you can build the docker container

    docker build -t beehive .

If you'd like to push the image up to docker.io so that you can use it elsewhere, you need
to namespace it with your docker.io username.

    docker build -t <username>/beehive .
    docker push <username>/beehive

## Running a container

Once you have the image on your machine, it's time to spin up an instance of it!
Of course if you built the container yourself without adding your username, leave out
the `<username>/` in this command.

    docker run --name beehive -d -p 8181:8181 <username>/beehive

The `--name` parameter will give your container a name so that you can easily reference it with future commands.
The `-d` flag specifies that the container should be run as a daemon.
The `-p` parameter tells docker to map port 8181 on the host machine to port 8181 in the container.
You can expose as many ports as is necessary. If you're running http server bees then you may need to
add additional flags so that those servers can be seen by your machine: `-p 8181:8181 -p 12345:12345 ... -p 34343:34343`

If ever you want to stop the container, run the following

    docker stop beehive

Then you can start it again with

    docker start beehive

#### Note

This container will store the `beehive.conf` file in a persistent volume.
As long as you use `docker stop` / `docker start` to stop/start the container
the configuration will persist.

If you'd like to have the container use an old config file, you can mount it as
a volume with `docker run`.

Suppose you had a config file stored in `/path/to/beehive.conf` then when running the container use

    docker run -d -p 8181:8181 -v /path/to/beehive.conf:/conf/beehive.conf <username>/beehive

This will tell docker to put your config file at `/conf/beehive.conf` within the container's filesystem.
Thus beehive will startup using your configuration file.
