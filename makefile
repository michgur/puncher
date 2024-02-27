run: build zip package deploy

build:
	GOARCH=amd64 GOOS=linux go build -o main main.go

zip:
	zip main.zip main sqls/* static/* design.settings.json puncher.db

package:
	aws cloudformation package --template-file sam.yaml --output-template-file output-sam.yaml --s3-bucket puncher-bucket

deploy:
	aws cloudformation deploy --template-file output-sam.yaml --stack-name puncher-stack --capabilities CAPABILITY_IAM