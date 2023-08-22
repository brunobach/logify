FROM golang:1.19 AS builder

# Defina o diretório de trabalho
WORKDIR /app

# Copie o código-fonte para o diretório de trabalho
COPY . .

# Compilar o aplicativo
RUN go build -o logify

# Construa a imagem final
FROM debian:buster-slim

# Copie o binário compilado do estágio anterior
COPY --from=builder /app/logify /usr/local/bin/
# Exponha a porta na qual o servidor da sua aplicação está ouvindo
EXPOSE 8080

# Comando para executar o aplicativo quando o container for iniciado
CMD ["logify"]