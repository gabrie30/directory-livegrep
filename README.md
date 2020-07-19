# directory-livegrep

Clone this directory and run `go build ./...`

Then run `./directory-livegrep <directory to index>`

Takes a single argument which is a path to a directory which has subdirectories of *bare* git repos. It will walk all subdirectories looking for bare git repositories, it will create a codeseach configuration and output a command for creating the livegrep index. It will also write a docker-compose.yaml file to use to run livegrep.

You must use git clone --mirror and only the https protocol (not ssh)