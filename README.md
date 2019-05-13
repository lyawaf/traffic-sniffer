# traffic-sniffer
## Build service
  ```
  docker-compose up
  cd cmd/waf
  go build .
  ./waf
  ```
## Drop database from mongo
1. Install mongodb client on your system:
(install mongo on ubuntu)[https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/]
2. Connect to database
```
mongo
```
3. Drop `streams` database:
```
use streams
db.dropDatabase()
```
