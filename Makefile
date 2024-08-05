dev:
	cd web && npx vite build --outDir /home/rosta/kleofas3/dist/web && cd .. && go build -o dist/kleofas3 . && cd dist && ./kleofas3 serve
build:
	go build -o dist/kleofas3 .
