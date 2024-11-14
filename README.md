# Go-Expert-Stress-test


# Build do projeto
docker build -t stresstest .

# Rodar o projeto exemple
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
