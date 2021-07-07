# Golastic

Golastic is a web API offering full text search on a books collection with Elasticsearch.

## Run Elasticsearch

### Setup docker

```sh
docker network create elastic

docker pull docker.elastic.co/elasticsearch/elasticsearch:7.13.2
docker pull docker.elastic.co/kibana/kibana:7.13.2
```

### Run ES and Kibana images

```sh
docker run --name es01-test --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.13.2
docker run --name kib01-test --net elastic -p 5601:5601 -e "ELASTICSEARCH_HOSTS=http://es01-test:9200" docker.elastic.co/kibana/kibana:7.13.2
```

### Stop and clean up

```sh
docker stop es01-test
docker stop kib01-test
```

```sh
docker network rm elastic
docker rm es01-test
docker rm kib01-test
```
