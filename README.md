
# Libertyproxybeat

**Welcome to Libertyproxybeat - beat that utilises IBMJMXConnectorREST to retrieve JMX metrics for a Websphere Liberty Profile Instance.**

This is still **development version** and I expect changes based on feedback.

This beat retrieves JMX metrics from a running Liberty Profile instance and sends them to Logstash or Elasticsearch.
JMX metrics are requested via 'JMX Proxy Servlet' configured and enabled in Tomcat for HTTP listener. JMX Proxy Servlet is a lightweight proxy to get and set the Tomcat internals.

## Liberty Profile configuration

In order to enable the JMX url required by this beat the following features need to be added
```xml
        <feature>restConnector-1.0</feature>
        <feature>monitor-1.0</feature>
```

You must also setup a user with with the <administrator-role> tag

Example of request for **HeapMemoryUsage**
```
http://127.0.0.1:8443/IBMJMXConnectorREST/mbeans/WebSphere%3Atype%3DJvmStats/attributes?attribute=UsedMemory
```

## Getting Started with Libertyproxybeat

Ensure that this folder is at the following location:
`${GOPATH}/github.com/ninjasftw

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Build

To build the binary for Libertyproxybeat run the command below. This will generate a binary in the same directory with the name Libertyproxybeat.

```
make
```


### Run

To run Jmxproxybeat with debugging output enabled, run:

```
./Libertyproxybeat -c Libertyproxybeat.yml -e -d "*"
```

### Example JSON output
```
{
  "_index": "libertyproxybeat-2016.04.20",
  "_type": "jmx",
  "_id": "AVQ0FOGeegQ15caFDGZ7",
  "_score": null,
  "_source": {
    "@timestamp": "2016-04-20T14:31:03.385Z",
    "bean": {
      "attribute": "usedMemory",
      "hostname": "127.0.0.1:8443",
      "name": "WebSphere:type=JvmStats",
      "value": 81920
    },
    "beat": {
      "hostname": "localhost",
      "name": "localhost"
    },
    "type": "jmx"
  }
```

### Package - not complete yet

To be able to package Libertyproxybeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Libertyproxybeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/libertyproxybeat.template.json and etc/libertyproxybeat.asciidoc

```
make update
```


### Cleanup

To clean  Libertyproxybeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Jmxproxybeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/ninjasftw
cd ${GOPATH}/github.com/ninjasftw
git clone https://github.com/ninjasftw/libertyproxybeat
```

### Origins
This project is based upon a fork of https://github.com/radoondas/jmxproxybeat that was created for retrieving JMX metrics from a Tomcat instance
