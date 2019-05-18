# microservice-apps-management

[![Build](https://travis-ci.com/Microkubes/microservice-apps-management.svg?token=UB5yzsLHNSbtjSYrGbWf&branch=master)](https://travis-ci.com/Microkubes/microservice-apps-management)
[![Test Coverage](https://api.codeclimate.com/v1/badges/84383c8e579c181760ed/test_coverage)](https://codeclimate.com/repos/59c2524fbed4f6028e000bb6/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/84383c8e579c181760ed/maintainability)](https://codeclimate.com/repos/59c2524fbed4f6028e000bb6/maintainability)

Microservice Apps Management

## Prerequisite
Create a project directory. Set GOPATH enviroment variable to that project. Add $GOPATH/bin to the $PATH
```
export GOPATH=/path/to/project-workspace
export PATH=$GOPATH/bin:$PATH
```
Install goa and goagen:
```
cd $GOPATH
go get -u github.com/keitaroinc/goa/...
```

## Compile and run the service:
Clone the repo:
```
cd $GOPATH/src
git clone https://github.com/Microkubes/microservice-apps-management.git /path/to/project-workspace/src/github.com/Microkubes/microservice-apps-management
```
Be sure to use the full domain name and resource path here (compatible with ```go get```).


Then compile and run:
```
cd /path/to/project-workspace/src/github.com/Microkubes/microservice-apps-management
go build -o apps-management
./apps-management
```

## Change the design
If you change the design then you should regenerate the files. Run:
```
cd /path/to/project-workspace/src/github.com/Microkubes/microservice-apps-management
go generate
```
**NOTE:** If the above command does not update the generated code per the changes in the design,
then run ```goagen bootstrap```:

```bash
goagen bootstrap -d github.com/Microkubes/microservice-apps-management/design -o .
```


Also, recompile the service and start it again:
```
go build -o apps-management
./apps-management
```

## Other changes, not related to the design
For all other changes that are not related to the design just recompile the service and start it again:
```
cd $GOPATH/src/github.com/Microkubes/microservice-apps-management
go build -o apps-management
./apps-management
```

## Tests
For testing we use controller_test.go files which call the generated test helpers which package that data into HTTP requests and calls the actual controller functions. The test helpers retrieve the written responses, deserialize them, validate the generated data structures (against the validations written in the design) and make them available to the tests. Run:
```
go test -v
```

## Set up MongoDB
Create apps-management database with default username and password.
See: [Set up MongoDB](https://github.com/Microkubes/jormungandr-infrastructure#mongodb--v346-)
```
export MS_DBNAME=apps-management
./mongo/run.sh
```
Then install mgo package:
```
cd $GOPATH
go get gopkg.in/mgo.v2
```

# Docker Builds

First, create a directory for the shh keys:
```bash
mkdir keys
```

Find a key that you'll use to acceess Microkubes organization on github. Then copy the
private key to the directory you created above. The build would use this key to
access ```Microkubes/microservice-tools``` repository.

```bash
cp ~/.ssh/id_rsa keys/
```

**WARNING!** Make sure you don't commit or push this key to the repository!

To build the docker image of the microservice, run the following command:
```bash
docker build -t apps-management-microservice .
```

Also, you can build docker image using Makefile. Run the following command:
```bash
make run ARGS="-e API_GATEWAY_URL=http://192.168.1.10:8001 -e MONGO_URL=192.168.1.10:27017"
```

# Running the microservice

To run the apps-management microservice you'll need to set up some ENV variables:

 * **SERVICE_CONFIG_FILE** - Location of the configuration JSON file
 * **API_GATEWAY_URL** - Kong API url (default: http://localhost:8001)
 * **MONGO_URL** - Host IP(example: 192.168.1.10:27017)
 * **MS_USERNAME** - Mongo username (default: restapi)
 * **MS_PASSWORD** - Mongo password (default: restapi)
 * **MS_DBNAME** - Mongo database name (default: apps-management)

Run the docker image:
```bash
docker run apps-management-microservice
```

## Check if the service is self-registering on Kong Gateway

First make sure you have started Kong. See [Jormungandr Infrastructure](https://github.com/Microkubes/jormungandr-infrastructure)
on how to set up Kong locally.

If you have Kong admin endpoint running on http://localhost:8001 , you're good to go.
Build and run the service:
```bash
go build -o apps-management
./apps-management
```

To access the apps-management service, then instead of calling the service on :8080 port,
make the call to Kong:

```bash
curl -v --header "Host: apps-management.services.jormugandr.org" http://localhost:8000/apps/1
```

You should see a log on the terminal running the service that it received and handled the request.

## Running with the docker image

Assuming that you have Kong and it is availabel on your host (ports: 8001 - admin, and 8000 - proxy) and
you have build the service docker image (microservice-apps-management), then you need to pass
the Kong URL as an ENV variable to the docker run. This is needed because by default
the service will try http://localhost:8001 inside the container and won't be able to connect to kong.

Find your host IP using ```ifconfig``` or ```ip addr```.
Assuming your host IP is 192.168.1.10, then run:

```bash
docker run -ti -e API_GATEWAY_URL=http://192.168.1.10:8001 -e MONGO_URL=192.168.1.10:27017 apps-management-microservice
```

Also, you can build and run docker image using Makefile. Run:
```bash
make run ARGS="-e API_GATEWAY_URL=http://192.168.1.10:8001 -e MONGO_URL=192.168.1.10:27017"
```

If there are no errors, on a different terminal try calling Kong on port :8000

```bash
curl -v --header "Host: apps-management.services.jormugandr.org" http://localhost:8000/apps/1
```

You should see output (log) in the container running the service.

# Service configuration

The service loads the gateway configuration from a JSON file /run/secrets/microservice_apps_management_config.json. To change the path set the
**SERVICE_CONFIG_FILE** env var.
Here's an example of a JSON configuration file:

```json
{
  "name": "apps-management-microservice",
  "port": 8080,
  "virtual_host": "apps-management.services.jormugandr.org",
  "hosts": ["localhost", "apps-management.services.jormugandr.org"],
  "weight": 10,
  "slots": 100
}
```

Configuration properties:
 * **name** - ```"apps-management-microservice"``` - the name of the service, do not change this.
 * **port** - ```8080``` - port on which the microservice is running.
 * **virtual_host** - ```"apps-management.services.jormugandr.org"``` domain name of the service group/cluster. Don't change if not sure.
 * **hosts** - list of valid hosts. Used for proxying and load balancing of the incoming request. You need to have at least the **virtual_host** in the list.
 * **weight** - instance weight - use for load balancing.
 * **slots** - maximal number of service instances under ```"apps-management.services.jormugandr.org"```.

## Contributing

For contributing to this repository or its documentation, see the [Contributing guidelines](CONTRIBUTING.md).