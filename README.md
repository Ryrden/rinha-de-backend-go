# Rinha de Back-end 2024 (Q1)

Código da minha aplicação para a Rinha de Back-end do @zanfranceschi.

Repositório da competição [aqui](https://github.com/zanfranceschi/rinha-de-backend-2024-q1).

## Tech

- Go 1.21
- Fiber (HTTP framework w/ fasthttp)
- pgx (SQL Driver)
- Fx (DI framework)
- Postgres
- Nginx

## Como rodar

1. Builde a imagem docker da aplicação com o comando `docker buildx build --platform linux/amd64 -t ryrden/rinha-de-backend-go:latest .`
2. Execute o script `./restart-container.ps1` Caso esteja no Linux, execute o comando `docker-compose up -d` na raiz do projeto.

## Minhas redes

- [LinkedIn](https://www.linkedin.com/in/ryan25/)
- [Website](https://ryrden.dev.br)
