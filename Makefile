run :
	@docker-compose -f config/compose.yml up -p simple-chat --build $(ARGS)
prod :
	@docker-compose -f config/deploy.yml up -p simple-chat --build -d
stop :
	@docker-compose -f config/compose.yml down -p simple-chat --remove-orphans
stop-prod :
	@docker-compose -f config/deploy.yml down -p simple-chat --remove-orphans
