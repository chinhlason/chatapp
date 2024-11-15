How to run this app
1. start docker-container => cd to project dir, exec this command : "docker-compose up -d"  
2. start the server : cd to /server, exec command "make g-up" to migrate the db; "go mod download" to get all the dependencies, cd to /cmd and exec "go run ."
3. start the client : cd to /client, exec "npm i" to get all the dependencies, then "npm start" to start your app at :3000
