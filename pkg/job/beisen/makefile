name=syncjob
version=vtemp0.0.4
devHost=192.168.200.20
devUser=ubuntu
repository=lowcode
dockerHost=qxcr.io

env:
#-- open go mod vendor --
	go mod vendor

docker-test: env
	cd ../../../ && \
	docker build -f ./pkg/job/beisen/Dockerfile -t  $(dockerHost)/$(repository)/$(name):$(version) .
	#docker push  $(dockerHost)/$(repository)/$(name):$(version)
