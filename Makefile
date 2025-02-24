run :
	@docker-compose -f config/compose.yml up --build $(ARGS)
prod :
	@docker-compose -f config/deploy.yml up --build -d
stop :
	@docker-compose -f config/compose.yml down --remove-orphans
stop-prod :
	@docker-compose -f config/deploy.yml down --remove-orphans
