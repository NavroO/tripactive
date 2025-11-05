dev:
	go run cmd/api/main.go

migration-add:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Please provide migration name using 'make migration-add NAME=create_users_table'"; \
	else \
		migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME); \
	fi
