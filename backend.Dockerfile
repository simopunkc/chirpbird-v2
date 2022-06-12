FROM golang:alpine
WORKDIR /app/backend
COPY backend /app/backend
RUN mv .env.example .env
RUN go version
RUN go get ./ && go build -o chirpbird && go mod download
CMD ["./chirpbird"]