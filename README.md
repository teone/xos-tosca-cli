# XOS TOSCA CLI

This is an helper around the new TOSCA engine. It let you easily submit TOSCA recipes.

_It is mainly an excgit use to play with Go and it is not in any way supported by the OpenCORD project_

## How to use it

Pull the docker container: `docker pull matteoscandolo/xos-cli`
or build it:
```
docker build -t matteoscandolo/xos-cli .
```

Run the docker container:
```
docker run --name xos-cli --net host --volume <tosca-recipe-folder>:/opt/tosca -d matteoscandolo/xos-cli
```

Use the CLI:
```
docker exec -it xos-cli ./xos-tosca-cli
```