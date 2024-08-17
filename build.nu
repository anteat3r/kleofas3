def build [] {
  npx vite build --outDir ../dist/web --emptyOutDir ../web
  go build -C .. -v -o dist/kleofas3 .; ./kleofas3 serve
}
