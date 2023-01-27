# test-kafka
## Comandos para configurar o kafka
# Acessar o kafka
	docker exec -it {name_container} /bin/sh ou ir no docker app e clicar CLI
# Dentro criar o topico com o comando
	kafka-topics --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic weather-data --config retention.ms=1000
		em caso de sucesso
			Created topic weather-data
# Alterar Topico
	kafka-topics --alter --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic weather-data --config retention.ms=1000
# Deletar topico
	kafka-topics --delete  --bootstrap-server localhost:9092 --topic weather-data
# Listar topicops
	kafka-topics --list --bootstrap-server localhost:9092
# Enviando msg
	kafka-console-producer --broker-list localhost:9092 --topic weather-data
		{"temperature": 72, "humidity": 40, "wind_speed": 15}
# Consumindo msg
	kafka-console-consumer --bootstrap-server localhost:9092 --topic weather-data --from-beginning

# Find em todo pc
	find . -type f -exec grep "topics.sh" '{}' \; -print