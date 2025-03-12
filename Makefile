run :
	@docker-compose -f config/compose.yml -p simple-chat up --build $(ARGS)
prod :
	@docker-compose -f config/deploy.yml -p simple-chat up --build -d
stop :
	@docker-compose -f config/compose.yml -p simple-chat down --remove-orphans
stop-prod :
	@docker-compose -f config/deploy.yml -p simple-chat down --remove-orphans
