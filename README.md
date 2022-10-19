### Clean Architecture (template)

###swagger
`http://localhost:80/docs/index.html`

### docker run

`sudo docker run -d -p 8081:80 -e JAEGER_AGENT_HOST=jaeger -e PGDATABASE=dbname -e PGPORT=5432 -e PGHOST=postgres -e PGPASSWORD=password -e PGUSER=user mensheninao/slurm:contact`
