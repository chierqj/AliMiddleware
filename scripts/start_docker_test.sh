contanier_id=`docker ps -a | grep chier | awk '{print $1}'`
docker stop $contanier_id
docker rm $contanier_id
docker run -td -p 3355:3355 -v /Users/chier/mygo/src/AliMiddleware/input:/root/input chier:latest