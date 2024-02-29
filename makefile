run: templ build zip package deploy

templ:
	templ generate

build:
	GOARCH=amd64 GOOS=linux go build -o main main.go

zip:
	zip main.zip main sqls/* static/**/* design.settings.json

package:
	aws cloudformation package --template-file sam.yaml --output-template-file output-sam.yaml --s3-bucket puncher-bucket

deploy:
	aws cloudformation deploy --template-file output-sam.yaml --stack-name puncher-stack --capabilities CAPABILITY_IAM
