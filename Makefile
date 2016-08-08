DB ?= wildcare_development

migrate:
	mkdir -p tmp
	cp db/migrate.sql tmp/migrate.sql
	sed -i "s/DATABASE_NAME/$(DB)/g" tmp/migrate.sql
	mysql -u root < tmp/migrate.sql
