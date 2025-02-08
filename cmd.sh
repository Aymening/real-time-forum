### Install docker 
###curl -fsSL https://get.docker.com/rootless | sh

## build image
docker build -t forum .

## run container
docker run -p 8228:8228 -t forum

# delete images and 
# docker system prune -a -f

## list containers
# docker ps -a

## list images 
# docker images
