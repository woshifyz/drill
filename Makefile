all:
	cd ./fe && npm run build && cd ../backend/static && statik -src ../../fe/build && cd .. && GOOS=linux GOARCH=amd64 go build
mac:
	cd ./fe && npm run build && cd ../backend/static && statik -src ../../fe/build && cd .. &&  go build
