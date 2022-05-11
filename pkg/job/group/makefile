name=organizations
version=pushdatav0.0.1
devHost=192.168.200.20
devUser=ubuntu
repository=lowcode
dockerHost=qxcr.xyz

env:
#-- open go mod vendor --
	go mod vendor

docker-test: env
	cd ../../../ && \
	docker build -f ./pkg/job/refactor/Dockerfile -t  $(dockerHost)/$(repository)/$(name):$(version) .
	docker push  $(dockerHost)/$(repository)/$(name):$(version)
