docker-up :
	cd docker && sudo docker-compose up -d

docker-down :
	cd docker && sudo docker-compose down

docker-build :
	sudo docker build -t horgh_replicator -f docker/Dockerfile_prod .

docker-run :
	sudo docker run -d -P horgh_replicator

### dev mode ###
start-dev : docker-up

restart-dev : docker-down docker-up

stop-dev : docker-down

### prod mode ###
build-prod : docker-build

start-prod : docker-run

stop-prod :
	sudo docker stop ${container} && sudo docker rm ${container}
