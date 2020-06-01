imgID=`docker images chier | grep chier | awk '{print $3}'`
docker login --username=15902906162 registry.cn-hangzhou.aliyuncs.com
docker tag $imgID registry.cn-hangzhou.aliyuncs.com/chier/middleware
docker push registry.cn-hangzhou.aliyuncs.com/chier/middleware