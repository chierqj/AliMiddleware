contanier_id=`docker ps -a | grep chier | awk '{print $1}'`
docker stop $contanier_id
docker rm $contanier_id