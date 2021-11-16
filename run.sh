docker-compose down
docker-compose build
docker-compose up --scale redis-sentinel=3 -d