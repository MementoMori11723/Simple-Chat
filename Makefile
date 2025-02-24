run :
	@docker-compose -f config/compose.yml up --build $(ARGS)

stop :
	@docker-compose -f config/compose.yml down --remove-orphans
