docker build -t forum .
echo "...........docker images.........................."
docker images
echo "....................................."
docker container run -p  8080:8081 --detach --name containersforum forum 
echo "...........docker ps -a.........................."
docker ps -a
echo "....................................."
docker exec -it containersforum /bin/bash
echo "....................................."
ls -l