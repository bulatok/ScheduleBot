start:
	export $(grep -v '#.*' .env | xargs) && go run main.go