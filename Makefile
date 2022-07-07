.PHONY: migrate
migrate: migrate-billing migrate-orders migrate-products

.PHONY: migrate-products
migrate-products:
	goose -dir=./products_service/migrations postgres "user=postgres password=postgres dbname=homework_3_products_service sslmode=disable" up

.PHONY: migrate-orders
migrate-orders:
	goose -dir=./orders_service/migrations postgres "user=postgres password=postgres dbname=homework_3_orders_service sslmode=disable" up

.PHONY: migrate-billing
migrate-billing:
	goose -dir=./billing_service/migrations postgres "user=postgres password=postgres dbname=homework_3_billing_service sslmode=disable" up

.PHONY: proto-gen-billing
proto-gen-billing:
	protoc -I ./billing_service/proto \
	-I ${GOPATH}/src/github.com/googleapis \
	--go_out=./billing_service/pkg --go_opt=paths=source_relative \
	--go-grpc_out=./billing_service/pkg --go-grpc_opt=paths=source_relative \
	./billing_service/proto/api/invoices.proto
