# Internal App Backend

## Installation
In order to start development on this project you need to follow the golang project structure and setup ( https://golang.org/doc/code.html ).

After setup please follow theese steps:

 1. Go to $GOPATH

 2. ```mkdir -p src/digiedu/iapp/main && cd src/digiedu/iapp/main```

 3. ```git init . && git remote add origin git@bitbucket.org:digital_education/iapp-backend.git```
 
 4. ```git pull origin master```

### Dependinces

1. Download go https://golang.org/dl/

2. ```sudo apt-get install build-essential``` ( for *make* )


## Build & Run
To build the project execute **```make build```** and this will create the binary files in $GOPATH/bin.
Also if you just want to test the application execute **```make run```**. For more informations about available tasks view *Makefile*

## Servers

**Broadcast me enabled**

iapp.digiedutm.com:9000

iapp.digiedutm.com:9001/info

**Broadcast me disabled**

iapp.digiedutm.com:9002

iapp.digiedutm.com:9003/info

### Hosting

Hosting is provided by amazon on eu-central-1b zone ( Frankfurt ) and EC2 instance name is ( DIGIEDU INTERNAL APP ) and id [i-098af5c8](https://eu-central-1.console.aws.amazon.com/ec2/v2/home?region=eu-central-1)

For machine access please contact cosmin.albulescu@digital-education.com
## Contribution
If you want to create a feature just create the branch from master with name convention **feature/XXX-some-description**. When the feature is done create a pull request in master.
This flow is also for **bug/XXX-some-description**.
Before create pull request make sure you pulled the master branch and all conflicts are merged.

Please request developer access to cosmin.albulescu@digital-education.com using the subject: **IAPP DEVELOPER**

### Bugs

If you found any bug please report it at https://bitbucket.org/digital_education/iapp-backend/issues