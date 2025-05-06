# run if container is already running
# docker-compose exec k6 sh
# k6 run /scripts/k6-load-test.js

# Run if container is not running
docker-compose run --rm k6 k6 run /scripts/k6-load-test.js
