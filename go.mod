module keykeeper

go 1.13

require (
	app v0.0.0
	config v0.0.0
	dbmgr v0.0.0
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	go.mongodb.org/mongo-driver v1.1.2
	gopkg.in/yaml.v2 v2.2.4
	gotest.tools v2.2.0+incompatible
	server v0.0.0

)

replace app => ./app

replace server => ./server

replace dbmgr => ./dbmgr

replace config => ./config
